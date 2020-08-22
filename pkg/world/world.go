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
	"github.com/juan-medina/goecs/pkg/view"
	"reflect"
	"sort"
)

type systemWithPriority struct {
	system   System
	priority int32
}

const defaultPriority = int32(0)

// World is a view.View that contains the entity.Entity and System of our ECS
type World struct {
	*view.View
	systems []systemWithPriority
	events  []interface{}
}

// String get a string representation of our World
func (wld World) String() string {
	result := ""

	result += fmt.Sprintf("World[view: %v, systems: [", wld.View)

	for _, s := range wld.systems {
		result += fmt.Sprintf("%s,", reflect.TypeOf(s).String())
	}

	result += "]"

	return result
}

// AddSystem adds the given System to the world
func (wld *World) AddSystem(sys System) {
	wld.systems = append(wld.systems, systemWithPriority{
		system:   sys,
		priority: defaultPriority,
	})
}

// AddSystemWithPriority adds the given System to the world
func (wld *World) AddSystemWithPriority(sys System, priority int32) {
	wld.systems = append(wld.systems, systemWithPriority{
		system:   sys,
		priority: priority,
	})
}

// sendEvents send the pending events to the System on the world
func (wld *World) sendEvents(delta float32, systems []systemWithPriority) error {
	// get all events for this hold
	for i, e := range wld.events {
		// range systems
		for _, s := range systems {
			// notify the event to the system
			if err := s.system.Notify(wld, e, delta); err != nil {
				return err
			}
		}
		//clear the event
		wld.events[i] = nil
	}

	//empty hold
	wld.events = wld.events[:0]

	return nil
}

func (wld World) getPriorityList() []systemWithPriority {
	result := make([]systemWithPriority, len(wld.systems))

	for i, v := range wld.systems {
		result[i] = v
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].priority > result[j].priority
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
	wld.events = append(wld.events, event)

	return nil
}

// New creates a new World
func New() *World {
	return &World{
		View:    view.New(),
		systems: make([]systemWithPriority, 0),
		events:  make([]interface{}, 0),
	}
}
