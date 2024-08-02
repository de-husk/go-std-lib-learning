package context

import (
	"fmt"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func TestWithCancel(t *testing.T) {
	defer goleak.VerifyNone(t)

	ctx, cancel := WithCancel(Background())

	if ctx.Err() != nil {
		t.Errorf("Err() should be nil, but was %v", ctx.Err())
	}

	if ctx.Done() == nil {
		t.Error("Done() channel should not be nil")
	}

	go func() {
		time.Sleep(1 * time.Millisecond)
		cancel()
	}()

	select {
	case <-time.After(10 * time.Millisecond):
		t.Error("Done should have been called before timeout")
	case <-ctx.Done():
		if ctx.Err() != ErrCanceled {
			t.Error("Err() should be ErrCanceled")
		}
	}
}

func TestWithCancelParent(t *testing.T) {
	defer goleak.VerifyNone(t)

	ctxParent, cancelParent := WithCancel(Background())

	ctx, _ := WithCancel(ctxParent)
	if ctx.Err() != nil {
		t.Errorf("Err() should be nil, but was %v", ctx.Err())
	}
	if ctx.Done() == nil {
		t.Error("Done() channel should not be nil")
	}

	ctx2, _ := WithCancel(ctxParent)
	if ctx2.Err() != nil {
		t.Errorf("Err() should be nil, but was %v", ctx.Err())
	}
	if ctx2.Done() == nil {
		t.Error("Done() channel should not be nil")
	}

	cancelParent()

	// Sleep to let Done cancelling propagate to children
	time.Sleep(1 * time.Millisecond)

	if ctxParent.Err() != ErrCanceled {
		t.Error("Err() should be ErrCanceled")
	}
	if ctx.Err() != ErrCanceled {
		t.Error("Err() should be ErrCanceled")
	}
	if ctx2.Err() != ErrCanceled {
		t.Error("Err() should be ErrCanceled")
	}
}

func TestWithValue(t *testing.T) {
	defer goleak.VerifyNone(t)

	// KV pair exists in current ctx:
	type testContextKey string

	ctx := WithValue(Background(), testContextKey("foo_key"), "bar_val")

	v := ctx.Value(testContextKey("foo_key"))
	if v != "bar_val" {
		t.Errorf("Expected %v, but got %v", "bar_val", v)
	}

	// KV pair exists in parent ctx:
	ctx, cancel := WithCancel(ctx)
	defer cancel()

	ctx2 := WithValue(ctx, testContextKey("foo_key2"), "bar_val2")

	v = ctx.Value(testContextKey("foo_key2"))
	if v != nil {
		t.Errorf("Expected nik, but got %v", v)
	}

	v = ctx2.Value(testContextKey("foo_key"))
	if v != "bar_val" {
		t.Errorf("Expected %v, but got %v", "bar_val", v)
	}

	v = ctx2.Value(testContextKey("foo_key2"))
	if v != "bar_val2" {
		t.Errorf("Expected %v, but got %v", "bar_val2", v)
	}
}

func TestWithTimeout(t *testing.T) {
	defer goleak.VerifyNone(t)

	var neverReady chan struct{}

	ctx, cancel := WithTimeout(Background(), 1*time.Millisecond)
	defer cancel()

	select {
	case <-neverReady:
		t.Error("should never happen")
	case <-ctx.Done():
		fmt.Println(ctx.Err())

		if ctx.Err() != ErrDeadlineExceeded {
			t.Errorf("Expected %v, but got %v", ErrDeadlineExceeded, ctx.Err())
		}
	}

	// Test with a timer that is a lot longer than this timeout
	tt := time.NewTimer(10 * time.Millisecond)

	ctx, cancel = WithTimeout(Background(), 1*time.Millisecond)
	defer cancel()

	select {
	case <-tt.C:
		t.Error("should never happen")
	case <-ctx.Done():
		fmt.Println(ctx.Err())

		if ctx.Err() != ErrDeadlineExceeded {
			t.Errorf("Expected %v, but got %v", ErrDeadlineExceeded, ctx.Err())
		}
	}
}

func TestWithDeadline(t *testing.T) {
	defer goleak.VerifyNone(t)

	// Test with deadline in the past:
	var neverReady chan struct{}

	ctx, cancel := WithDeadline(Background(), time.Now().Add(-1*time.Hour))
	defer cancel()

	select {
	case <-neverReady:
		t.Error("should never happen")
	case <-ctx.Done():
		fmt.Println(ctx.Err())

		if ctx.Err() != ErrDeadlineExceeded {
			t.Errorf("Expected %v, but got %v", ErrDeadlineExceeded, ctx.Err())
		}
	}

	// Test with deadline in the future:
	ctx, cancel = WithDeadline(Background(), time.Now().Add(10*time.Millisecond))
	defer cancel()

	select {
	case <-neverReady:
		t.Error("should never happen")
	case <-ctx.Done():
		fmt.Println(ctx.Err())

		if ctx.Err() != ErrDeadlineExceeded {
			t.Errorf("Expected %v, but got %v", ErrDeadlineExceeded, ctx.Err())
		}
	}

}
