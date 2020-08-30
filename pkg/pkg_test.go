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

package pkg_test

import (
	"fmt"
	"github.com/juan-medina/goecs/pkg/entity"
	"github.com/juan-medina/goecs/pkg/world"
	"reflect"
)

// Simple Usage
func Example() {
	// creates the world
	wld := world.New()

	// ad our movement system
	wld.AddSystem(&MovementSystem{})

	// add a first entity
	wld.Add(entity.New(
		Pos{X: 0, Y: 0},
		Vel{X: 1, Y: 1},
	))

	// this entity shouldn't be updated
	wld.Add(entity.New(
		Pos{X: 2, Y: 2},
	))

	// add a third entity
	wld.Add(entity.New(
		Pos{X: 2, Y: 2},
		Vel{X: 3, Y: 3},
	))

	// ask the world to update
	_ = wld.Update(0.5)

	// print the world
	fmt.Printf("world: %v", wld)
}

// MovementSystem is a simple movement system
type MovementSystem struct{}

// Update the ECS world.World
func (m *MovementSystem) Update(wld *world.World, delta float32) error {
	for it := wld.Iterator(PosType, VelType); it.HasNext(); {
		ent := it.Value()
		pos := ent.Get(PosType).(Pos)
		vel := ent.Get(VelType).(Vel)

		pos.X += vel.X * delta
		pos.Y += vel.Y * delta

		ent.Set(pos)
	}

	return nil
}

// Notify get call when events trigger
func (m *MovementSystem) Notify(world *world.World, event interface{}, delta float32) error {
	return nil
}

// Pos represent a 2D position
type Pos struct {
	X float32
	Y float32
}

// PosType is the reflect.Type of Pos
var PosType = reflect.TypeOf(Pos{})

// Vel represent a 2D velocity
type Vel struct {
	X float32
	Y float32
}

// VelType is the reflect.Type of Vel
var VelType = reflect.TypeOf(Vel{})
