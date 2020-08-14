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

package view

import (
	"fmt"
	"github.com/juan-medina/goecs/pkg/entitiy"
	"reflect"
)

type View struct {
	entities []*entitiy.Entity
}

func (vw View) String() string {
	return fmt.Sprintf("View{entities: %v}", vw.entities)
}

func (vw *View) Add(ent *entitiy.Entity) *entitiy.Entity {
	vw.entities = append(vw.entities, ent)
	return ent
}

func (vw *View) Remove(ent *entitiy.Entity) {
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

func (vw View) Size() int {
	return len(vw.entities)
}

func (vw View) Entities(rtypes ...reflect.Type) []*entitiy.Entity {
	result := make([]*entitiy.Entity, 0)

	for _, v := range vw.entities {
		if v.Contains(rtypes...) {
			result = append(result, v)
		}
	}

	return result
}

func (vw View) Entity(rtypes ...reflect.Type) *entitiy.Entity {
	entities := vw.Entities(rtypes...)
	if len(entities) == 0 {
		return nil
	} else {
		return entities[0]
	}
}

func (vw View) SubView(rtypes ...reflect.Type) *View {
	return fromEntities(vw.Entities(rtypes...))
}

func New() *View {
	return &View{
		entities: make([]*entitiy.Entity, 0),
	}
}

func fromEntities(entities []*entitiy.Entity) *View {
	view := New()

	view.entities = entities

	return view
}
