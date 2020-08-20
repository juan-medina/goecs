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

package entity

import (
	"fmt"
	"reflect"
)

type Entity struct {
	components map[reflect.Type]interface{}
}

func (ent *Entity) String() string {
	var result = ""

	for _, v := range ent.components {
		if result != "" {
			result += ","
		}
		result += fmt.Sprintf("%s%v", reflect.TypeOf(v).String(), v)
	}

	return "Entity{" + result + "}"
}

func New(components ...interface{}) *Entity {
	ent := Entity{
		components: make(map[reflect.Type]interface{}),
	}

	for _, v := range components {
		ent.Add(v)
	}

	return &ent
}

func (ent *Entity) Add(component interface{}) *Entity {
	ent.components[reflect.TypeOf(component)] = component
	return ent
}

func (ent *Entity) Set(component interface{}) *Entity {
	return ent.Add(component)
}

func (ent Entity) Get(rtype reflect.Type) interface{} {
	return ent.components[rtype]
}

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
