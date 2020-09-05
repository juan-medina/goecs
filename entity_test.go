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

func TestNewEntity(t *testing.T) {
	ent1 := goecs.NewEntity().Add(Pos{
		X: 1,
		Y: 1,
	}).Add(Vel{
		X: 2,
		Y: 2,
	})

	gotPos := ent1.Get(PosType).(Pos)
	gotVel := ent1.Get(VelType).(Vel)

	expectPos := Pos{
		X: 1,
		Y: 1,
	}

	expectVel := Vel{
		X: 2,
		Y: 2,
	}

	if !reflect.DeepEqual(gotPos, expectPos) {
		t.Fatalf("fail to Get Pos: got %v, want %v", gotPos, expectPos)
	}

	if !reflect.DeepEqual(gotVel, expectVel) {
		t.Fatalf("fail to Get Vel: got %v, want %v", gotVel, expectVel)
	}
}

func TestNew_With_Components(t *testing.T) {
	ent1 := goecs.NewEntity(Pos{
		X: 1,
		Y: 1,
	}, Vel{
		X: 2,
		Y: 2,
	})

	gotPos := ent1.Get(PosType).(Pos)
	gotVel := ent1.Get(VelType).(Vel)

	expectPos := Pos{
		X: 1,
		Y: 1,
	}

	expectVel := Vel{
		X: 2,
		Y: 2,
	}

	if !reflect.DeepEqual(gotPos, expectPos) {
		t.Fatalf("fail to Get Pos: got %v, want %v", gotPos, expectPos)
	}

	if !reflect.DeepEqual(gotVel, expectVel) {
		t.Fatalf("fail to Get Vel: got %v, want %v", gotVel, expectVel)
	}
}

func TestEntity_Get(t *testing.T) {
	ent1 := goecs.NewEntity().Add(Pos{
		X: 1,
		Y: 1,
	})

	expectPos := Pos{
		X: 1,
		Y: 1,
	}

	gotPos, ok := ent1.Get(PosType).(Pos)

	if !ok {
		t.Fatalf("expect get is ok, but it wasn't")
	}

	if !reflect.DeepEqual(gotPos, expectPos) {
		t.Fatalf("fail to Get Pos: got %v, want %v", gotPos, expectPos)
	}

	gotVel, ok := ent1.Get(VelType).(Vel)

	if ok {
		t.Fatalf("expect cast is not ok, but it was")
	}

	expectVel := Vel{}

	if !reflect.DeepEqual(gotVel, expectVel) {
		t.Fatalf("expect to get empty but got %v, want %v", gotVel, expectVel)
	}

}

func TestEntity_Contains(t *testing.T) {
	ent1 := goecs.NewEntity().Add(Pos{
		X: 1,
		Y: 1,
	}).Add(Vel{
		X: 2,
		Y: 2,
	})

	ent2 := goecs.NewEntity().Add(Pos{
		X: 1,
		Y: 1,
	})

	t.Run("expect to Contains pos and vel", func(t *testing.T) {
		if !ent1.Contains(PosType, VelType) {
			t.Fatalf("error, does not container them")
		}
	})

	t.Run("expect to not Contains pos and vel", func(t *testing.T) {
		if ent2.Contains(PosType, VelType) {
			t.Fatalf("error, it Contains them")
		}
	})

	t.Run("expect to Contains pos", func(t *testing.T) {
		if !ent1.Contains(PosType) {
			t.Fatalf("error, does not contain it")
		}
	})

	t.Run("expect to not Contains vel", func(t *testing.T) {
		if ent2.Contains(VelType) {
			t.Fatalf("error, its Contains it")
		}
	})
}

func TestEntity_NoContains(t *testing.T) {
	ent1 := goecs.NewEntity().Add(Pos{
		X: 1,
		Y: 1,
	}).Add(Vel{
		X: 2,
		Y: 2,
	})

	ent2 := goecs.NewEntity().Add(Pos{
		X: 1,
		Y: 1,
	})

	var want bool
	var got bool

	t.Run("expect ent1 to contains pos and vel", func(t *testing.T) {
		got = ent1.NotContains(PosType, VelType)
		want = false
		if got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("expect ent2 to not Contains pos and vel", func(t *testing.T) {
		got := ent2.NotContains(PosType, VelType)
		want = false
		if got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("expect ent1 to Contains pos", func(t *testing.T) {
		got := ent1.NotContains(PosType)
		want = false
		if got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("expect ent1 to Contains vel", func(t *testing.T) {
		got := ent1.NotContains(VelType)
		want = false
		if got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("expect ent12 to Contains pos", func(t *testing.T) {
		got := ent2.NotContains(PosType)
		want = false
		if got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("expect ent2 to not Contains vel", func(t *testing.T) {
		got := ent2.NotContains(VelType)
		want = true
		if got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
}

func TestEntity_Set(t *testing.T) {
	ent1 := goecs.NewEntity().Add(Pos{
		X: 1,
		Y: 1,
	}).Add(Vel{
		X: 2,
		Y: 2,
	})

	pos := ent1.Get(PosType).(Pos)

	pos.X = 3
	pos.Y = 3

	ent1.Set(pos)

	got := ent1.Get(PosType).(Pos)

	if pos != got {
		t.Fatalf("erro changin component value got %v, expect %v", got, pos)
	}
}

func TestEntity_String(t *testing.T) {
	ent := goecs.NewEntity().Add(Pos{X: 0, Y: 0}).Add(Vel{X: 1, Y: 1})

	s := ent.String()

	if len(s) == 0 {
		t.Fatalf("shoudl get string, got empty")
	}
}

func TestEntity_Id(t *testing.T) {
	ent1 := goecs.NewEntity().Add(Pos{X: 0, Y: 0}).Add(Vel{X: 1, Y: 1})
	ent2 := goecs.NewEntity().Add(Pos{X: 0, Y: 0}).Add(Vel{X: 1, Y: 1})

	if !(ent1.ID() <= ent2.ID()) {
		t.Fatalf("expect ent1 to have bigger id than ent 2, ent1 %v ent2 %v", ent1, ent2)
	}
}

func TestEntity_Remove(t *testing.T) {
	ent := goecs.NewEntity().Add(Pos{X: 0, Y: 0}).Add(Vel{X: 1, Y: 1})

	var got bool
	var want bool

	t.Run("expect ent to not Contains pol", func(t *testing.T) {
		ent.Remove(PosType)
		got = ent.NotContains(PosType)
		want = true
		if got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("expect ent to Contains vel", func(t *testing.T) {
		got = ent.Contains(VelType)
		want = true
		if got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("expect ent to not Contains vel", func(t *testing.T) {
		ent.Remove(VelType)
		got = ent.NotContains(PosType)
		want = true
		if got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("expect not to fail", func(t *testing.T) {
		ent.Remove(VelType)
	})
}
