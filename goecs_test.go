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
	"fmt"
	"github.com/juan-medina/goecs"
	"reflect"
)

// Simple Usage
func Example() {
	// creates the world
	world := goecs.New()

	// add our movement system
	world.AddSystem(MovementSystem)

	// add a listener
	world.AddListener(ChangePostListener)

	// add a first entity
	world.AddEntity(
		Pos{X: 0, Y: 0},
		Vel{X: 1, Y: 1},
	)

	// this entity shouldn't be updated
	world.AddEntity(
		Pos{X: 2, Y: 2},
	)

	// add a third entity
	world.AddEntity(
		Pos{X: 2, Y: 2},
		Vel{X: 3, Y: 3},
	)

	// ask the world to update
	if err := world.Update(0.5); err != nil {
		fmt.Printf("error on update %v\n", err)
	}

	// print the world
	fmt.Println()
	fmt.Println("Pos:")
	for it := world.Iterator(PosType); it != nil; it = it.Next() {
		fmt.Println(it.Value().Get(PosType))
	}
	// Output:
	// pos change from {0 0} to {0.5 0.5}
	// pos change from {2 2} to {3.5 3.5}
	//
	// Pos:
	// {0.5 0.5}
	// {2 2}
	// {3.5 3.5}
}

// MovementSystem is a simple movement system
func MovementSystem(world *goecs.World, delta float32) error {
	// get all the entities that we need to update, only if they have Pos & Vel
	for it := world.Iterator(PosType, VelType); it != nil; it = it.Next() {
		// get the values
		ent := it.Value()
		pos := ent.Get(PosType).(Pos)
		vel := ent.Get(VelType).(Vel)

		// calculate new pos
		npos := Pos{
			X: pos.X + (vel.X * delta),
			Y: pos.Y + vel.Y*delta,
		}

		// signal the change
		if err := world.Signal(PosChangeSignal{From: pos, To: npos}); err != nil {
			return err
		}

		// set the new pos
		ent.Set(npos)
	}

	return nil
}

// ChangePostListener listen to PosChangeSignal
func ChangePostListener(world *goecs.World, signal interface{}, delta float32) error {
	switch s := signal.(type) {
	case PosChangeSignal:
		// print the change
		fmt.Printf("pos change from %v to %v\n", s.From, s.To)
	}
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

// PosChangeSignal is a signal that a Pos has change
type PosChangeSignal struct {
	From Pos
	To   Pos
}
