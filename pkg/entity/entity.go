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

// Package entity  contains the objects to manage the Entity objects in a ECS
package entity

import (
	"fmt"
	"reflect"
)

var (
	lastID = int64(0)
)

// Entity represents a instance of an object in a ECS
type Entity struct {
	id         int64
	components map[reflect.Type]interface{}
}

// ID : get the unique id for this Entity
func (ent Entity) ID() int64 {
	return ent.id
}

// String : get a string representation of an Entity
func (ent Entity) String() string {
	var result = fmt.Sprintf("id{%d} ", ent.id)

	for _, v := range ent.components {
		if result != "" {
			result += ","
		}
		result += fmt.Sprintf("%s%v", reflect.TypeOf(v).String(), v)
	}

	return "Entity{" + result + "}"
}

// New creates a new Entity giving a set of varg components
func New(components ...interface{}) *Entity {
	ent := Entity{
		id:         lastID,
		components: make(map[reflect.Type]interface{}),
	}

	for _, v := range components {
		ent.Add(v)
	}

	lastID++
	return &ent
}

// Add a new component into an Entity
func (ent *Entity) Add(component interface{}) *Entity {
	ent.components[reflect.TypeOf(component)] = component
	return ent
}

// Set a new component into an Entity
func (ent *Entity) Set(component interface{}) *Entity {
	return ent.Add(component)
}

// Get the component of the given reflect.Type
func (ent Entity) Get(rtype reflect.Type) interface{} {
	return ent.components[rtype]
}

// Contains check that the Entity has the given varg reflect.Type
func (ent Entity) Contains(rtypes ...reflect.Type) bool {
	var contains = true

	for _, t := range rtypes {
		if _, ok := ent.components[t]; !ok {
			contains = false
			break
		}
	}

	return contains
}
