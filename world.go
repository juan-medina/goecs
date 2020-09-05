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
	signalsInitialCapacity   = 100
	signalsCapacityGrow      = signalsInitialCapacity / 4
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
	systems   sparse.Slice // sparse.Slice of systemWithPriority
	listeners sparse.Slice // sparse.Slice of listenerWithPriority
	signals   sparse.Slice // sparse.Slice of signals
	toSend    sparse.Slice // sparse.Slice of signals is a copy to signals to be send
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
	// sort systems, is better to keep them sorted that sort them on update
	world.systems.Sort(world.sortSystemsByPriority)
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
	// sort listeners, is better to keep them sorted that sort them on update
	world.listeners.Sort(world.sortListenersByPriority)
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

	var l listenerWithPriority
	// get thee signals to send
	for ite := world.toSend.Iterator(); ite != nil; ite = ite.Next() {
		// range systems
		for itl := world.listeners.Iterator(); itl != nil; itl = itl.Next() {
			// get the listener
			l = itl.Value().(listenerWithPriority)
			// notify the signal to the listener
			if err := l.listener(world, ite.Value(), delta); err != nil {
				return err
			}
		}
	}

	// clear the signals to be send
	world.toSend.Clear()

	return nil
}

func (world World) sortSystemsByPriority(a interface{}, b interface{}) bool {
	first := a.(systemWithPriority)
	second := b.(systemWithPriority)
	if first.priority == second.priority {
		return first.id < second.id
	}
	return first.priority > second.priority
}

func (world World) sortListenersByPriority(a interface{}, b interface{}) bool {
	first := a.(listenerWithPriority)
	second := b.(listenerWithPriority)
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
	world.listeners.Clear()
	world.signals.Clear()
	world.toSend.Clear()
	world.View.Clear()
}

// New creates a new World
func New() *World {
	return &World{
		View:      NewView(),
		systems:   sparse.NewSlice(systemsInitialCapacity, systemsCapacityGrow),
		listeners: sparse.NewSlice(listenersInitialCapacity, listenersCapacityGrow),
		signals:   sparse.NewSlice(signalsInitialCapacity, signalsCapacityGrow),
		toSend:    sparse.NewSlice(signalsInitialCapacity, signalsCapacityGrow),
	}
}
