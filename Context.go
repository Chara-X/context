package context

import (
	"context"
	"time"
)

type Context interface {
	Value(key interface{}) interface{}
	Done() <-chan struct{}
}
type backgroundCtx struct{ c context.Context }

func Background() Context {
	if Reference {
		return &backgroundCtx{c: context.Background()}
	}
	return &backgroundCtx{}
}
func (c backgroundCtx) Done() <-chan struct{} {
	if Reference {
		return c.c.Done()
	}
	return nil
}
func (c backgroundCtx) Value(key interface{}) interface{} {
	if Reference {
		return c.c.Value(key)
	}
	return nil
}

type valueCtx struct {
	c context.Context
	Context
	key, val interface{}
}

func WithValue(parent Context, key, val interface{}) Context {
	if Reference {
		return &valueCtx{c: context.WithValue(parent.(context.Context), key, val)}
	}
	return &valueCtx{Context: parent, key: key, val: val}
}
func (c *valueCtx) Value(key interface{}) interface{} {
	if Reference {
		return c.c.Value(key)
	}
	if c.key == key {
		return c.val
	}
	return c.Context.Value(key)
}

type cancelCtx struct {
	c context.Context
	Context
	done chan struct{}
}

func WithCancel(parent Context, timeout time.Duration) (ctx Context, cancel func()) {
	if Reference {
		var c, cancel = context.WithCancel(parent.(context.Context))
		return &cancelCtx{c: c}, cancel
	}
	var c = &cancelCtx{Context: parent, done: make(chan struct{})}
	go func() {
		select {
		case <-parent.Done():
			close(c.done)
		case <-time.After(timeout):
			close(c.done)
		case <-c.done:
		}
	}()
	return c, func() { close(c.done) }
}
func (c *cancelCtx) Done() <-chan struct{} {
	if Reference {
		return c.c.Done()
	}
	return c.done
}
