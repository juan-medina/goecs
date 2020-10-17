# goecs
A simple Go [ECS](https://en.wikipedia.org/wiki/Entity_component_system)

[![License: Apache2](https://img.shields.io/badge/license-Apache%202-blue.svg)](/LICENSE)
[![version](https://img.shields.io/github/v/tag/juan-medina/goecs?label=version)](https://pkg.go.dev/mod/github.com/juan-medina/goecs?tab=versions)
[![godoc](https://godoc.org/github.com/juan-medina/goecs?status.svg)](https://pkg.go.dev/mod/github.com/juan-medina/goecs)
[![Build Status](https://travis-ci.com/juan-medina/goecs.svg?branch=main)](https://travis-ci.com/juan-medina/goecs)
[![codecov](https://codecov.io/gh/juan-medina/goecs/branch/main/graph/badge.svg)](https://codecov.io/gh/juan-medina/goecs)
[![conduct](https://img.shields.io/badge/code%20of%20conduct-contributor%20covenant%202.0-purple.svg?style=flat-square)](https://www.contributor-covenant.org/version/2/0/code_of_conduct/)



## Info
Entity–component–system (ECS) is an architectural patter that follows the composition over inheritance principle that allows greater flexibility in defining entities where every object in a world.

Every entity consists of one or more components which contains data or state. Therefore, the behavior of an entity can be changed at runtime by systems that add, remove or mutate components.

This eliminates the ambiguity problems of deep and wide inheritance hierarchies that are difficult to understand, maintain and extend.

Common ECS approaches are highly compatible and often combined with data-oriented design techniques.

For a more in deep read on this topic I could recommend this [article](https://medium.com/ingeniouslysimple/entities-components-and-systems-89c31464240d).

## Example

[Run it on the Go Playground](https://play.golang.org/p/l3wgiqeeBf6)
```go
package main

import (
	"fmt"
	"github.com/juan-medina/goecs"
)

// Simple Usage
func main() {
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

```
This will output:
```
World:
Id: 1, Pos: {0 0}, Vel: {2 2}/s
Id: 2, Pos: {6 6}
Id: 3, Pos: {8 8}, Vel: {4 4}/s

updating world for halve a second:
pos change for id: 1, from Pos{0 0} to Pos{1 1}
pos change for id: 3, from Pos{8 8} to Pos{10 10}

World:
Id: 1, Pos: {1 1}, Vel: {2 2}/s
Id: 2, Pos: {6 6}
Id: 3, Pos: {10 10}, Vel: {4 4}/s
```

## Installation

```bash
go get -v -u github.com/juan-medina/goecs
```

## License

```text
    Copyright (C) 2020 Juan Medina

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
```
