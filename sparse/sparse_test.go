/*
 * Copyright (c) 2020 Juan Medina.
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in
 *  all copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 *  THE SOFTWARE.
 */

package sparse

import (
	"errors"
	"reflect"
	"testing"
)

func expectCapacityGrow(t *testing.T, sl *slice, capacity, grow int) {
	t.Helper()
	expect := capacity
	got := sl.capacity

	if got != expect {
		t.Fatalf("error capacity got %d, want %d", got, expect)
	}

	expect = grow
	got = sl.grow

	if got != expect {
		t.Fatalf("error grow got %d, want %d", got, expect)
	}
}

func expectFound(t *testing.T, sl *slice, ref ...interface{}) {
	t.Helper()

	for _, r := range ref {
		var expect error = nil
		_, got := sl.find(r)

		if !errors.Is(got, expect) {
			t.Fatalf("error find %v : got %v, want %v", r, got, expect)
		}
	}
}

func expectNotFound(t *testing.T, sl *slice, ref ...interface{}) {
	t.Helper()

	for _, r := range ref {
		var expect = ErrItemNotFound
		_, got := sl.find(r)

		if !errors.Is(got, expect) {
			t.Fatalf("error find %v : got %v, want %v", r, got, expect)
		}
	}
}

func expectSize(t *testing.T, sl *slice, size int) {
	t.Helper()
	expect := size
	got := sl.Size()

	if got != expect {
		t.Fatalf("error in size got %d, want %d", got, expect)
	}
}

func TestNewSlice(t *testing.T) {
	sl := NewSlice(5, 3).(*slice)

	expectCapacityGrow(t, sl, 5, 3)
}

func TestSlice_find(t *testing.T) {

	sl := NewSlice(3, 2).(*slice)

	sl.Add(1)

	expectFound(t, sl, 1)
	expectNotFound(t, sl, 2)
}

func TestSlice_Add(t *testing.T) {
	sl := NewSlice(3, 2).(*slice)

	sl.Add(1)
	sl.Add(2)
	sl.Add(3)

	expectCapacityGrow(t, sl, 3, 2)

	expectFound(t, sl, 1, 2, 3)
	expectNotFound(t, sl, 4, 5, 6)

	sl.Add(4)

	expectCapacityGrow(t, sl, 5, 2)

	expectFound(t, sl, 1, 2, 3, 4)
	expectNotFound(t, sl, 5, 6)

	sl.Add(5)

	expectCapacityGrow(t, sl, 5, 2)

	expectFound(t, sl, 1, 2, 3, 4, 5)
	expectNotFound(t, sl, 6)

	sl.Add(6)

	expectCapacityGrow(t, sl, 7, 2)

	expectFound(t, sl, 1, 2, 3, 4, 5, 6)
}

func TestSlice_Remove(t *testing.T) {
	sl := NewSlice(3, 2).(*slice)

	sl.Add(1)
	sl.Add(2)
	sl.Add(3)

	expectCapacityGrow(t, sl, 3, 2)
	expectFound(t, sl, 1, 2, 3)

	_ = sl.Remove(2)

	expectCapacityGrow(t, sl, 3, 2)
	expectFound(t, sl, 1, 3)
	expectNotFound(t, sl, 2)

	sl.Add(2)
	expectCapacityGrow(t, sl, 3, 2)
	expectFound(t, sl, 1, 2, 3)

	sl.Add(4)
	sl.Add(5)
	sl.Add(6)

	expectCapacityGrow(t, sl, 7, 2)
	expectFound(t, sl, 1, 2, 3, 4, 5, 6)

	_ = sl.Remove(3)
	_ = sl.Remove(5)

	sl.Add(7)
	sl.Add(8)
	sl.Add(9)
	sl.Add(10)

	expectCapacityGrow(t, sl, 9, 2)
	expectFound(t, sl, 1, 2, 4, 6, 7, 8, 9, 10)
	expectNotFound(t, sl, 3, 5)

	err := sl.Remove(11)

	if !errors.Is(err, ErrItemNotFound) {
		t.Fatalf("error find : got %v, want %v", err, ErrItemNotFound)
	}

	err = sl.Remove(9)

	if !errors.Is(err, nil) {
		t.Fatalf("error find : got %v, want %v", err, nil)
	}
}

func TestSlice_Iterator(t *testing.T) {
	sl := NewSlice(3, 2).(*slice)

	sl.Add(1)
	sl.Add(2)
	sl.Add(3)

	got := make([]int, 0)

	for it := sl.Iterator(); it != nil; it = it.Next() {
		got = append(got, it.Value().(int))
	}

	expect := []int{1, 2, 3}

	if !reflect.DeepEqual(got, expect) {
		t.Fatalf("error for each got %d, want %d", got, expect)
	}

	_ = sl.Remove(2)

	got = make([]int, 0)

	for it := sl.Iterator(); it != nil; it = it.Next() {
		got = append(got, it.Value().(int))
	}

	expect = []int{1, 3}

	if !reflect.DeepEqual(got, expect) {

	}

	_ = sl.Remove(1)
	_ = sl.Remove(3)

	got = make([]int, 0)

	for it := sl.Iterator(); it != nil; it = it.Next() {
		got = append(got, it.Value().(int))
	}

	expect = []int{}

	if !reflect.DeepEqual(got, expect) {

	}
}

func TestSlice_Size(t *testing.T) {

	sl := NewSlice(3, 2).(*slice)

	sl.Add(1)
	sl.Add(2)
	sl.Add(3)

	expectSize(t, sl, 3)

	_ = sl.Remove(2)
	expectSize(t, sl, 2)

	_ = sl.Remove(3)
	expectSize(t, sl, 1)

	_ = sl.Remove(1)
	expectSize(t, sl, 0)

	sl.Add(4)
	expectSize(t, sl, 1)
}

func TestSlice_Clear(t *testing.T) {
	sl := NewSlice(3, 2).(*slice)

	sl.Add(1)
	sl.Add(2)
	sl.Add(3)
	sl.Add(4)

	_ = sl.Remove(2)

	sl.Clear()

	expectSize(t, sl, 0)
}
