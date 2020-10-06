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
	"github.com/juan-medina/goecs/sparse"
	"reflect"
	"runtime"
)

// System get invoke with Update() from a World
type System func(world *World, delta float32) error

// systemRegistration hold the registration of a system
type systemRegistration struct {
	system   System // system registered
	priority int32  // priority for this system
	id       int64  // this system id
}

// Systems manage registration of systems
type Systems struct {
	registrations      sparse.Slice // registrations of System
	lastRegistrationID int64        // lastRegistrationID is the id of the last registration
}

// Register adds a new registration with a given priority
func (sys *Systems) Register(system System, priority int32) {
	// increment the id
	sys.lastRegistrationID++
	// add the registration
	sys.registrations.Add(systemRegistration{
		id:       sys.lastRegistrationID,
		system:   system,
		priority: priority,
	})
	// keep the registration sorted
	sys.registrations.Sort(sys.sortSystemByPriority)
}

// sortSystemByPriority sorts by systemRegistration priority, if equal by id
func (sys *Systems) sortSystemByPriority(a interface{}, b interface{}) bool {
	first := a.(systemRegistration)
	second := b.(systemRegistration)
	if first.priority == second.priority {
		return first.id < second.id
	}
	return first.priority > second.priority
}

// Update the systems
func (sys *Systems) Update(world *World, delta float32) error {
	var err error
	// go trough al registrations
	for it := sys.registrations.Iterator(); it != nil; it = it.Next() {
		// get the value
		sr := it.Value().(systemRegistration)
		//invoke the system, if error return
		if err = sr.system(world, delta); err != nil {
			return err
		}
	}
	return nil
}

// Clear the systems
func (sys *Systems) Clear() {
	sys.registrations.Clear()
}

// String returns the string representation of the systems
func (sys Systems) String() string {
	str := ""
	for it := sys.registrations.Iterator(); it != nil; it = it.Next() {
		l := it.Value().(systemRegistration)
		if str != "" {
			str += ","
		}
		name := runtime.FuncForPC(reflect.ValueOf(l.system).Pointer()).Name()
		str += name
	}
	return str
}

// NewSystems creates a new Systems
func NewSystems(systems int) *Systems {
	return &Systems{
		registrations: sparse.NewSlice(systems),
	}
}
