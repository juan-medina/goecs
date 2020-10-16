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
	"reflect"
)

// Entity represents a instance of an object in a ECS
type Entity struct {
	id         uint64
	components map[ComponentType]Component
}

// ID : get the unique id for this Entity
func (ent Entity) ID() uint64 {
	return ent.id
}

// String : get a string representation of an Entity
func (ent Entity) String() string {
	var result = fmt.Sprintf("id{%d}", ent.id)

	for _, v := range ent.components {
		if result != "" {
			result += ","
		}
		result += fmt.Sprintf("%s%v", reflect.TypeOf(v), v)
	}

	return "Entity{" + result + "}"
}

// NewEntity creates a new Entity giving a set of varg components
func NewEntity(ID uint64, components ...Component) *Entity {
	ent := Entity{
		id:         ID,
		components: make(map[ComponentType]Component),
	}

	for _, v := range components {
		ent.Add(v)
	}
	return &ent
}

// Add a new component into an Entity
func (ent *Entity) Add(component Component) *Entity {
	ent.components[component.Type()] = component
	return ent
}

// Set a new component into an Entity
func (ent *Entity) Set(component Component) *Entity {
	return ent.Add(component)
}

// Get the component of the given ComponentType
func (ent Entity) Get(ctype ComponentType) Component {
	return ent.components[ctype]
}

// Remove the component of the given ComponentType
func (ent *Entity) Remove(ctype ComponentType) {
	delete(ent.components, ctype)
}

// Contains check that the Entity has the given varg ComponentType
func (ent Entity) Contains(types ...ComponentType) bool {
	var contains = true

	for _, t := range types {
		if _, ok := ent.components[t]; !ok {
			contains = false
			break
		}
	}

	return contains
}

// NotContains check that the Entity has not the given varg ComponentType
func (ent Entity) NotContains(types ...ComponentType) bool {
	var noContains = true

	for _, t := range types {
		if _, ok := ent.components[t]; ok {
			noContains = false
			break
		}
	}

	return noContains
}
