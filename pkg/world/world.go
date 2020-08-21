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
)

type eventHold []interface{}

// World is a view.View that contains the entity.Entity and System of our ECS
type World struct {
	*view.View
	systems map[string][]System
	events  map[string]eventHold
}

const (
	defaultSystemGroup = "DEFAULT_GROUP"
)

// String get a string representation of our World
func (wld World) String() string {
	result := ""

	result += fmt.Sprintf("World[view: %v, systems: [", wld.View)

	for g := range wld.systems {
		result += fmt.Sprintf("%s:[", g)
		for _, s := range wld.systems[g] {
			result += fmt.Sprintf("%s,", reflect.TypeOf(s).String())
		}
		result += "],"
	}
	result += "]"

	return result
}

// AddSystemToGroup adds the given System to a given group
func (wld *World) AddSystemToGroup(sys System, group string) {
	if _, ok := wld.systems[group]; !ok {
		wld.systems[group] = make([]System, 0)
	}
	wld.systems[group] = append(wld.systems[group], sys)
}

// AddSystem adds the given System to the default group
func (wld *World) AddSystem(sys System) {
	wld.AddSystemToGroup(sys, defaultSystemGroup)
}

// sendGroupEvents send the pending events to the System on the given group
func (wld *World) sendGroupEvents(group string, delta float32) error {
	//if this group has a event hold
	if h, ok := wld.events[group]; ok {
		// get all events for this hold
		for i, e := range h {
			// range systems on this group
			for _, s := range wld.systems[group] {
				// notify the event to the system
				if err := s.Notify(wld, e, delta); err != nil {
					return err
				}
			}
			//clear the event
			wld.events[group][i] = nil
		}
		//empty hold
		wld.events[group] = wld.events[group][:0]
	}

	return nil
}

// UpdateGroup ask to update the System on the given group, and send the pending events
func (wld *World) UpdateGroup(group string, delta float32) error {
	for _, s := range wld.systems[group] {
		if err := s.Update(wld, delta); err != nil {
			return err
		}
	}

	if err := wld.sendGroupEvents(group, delta); err != nil {
		return err
	}

	return nil
}

// Update ask to update the System on the default group and send the pending events
func (wld *World) Update(delta float32) error {
	return wld.UpdateGroup(defaultSystemGroup, delta)
}

// NotifyGroup add an event to be sent to the given group
func (wld *World) NotifyGroup(group string, event interface{}) error {
	// if this group has not a event hold for this group create it
	if _, ok := wld.events[group]; !ok {
		wld.events[group] = make([]interface{}, 0)
	}
	// add the event
	wld.events[group] = append(wld.events[group], event)

	return nil
}

// Notify add an event to be sent to the default group
func (wld *World) Notify(event interface{}) error {
	return wld.NotifyGroup(defaultSystemGroup, event)
}

// New creates a new World
func New() *World {
	return &World{
		View:    view.New(),
		systems: make(map[string][]System, 0),
		events:  make(map[string]eventHold, 0),
	}
}
