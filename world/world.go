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

// Package world allows to create an use our ECS World
package world

import (
	"fmt"
	"github.com/juan-medina/goecs/sparse"
	"github.com/juan-medina/goecs/view"
	"reflect"
	"runtime"
	"sort"
)

// System get invoke with Update() from a World
type System func(wld *World, delta float32) error

// Listener the get notified that a new signal has been received by World.Signal
type Listener func(wld *World, signal interface{}, delta float32) error

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

// World is a view.View that contains the entity.Entity and System of our ECS
type World struct {
	*view.View
	systems   sparse.Slice
	listeners sparse.Slice
	events    sparse.Slice
}

// String get a string representation of our World
func (wld World) String() string {
	result := ""

	result += fmt.Sprintf("World{view: %v, systems: [", wld.View)

	str := ""
	for it := wld.systems.Iterator(); it != nil; it = it.Next() {
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
	for it := wld.listeners.Iterator(); it != nil; it = it.Next() {
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
func (wld *World) AddSystem(sys System) {
	wld.AddSystemWithPriority(sys, defaultPriority)
}

// AddSystemWithPriority adds the given System to the world
func (wld *World) AddSystemWithPriority(sys System, priority int32) {
	wld.systems.Add(systemWithPriority{
		system:   sys,
		priority: priority,
		id:       lastSystemID,
	})
	lastSystemID++
}

// Listen adds the given System to the world
func (wld *World) Listen(lis Listener) {
	wld.ListenWithPriority(lis, defaultPriority)
}

// ListenWithPriority adds the given System to the world
func (wld *World) ListenWithPriority(lis Listener, priority int32) {
	wld.listeners.Add(listenerWithPriority{
		listener: lis,
		priority: priority,
		id:       lastListenerID,
	})
	lastListenerID++
}

// sendEvents send the pending events to the System on the world
func (wld *World) sendEvents(delta float32) error {
	// for hold a copy of the events
	events := make([]interface{}, wld.events.Size())

	// get all events for this hold
	i := 0
	for it := wld.events.Iterator(); it != nil; it = it.Next() {
		events[i] = it.Value()
		i++
	}

	// clear the hold
	wld.events.Clear()

	// get the listener list
	listeners := wld.getListenersPriorityList()

	for _, e := range events {
		// range systems
		for _, l := range listeners {
			// notify the event to the system
			if err := l.listener(wld, e, delta); err != nil {
				return err
			}
		}
	}

	return nil
}

func (wld World) getSystemsPriorityList() []systemWithPriority {
	result := make([]systemWithPriority, wld.systems.Size())

	i := 0
	for it := wld.systems.Iterator(); it != nil; it = it.Next() {
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

func (wld World) getListenersPriorityList() []listenerWithPriority {
	result := make([]listenerWithPriority, wld.listeners.Size())

	i := 0
	for it := wld.listeners.Iterator(); it != nil; it = it.Next() {
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
func (wld *World) Update(delta float32) error {
	pl := wld.getSystemsPriorityList()
	for _, s := range pl {
		if err := s.system(wld, delta); err != nil {
			return err
		}
	}

	if err := wld.sendEvents(delta); err != nil {
		return err
	}

	return nil
}

// Signal signal an event to be sent
func (wld *World) Signal(event interface{}) error {
	// add the event
	wld.events.Add(event)

	return nil
}

// Clear removes all world.System and entity.Entity from the world.World
func (wld *World) Clear() {
	wld.systems.Clear()
	wld.View.Clear()
}

// New creates a new World
func New() *World {
	return &World{
		View:      view.New(),
		systems:   sparse.NewSlice(systemsInitialCapacity, systemsCapacityGrow),
		listeners: sparse.NewSlice(listenersInitialCapacity, listenersCapacityGrow),
		events:    sparse.NewSlice(eventsInitialCapacity, eventsCapacityGrow),
	}
}
