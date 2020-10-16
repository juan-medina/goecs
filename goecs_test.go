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
)

// Simple Usage
func Example() {
	// creates the world
	world := goecs.Default()

	// add our movement system
	world.AddSystem(MovementSystem)

	// add a listener
	world.AddListener(ChangePostListener, PosChangeSignalType)

	// add a first entity
	world.AddEntity(
		Pos{X: 0, Y: 0},
		Vel{X: 2, Y: 2},
	)

	// this entity shouldn't be updated
	world.AddEntity(
		Pos{X: 6, Y: 6},
	)

	// add a third entity
	world.AddEntity(
		Pos{X: 8, Y: 8},
		Vel{X: 4, Y: 4},
	)

	// print the world
	PrintWorld(world)
	fmt.Println()
	fmt.Println("updating world for halve a second:")

	// ask the world to update
	if err := world.Update(0.5); err != nil {
		fmt.Printf("error on update %v\n", err)
	}

	// print the world
	fmt.Println()
	PrintWorld(world)

	// Output:
	// World:
	// Id: 0, Pos: {0 0}, Vel: {2 2}/s
	// Id: 1, Pos: {6 6}
	// Id: 2, Pos: {8 8}, Vel: {4 4}/s
	//
	// updating world for halve a second:
	// pos change for id: 0, from Pos{0 0} to Pos{1 1}
	// pos change for id: 2, from Pos{8 8} to Pos{10 10}
	//
	// World:
	// Id: 0, Pos: {1 1}, Vel: {2 2}/s
	// Id: 1, Pos: {6 6}
	// Id: 2, Pos: {10 10}, Vel: {4 4}/s
}

// PrintWorld prints the content of our world
func PrintWorld(world *goecs.World) {
	fmt.Println("World:")
	for it := world.Iterator(); it != nil; it = it.Next() {
		ent := it.Value()
		id := ent.ID()
		pos := ent.Get(PosType).(Pos)
		if ent.Contains(VelType) {
			vel := ent.Get(VelType).(Vel)
			fmt.Printf("Id: %d, Pos: %v, Vel: %v/s\n", id, pos, vel)
		} else {
			fmt.Printf("Id: %d, Pos: %v\n", id, pos)
		}
	}
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
			X: pos.X + vel.X*delta,
			Y: pos.Y + vel.Y*delta,
		}

		// set the new pos
		ent.Set(npos)

		// signal the change
		world.Signal(PosChangeSignal{ID: ent.ID(), From: pos, To: npos})
	}

	return nil
}

// ChangePostListener listen to PosChangeSignal
func ChangePostListener(world *goecs.World, signal goecs.Component, delta float32) error {
	switch s := signal.(type) {
	case PosChangeSignal:
		// print the change
		fmt.Printf("pos change for id: %d, from Pos%v to Pos%v\n", s.ID, s.From, s.To)
	}
	return nil
}

// PosType is the ComponentType of Pos
var PosType = goecs.NewComponentType()

// Pos represent a 2D position
type Pos struct {
	X float32
	Y float32
}

// Type will return Pos goecs.ComponentType
func (p Pos) Type() goecs.ComponentType {
	return PosType
}

// VelType is the ComponentType of Vel
var VelType = goecs.NewComponentType()

// Vel represent a 2D velocity
type Vel struct {
	X float32
	Y float32
}

// Type will return Vel goecs.ComponentType
func (v Vel) Type() goecs.ComponentType {
	return VelType
}

// PosChangeSignalType is the type of the PosChangeSignal
var PosChangeSignalType = goecs.NewComponentType()

// PosChangeSignal is a signal that a Pos has change
type PosChangeSignal struct {
	ID   uint64
	From Pos
	To   Pos
}

// Type will return PosChangeSignal goecs.ComponentType
func (p PosChangeSignal) Type() goecs.ComponentType {
	return PosChangeSignalType
}
