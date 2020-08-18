package world

import (
	"errors"
	"github.com/juan-medina/goecs/pkg/entitiy"
	"github.com/juan-medina/goecs/pkg/system"
	"github.com/juan-medina/goecs/pkg/view"
	"reflect"
	"testing"
)

type Pos struct {
	x float64
	y float64
}

var PosType = reflect.TypeOf(Pos{})

type Vel struct {
	x float64
	y float64
}

var VelType = reflect.TypeOf(Vel{})

type HMovementSystem struct{}

func (m *HMovementSystem) Update(view *view.View, _ float64) error {
	for _, v := range view.Entities(PosType, VelType) {
		pos := v.Get(PosType).(Pos)
		vel := v.Get(VelType).(Vel)

		pos.x += vel.x

		v.Set(pos)
	}
	return nil
}

type VMovementSystem struct{}

func (m *VMovementSystem) Update(view *view.View, _ float64) error {
	for _, v := range view.Entities(PosType, VelType) {
		pos := v.Get(PosType).(Pos)
		vel := v.Get(VelType).(Vel)

		pos.y += vel.y

		v.Set(pos)
	}
	return nil
}

var errFailure = errors.New("failure")

type FailureSystem struct{}

func (f FailureSystem) Update(_ *view.View, _ float64) error {
	return errFailure
}

func TestWorld_Update(t *testing.T) {
	t.Run("update single system should work", func(t *testing.T) {
		world := New()
		world.AddSystem(&HMovementSystem{})

		world.Add(entitiy.New().Add(Pos{x: 0, y: 0}).Add(Vel{x: 1, y: 1}))
		world.Add(entitiy.New().Add(Pos{x: 2, y: 2}))
		world.Add(entitiy.New().Add(Pos{x: 3, y: 3}).Add(Vel{x: 4, y: 4}))

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

		world.Add(entitiy.New().Add(Pos{x: 0, y: 0}).Add(Vel{x: 1, y: 1}))
		world.Add(entitiy.New().Add(Pos{x: 2, y: 2}))
		world.Add(entitiy.New().Add(Pos{x: 3, y: 3}).Add(Vel{x: 4, y: 4}))

		_ = world.Update(0)

		expectPositions(t, world, []Pos{
			{x: 1, y: 1},
			{x: 2, y: 2},
			{x: 7, y: 7},
		})
	})

	t.Run("update with special groups", func(t *testing.T) {
		world := New()
		world.AddSystemToGroup(&HMovementSystem{}, "special group")
		world.AddSystem(&VMovementSystem{})

		world.Add(entitiy.New().Add(Pos{x: 0, y: 0}).Add(Vel{x: 1, y: 1}))
		world.Add(entitiy.New().Add(Pos{x: 2, y: 2}))
		world.Add(entitiy.New().Add(Pos{x: 3, y: 3}).Add(Vel{x: 4, y: 4}))

		_ = world.Update(0)

		expectPositions(t, world, []Pos{
			{x: 0, y: 1},
			{x: 2, y: 2},
			{x: 3, y: 7},
		})
	})

	t.Run("update only to special groups", func(t *testing.T) {
		world := New()
		world.AddSystemToGroup(&HMovementSystem{}, "special group")

		world.Add(entitiy.New().Add(Pos{x: 0, y: 0}).Add(Vel{x: 1, y: 1}))
		world.Add(entitiy.New().Add(Pos{x: 2, y: 2}))
		world.Add(entitiy.New().Add(Pos{x: 3, y: 3}).Add(Vel{x: 4, y: 4}))

		_ = world.Update(0)

		expectPositions(t, world, []Pos{
			{x: 0, y: 0},
			{x: 2, y: 2},
			{x: 3, y: 3},
		})
	})
}

func TestWorld_UpdateGroup(t *testing.T) {
	t.Run("update single system should work", func(t *testing.T) {
		world := New()
		world.AddSystem(&HMovementSystem{})

		world.Add(entitiy.New().Add(Pos{x: 0, y: 0}).Add(Vel{x: 1, y: 1}))
		world.Add(entitiy.New().Add(Pos{x: 2, y: 2}))
		world.Add(entitiy.New().Add(Pos{x: 3, y: 3}).Add(Vel{x: 4, y: 4}))

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

		world.Add(entitiy.New().Add(Pos{x: 0, y: 0}).Add(Vel{x: 1, y: 1}))
		world.Add(entitiy.New().Add(Pos{x: 2, y: 2}))
		world.Add(entitiy.New().Add(Pos{x: 3, y: 3}).Add(Vel{x: 4, y: 4}))

		_ = world.UpdateGroup(system.DefaultGroup, 0)

		expectPositions(t, world, []Pos{
			{x: 1, y: 1},
			{x: 2, y: 2},
			{x: 7, y: 7},
		})
	})

	t.Run("update with special groups", func(t *testing.T) {
		world := New()
		world.AddSystemToGroup(&HMovementSystem{}, "special group")
		world.AddSystem(&VMovementSystem{})

		world.Add(entitiy.New().Add(Pos{x: 0, y: 0}).Add(Vel{x: 1, y: 1}))
		world.Add(entitiy.New().Add(Pos{x: 2, y: 2}))
		world.Add(entitiy.New().Add(Pos{x: 3, y: 3}).Add(Vel{x: 4, y: 4}))

		_ = world.UpdateGroup("special group", 0)

		expectPositions(t, world, []Pos{
			{x: 1, y: 0},
			{x: 2, y: 2},
			{x: 7, y: 3},
		})
	})

	t.Run("update only to special groups", func(t *testing.T) {
		world := New()
		world.AddSystemToGroup(&HMovementSystem{}, "special group")

		world.Add(entitiy.New().Add(Pos{x: 0, y: 0}).Add(Vel{x: 1, y: 1}))
		world.Add(entitiy.New().Add(Pos{x: 2, y: 2}))
		world.Add(entitiy.New().Add(Pos{x: 3, y: 3}).Add(Vel{x: 4, y: 4}))

		_ = world.UpdateGroup("special group", 0)

		expectPositions(t, world, []Pos{
			{x: 1, y: 0},
			{x: 2, y: 2},
			{x: 7, y: 3},
		})
	})

	t.Run("update only to no groups", func(t *testing.T) {
		world := New()
		world.AddSystemToGroup(&HMovementSystem{}, "special group")

		world.Add(entitiy.New().Add(Pos{x: 0, y: 0}).Add(Vel{x: 1, y: 1}))
		world.Add(entitiy.New().Add(Pos{x: 2, y: 2}))
		world.Add(entitiy.New().Add(Pos{x: 3, y: 3}).Add(Vel{x: 4, y: 4}))

		_ = world.UpdateGroup(system.DefaultGroup, 0)

		expectPositions(t, world, []Pos{
			{x: 0, y: 0},
			{x: 2, y: 2},
			{x: 3, y: 3},
		})
	})
}

func expectPositions(t *testing.T, world *World, want []Pos) {
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

	world.AddSystemToGroup(&HMovementSystem{}, "special")

	world.Add(entitiy.New().Add(Pos{x: 0, y: 0}).Add(Vel{x: 1, y: 1}))
	world.Add(entitiy.New().Add(Pos{x: 2, y: 2}))
	world.Add(entitiy.New().Add(Pos{x: 3, y: 3}).Add(Vel{x: 4, y: 4}))

	s := world.String()

	if len(s) == 0 {
		t.Fatalf("shoudl get string, got empty")
	}
}

func TestWorld_Update_Error(t *testing.T) {
	world := New()
	world.AddSystem(&FailureSystem{})
	world.AddSystem(&HMovementSystem{})
	world.AddSystem(&VMovementSystem{})

	world.AddSystemToGroup(&HMovementSystem{}, "special")

	world.Add(entitiy.New().Add(Pos{x: 0, y: 0}).Add(Vel{x: 1, y: 1}))

	e := world.Update(0)

	if !errors.Is(e, errFailure) {
		t.Fatalf("shoudl get failure but got %v", e)
	}

	expectPositions(t, world, []Pos{
		{x: 0, y: 0},
	})

	e = world.UpdateGroup("special", 0)

	if e != nil {
		t.Fatalf("shoudl not get error but got %v", e)
	}

	expectPositions(t, world, []Pos{
		{x: 1, y: 0},
	})
}
