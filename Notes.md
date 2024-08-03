# Notes

## Context

### Efficiently propagating cancels from parents to children
 At first, my implementation for cancellable contexts each had a single goroutine running that looked like this for contexts made via WithCancel():
```
go func() {
    select {
    case <-ctx.parent.Done():
        ctx.close()
    case <-ctx.Done():
        // avoids go routine leak
    }
}()
```

where `ctx` type looked like:
```
type cancelCtx struct {
	parent Context

	mu sync.Mutex

	done chan struct{}
}
...

type Context interface {...}
```

So, each WithCancel() context had a pointer up to it's direct parent, and when the parent context was cancelled, all children would cancel their contexts through that goroutine seen above.

However, after I read the std lib context code, I noticed that the real cancelCtx looks like this:
```
// A cancelCtx can be canceled. When canceled, it also cancels any children
// that implement canceler.
type cancelCtx struct {
	Context

	mu       sync.Mutex            // protects following fields
	done     atomic.Value          // of chan struct{}, created lazily, closed by first cancel call
	children map[canceler]struct{} // set to nil by the first cancel call
	err      error                 // set to non-nil by the first cancel call
	cause    error                 // set to non-nil by the first cancel call
}
```

Very cool to see that the cancellable contexts, not only store a pointer to their parent (the embedded `cancelCtx.Context`), but they also store a list of their direct children contexts in `cancelCtx.children`. 

This allows them to entirely SKIP having the extra goroutine running per `cancelCtx`, just sitting there waiting to see if the parent gets cancelled.

Each `cancelCtx` just direcly calls `child.cancel()` for each child  when cancelling itself in `cancel()`:
```
for child := range c.children {
    child.cancel(false, err, cause)
}
```

Very cool inversion of control from what I naturally thought of. This also forces a strict ordering of cancelation which is really nice.


### RWMutexes are slow

Mutexes are the only locks used in the context std lib file.

After doing some digging it seems RWMutex [is slower than Mutex](https://github.com/golang/go/issues/17973). Im not sure if thats the only reason, it probably has to do with the usecase too, eg: there are probably more writes to cancelCtx.cause/err than reads.


## List

* Doubly linked list in std lib is implemented as a ring under the hood, which makes the code cleaner and makes it so we only need to store one Element pointer in List instead of 2 (first and last nodes)