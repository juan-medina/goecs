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

// Listener that get notified that a new signal has been received by World.Signal
type Listener func(world *World, signal interface{}, delta float32) error

// subscription hold the information of listener subscribed to signals with a priority and id
type subscription struct {
	listener Listener       // listener for this subscription
	signals  []reflect.Type // signals that we are subscribed to
	priority int32          // priority of this subscription
	id       int64          // id of the subscription
}

// Subscriptions manage subscriptions of Listeners to signals
type Subscriptions struct {
	subscriptions      sparse.Slice // subscriptions is an sparse.Slice of subscriptions
	lastSubscriptionID int64        // lastSubscriptionID is the last subscription id
	signals            sparse.Slice // sparse.Slice of signals
	toSend             sparse.Slice // sparse.Slice of signals is a copy to signals to be send
}

// Subscribe adds a new subscription given a priority and set of signals types
func (subs *Subscriptions) Subscribe(listener Listener, priority int32, signals ...reflect.Type) {
	// increment the id
	subs.lastSubscriptionID++
	// add the subscription
	subs.subscriptions.Add(subscription{
		id:       subs.lastSubscriptionID,
		listener: listener,
		signals:  signals,
		priority: priority,
	})
	// keep the subscriptions sorted
	subs.subscriptions.Sort(subs.sortSubsByPriority)
}

// Signal adds a signal to to be sent
func (subs *Subscriptions) Signal(signal interface{}) {
	// add the signal
	subs.signals.Add(signal)
}

// sortSubsByPriority sorts by subscription priority, if equal by id
func (subs Subscriptions) sortSubsByPriority(a, b interface{}) bool {
	first := a.(subscription)
	second := b.(subscription)
	if first.priority == second.priority {
		return first.id < second.id
	}
	return first.priority > second.priority
}

// Update send the pending signals to the listeners on the world
func (subs *Subscriptions) Update(world *World, delta float32) error {
	// avoid to copy empty signals
	if subs.signals.Size() == 0 {
		return nil
	}
	// replace the signals to send, so we do not send the signals triggered by the current signals
	subs.signals.Replace(subs.toSend)

	// clear the hold so new signals will be here
	subs.signals.Clear()

	var err error
	// get thee signals to send
	for ite := subs.toSend.Iterator(); ite != nil; ite = ite.Next() {
		if err = subs.process(world, ite.Value(), delta); err != nil {
			return err
		}
	}

	// clear the signals to be send
	subs.toSend.Clear()

	return nil
}

// process the subscriptions for this signal
func (subs Subscriptions) process(world *World, signal interface{}, delta float32) error {
	var err error
	// get the signal type
	signalType := reflect.TypeOf(signal)
	// iterate trough the subscriptions
	for it := subs.subscriptions.Iterator(); it != nil; it = it.Next() {
		// get te subscription value
		sub := it.Value().(subscription)
		// go to the signal that this subscription is listen to
		for _, t := range sub.signals {
			// if we listen to this signal type
			if t == signalType {
				// notify the listener, return error if happen
				if err = sub.listener(world, signal, delta); err != nil {
					return err
				}
				// we do not need to iterate further for this subscription
				break
			}
		}
	}
	// no error happens
	return nil
}

// Clear the subscriptions & signals
func (subs *Subscriptions) Clear() {
	subs.subscriptions.Clear()
	subs.signals.Clear()
	subs.toSend.Clear()
}

// String returns the string representation of the subscriptions
func (subs Subscriptions) String() string {
	str := ""
	for it := subs.subscriptions.Iterator(); it != nil; it = it.Next() {
		l := it.Value().(subscription)
		if str != "" {
			str += ","
		}
		name := runtime.FuncForPC(reflect.ValueOf(l.listener).Pointer()).Name()
		signals := ""
		for _, v := range l.signals {
			if signals != "" {
				signals += ","
			}
			signals += v.Name()
		}
		str += fmt.Sprintf("{listener: %s, signals: {%s}}", name, signals)
	}
	return str
}

// NewSubscriptions creates a new Subscriptions
func NewSubscriptions(listeners, signals int) *Subscriptions {
	return &Subscriptions{
		subscriptions: sparse.NewSlice(listeners),
		signals:       sparse.NewSlice(signals),
		toSend:        sparse.NewSlice(signals),
	}
}
