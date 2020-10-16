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

package goecs

import (
	"errors"
	"fmt"
	"sort"
)

var (
	// ErrEntityNotFound is the error when we could not find an viewItem
	ErrEntityNotFound = errors.New("entity not found")
)

// View represent a set of Entity objects
type View struct {
	capacity int
	grow     int
	items    []*Entity
	size     int
	lastID   EntityID
}

// Iterator allow to iterate trough the View
type Iterator struct {
	data    *View
	current int
	filter  []ComponentType
}

// Next return a Iterator to the next Entity
func (ei *Iterator) Next() *Iterator {
	for i := ei.current + 1; i < len(ei.data.items); i++ {
		item := ei.data.items[i]
		if item != nil {
			if item.Contains(ei.filter...) {
				ei.current = i
				return ei
			}
		}
	}
	return nil
}

// first return a Iterator to the first Entity
func (ei *Iterator) first() *Iterator {
	for i := ei.current + 1; i < len(ei.data.items); i++ {
		item := ei.data.items[i]
		if item != nil {
			if item.Contains(ei.filter...) {
				ei.current = i
				return ei
			}
		}
	}
	return nil
}

// Value returns the value of the current Iterator
func (ei *Iterator) Value() *Entity {
	return ei.data.items[ei.current]
}

// AddEntity a Entity instance to a View given it components
func (v *View) AddEntity(data ...Component) EntityID {
	v.lastID++
	ent := NewEntity(v.lastID, data...)
	for i, si := range v.items {
		if si == nil {
			v.items[i] = ent
			v.size++
			return ent.ID()
		}
	}

	v.growCapacity()
	v.items[v.size] = ent
	v.size++
	return ent.ID()
}

// Remove a Entity from a View
func (v *View) Remove(id EntityID) error {
	if i, err := v.find(id); err == nil {
		v.items[i] = nil
		v.size--
	} else {
		return err
	}
	return nil
}

// Get a Entity from a View giving it EntityID
func (v *View) Get(id EntityID) (*Entity, error) {
	var err error = nil
	var i int
	if i, err = v.find(id); err == nil {
		return v.items[i], nil
	}
	return nil, err
}

// Clear removes all Entity from the View
func (v *View) Clear() {
	for i := 0; i < v.capacity; i++ {
		v.items[i] = nil
	}
	v.size = 0
}

// Size is the number of Entity in this View
func (v View) Size() int {
	return v.size
}

// Iterator return an view.Iterator for the given varg ComponentType
func (v *View) Iterator(types ...ComponentType) *Iterator {
	it := Iterator{
		data:    v,
		current: -1,
		filter:  types,
	}
	return it.first()
}

// growCapacity increases the View capacity
func (v *View) growCapacity() {
	v.capacity += v.grow
	v.items = append(v.items, make([]*Entity, v.grow)...)
	v.grow = (v.capacity >> 2) + 1 // next grow will be 25% + 1
}

// find a Entity position in a View giving it EntityID
func (v View) find(id EntityID) (int, error) {
	for i, si := range v.items {
		if si != nil {
			if si.ID() == id {
				return i, nil
			}
		}
	}
	return 0, ErrEntityNotFound
}

// Sort the entities in place with a less function
func (v *View) Sort(less func(a, b *Entity) bool) {
	sort.Slice(v.items, func(i, j int) bool {
		a := v.items[i]
		b := v.items[j]
		if a == nil {
			return false
		} else if b == nil {
			return true
		} else {
			return less(a, b)
		}
	})
}

// String get a string representation of a View
func (v View) String() string {
	str := ""
	for it := v.Iterator(); it != nil; it = it.Next() {
		if str != "" {
			str += ","
		}
		ent := it.Value()
		str += ent.String()
	}

	return fmt.Sprintf("View{entities: [%v]}", str)
}

// NewView creates a new empty View with a given capacity
func NewView(capacity int) *View {
	slice := View{
		items:    make([]*Entity, capacity),
		capacity: capacity,
		grow:     capacity, // first grow will double capacity
		size:     0,
	}
	return &slice
}
