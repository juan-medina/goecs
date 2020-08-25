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

// Package view allows to get a View of entity.Entity objects
package view

import (
	"fmt"
	"github.com/juan-medina/goecs/pkg/entity"
	"github.com/juan-medina/goecs/pkg/sparse"
	"reflect"
)

const (
	entitiesCapacity = 2000
	entitiesGrow     = entitiesCapacity / 4
)

// Iterator for view
type Iterator interface {
	// HasNext returns true if we have more items
	HasNext() bool
	// Value returns the current item value
	Value() *entity.Entity
}

// View represent a set of entity.Entity objects
type View struct {
	entities sparse.Slice
}

// String get a string representation of a View
func (vw View) String() string {
	return fmt.Sprintf("View{entities: %v}", vw.entities)
}

// Add a entity.Entity instance to a View
func (vw *View) Add(ent *entity.Entity) *entity.Entity {
	vw.entities.Add(ent)
	return ent
}

// Remove a entity.Entity from a View
func (vw *View) Remove(ent *entity.Entity) error {
	return vw.entities.Remove(ent)
}

// Size of entity.Entity in the View
func (vw View) Size() int {
	return vw.entities.Size()
}

type viewIterator struct {
	eit    sparse.Iterator
	filter []reflect.Type
}

func (vi *viewIterator) HasNext() bool {
	for vi.eit.HasNext() {
		val := vi.Value()
		if val.Contains(vi.filter...) {
			return true
		}
	}
	return false
}

func (vi *viewIterator) Value() *entity.Entity {
	return vi.eit.Value().(*entity.Entity)
}

// Iterator return an view.Iterator for the given varg reflect.Type
func (vw *View) Iterator(rtypes ...reflect.Type) Iterator {
	return &viewIterator{
		eit:    vw.entities.Iterator(),
		filter: rtypes,
	}
}

// Clear removes all entity.Entity from the view.View
func (vw *View) Clear() {
	vw.entities.Clear()
}

// New creates a new empty View
func New() *View {
	return &View{
		entities: sparse.NewSlice(entitiesCapacity, entitiesGrow),
	}
}
