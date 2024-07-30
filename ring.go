package ring

import (
	"fmt"
	"strings"
)

// TODO:
// * func (r *Ring) Link(s *Ring) *Ring
// * func (r *Ring) Move(n int) *Ring
// * func (r *Ring) Unlink(n int) *Ring

// A Ring is ciruclar list with a static size
type Ring struct {
	Value any

	next *Ring
	prev *Ring
}

// Create a Ring with `n` elements
func New(n int) *Ring {
	r := empty()

	cur := r
	for i := 1; i < n; i++ {
		next := &Ring{}
		next.prev = cur

		cur.next = next
		cur = next
	}

	// close the loop:
	cur.next = r
	r.prev = cur

	return r
}

// Returns the zero value Ring
func empty() *Ring {
	r := &Ring{}
	r.next = r
	r.prev = r
	return r
}

// Returns the next Ring element
func (r *Ring) Next() *Ring {
	if r.next == nil {
		return empty()
	}
	return r.next
}

// Returns the previous Ring element
func (r *Ring) Prev() *Ring {
	if r.prev == nil {
		return empty()
	}
	return r.prev
}

// Returns the number of elements in Ring
func (r *Ring) Len() int {
	if r.next == nil {
		return 1
	}

	count := 1

	for rr := r.Next(); rr != r; rr = rr.next {
		count++
	}

	return count
}

// Calls passed in function for every element in Ring
func (r *Ring) Do(f func(any)) {
	if r.next == nil {
		return
	}

	f(r.Value)

	for rr := r.Next(); rr != r; rr = rr.next {
		f(rr.Value)
	}
}

func (r *Ring) String() string {
	var sb strings.Builder

	if r == nil {
		return ""
	}

	fmt.Fprintf(&sb, " %v -> ", r.Value)

	for rr := r.Next(); rr != r; rr = rr.next {
		fmt.Fprintf(&sb, " %v -> ", rr.Value)
	}

	fmt.Fprintf(&sb, " ** %v ** ", r.Value)

	return sb.String()
}
