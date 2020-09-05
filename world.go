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
	"sort"
)

// System get invoke with Update() from a World
type System func(world *World, delta float32) error

// Listener the get notified that a new signal has been received by World.Signal
type Listener func(world *World, signal interface{}, delta float32) error

type systemWithPriority struct {
	system   System
	priority int32
	id       int64
}

type listenerWithPriority struct {
	listener Listener
	priority int32
	id       int64
}

const (
	defaultPriority          = int32(0)
	eventsInitialCapacity    = 100
	eventsCapacityGrow       = eventsInitialCapacity / 4
	systemsInitialCapacity   = 100
	systemsCapacityGrow      = systemsInitialCapacity / 4
	listenersInitialCapacity = 100
	listenersCapacityGrow    = systemsInitialCapacity / 4
)

var (
	lastSystemID   = int64(0)
	lastListenerID = int64(0)
)

// World is a view.View that contains the Entity and System of our ECS
type World struct {
	*View
	systems   sparse.Slice
	listeners sparse.Slice
	events    sparse.Slice
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

	result += str + " listeners: ["

	str = ""
	for it := world.listeners.Iterator(); it != nil; it = it.Next() {
		l := it.Value().(listenerWithPriority)
		if str != "" {
			str += ","
		}
		name := runtime.FuncForPC(reflect.ValueOf(l.listener).Pointer()).Name()
		str += fmt.Sprintf("{%s}", name)
	}
	str += "]"

	result += str + "}"

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
}

// AddListener adds the given Listener to the world
func (world *World) AddListener(lis Listener) {
	world.AddListenerWithPriority(lis, defaultPriority)
}

// AddListenerWithPriority adds the given Listener to the world with a priority
func (world *World) AddListenerWithPriority(lis Listener, priority int32) {
	world.listeners.Add(listenerWithPriority{
		listener: lis,
		priority: priority,
		id:       lastListenerID,
	})
	lastListenerID++
}

// sendEvents send the pending events to the System on the world
func (world *World) sendEvents(delta float32) error {
	// for hold a copy of the events
	events := make([]interface{}, world.events.Size())

	// get all events for this hold
	i := 0
	for it := world.events.Iterator(); it != nil; it = it.Next() {
		events[i] = it.Value()
		i++
	}

	// clear the hold
	world.events.Clear()

	// get the listener list
	listeners := world.getListenersPriorityList()

	for _, e := range events {
		// range systems
		for _, l := range listeners {
			// notify the event to the system
			if err := l.listener(world, e, delta); err != nil {
				return err
			}
		}
	}

	return nil
}

func (world World) getSystemsPriorityList() []systemWithPriority {
	result := make([]systemWithPriority, world.systems.Size())

	i := 0
	for it := world.systems.Iterator(); it != nil; it = it.Next() {
		result[i] = it.Value().(systemWithPriority)
		i++
	}

	sort.Slice(result, func(i, j int) bool {
		first := result[i]
		second := result[j]
		if first.priority == second.priority {
			return first.id < second.id
		}
		return first.priority > second.priority
	})

	return result
}

func (world World) getListenersPriorityList() []listenerWithPriority {
	result := make([]listenerWithPriority, world.listeners.Size())

	i := 0
	for it := world.listeners.Iterator(); it != nil; it = it.Next() {
		result[i] = it.Value().(listenerWithPriority)
		i++
	}

	sort.Slice(result, func(i, j int) bool {
		first := result[i]
		second := result[j]
		if first.priority == second.priority {
			return first.id < second.id
		}
		return first.priority > second.priority
	})

	return result
}

// Update ask to update the System send the pending events
func (world *World) Update(delta float32) error {
	pl := world.getSystemsPriorityList()
	for _, s := range pl {
		if err := s.system(world, delta); err != nil {
			return err
		}
	}

	if err := world.sendEvents(delta); err != nil {
		return err
	}

	return nil
}

// Signal signal an event to be sent
func (world *World) Signal(event interface{}) error {
	// add the event
	world.events.Add(event)

	return nil
}

// Clear removes all world.System and Entity from the World
func (world *World) Clear() {
	world.systems.Clear()
	world.View.Clear()
}

// New creates a new World
func New() *World {
	return &World{
		View:      NewView(),
		systems:   sparse.NewSlice(systemsInitialCapacity, systemsCapacityGrow),
		listeners: sparse.NewSlice(listenersInitialCapacity, listenersCapacityGrow),
		events:    sparse.NewSlice(eventsInitialCapacity, eventsCapacityGrow),
	}
}
