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

package entitiy

import (
	"fmt"
	"reflect"
)

type Entity struct {
	components map[reflect.Type]interface{}
}

func (e *Entity) String() string {
	var result = ""

	for _, v := range e.components {
		if result != "" {
			result += ","
		}
		result += fmt.Sprintf("%s%v", reflect.TypeOf(v).String(), v)
	}

	return "Entity{" + result + "}"
}

func New() *Entity {
	return &Entity{
		components: make(map[reflect.Type]interface{}),
	}
}

func (e *Entity) Add(i interface{}) *Entity {
	e.components[reflect.TypeOf(i)] = i
	return e
}

func (e *Entity) Set(i interface{}) *Entity {
	return e.Add(i)
}

func (e Entity) Get(t reflect.Type) interface{} {
	return e.components[t]
}

func (e Entity) Contains(ts ...reflect.Type) bool {
	var contains = true

	for _, t := range ts {
		if _, ok := e.components[t]; !ok {
			contains = false
			break
		}
	}

	return contains
}
