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
	"github.com/juan-medina/goecs/sparse"
	"reflect"
	"runtime"
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

var (
	lastSystemID = int64(0)
)

// World is a view.View that contains the Entity and System of our ECS
type World struct {
	*View
	systems       sparse.Slice   // sparse.Slice of systemWithPriority
	subscriptions *Subscriptions // subscriptions of Listeners to signals
	signals       sparse.Slice   // sparse.Slice of signals
	toSend        sparse.Slice   // sparse.Slice of signals is a copy to signals to be send
}

// String get a string representation of our World
func (world World) String() string {
	result := ""

	result += fmt.Sprintf("World{view: %v, systems: [", world.View)

	str := ""
	for it := world.systems.Iterator(); it != nil; it = it.Next() {
		s := it.Value().(systemWithPriority)
		if str != "" {
			str += ","
		}
		name := runtime.FuncForPC(reflect.ValueOf(s.system).Pointer()).Name()
		str += fmt.Sprintf("{%s}", name)
	}
	str += "]"

	result += str + " subscriptions: [" + world.subscriptions.String() + "]"

	result += "}"

	return result
}

// AddSystem adds the given System to the world
func (world *World) AddSystem(sys System) {
	world.AddSystemWithPriority(sys, defaultPriority)
}

// AddSystemWithPriority adds the given System to the world with a priority
func (world *World) AddSystemWithPriority(sys System, priority int32) {
	world.systems.Add(systemWithPriority{
		system:   sys,
		priority: priority,
		id:       lastSystemID,
	})
	lastSystemID++
	// sort systems, is better to keep them sorted that sort them on update
	world.systems.Sort(world.sortSystemsByPriority)
}

// AddListener adds the given Listener to the world
func (world *World) AddListener(lis Listener, signals ...reflect.Type) {
	world.AddListenerWithPriority(lis, defaultPriority, signals...)
}

// AddListenerWithPriority adds the given Listener to the world with a priority
func (world *World) AddListenerWithPriority(lis Listener, priority int32, signals ...reflect.Type) {
	world.subscriptions.Subscribe(lis, priority, signals...)
}

// sendSignals send the pending signals to the System on the world
func (world *World) sendSignals(delta float32) error {
	// avoid to copy empty signals
	if world.signals.Size() == 0 {
		return nil
	}
	// replace the signals to send, so we do not send the signals triggered by the current signals
	world.signals.Replace(world.toSend)

	// clear the hold so new signals will be here
	world.signals.Clear()

	var err error
	// get thee signals to send
	for ite := world.toSend.Iterator(); ite != nil; ite = ite.Next() {
		if err = world.subscriptions.Process(world, ite.Value(), delta); err != nil {
			return err
		}
	}

	// clear the signals to be send
	world.toSend.Clear()

	return nil
}

func (world World) sortSystemsByPriority(a, b interface{}) bool {
	first := a.(systemWithPriority)
	second := b.(systemWithPriority)
	if first.priority == second.priority {
		return first.id < second.id
	}
	return first.priority > second.priority
}

// Update ask to update the System send the pending signals
func (world *World) Update(delta float32) error {
	var s systemWithPriority
	// go trough the systems
	for it := world.systems.Iterator(); it != nil; it = it.Next() {
		s = it.Value().(systemWithPriority)
		if err := s.system(world, delta); err != nil {
			return err
		}
	}

	if err := world.sendSignals(delta); err != nil {
		return err
	}

	return nil
}

// Signal to be sent
func (world *World) Signal(signal interface{}) error {
	// add the signal
	world.signals.Add(signal)

	return nil
}

// Clear removes all System, Listener, Signals and Entity from the World
func (world *World) Clear() {
	world.systems.Clear()
	world.subscriptions.Clear()
	world.signals.Clear()
	world.toSend.Clear()
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
		systems:       sparse.NewSlice(systems),
		subscriptions: NewSubscriptions(listeners),
		signals:       sparse.NewSlice(signals),
		toSend:        sparse.NewSlice(signals),
	}
}
