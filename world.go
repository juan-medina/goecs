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

package goecs

import (
	"fmt"
	"reflect"
)

const (
	defaultPriority = int32(0)
)

// Default values for Default()
const (
	DefaultSignalsInitialCapacity   = 20   // Default Signals initial capacity
	DefaultSystemsInitialCapacity   = 50   // Default System initial capacity
	DefaultListenersInitialCapacity = 50   // Default Listener initial capacity
	DefaultEntitiesInitialCapacity  = 2000 // Default Entity initial capacity
)

// World is a view.View that contains the Entity and System of our ECS
type World struct {
	*View
	systems       *Systems       // systems registration of System
	subscriptions *Subscriptions // subscriptions of Listener to signals
}

// String get a string representation of our World
func (world World) String() string {
	result := ""

	result += fmt.Sprintf("World{view: %v, systems: [", world.View)

	result += world.systems.String() + "], "

	result += " subscriptions: [" + world.subscriptions.String() + "]"

	result += "}"

	return result
}

// AddSystem adds the given System to the world
func (world *World) AddSystem(sys System) {
	world.AddSystemWithPriority(sys, defaultPriority)
}

// AddSystemWithPriority adds the given System to the world with a priority
func (world *World) AddSystemWithPriority(sys System, priority int32) {
	world.systems.Register(sys, priority)
}

// AddListener adds the given Listener to the world
func (world *World) AddListener(lis Listener, signals ...reflect.Type) {
	world.AddListenerWithPriority(lis, defaultPriority, signals...)
}

// AddListenerWithPriority adds the given Listener to the world with a priority
func (world *World) AddListenerWithPriority(lis Listener, priority int32, signals ...reflect.Type) {
	world.subscriptions.Subscribe(lis, priority, signals...)
}

// Update ask to update the System and send the signals
func (world *World) Update(delta float32) error {
	// update the systems
	if err := world.systems.Update(world, delta); err != nil {
		return err
	}

	// update the subscriptions
	if err := world.subscriptions.Update(world, delta); err != nil {
		return err
	}
	return nil
}

// Signal to be sent
func (world *World) Signal(signal interface{}) {
	world.subscriptions.Signal(signal)
}

// Clear removes all System, Listener, Subscriptions and Entity from the World
func (world *World) Clear() {
	world.systems.Clear()
	world.subscriptions.Clear()
	world.View.Clear()
}

// Default creates a default World with a initial capacity
//  const (
//  	DefaultSignalsInitialCapacity   = 20   // Default Signals initial capacity
//  	DefaultSystemsInitialCapacity   = 50   // Default System initial capacity
//  	DefaultListenersInitialCapacity = 50   // Default Listener initial capacity
//  	DefaultEntitiesInitialCapacity  = 2000 // Default Entity initial capacity
//  )
func Default() *World {
	return New(
		DefaultEntitiesInitialCapacity,
		DefaultSystemsInitialCapacity,
		DefaultListenersInitialCapacity,
		DefaultSignalsInitialCapacity,
	)
}

// New creates World with a giving initial capacity of entities, systems, listeners and signals
//
// Since those elements are sparse.Slice the will grow dynamically
func New(entities, systems, listeners, signals int) *World {
	return &World{
		View:          NewView(entities),
		systems:       NewSystems(systems),
		subscriptions: NewSubscriptions(listeners, signals),
	}
}
