package app

import (
	"context"
	"fmt"

	"github.com/maxence-charriere/go-app/v6/pkg/errors"
)

type elem struct {
	attrs         map[string]string
	body          []UI
	ctx           context.Context
	ctxCancel     func()
	eventHandlers map[string]elemEventHandler
	jsvalue       Value
	parentElem    UI
	selfClosing   bool
	tag           string
}

func (e *elem) Kind() Kind {
	return HTML
}

func (e *elem) JSValue() Value {
	return e.jsvalue
}

func (e *elem) Mounted() bool {
	return e.ctx != nil && e.jsvalue != nil
}

func (e *elem) parent() UI {
	return e.parentElem
}

func (e *elem) setParent(p UI) {
	e.parentElem = p
}

func (e *elem) children() []UI {
	return e.body
}

func (e *elem) appendChild(c UI) {
	panic("not implemented")
}

func (e *elem) removeChild(c UI) {
	panic("not implemented")
}

func (e *elem) mount() error {
	panic("not implemented")
}

func (e *elem) update(n UI) error {
	panic("not implemented")
}

func (e *elem) dismount() {
	panic("not implemented")
}

func (e *elem) setAttr(k string, v interface{}) {
	if e.attrs == nil {
		e.attrs = make(map[string]string)
	}

	switch k {
	case "style":
		s := e.attrs[k] + toString(v) + ";"
		e.attrs[k] = s
		return
	}

	switch v := v.(type) {
	case bool:
		if !v {
			delete(e.attrs, k)
			return
		}
		e.attrs[k] = ""

	default:
		e.attrs[k] = toString(v)
	}
}

func (e *elem) setEventHandler(k string, h EventHandler) {
	if e.eventHandlers == nil {
		e.eventHandlers = make(map[string]elemEventHandler)
	}

	e.eventHandlers[k] = elemEventHandler{
		event: k,
		value: h,
	}
}

func (e *elem) setBody(self UI, body ...UI) {
	if e.selfClosing {
		panic(errors.New("setting html element body failed").
			Tag("reason", "self closing element can't have children").
			Tag("tag", e.tag),
		)
	}

	body = FilterUIElems(body...)
	for _, n := range body {
		n.setParent(self)
	}
	e.body = body
}

type elemEventHandler struct {
	event   string
	jsvalue Func
	value   EventHandler
}

func (h elemEventHandler) equal(o elemEventHandler) bool {
	return h.event == o.event &&
		fmt.Sprintf("%p", h.value) == fmt.Sprintf("%p", o.value)
}