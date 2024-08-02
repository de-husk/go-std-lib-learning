package ring

import (
	"fmt"
	"testing"
)

func TestEmpty(t *testing.T) {
	var r Ring

	if r.Len() != 1 {
		t.Errorf("Zero value ring should have 1 element, but got: %d", r.Len())
	}

	if r.Value != nil {
		t.Errorf("Expected Value to be nil, but got: %s", r.Value)
	}
}

func TestNew(t *testing.T) {
	r := New(5)

	if r.Len() != 5 {
		t.Errorf("Ring should have 5 elements, but got: %d", r.Len())
	}

	for i := 0; i < 5; i++ {
		if r == nil {
			t.Errorf("Expected r to be non nil")
			break
		}

		if r.Value != nil {
			t.Errorf("Expected Value to be nil, but got: %s", r.Value)
		}
		r = r.Next()
	}

	// A Ring with 1 element will be a loop of 1
	// with the next and prev pointing to itself:
	r = New(1)

	if r.Len() != 1 {
		t.Errorf("Ring should have 1 elements, but got: %d", r.Len())
	}

	if r.Next() != r {
		t.Error("r.Next() is not equal to r")
	}

	if r.Prev() != r {
		t.Error("r.Prev() is not equal to r")
	}
}

func TestDo(t *testing.T) {
	r := New(11)

	if r.Len() != 11 {
		t.Errorf("Ring should have 11 elements, but got: %d", r.Len())
	}

	count := 0
	r.Do(func(p any) {
		count++
	})

	if count != 11 {
		t.Errorf("Expected count to be 11, but got: %d", count)
	}
}

func TestNext(t *testing.T) {
	r := New(4)
	for i := 0; i < r.Len(); i++ {
		r.Value = i
		r = r.Next()
	}

	fmt.Println(r)

	if r.Value != 0 {
		t.Errorf("Expected %d -  got %d", 0, r.Value)
	}

	if r.Next().Value != 1 {
		t.Errorf("Expected %d -  got %d", 1, r.Next().Value)
	}

	count := 0
	r.Do(func(p any) {
		count += p.(int)
	})

	expected := ((r.Len() - 1) * (r.Len() - 1 + 1)) / 2

	if count != expected {
		t.Errorf("Expected %d -  got %d", expected, count)
	}
}

func TestPrev(t *testing.T) {
	r := New(42)
	for i := 0; i < r.Len(); i++ {
		r.Value = i
		r = r.Next()
	}

	if r.Value != 0 {
		t.Errorf("Expected %d -  got %d", 0, r.Value)
	}

	if r.Prev().Value != 41 {
		t.Errorf("Expected %d -  got %d", 41, r.Prev().Value)
	}

	count := 0
	r.Do(func(p any) {
		count += p.(int)
	})

	expected := ((r.Len() - 1) * (r.Len() - 1 + 1)) / 2

	if count != expected {
		t.Errorf("Expected %d -  got %d", expected, count)
	}
}
