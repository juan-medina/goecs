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
	"github.com/juan-medina/goecs/pkg/sparse"
	"github.com/juan-medina/goecs/pkg/view"
	"reflect"
	"sort"
)

type systemWithPriority struct {
	system   System
	priority int32
	id       int64
}

const (
	defaultPriority        = int32(0)
	eventsInitialCapacity  = 10
	eventsCapacityGrow     = eventsInitialCapacity / 4
	systemsInitialCapacity = 100
	systemsCapacityGrow    = systemsInitialCapacity / 4
)

var (
	lastID = int64(0)
)

// World is a view.View that contains the entity.Entity and System of our ECS
type World struct {
	*view.View
	systems sparse.Slice
	events  sparse.Slice
}

// String get a string representation of our World
func (wld World) String() string {
	result := ""

	result += fmt.Sprintf("World[view: %v, systems: [", wld.View)

	for it := wld.systems.Iterator(); it.HasNext(); {
		s := it.Value().(systemWithPriority)
		result += fmt.Sprintf("%s,", reflect.TypeOf(s.system).String())
	}

	result += "]"

	return result
}

// AddSystem adds the given System to the world
func (wld *World) AddSystem(sys System) {
	wld.AddSystemWithPriority(sys, defaultPriority)
}

// RemoveSystem deletes the given System from the world
func (wld *World) RemoveSystem(sys System) error {
	for it := wld.systems.Iterator(); it.HasNext(); {
		s := it.Value().(systemWithPriority)
		if s.system == sys {
			return wld.systems.Remove(s)
		}
	}

	return sparse.ErrItemNotFound
}

// AddSystemWithPriority adds the given System to the world
func (wld *World) AddSystemWithPriority(sys System, priority int32) {
	wld.systems.Add(systemWithPriority{
		system:   sys,
		priority: priority,
		id:       lastID,
	})
	lastID++
}

// sendEvents send the pending events to the System on the world
func (wld *World) sendEvents(delta float32, systems []systemWithPriority) error {
	// get all events for this hold
	for it := wld.events.Iterator(); it.HasNext(); {
		e := it.Value()
		// range systems
		for _, s := range systems {
			// notify the event to the system
			if err := s.system.Notify(wld, e, delta); err != nil {
				return err
			}
		}
	}

	//empty hold
	wld.events.Clear()

	return nil
}

func (wld World) getPriorityList() []systemWithPriority {
	result := make([]systemWithPriority, wld.systems.Size())

	i := 0
	for it := wld.systems.Iterator(); it.HasNext(); {
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

// Update ask to update the System send the pending events
func (wld *World) Update(delta float32) error {
	pl := wld.getPriorityList()
	for _, s := range pl {
		if err := s.system.Update(wld, delta); err != nil {
			return err
		}
	}

	if err := wld.sendEvents(delta, pl); err != nil {
		return err
	}

	return nil
}

// Notify add an event to be sent
func (wld *World) Notify(event interface{}) error {
	// add the event
	wld.events.Add(event)

	return nil
}

// New creates a new World
func New() *World {
	return &World{
		View:    view.New(),
		systems: sparse.NewSlice(systemsInitialCapacity, systemsCapacityGrow),
		events:  sparse.NewSlice(eventsInitialCapacity, eventsCapacityGrow),
	}
}
