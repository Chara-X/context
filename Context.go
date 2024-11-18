package context

import "context"

type Context interface {
	Done() <-chan struct{}
	Value(key any) any
}
type backgroundCtx struct {
	ctx context.Context
}

func Background() Context {
	if Reference {
		return context.Background()
	}
	return &backgroundCtx{ctx: context.Background()}
}
func (c backgroundCtx) Done() <-chan struct{} {
	if Reference {
		return c.ctx.Done()
	}
	return nil
}
func (c backgroundCtx) Value(key any) any {
	if Reference {
		return c.ctx.Value(key)
	}
	return nil
}

type valueCtx struct {
	Context
	c        context.Context
	key, val any
}

func WithValue(parent Context, key, val any) Context {
	if Reference {
		return context.WithValue(parent.(context.Context), key, val)
	}
	return &valueCtx{Context: parent, key: key, val: val}
}
func (c *valueCtx) Value(key any) any {
	if Reference {
		return c.c.Value(key)
	}
	if c.key == key {
		return c.val
	}
	return c.Context.Value(key)
}

type cancelCtx struct {
	Context
	c    context.Context
	done chan struct{}
}

func WithCancel(parent Context) (Context, func()) {
	if Reference {
		return context.WithCancel(parent.(context.Context))
	}
	var c = &cancelCtx{c: parent.(context.Context)}
	go func() {
		select {
		case <-parent.Done():
			close(c.done)
		case <-c.Done():
		}
	}()
	return c, func() { close(c.done) }
}
