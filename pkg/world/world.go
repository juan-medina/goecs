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

package world

import (
	"fmt"
	"github.com/juan-medina/goecs/pkg/view"
	"reflect"
)

type World struct {
	*view.View
	systems map[string][]System
}

const (
	defaultSystemGroup = "DEFAULT_GROUP"
)

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

func (wld *World) AddSystemToGroup(sys System, group string) {
	if _, ok := wld.systems[group]; !ok {
		wld.systems[group] = make([]System, 0)
	}
	wld.systems[group] = append(wld.systems[group], sys)
}

func (wld *World) AddSystem(sys System) {
	wld.AddSystemToGroup(sys, defaultSystemGroup)
}

func (wld *World) UpdateGroup(group string, delta float64) error {
	for _, s := range wld.systems[group] {
		if err := s.Update(wld, delta); err != nil {
			return err
		}
	}
	return nil
}

func (wld *World) Update(delta float64) error {
	return wld.UpdateGroup(defaultSystemGroup, delta)
}

func (wld *World) NotifyGroup(group string, event interface{}, delta float64) error {
	for _, s := range wld.systems[group] {
		if err := s.Notify(wld, event, delta); err != nil {
			return err
		}
	}
	return nil
}

func (wld *World) Notify(event interface{}, delta float64) error {
	return wld.NotifyGroup(defaultSystemGroup, event, delta)
}

func New() *World {
	return &World{
		View:    view.New(),
		systems: make(map[string][]System, 0),
	}
}
