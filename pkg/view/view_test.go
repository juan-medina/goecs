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
	"github.com/juan-medina/goecs/pkg/entity"
	"reflect"
	"testing"
)

type Pos struct {
	x float32
	y float32
}

var PosType = reflect.TypeOf(Pos{})

type Vel struct {
	x float32
	y float32
}

var VelType = reflect.TypeOf(Vel{})

type GameObject struct {
	name string
}

var GameObjectType = reflect.TypeOf(GameObject{})

func TestNew(t *testing.T) {
	view := New()
	got := view.Size()
	expect := 0

	if got != expect {
		t.Fatalf("error on view size got %d, want %d", got, expect)
	}
}

func TestView_Add(t *testing.T) {
	view := New()

	view.Add(entity.New().Add(Pos{
		x: 1,
		y: 1,
	}).Add(Vel{
		x: 2,
		y: 2,
	}))

	got := view.Size()
	expect := 1

	if got != expect {
		t.Fatalf("error on view size got %d, want %d", got, expect)
	}
}

func TestView_Size(t *testing.T) {
	view := New()
	ent1 := entity.New().Add(Pos{
		x: 1,
		y: 1,
	}).Add(Vel{
		x: 2,
		y: 2,
	})

	view.Add(ent1)

	ent2 := entity.New().Add(Pos{
		x: 1,
		y: 1,
	})

	view.Add(ent2)

	got := view.Size()
	expect := 2

	if got != expect {
		t.Fatalf("error on view size got %d, want %d", got, expect)
	}
}

func entitiesEqual(a, b []*entity.Entity) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if !reflect.DeepEqual(v, b[i]) {
			return false
		}
	}
	return true
}

func TestView_Iterator(t *testing.T) {
	view := New()

	ent1 := entity.New().Add(Pos{
		x: 1,
		y: 1,
	}).Add(Vel{
		x: 2,
		y: 2,
	})

	view.Add(ent1)

	ent2 := entity.New().Add(Pos{
		x: 1,
		y: 1,
	})

	view.Add(ent2)
	type testCase struct {
		name   string
		params []reflect.Type
		expect []*entity.Entity
	}
	var cases = []testCase{
		{
			name:   "should get ent1 asking for pos and vel",
			params: []reflect.Type{PosType, VelType},
			expect: []*entity.Entity{ent1},
		},
		{
			name:   "should get ent1 and ent2 asking only for pos",
			params: []reflect.Type{PosType},
			expect: []*entity.Entity{ent1, ent2},
		},
		{
			name:   "should get no entities with non existing component",
			params: []reflect.Type{GameObjectType},
			expect: []*entity.Entity{},
		},
		{
			name:   "should get ent1 asking for only vel",
			params: []reflect.Type{VelType},
			expect: []*entity.Entity{ent1},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			result := make([]*entity.Entity, 0)
			for it := view.Iterator(tt.params...); it.HasNext(); {
				value := it.Value()
				result = append(result, value)
			}

			if !entitiesEqual(result, tt.expect) {
				t.Fatalf("error on get entities got %v, want %v", result, tt.expect)
			}

		})
	}
}

func TestView_Remove(t *testing.T) {
	view := New()

	ent1 := view.Add(entity.New().Add(Pos{
		x: 1,
		y: 1,
	}).Add(Vel{
		x: 2,
		y: 2,
	}))

	ent2 := view.Add(entity.New().Add(Pos{
		x: 1,
		y: 1,
	}))

	view.Add(entity.New().Add(Pos{
		x: 1,
		y: 1,
	}).Add(Vel{
		x: 2,
		y: 2,
	}))

	_ = view.Remove(ent2)

	got := view.Size()
	expect := 2

	if got != expect {
		t.Fatalf("error on view size got %d, want %d", got, expect)
	}

	_ = view.Remove(ent1)

	got = view.Size()
	expect = 1

	if got != expect {
		t.Fatalf("error on view size got %d, want %d", got, expect)
	}
}

func TestView_String(t *testing.T) {
	view := New()

	view.Add(entity.New().Add(Pos{x: 0, y: 0}).Add(Vel{x: 1, y: 1}))
	view.Add(entity.New().Add(Pos{x: 2, y: 2}))
	view.Add(entity.New().Add(Pos{x: 3, y: 3}).Add(Vel{x: 4, y: 4}))

	s := view.String()

	if len(s) == 0 {
		t.Fatalf("shoudl get string, got empty")
	}
}

/*
func TestView_Filter(t *testing.T) {
	view := New()

	view.Add(entity.New().Add(Pos{x: 0, y: 0}).Add(Vel{x: 1, y: 1}))
	view.Add(entity.New().Add(Pos{x: 2, y: 2}))
	view.Add(entity.New().Add(Pos{x: 3, y: 3}).Add(Vel{x: 4, y: 4}))

	filtered := view.Filter(func(e *entity.Entity) bool {
		pos := e.Get(PosType).(Pos)
		return pos.x > 0
	})

	got := len(filtered)
	expect := 2

	if got != expect {
		t.Fatalf("error on view filter got %d, want %d", got, expect)
	}
}
*/
