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

// Package sparse provides the creation of sparse.Slice via sparse.NewSlice
//
// Sparse items are design to automatically grow and reuse memory so they try to
// minimize the GC pauses
package sparse

import (
	"errors"
	"sort"
)

var (
	// ErrItemNotFound is the error when we could not find an item
	ErrItemNotFound = errors.New("item not found")
)

// Iterator allows to move trough the items in sparse with a for-loop
type Iterator interface {
	// Next returns the next element, nil if we haven't got more
	Next() Iterator
	// Value returns the current item value
	Value() interface{}
}

// Slice is an slice that contains interfaces, we need to specify a initial capacity
//
// When remove at item it marked as free so the next added item will reuse the free slot
//
// Is dynamically sized, so first grow will double it current capacity, consecutive grows will be 25% of the current
// capacity
//
// The implementation use go standard slices so when grow it probably do not need to be reallocated
type Slice interface {
	// Add a new item to the slice
	Add(ref interface{})
	// Remove a item in the slice
	Remove(ref interface{}) error
	// Clear al the items in the slice
	Clear()
	// Size return the number of items in this slice
	Size() int
	// Iterator returns a new sparse.Iterator for sparse.Slice
	Iterator() Iterator
	// AssureCapacity grows the Slice until it has at least the desired capacity
	AssureCapacity(capacity int)
	// Sort a sparse.Slice in-place using a less function
	Sort(less func(a interface{}, b interface{}) bool)
	// Copy makes a copy of this Slice into dest
	Copy(dest Slice)
	// Replace replace dest contents with this Slice contents
	Replace(dest Slice)
}

type item struct {
	ref   interface{}
	valid bool
}

type slice struct {
	capacity int
	grow     int
	items    []item
	size     int
}

type sliceIterator struct {
	data    slice
	current int
}

func (si *sliceIterator) Next() Iterator {
	for i := si.current + 1; i < len(si.data.items); i++ {
		item := si.data.items[i]
		if item.valid {
			si.current = i
			return si
		}
	}
	return nil
}

func (si *sliceIterator) Value() interface{} {
	return si.data.items[si.current].ref
}

func (ss *slice) Add(ref interface{}) {
	for i, si := range ss.items {
		if !si.valid {
			ss.items[i].ref = ref
			ss.items[i].valid = true
			ss.size++
			return
		}
	}

	ss.growCapacity()
	ss.items[ss.size].ref = ref
	ss.items[ss.size].valid = true
	ss.size++
}

func (ss *slice) Remove(ref interface{}) error {
	if i, err := ss.find(ref); err == nil {
		ss.items[i].valid = false
		ss.items[i].ref = nil
		ss.size--
	} else {
		return err
	}
	return nil
}

func (ss *slice) Clear() {
	for i := 0; i < ss.capacity; i++ {
		ss.items[i].valid = false
		ss.items[i].ref = nil
	}
	ss.size = 0
}

func (ss slice) Size() int {
	return ss.size
}

func (ss slice) Iterator() Iterator {
	it := sliceIterator{
		data:    ss,
		current: -1,
	}
	return it.Next()
}

func (ss *slice) growCapacity() {
	ss.capacity += ss.grow
	ss.items = append(ss.items, make([]item, ss.grow)...)
	ss.grow = (ss.capacity >> 2) + 1 // next grow will be 25% + 1
}

func (ss slice) find(ref interface{}) (int, error) {
	for i, si := range ss.items {
		if si.valid {
			if si.ref == ref {
				return i, nil
			}
		}
	}
	return 0, ErrItemNotFound
}

// Sort a sparse.Slice in-place using a less function
func (ss *slice) Sort(less func(a, b interface{}) bool) {
	sort.Slice(ss.items, func(i, j int) bool {
		a := ss.items[i]
		b := ss.items[j]
		if !a.valid {
			return false
		} else if !b.valid {
			return true
		} else {
			return less(a.ref, b.ref)
		}
	})
}

// AssureCapacity grows the slice until the desired capacity is meet
func (ss *slice) AssureCapacity(capacity int) {
	for ss.capacity < capacity {
		ss.growCapacity()
	}
}

// Copy makes a copy of this Slice into dest
func (ss slice) Copy(dest Slice) {
	dest.AssureCapacity(ss.capacity)
	for it := ss.Iterator(); it != nil; it = it.Next() {
		dest.Add(it.Value())
	}
}

// Replace replace dest content with this Slice contents
func (ss slice) Replace(dest Slice) {
	dest.Clear()
	ss.Copy(dest)
}

// NewSlice creates a new sparse.Slice with the given capacity
func NewSlice(capacity int) Slice {
	slice := slice{
		items:    make([]item, capacity),
		capacity: capacity,
		grow:     capacity, // first grow will double capacity
		size:     0,
	}
	return &slice
}
