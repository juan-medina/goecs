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

// Listener the get notified that a new signal has been received by World.Signal
type Listener func(world *World, signal interface{}, delta float32) error

// lastSubscriptionID to be used when sorted subscriptions
var lastSubscriptionID = int64(0)

// subscription hold the information of listener subscription
type subscription struct {
	listener Listener       // listener for this subscription
	signals  []reflect.Type // signals that we are subscribed to
	priority int32          // priority of this subscription
	id       int64          // id of the subscription
}

// newSubscription create a new subscription
func newSubscription(listener Listener, priority int32, signals ...reflect.Type) *subscription {
	lastSubscriptionID++
	sub := &subscription{
		listener: listener,
		signals:  signals,
		priority: priority,
		id:       lastSubscriptionID,
	}
	return sub
}

// Subscriptions manage subscriptions of Listeners to signals
type Subscriptions struct {
	subscriptions sparse.Slice
}

// Subscribe adds a new subscription
func (subs *Subscriptions) Subscribe(listener Listener, priority int32, signals ...reflect.Type) {
	sub := newSubscription(listener, priority, signals...)
	subs.subscriptions.Add(sub)
	subs.subscriptions.Sort(subs.sortSubsByPriority)
}

// sortSubsByPriority sorts by subscription priority, if equal by id
func (subs Subscriptions) sortSubsByPriority(a, b interface{}) bool {
	first := a.(*subscription)
	second := b.(*subscription)
	if first.priority == second.priority {
		return first.id < second.id
	}
	return first.priority > second.priority
}

// Process the subscriptions
func (subs Subscriptions) Process(world *World, signal interface{}, delta float32) error {
	var err error
	signalType := reflect.TypeOf(signal)
	for it := subs.subscriptions.Iterator(); it != nil; it = it.Next() {
		sub := it.Value().(*subscription)
		for _, t := range sub.signals {
			if t == signalType {
				if err = sub.listener(world, signal, delta); err != nil {
					return err
				}
				break
			}
		}
	}
	return nil
}

// Clear the subscriptions
func (subs *Subscriptions) Clear() {
	subs.subscriptions.Clear()
}

// String returns the string representation of the subscriptions
func (subs Subscriptions) String() string {
	str := ""
	for it := subs.subscriptions.Iterator(); it != nil; it = it.Next() {
		l := it.Value().(*subscription)
		if str != "" {
			str += ","
		}
		name := runtime.FuncForPC(reflect.ValueOf(l.listener).Pointer()).Name()
		str += fmt.Sprintf("{%s}", name)
	}
	return str
}

// NewSubscriptions creates a new Subscriptions
func NewSubscriptions(listeners int) *Subscriptions {
	return &Subscriptions{
		subscriptions: sparse.NewSlice(listeners),
	}
}
