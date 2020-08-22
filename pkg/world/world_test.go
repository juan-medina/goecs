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

package world

import (
	"errors"
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

type resetEvent struct{}

type HMovementSystem struct{}

func (m *HMovementSystem) Notify(world *World, e interface{}, _ float32) error {
	switch e.(type) {
	case resetEvent:
		for _, v := range world.Entities(PosType, VelType) {
			pos := v.Get(PosType).(Pos)
			pos.x = 0
			v.Set(pos)
		}
	}
	return nil
}

func (m *HMovementSystem) Update(world *World, _ float32) error {
	for _, v := range world.Entities(PosType, VelType) {
		pos := v.Get(PosType).(Pos)
		vel := v.Get(VelType).(Vel)

		pos.x += vel.x

		v.Set(pos)
	}
	return nil
}

type VMovementSystem struct{}

func (m *VMovementSystem) Notify(world *World, e interface{}, _ float32) error {
	switch e.(type) {
	case resetEvent:
		for _, v := range world.Entities(PosType, VelType) {
			pos := v.Get(PosType).(Pos)
			pos.y = 0
			v.Set(pos)
		}
	}
	return nil
}
func (m *VMovementSystem) Update(world *World, _ float32) error {
	for _, v := range world.Entities(PosType, VelType) {
		pos := v.Get(PosType).(Pos)
		vel := v.Get(VelType).(Vel)

		pos.y += vel.y

		v.Set(pos)
	}
	return nil
}

var errFailure = errors.New("failure")

type FailureUpdateSystem struct{}

func (f *FailureUpdateSystem) Notify(_ *World, _ interface{}, _ float32) error {
	return nil
}

func (f FailureUpdateSystem) Update(_ *World, _ float32) error {
	return errFailure
}

type FailureNotifySystem struct{}

func (f *FailureNotifySystem) Notify(_ *World, _ interface{}, _ float32) error {
	return errFailure
}

func (f FailureNotifySystem) Update(_ *World, _ float32) error {
	return nil
}
func TestWorld_Update(t *testing.T) {
	t.Run("update single system should work", func(t *testing.T) {
		world := New()
		world.AddSystem(&HMovementSystem{})

		world.Add(entity.New().Add(Pos{x: 0, y: 0}).Add(Vel{x: 1, y: 1}))
		world.Add(entity.New().Add(Pos{x: 2, y: 2}))
		world.Add(entity.New().Add(Pos{x: 3, y: 3}).Add(Vel{x: 4, y: 4}))

		_ = world.Update(0)

		expectPositions(t, world, []Pos{
			{x: 1, y: 0},
			{x: 2, y: 2},
			{x: 7, y: 3},
		})
	})

	t.Run("update multiple systems should work", func(t *testing.T) {
		world := New()
		world.AddSystem(&HMovementSystem{})
		world.AddSystem(&VMovementSystem{})

		world.Add(entity.New().Add(Pos{x: 0, y: 0}).Add(Vel{x: 1, y: 1}))
		world.Add(entity.New().Add(Pos{x: 2, y: 2}))
		world.Add(entity.New().Add(Pos{x: 3, y: 3}).Add(Vel{x: 4, y: 4}))

		_ = world.Update(0)

		expectPositions(t, world, []Pos{
			{x: 1, y: 1},
			{x: 2, y: 2},
			{x: 7, y: 7},
		})
	})
}

func TestWorld_UpdateGroup(t *testing.T) {
	t.Run("update single system should work", func(t *testing.T) {
		world := New()
		world.AddSystem(&HMovementSystem{})

		world.Add(entity.New().Add(Pos{x: 0, y: 0}).Add(Vel{x: 1, y: 1}))
		world.Add(entity.New().Add(Pos{x: 2, y: 2}))
		world.Add(entity.New().Add(Pos{x: 3, y: 3}).Add(Vel{x: 4, y: 4}))

		_ = world.Update(0)

		expectPositions(t, world, []Pos{
			{x: 1, y: 0},
			{x: 2, y: 2},
			{x: 7, y: 3},
		})
	})

	t.Run("update multiple systems should work", func(t *testing.T) {
		world := New()
		world.AddSystem(&HMovementSystem{})
		world.AddSystem(&VMovementSystem{})

		world.Add(entity.New().Add(Pos{x: 0, y: 0}).Add(Vel{x: 1, y: 1}))
		world.Add(entity.New().Add(Pos{x: 2, y: 2}))
		world.Add(entity.New().Add(Pos{x: 3, y: 3}).Add(Vel{x: 4, y: 4}))

		_ = world.Update(0)

		expectPositions(t, world, []Pos{
			{x: 1, y: 1},
			{x: 2, y: 2},
			{x: 7, y: 7},
		})
	})
}

func expectPositions(t *testing.T, world *World, want []Pos) {
	t.Helper()
	entities := world.Entities(PosType)
	got := make([]Pos, 0)
	for _, v := range entities {
		got = append(got, v.Get(PosType).(Pos))
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestWorld_String(t *testing.T) {
	world := New()
	world.AddSystem(&HMovementSystem{})
	world.AddSystem(&VMovementSystem{})

	world.Add(entity.New().Add(Pos{x: 0, y: 0}).Add(Vel{x: 1, y: 1}))
	world.Add(entity.New().Add(Pos{x: 2, y: 2}))
	world.Add(entity.New().Add(Pos{x: 3, y: 3}).Add(Vel{x: 4, y: 4}))

	s := world.String()

	if len(s) == 0 {
		t.Fatalf("shoudl get string, got empty")
	}
}

func TestWorld_Update_Error(t *testing.T) {
	world := New()
	world.AddSystem(&FailureUpdateSystem{})
	world.AddSystem(&HMovementSystem{})
	world.AddSystem(&VMovementSystem{})

	world.Add(entity.New().Add(Pos{x: 0, y: 0}).Add(Vel{x: 1, y: 1}))

	e := world.Update(0)

	if !errors.Is(e, errFailure) {
		t.Fatalf("shoudl get failure but got %v", e)
	}

	expectPositions(t, world, []Pos{
		{x: 0, y: 0},
	})
}

func TestWorld_Notify(t *testing.T) {
	world := New()
	world.AddSystem(&HMovementSystem{})
	world.AddSystem(&VMovementSystem{})

	world.Add(entity.New().Add(Pos{x: 0, y: 0}).Add(Vel{x: 1, y: 1}))
	world.Add(entity.New().Add(Pos{x: 2, y: 2}))
	world.Add(entity.New().Add(Pos{x: 3, y: 3}).Add(Vel{x: 4, y: 4}))

	_ = world.Update(0)

	expectPositions(t, world, []Pos{
		{x: 1, y: 1},
		{x: 2, y: 2},
		{x: 7, y: 7},
	})

	_ = world.Notify(resetEvent{})
	_ = world.Update(0)

	expectPositions(t, world, []Pos{
		{x: 0, y: 0},
		{x: 2, y: 2},
		{x: 0, y: 0},
	})

	_ = world.Update(0)

	expectPositions(t, world, []Pos{
		{x: 1, y: 1},
		{x: 2, y: 2},
		{x: 4, y: 4},
	})
}

func TestWorld_Notify_Error(t *testing.T) {
	world := New()
	world.AddSystem(&FailureNotifySystem{})
	world.AddSystem(&HMovementSystem{})
	world.AddSystem(&VMovementSystem{})

	world.Add(entity.New().Add(Pos{x: 0, y: 0}).Add(Vel{x: 1, y: 1}))
	world.Add(entity.New().Add(Pos{x: 2, y: 2}))
	world.Add(entity.New().Add(Pos{x: 3, y: 3}).Add(Vel{x: 4, y: 4}))

	e := world.Update(0)

	if e != nil {
		t.Fatalf("shoudl not get error but got %v", e)
	}

	expectPositions(t, world, []Pos{
		{x: 1, y: 1},
		{x: 2, y: 2},
		{x: 7, y: 7},
	})

	e = world.Notify(resetEvent{})

	if e != nil {
		t.Fatalf("shoudl not get error but got %v", e)
	}

	e = world.Update(0)

	if !errors.Is(e, errFailure) {
		t.Fatalf("shoudl get failure but got %v", e)
	}

	expectPositions(t, world, []Pos{
		{x: 2, y: 2},
		{x: 2, y: 2},
		{x: 11, y: 11},
	})
}
