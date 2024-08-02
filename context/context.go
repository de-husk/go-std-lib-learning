package context

import (
	"errors"
	"reflect"
	"sync"
	"time"
)

var ErrCanceled = errors.New("context canceled")

var ErrDeadlineExceeded error = errors.New("deadline exceeded")

type Context interface {
	// Deadline returns the time when the work tied to this context should be cancelled.
	Deadline() (deadline time.Time, ok bool)

	// Done returns a closed channel when the work tied to this context should be cancelled.
	Done() <-chan struct{}

	// Err returns nil when Done isn't closed yet, else it returns a non-nil error
	Err() error

	// Value returns the value associated with this context for Key, or nil if no associated context for Key
	Value(key any) any
}

// CancelFunc closes the returned Done channel when called the first time.
// Subsequent calls do nothing.
type CancelFunc func()

// Background returns a non-nil. empty Context.
// It is never cancelled, has no values, and has no deadline.
// It is typically used by the main function, initialization, and tests,
// and as the top-level Context for incoming requests.

type backgroundCtx struct{}

func (backgroundCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (backgroundCtx) Done() <-chan struct{} {
	return nil
}

func (backgroundCtx) Err() error {
	return nil
}

func (backgroundCtx) Value(key any) any {
	return nil
}

func Background() Context {
	return backgroundCtx{}
}

func TODO() Context {
	return backgroundCtx{}
}

type cancelCtx struct {
	parent Context

	mu sync.Mutex

	done chan struct{}
}

func (c *cancelCtx) Done() <-chan struct{} {
	return c.done
}

func (c *cancelCtx) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

func (c *cancelCtx) Err() error {
	select {
	case _, ok := <-c.done:
		if !ok {
			return ErrCanceled
		}
	default:
	}

	return nil
}

func (c *cancelCtx) Value(key any) any {
	// cancelCtx doesnt store any kvs but the parents might:
	return c.parent.Value(key)
}

func (c *cancelCtx) close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case <-c.done:
	default:
		close(c.done)
	}
}

// WithCancel closes Done channel when CancelFunc is called,
// or when the parent context's Done channel is closed, whichever happens first
func WithCancel(parent Context) (Context, CancelFunc) {
	if parent == nil {
		panic("cannot create cancel context with a nil parent")
	}

	ctx := &cancelCtx{
		parent: parent,
		done:   make(chan struct{}),
	}

	// TODO: Can we avoid running a go routine per cancellable context?
	go func() {
		select {
		case <-ctx.parent.Done():
			ctx.close()
		case <-ctx.Done():
			// avoids go routine leak
		}
	}()

	return ctx, ctx.close
}

type valueCtx struct {
	// NOTE: We embed the interface here, so that we dont have to create dummy wrapper implementations of
	// the interface methods (Done(), Err(), etc) that just use `v.parent.Done()` etc
	Context

	key, value any
}

func (v *valueCtx) Value(key any) any {
	if key == nil {
		panic("key cannot be nil")
	}

	if v.key == key {
		return v.value
	}

	// look in parent:
	return v.Context.Value(key)
}

// WithValue returns a copy of parent with key associated with val.
// The provided key must be comparable and should not be of type string or any other built-in type to avoid collisions between packages using context. Users of WithValue should define their own types for keys. To avoid allocating when assigning to an interface{}, context keys often have concrete type struct{}. Alternatively, exported context key variables' static type should be a pointer or interface.
func WithValue(parent Context, key, val any) Context {
	if parent == nil {
		panic("cannot create a new context kv with a nil parent")
	}

	if key == nil {
		panic("key cannot be nil")
	}

	if !reflect.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}

	return &valueCtx{parent, key, val}
}

type deadlineCtx struct {
	cancelCtx

	err error // locked under cancelCtx.mu

	deadline time.Time
}

func (d *deadlineCtx) Deadline() (deadline time.Time, ok bool) {
	return d.deadline, true
}

func (d *deadlineCtx) Err() error {
	d.mu.Lock()
	err := d.err
	d.mu.Unlock()

	return err
}

func (d *deadlineCtx) close() {
	d.closeErr(ErrCanceled)
}

func (d *deadlineCtx) closeErr(err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	select {
	case <-d.done:
	default:
		close(d.done)
		d.err = err
	}
}

func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	if parent == nil {
		panic("cannot create a new timeout context with a nil parent")
	}

	d := time.Now().Add(timeout)

	ctx := &deadlineCtx{
		cancelCtx{
			parent: parent,
			done:   make(chan struct{}),
		},
		nil,
		d,
	}

	t := time.NewTimer(timeout)

	// TODO: Can we avoid running a go routine per cancellable context?
	go func() {
		select {
		case <-ctx.parent.Done():
			ctx.closeErr(ErrCanceled)
		case <-t.C:
			ctx.closeErr(ErrDeadlineExceeded)
		case <-ctx.Done():
		}
	}()

	return ctx, ctx.close
}

// WithDeadline closes Done channel when deadline expires,
// or when the parent context's Done channel is closed, whichever happens first
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
	if parent == nil {
		panic("cannot create a new timeout context with a nil parent")
	}

	n := time.Now()

	if d.Before(n) {
		// d has already expired, return a cancelled context
		ctx := &deadlineCtx{
			cancelCtx{
				parent: parent,
				done:   make(chan struct{}),
			},
			nil,
			d,
		}
		close(ctx.done)
		ctx.err = ErrDeadlineExceeded
		return ctx, ctx.close
	}

	tt := d.Sub(n)

	ctx := &deadlineCtx{
		cancelCtx{
			parent: parent,
			done:   make(chan struct{}),
		},
		nil,
		d,
	}

	t := time.NewTimer(tt)

	// TODO: Can we avoid running a go routine per cancellable context?
	go func() {
		select {
		case <-ctx.parent.Done():
			ctx.closeErr(ErrCanceled)
		case <-t.C:
			ctx.closeErr(ErrDeadlineExceeded)
		case <-ctx.Done():
		}
	}()

	return ctx, ctx.close
}
