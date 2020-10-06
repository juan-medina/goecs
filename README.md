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

[Run it on the Go Playground](https://play.golang.org/p/9urqRvY8u-v)
```go
package main

import (
	"fmt"
	"github.com/juan-medina/goecs"
	"reflect"
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
		world.Signal(PosChangeSignal{From: pos, To: npos})

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

// PosChangeSignalType is the type of the PosChangeSignal
var PosChangeSignalType = reflect.TypeOf(PosChangeSignal{})
```
This will output:

```
pos change from {0 0} to {0.5 0.5}
pos change from {2 2} to {3.5 3.5}

Pos:
{0.5 0.5}
{2 2}
{3.5 3.5}
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
