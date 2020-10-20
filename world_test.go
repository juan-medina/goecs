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
	"errors"
	"fmt"
	"github.com/juan-medina/goecs"
	"reflect"
	"testing"
)

var resetSignalEventType = goecs.NewComponentType()

type resetSignalEvent struct{}

func (r resetSignalEvent) Type() goecs.ComponentType {
	return resetSignalEventType
}

var dummySignalType = goecs.NewComponentType()

type dummySignal struct{}

func (d dummySignal) Type() goecs.ComponentType {
	return dummySignalType
}

func ResetHListener(world *goecs.World, e goecs.Component, _ float32) error {
	switch e.(type) {
	case resetSignalEvent:
		for it := world.Iterator(PosType, VelType); it != nil; it = it.Next() {
			v := it.Value()
			pos := v.Get(PosType).(Pos)

			pos.X = 0
			v.Set(pos)
		}
	}
	return nil
}

func HMovementSystem(world *goecs.World, _ float32) error {
	for it := world.Iterator(PosType, VelType); it != nil; it = it.Next() {
		v := it.Value()
		pos := v.Get(PosType).(Pos)
		vel := v.Get(VelType).(Vel)

		pos.X += vel.X
		v.Set(pos)
	}

	return nil
}

func ResetVListener(world *goecs.World, e goecs.Component, _ float32) error {
	switch e.(type) {
	case resetSignalEvent:
		for it := world.Iterator(PosType, VelType); it != nil; it = it.Next() {
			v := it.Value()
			pos := v.Get(PosType).(Pos)

			pos.Y = 0
			v.Set(pos)
		}
	}
	return nil
}

func VMovementSystem(world *goecs.World, _ float32) error {
	for it := world.Iterator(PosType, VelType); it != nil; it = it.Next() {
		v := it.Value()
		pos := v.Get(PosType).(Pos)
		vel := v.Get(VelType).(Vel)

		pos.Y += vel.Y
		v.Set(pos)
	}
	return nil
}

var errFailure = errors.New("failure")

func FailureSystem(_ *goecs.World, _ float32) error {
	return errFailure
}

func FailureListener(_ *goecs.World, _ goecs.Component, _ float32) error {
	return errFailure
}

func TestWorld_Update(t *testing.T) {
	t.Run("update single system should work", func(t *testing.T) {
		world := goecs.Default()
		world.AddSystem(HMovementSystem)

		world.AddEntity(Pos{X: 0, Y: 0}, Vel{X: 1, Y: 1})
		world.AddEntity(Pos{X: 2, Y: 2})
		world.AddEntity(Pos{X: 3, Y: 3}, Vel{X: 4, Y: 4})

		_ = world.Update(0)

		expectWorldPositions(t, world, []Pos{
			{X: 1, Y: 0},
			{X: 2, Y: 2},
			{X: 7, Y: 3},
		})
	})

	t.Run("update multiple systems should work", func(t *testing.T) {
		world := goecs.Default()
		world.AddSystem(HMovementSystem)
		world.AddSystem(VMovementSystem)

		world.AddEntity(Pos{X: 0, Y: 0}, Vel{X: 1, Y: 1})
		world.AddEntity(Pos{X: 2, Y: 2})
		world.AddEntity(Pos{X: 3, Y: 3}, Vel{X: 4, Y: 4})

		_ = world.Update(0)

		expectWorldPositions(t, world, []Pos{
			{X: 1, Y: 1},
			{X: 2, Y: 2},
			{X: 7, Y: 7},
		})
	})
}

func TestWorld_UpdateGroup(t *testing.T) {
	t.Run("update single system should work", func(t *testing.T) {
		world := goecs.Default()
		world.AddSystem(HMovementSystem)

		world.AddEntity(Pos{X: 0, Y: 0}, Vel{X: 1, Y: 1})
		world.AddEntity(Pos{X: 2, Y: 2})
		world.AddEntity(Pos{X: 3, Y: 3}, Vel{X: 4, Y: 4})

		_ = world.Update(0)

		expectWorldPositions(t, world, []Pos{
			{X: 1, Y: 0},
			{X: 2, Y: 2},
			{X: 7, Y: 3},
		})
	})

	t.Run("update multiple systems should work", func(t *testing.T) {
		world := goecs.Default()
		world.AddSystem(HMovementSystem)
		world.AddSystem(VMovementSystem)

		world.AddEntity(Pos{X: 0, Y: 0}, Vel{X: 1, Y: 1})
		world.AddEntity(Pos{X: 2, Y: 2})
		world.AddEntity(Pos{X: 3, Y: 3}, Vel{X: 4, Y: 4})

		_ = world.Update(0)

		expectWorldPositions(t, world, []Pos{
			{X: 1, Y: 1},
			{X: 2, Y: 2},
			{X: 7, Y: 7},
		})
	})
}

func expectWorldPositions(t *testing.T, world *goecs.World, want []Pos) {
	t.Helper()
	expectViewPositions(t, world.View, want)
}

func TestWorld_String(t *testing.T) {
	world := goecs.Default()
	world.AddSystem(HMovementSystem)
	world.AddSystem(VMovementSystem)
	world.AddListener(ResetHListener, resetSignalEventType)
	world.AddListener(ResetVListener, resetSignalEventType)

	world.AddEntity(Pos{X: 0, Y: 0}, Vel{X: 1, Y: 1})
	world.AddEntity(Pos{X: 2, Y: 2})
	world.AddEntity(Pos{X: 3, Y: 3}, Vel{X: 4, Y: 4})

	s := world.String()
	fmt.Println(s)

	if len(s) == 0 {
		t.Fatalf("shoudl get string, got empty")
	}
}

func TestWorld_Update_Error(t *testing.T) {
	world := goecs.Default()
	world.AddSystem(FailureSystem)
	world.AddSystem(HMovementSystem)
	world.AddSystem(VMovementSystem)

	world.AddEntity(Pos{X: 0, Y: 0}, Vel{X: 1, Y: 1})

	e := world.Update(0)

	if !errors.Is(e, errFailure) {
		t.Fatalf("shoudl get failure but got %v", e)
	}

	expectWorldPositions(t, world, []Pos{
		{X: 0, Y: 0},
	})
}

func TestWorld_Signal(t *testing.T) {
	world := goecs.Default()
	world.AddSystem(HMovementSystem)
	world.AddSystem(VMovementSystem)
	world.AddListener(ResetHListener, resetSignalEventType)
	world.AddListener(ResetVListener, resetSignalEventType)

	world.AddEntity(Pos{X: 0, Y: 0}, Vel{X: 1, Y: 1})
	world.AddEntity(Pos{X: 2, Y: 2})
	world.AddEntity(Pos{X: 3, Y: 3}, Vel{X: 4, Y: 4})

	_ = world.Update(0)

	expectWorldPositions(t, world, []Pos{
		{X: 1, Y: 1},
		{X: 2, Y: 2},
		{X: 7, Y: 7},
	})

	world.Signal(resetSignalEvent{})
	_ = world.Update(0)

	expectWorldPositions(t, world, []Pos{
		{X: 0, Y: 0},
		{X: 2, Y: 2},
		{X: 0, Y: 0},
	})

	_ = world.Update(0)

	expectWorldPositions(t, world, []Pos{
		{X: 1, Y: 1},
		{X: 2, Y: 2},
		{X: 4, Y: 4},
	})
}

var numSignalType = goecs.NewComponentType()

type nunSignal struct {
	num int
}

func (s nunSignal) Type() goecs.ComponentType {
	return numSignalType
}

func TestWorld_SignalMultiple(t *testing.T) {
	world := goecs.Default()

	sum := 0
	world.AddListener(func(world *goecs.World, e goecs.Component, _ float32) error {
		switch n := e.(type) {
		case nunSignal:
			sum += n.num
		}
		return nil
	}, numSignalType)

	world.Signal(nunSignal{num: 1})
	world.Signal(nunSignal{num: 2})
	world.Signal(nunSignal{num: 3})
	world.Signal(nunSignal{num: 4})

	_ = world.Update(0)

	got := sum
	expect := 10

	if got != expect {
		t.Fatalf("error on testing multiple signals got %d , want %d", got, expect)
	}
}

func TestWorld_Signal_Error(t *testing.T) {
	world := goecs.Default()
	world.AddListener(FailureListener, resetSignalEventType)
	world.AddSystem(HMovementSystem)
	world.AddSystem(VMovementSystem)

	world.AddEntity(Pos{X: 0, Y: 0}, Vel{X: 1, Y: 1})
	world.AddEntity(Pos{X: 2, Y: 2})
	world.AddEntity(Pos{X: 3, Y: 3}, Vel{X: 4, Y: 4})

	e := world.Update(0)

	if e != nil {
		t.Fatalf("shoudl not get error but got %v", e)
	}

	expectWorldPositions(t, world, []Pos{
		{X: 1, Y: 1},
		{X: 2, Y: 2},
		{X: 7, Y: 7},
	})

	world.Signal(resetSignalEvent{})

	e = world.Update(0)

	if !errors.Is(e, errFailure) {
		t.Fatalf("shoudl get failure but got %v", e)
	}

	expectWorldPositions(t, world, []Pos{
		{X: 2, Y: 2},
		{X: 2, Y: 2},
		{X: 11, Y: 11},
	})
}

var systemCalls []string

func systemA(_ *goecs.World, _ float32) error {
	systemCalls = append(systemCalls, "update a")
	return nil
}
func listenerA(_ *goecs.World, _ goecs.Component, _ float32) error {
	systemCalls = append(systemCalls, "notify a")
	return nil
}

func systemB(_ *goecs.World, _ float32) error {
	systemCalls = append(systemCalls, "update b")
	return nil
}
func listenerB(_ *goecs.World, _ goecs.Component, _ float32) error {
	systemCalls = append(systemCalls, "notify b")
	return nil
}

func TestWorld_AddSystemWithPriority(t *testing.T) {
	type testCase struct {
		name   string
		setup  func(world *goecs.World)
		expect []string
	}
	for _, tc := range []testCase{
		{
			name: "without priority",
			setup: func(world *goecs.World) {
				world.AddSystem(systemA)
				world.AddSystem(systemB)
				world.AddListener(listenerA, dummySignalType)
				world.AddListener(listenerB, dummySignalType)
			},
			expect: []string{
				"update a",
				"update b",
				"notify a",
				"notify b",
			},
		},
		{
			name: "with priority",
			setup: func(world *goecs.World) {
				world.AddSystem(systemA)
				world.AddListener(listenerA, dummySignalType)
				world.AddSystemWithPriority(systemB, 100)
				world.AddListenerWithPriority(listenerB, 100, dummySignalType)
			},
			expect: []string{
				"update b",
				"update a",
				"notify b",
				"notify a",
			},
		},
		{
			name: "with priority inverted",
			setup: func(world *goecs.World) {
				world.AddSystem(systemB)
				world.AddListener(listenerB, dummySignalType)
				world.AddSystemWithPriority(systemA, -100)
				world.AddListenerWithPriority(listenerA, -100, dummySignalType)
			},
			expect: []string{
				"update b",
				"update a",
				"notify b",
				"notify a",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			systemCalls = make([]string, 0)
			world := goecs.Default()

			tc.setup(world)

			world.Signal(dummySignal{})
			_ = world.Update(0)

			if !reflect.DeepEqual(systemCalls, tc.expect) {
				t.Fatalf("got %v, want %v", systemCalls, tc.expect)
			}
		})
	}
}

func TestWorld_Clear(t *testing.T) {
	world := goecs.Default()

	world.AddSystem(HMovementSystem)
	world.AddSystem(VMovementSystem)

	world.AddEntity(Pos{X: 0, Y: 0}, Vel{X: 1, Y: 1})
	world.AddEntity(Pos{X: 2, Y: 2})
	world.AddEntity(Pos{X: 3, Y: 3}, Vel{X: 4, Y: 4})

	world.Clear()

	got := world.Size()

	expect := 0

	if got != expect {
		t.Fatalf("error on world clear got %d entities, want %d entities", got, expect)
	}
}

type constantVelocity struct {
	vel Vel
}

func (cv *constantVelocity) system(world *goecs.World, _ float32) error {
	for it := world.Iterator(PosType); it != nil; it = it.Next() {
		v := it.Value()
		pos := v.Get(PosType).(Pos)
		pos.X += cv.vel.X
		pos.Y += cv.vel.Y
		v.Set(pos)
	}

	return nil
}

func newConstantVelocity(vel Vel) constantVelocity {
	return constantVelocity{
		vel: vel,
	}
}

func Test_SystemsInStruct(t *testing.T) {
	world := goecs.Default()

	world.AddEntity(Pos{X: 0, Y: 0})
	world.AddEntity(Pos{X: 2, Y: 2})
	world.AddEntity(Pos{X: 3, Y: 3})

	cv := newConstantVelocity(Vel{X: 1, Y: 2})
	world.AddSystem(cv.system)

	_ = world.Update(0)

	expectWorldPositions(t, world, []Pos{
		{X: 1, Y: 2},
		{X: 3, Y: 4},
		{X: 4, Y: 5},
	})
}

func TestWorld_Sort(t *testing.T) {
	world := goecs.Default()

	world.AddEntity(Pos{X: 3, Y: -3}, Vel{X: 4, Y: 4})
	world.AddEntity(Pos{X: 0, Y: 0}, Vel{X: 1, Y: 1})
	world.AddEntity(Pos{X: 2, Y: -2})

	world.Sort(sortByPosX)

	expectWorldPositions(t, world, []Pos{
		{X: 0, Y: 0},
		{X: 2, Y: -2},
		{X: 3, Y: -3},
	})

	world.Sort(sortByPosY)

	expectWorldPositions(t, world, []Pos{
		{X: 3, Y: -3},
		{X: 2, Y: -2},
		{X: 0, Y: 0},
	})
}

func Test_Resources(t *testing.T) {
	world := goecs.Default()

	world.AddEntity(Pos{X: 3, Y: -3}, Vel{X: 4, Y: 4})
	world.AddEntity(Pos{X: 0, Y: 0}, Vel{X: 1, Y: 1})
	world.AddEntity(Pos{X: 2, Y: -2})

	idScore := world.AddResource(score{points: 100})

	id := world.FindResource(scoreType)

	if id != idScore {
		t.Fatalf("error on get resource got id %v, want id %v", id, idScore)
	}

	want := score{points: 100}
	got := world.GetResource(idScore).Get(scoreType).(score)

	if got != want {
		t.Fatalf("error on get resource got %v, want %v", got, want)
	}

	id = world.FindResource(PosType)

	if id != 0 {
		t.Fatalf("error on get resource got id %v, want id %v", id, 0)
	}

	err := world.RemoveResource(idScore)

	if err != nil {
		t.Fatalf("error on remove resource got  %v, want nil", err)
	}

	id = world.FindResource(scoreType)

	if id != 0 {
		t.Fatalf("error get remove resource got %v, want 0", id)
	}
}

type score struct {
	points int
}

func (s score) Type() goecs.ComponentType {
	return scoreType
}

var scoreType = goecs.NewComponentType()
