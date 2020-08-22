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
	"reflect"
)

// View represent a set of entity.Entity objects
type View struct {
	entities []*entity.Entity
}

// String get a string representation of a View
func (vw View) String() string {
	return fmt.Sprintf("View{entities: %v}", vw.entities)
}

// Add a entity.Entity instance to a View
func (vw *View) Add(ent *entity.Entity) *entity.Entity {
	vw.entities = append(vw.entities, ent)
	return ent
}

// Remove a entity.Entity from a View
func (vw *View) Remove(ent *entity.Entity) {
	i := 0
	for _, v := range vw.entities {
		if v != ent {
			vw.entities[i] = v
			i++
		}
	}

	// Prevent memory leak by erasing truncated values
	for j := i; j < len(vw.entities); j++ {
		vw.entities[j] = nil
	}

	vw.entities = vw.entities[:i]
}

// Size of entity.Entity in the View
func (vw View) Size() int {
	return len(vw.entities)
}

// Entities return slice of entity.Entity for the given varg reflect.Type
func (vw View) Entities(rtypes ...reflect.Type) []*entity.Entity {
	result := make([]*entity.Entity, len(vw.entities))

	count := 0
	for _, v := range vw.entities {
		if v.Contains(rtypes...) {
			result[count] = v
			count++
		}
	}
	result = result[:count]
	return result
}

// Entity return the first entity.Entity for the given varg reflect.Type in the View, nil
// if they are none
func (vw View) Entity(rtypes ...reflect.Type) *entity.Entity {
	if entities := vw.Entities(rtypes...); len(entities) != 0 {
		return entities[0]
	}
	return nil
}

// SubView generate a new View from the given varg reflect.Type
func (vw View) SubView(rtypes ...reflect.Type) *View {
	return fromEntities(vw.Entities(rtypes...))
}

// Filter return a set of entity.Entity trough a filter function
func (vw View) Filter(filter func(*entity.Entity) bool) []*entity.Entity {
	result := make([]*entity.Entity, len(vw.entities))

	count := 0
	for _, v := range vw.entities {
		if filter(v) {
			result[count] = v
			count++
		}
	}
	result = result[:count]
	return result
}

// New creates a new empty View
func New() *View {
	return &View{
		entities: make([]*entity.Entity, 0),
	}
}

// fromEntities creates a View with the given slice of entity.Entity
func fromEntities(entities []*entity.Entity) *View {
	view := New()

	view.entities = entities

	return view
}
