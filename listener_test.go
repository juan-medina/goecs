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
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type listenerCall struct {
	listener string
	signal   string
}

var listenersCalls = make([]listenerCall, 0)

func addCall(signal Component) {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc).Name()
	fn = strings.Replace(fn, "github.com/juan-medina/goecs.", "", -1)
	tn := reflect.TypeOf(signal).String()
	tn = strings.Replace(tn, "goecs.", "", -1)
	listenersCalls = append(listenersCalls, listenerCall{
		listener: fn,
		signal:   tn,
	})
}

func listenerA(_ *World, signal Component, _ float32) error {
	addCall(signal)
	return nil
}

func listenerB(_ *World, signal Component, _ float32) error {
	addCall(signal)
	return nil
}

func listenerC(_ *World, signal Component, _ float32) error {
	addCall(signal)
	return nil
}

func TestSubscriptions(t *testing.T) {
	subs := NewSubscriptions(10, 10)

	subs.Subscribe(listenerA, 0, signal1Type, signal2Type)
	subs.Subscribe(listenerB, 0, signal2Type)
	subs.Subscribe(listenerC, 0, signal1Type, signal3Type)

	type testCase struct {
		name    string
		signals []interface{}
		expect  []listenerCall
	}

	cases := []testCase{
		{
			name:    "signal1",
			signals: []interface{}{signal1{}},
			expect: []listenerCall{
				{
					listener: "listenerA",
					signal:   "signal1",
				},
				{
					listener: "listenerC",
					signal:   "signal1",
				},
			},
		},
		{
			name:    "signal2",
			signals: []interface{}{signal2{}},
			expect: []listenerCall{
				{
					listener: "listenerA",
					signal:   "signal2",
				},
				{
					listener: "listenerB",
					signal:   "signal2",
				},
			},
		},
		{
			name:    "signal3",
			signals: []interface{}{signal3{}},
			expect: []listenerCall{
				{
					listener: "listenerC",
					signal:   "signal3",
				},
			},
		},
		{
			name:    "signal1, signal2",
			signals: []interface{}{signal1{}, signal2{}},
			expect: []listenerCall{
				{
					listener: "listenerA",
					signal:   "signal1",
				},
				{
					listener: "listenerC",
					signal:   "signal1",
				},
				{
					listener: "listenerA",
					signal:   "signal2",
				},
				{
					listener: "listenerB",
					signal:   "signal2",
				},
			},
		},
		{
			name:    "signal1, signal3",
			signals: []interface{}{signal1{}, signal3{}},
			expect: []listenerCall{
				{
					listener: "listenerA",
					signal:   "signal1",
				},
				{
					listener: "listenerC",
					signal:   "signal1",
				},
				{
					listener: "listenerC",
					signal:   "signal3",
				},
			},
		},
		{
			name:    "signal2, signal3",
			signals: []interface{}{signal2{}, signal3{}},
			expect: []listenerCall{
				{
					listener: "listenerA",
					signal:   "signal2",
				},
				{
					listener: "listenerB",
					signal:   "signal2",
				},
				{
					listener: "listenerC",
					signal:   "signal3",
				},
			},
		},
		{
			name:    "signal2, signal1",
			signals: []interface{}{signal2{}, signal1{}},
			expect: []listenerCall{
				{
					listener: "listenerA",
					signal:   "signal2",
				},
				{
					listener: "listenerB",
					signal:   "signal2",
				},
				{
					listener: "listenerA",
					signal:   "signal1",
				},
				{
					listener: "listenerC",
					signal:   "signal1",
				},
			},
		},
		{
			name:    "signal1, signal2, signal3",
			signals: []interface{}{signal1{}, signal2{}, signal3{}},
			expect: []listenerCall{
				{
					listener: "listenerA",
					signal:   "signal1",
				},
				{
					listener: "listenerC",
					signal:   "signal1",
				},
				{
					listener: "listenerA",
					signal:   "signal2",
				},
				{
					listener: "listenerB",
					signal:   "signal2",
				},
				{
					listener: "listenerC",
					signal:   "signal3",
				},
			},
		},
		{
			name:    "signal3, signal2, signal1",
			signals: []interface{}{signal3{}, signal2{}, signal1{}},
			expect: []listenerCall{
				{
					listener: "listenerC",
					signal:   "signal3",
				},
				{
					listener: "listenerA",
					signal:   "signal2",
				},
				{
					listener: "listenerB",
					signal:   "signal2",
				},
				{
					listener: "listenerA",
					signal:   "signal1",
				},
				{
					listener: "listenerC",
					signal:   "signal1",
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			listenersCalls = make([]listenerCall, 0)
			for _, s := range c.signals {
				subs.Signal(s)
			}
			_ = subs.Update(nil, 0)
			if !reflect.DeepEqual(listenersCalls, c.expect) {
				t.Fatalf("got %v, want %v", listenersCalls, c.expect)
			}
		})
	}

}

var signal1Type = NewComponentType()

type signal1 struct{}

func (s signal1) Type() ComponentType {
	return signal1Type
}

var signal2Type = NewComponentType()

type signal2 struct{}

func (s signal2) Type() ComponentType {
	return signal2Type
}

var signal3Type = NewComponentType()

type signal3 struct{}

func (s signal3) Type() ComponentType {
	return signal3Type
}

func TestSubscriptions_String(t *testing.T) {
	subs := NewSubscriptions(10, 10)

	subs.Subscribe(listenerA, 0, signal1Type, signal2Type)
	subs.Subscribe(listenerB, 0, signal2Type)
	subs.Subscribe(listenerC, 0, signal1Type, signal3Type)

	str := subs.String()

	if str == "" {
		t.Fatal("got empty, expect an string")
	}
}
