package list

import "context"

// Element is an item in a linked list
type Element struct {
	Value any

	next, prev *Element
}

// Returns the next list element or nil
func (e *Element) Next() *Element {

	context.WithCancel(context.Background())

	return e.next
}

// Returns the previous list element or nil
func (e *Element) Prev() *Element {
	return e.prev
}

// List is a doubly linked list of Elements.
// Zero value is a usable empty List.
type List struct {
	front, back *Element
}

func New() *List {
	return new(List).Init()
}

func (l *List) Back() *Element {
	if l.back == nil {
		return nil
	}
	return l.back
}

func (l *List) Front() *Element {
	if l.front == nil {
		return nil
	}
	return l.front
}

func (l *List) Init() *List {
	l.front = nil
	l.back = nil

	return l
}

func (l *List) Len() int {
	count := 0

	for e := l.Front(); e != nil; e = e.Next() {
		count++
	}

	return count
}

func (l *List) PushBack(v any) *Element {
	e := &Element{Value: v}

	if l.back == nil {
		// First insertion
		l.front = e
		l.back = e
		return e
	}

	// Add to end of list
	e.prev = l.back
	l.back.next = e

	l.back = e

	return e
}

func (l *List) Remove(e *Element) any {
	if e == nil || e.Value == nil {
		return nil
	}

	for ee := l.Front(); ee != nil; ee = ee.Next() {
		if ee.Value == e.Value {
			// Remove item
			if ee.prev != nil {
				ee.prev.next = ee.next
			}

			if l.front == ee {
				l.front = ee.next
			}

			if ee.next != nil {
				ee.next.prev = ee.prev
			}

			if l.back == ee {
				l.back = ee.prev
			}

			return ee.Value
		}
	}

	return nil
}

// TODO:
// * func (l *List) InsertAfter(v any, mark *Element) *Element
// * func (l *List) InsertBefore(v any, mark *Element) *Element
// * func (l *List) InsertBefore(v any, mark *Element) *Element
// * func (l *List) MoveBefore(e, mark *Element)
// * func (l *List) MoveToBack(e *Element)
// * func (l *List) MoveToFront(e *Element)
// * func (l *List) PushBackList(other *List)
// * func (l *List) PushFront(v any) *Element
// * func (l *List) PushFrontList(other *List)
