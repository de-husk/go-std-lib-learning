package list

import (
	"testing"
)

func TestEmptyList(t *testing.T) {
	// zero value of a list is still a useable, empty list:
	var l List

	if l.Len() != 0 {
		t.Errorf("Expected 0 but got: %d", l.Len())
	}

	if l.Front() != nil {
		t.Errorf("Expected nil but got %v", l.Front())
	}

	l.PushBack(1)

	v := l.Front().Value
	if v != 1 {
		t.Errorf("Expected 1 but got: %d", v)
	}

	v = l.Back().Value
	if v != 1 {
		t.Errorf("Expected 1 but got: %d", v)
	}
}

func TestPushBack(t *testing.T) {
	l := New()

	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)

	if l.Len() != 3 {
		t.Errorf("Expected 3 but got: %d", l.Len())
	}

	v := l.Front().Value
	if v != 1 {
		t.Errorf("Expected 1 but got: %d", v)
	}

	v = l.Back().Value
	if v != 3 {
		t.Errorf("Expected 1 but got: %d", v)
	}
}

func TestInit(t *testing.T) {
	l := New()

	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)

	if l.Len() != 3 {
		t.Errorf("Expected 3 but got: %d", l.Len())
	}

	l.Init()

	if l.Len() != 0 {
		t.Errorf("Expected 0 but got: %d", l.Len())
	}

	v := l.Front()
	if v != nil {
		t.Errorf("Expected nil but got: %v", v)
	}

	v = l.Back()
	if v != nil {
		t.Errorf("Expected nil but got: %v", v)
	}
}

func TestRemove(t *testing.T) {
	l := New()

	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)

	if l.Len() != 3 {
		t.Errorf("Expected 3 but got: %d", l.Len())
	}

	v := l.Remove(&Element{Value: 2})
	if v != 2 {
		t.Errorf("Expected 2 but got: %s", v)
	}

	if l.Len() != 2 {
		t.Errorf("Expected 2 but got: %d", l.Len())
	}

	v = l.Remove(l.Front())
	if v != 1 {
		t.Errorf("Expected 1 but got: %s", v)
	}

	if l.Len() != 1 {
		t.Errorf("Expected 2 but got: %d", l.Len())
	}

	v = l.Remove(l.Back())
	if v != 3 {
		t.Errorf("Expected 1 but got: %s", v)
	}

	if l.Len() != 0 {
		t.Errorf("Expected 2 but got: %d", l.Len())
	}
}
