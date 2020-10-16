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
	"fmt"
	"github.com/juan-medina/goecs/sparse"
)

// Iterator for view
type Iterator interface {
	// Next returns next element
	Next() Iterator
	// Value returns the current item value
	Value() *Entity
}

// View represent a set of Entity objects
type View struct {
	entities sparse.Slice
	lastID   uint64
}

// String get a string representation of a View
func (vw View) String() string {
	str := ""
	for it := vw.Iterator(); it != nil; it = it.Next() {
		if str != "" {
			str += ","
		}
		ent := it.Value()
		str += ent.String()
	}

	return fmt.Sprintf("View{entities: [%v]}", str)
}

// AddEntity a Entity instance to a View given it components
func (vw *View) AddEntity(components ...Component) *Entity {
	ent := NewEntity(vw.lastID, components...)
	vw.entities.Add(ent)
	vw.lastID++
	return ent
}

// Remove a Entity from a View
func (vw *View) Remove(ent *Entity) error {
	return vw.entities.Remove(ent)
}

// Size of Entity in the View
func (vw View) Size() int {
	return vw.entities.Size()
}

type viewIterator struct {
	view   *View
	eit    sparse.Iterator
	filter []ComponentType
}

func (vi *viewIterator) Next() Iterator {
	for vi.eit = vi.eit.Next(); vi.eit != nil; vi.eit = vi.eit.Next() {
		val := vi.Value()
		if val.Contains(vi.filter...) {
			return vi
		}
	}
	return nil
}

func (vi *viewIterator) first() Iterator {
	for it := vi.view.entities.Iterator(); it != nil; it = it.Next() {
		val := it.Value().(*Entity)
		if val.Contains(vi.filter...) {
			vi.eit = it
			return vi
		}
	}
	return nil
}

func (vi *viewIterator) Value() *Entity {
	return vi.eit.Value().(*Entity)
}

// Iterator return an view.Iterator for the given varg ComponentType
func (vw *View) Iterator(types ...ComponentType) Iterator {
	it := viewIterator{
		view:   vw,
		eit:    nil,
		filter: types,
	}
	return it.first()
}

// Clear removes all Entity from the view.View
func (vw *View) Clear() {
	vw.entities.Clear()
}

// Sort the entities in place with a less function
func (vw *View) Sort(less func(a, b *Entity) bool) {
	vw.entities.Sort(func(a interface{}, b interface{}) bool {
		return less(a.(*Entity), b.(*Entity))
	})
}

// NewView creates a new empty View with a given capacity
func NewView(capacity int) *View {
	return &View{
		entities: sparse.NewSlice(capacity),
	}
}
