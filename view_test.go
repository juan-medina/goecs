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

package goecs_test

import (
	"github.com/juan-medina/goecs"
	"reflect"
	"testing"
)

type GameObject struct {
	name string
}

var GameObjectType = reflect.TypeOf(GameObject{})

func expectViewPositions(t *testing.T, view *goecs.View, want []Pos) {
	t.Helper()
	got := make([]Pos, 0)
	for it := view.Iterator(PosType); it != nil; it = it.Next() {
		v := it.Value()
		got = append(got, v.Get(PosType).(Pos))
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestNewView(t *testing.T) {
	view := goecs.NewView()
	got := view.Size()
	expect := 0

	if got != expect {
		t.Fatalf("error on view size got %d, want %d", got, expect)
	}
}

func TestView_Add(t *testing.T) {
	view := goecs.NewView()

	view.AddEntity(
		Pos{X: 1, Y: 1},
		Vel{X: 2, Y: 2},
	)

	got := view.Size()
	expect := 1

	if got != expect {
		t.Fatalf("error on view size got %d, want %d", got, expect)
	}
}

func TestView_Size(t *testing.T) {
	view := goecs.NewView()

	view.AddEntity(
		Pos{X: 1, Y: 1},
		Vel{X: 2, Y: 2},
	)

	view.AddEntity(
		Pos{X: 1, Y: 1},
	)

	got := view.Size()
	expect := 2

	if got != expect {
		t.Fatalf("error on view size got %d, want %d", got, expect)
	}
}

func entitiesEqual(a, b []*goecs.Entity) bool {
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
	view := goecs.NewView()

	ent1 := view.AddEntity(
		Pos{X: 1, Y: 1},
		Vel{X: 2, Y: 2},
	)

	ent2 := view.AddEntity(
		Pos{X: 1, Y: 1},
	)

	type testCase struct {
		name   string
		params []reflect.Type
		expect []*goecs.Entity
	}
	var cases = []testCase{
		{
			name:   "should get ent1 asking for pos and vel",
			params: []reflect.Type{PosType, VelType},
			expect: []*goecs.Entity{ent1},
		},
		{
			name:   "should get ent1 and ent2 asking only for pos",
			params: []reflect.Type{PosType},
			expect: []*goecs.Entity{ent1, ent2},
		},
		{
			name:   "should get no entities with non existing component",
			params: []reflect.Type{GameObjectType},
			expect: []*goecs.Entity{},
		},
		{
			name:   "should get ent1 asking for only vel",
			params: []reflect.Type{VelType},
			expect: []*goecs.Entity{ent1},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			result := make([]*goecs.Entity, 0)
			for it := view.Iterator(tt.params...); it != nil; it = it.Next() {
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
	view := goecs.NewView()

	ent1 := view.AddEntity(
		Pos{X: 1, Y: 1},
		Vel{X: 2, Y: 2},
	)

	ent2 := view.AddEntity(
		Pos{X: 1, Y: 1},
	)

	view.AddEntity(
		Pos{X: 1, Y: 1},
		Vel{X: 2, Y: 2},
	)

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
	view := goecs.NewView()

	view.AddEntity(Pos{X: 0, Y: 0}, Vel{X: 1, Y: 1})
	view.AddEntity(Pos{X: 2, Y: 2})
	view.AddEntity(Pos{X: 3, Y: 3}, Vel{X: 4, Y: 4})

	s := view.String()

	if len(s) == 0 {
		t.Fatalf("shoudl get string, got empty")
	}
}

func TestView_Clear(t *testing.T) {
	view := goecs.NewView()

	view.AddEntity(Pos{X: 3, Y: 3}).Add(Vel{X: 4, Y: 4})
	view.AddEntity(Pos{X: 0, Y: 0}).Add(Vel{X: 1, Y: 1})
	view.AddEntity(Pos{X: 2, Y: 2})

	view.Clear()

	got := view.Size()

	expect := 0

	if got != expect {
		t.Fatalf("error on view clear got %d, want %d", got, expect)
	}
}

func sortByPosX(a, b *goecs.Entity) bool {
	posA := a.Get(PosType).(Pos)
	posB := b.Get(PosType).(Pos)
	return posA.X < posB.X
}
func sortByPosY(a, b *goecs.Entity) bool {
	posA := a.Get(PosType).(Pos)
	posB := b.Get(PosType).(Pos)
	return posA.Y < posB.Y
}

func TestView_Sort(t *testing.T) {
	view := goecs.NewView()

	view.AddEntity(Pos{X: 3, Y: -3}).Add(Vel{X: 4, Y: 4})
	view.AddEntity(Pos{X: 0, Y: 0}).Add(Vel{X: 1, Y: 1})
	view.AddEntity(Pos{X: 2, Y: -2})

	view.Sort(sortByPosX)

	expectViewPositions(t, view, []Pos{
		{X: 0, Y: 0},
		{X: 2, Y: -2},
		{X: 3, Y: -3},
	})

	view.Sort(sortByPosY)

	expectViewPositions(t, view, []Pos{
		{X: 3, Y: -3},
		{X: 2, Y: -2},
		{X: 0, Y: 0},
	})
}
