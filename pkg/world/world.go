package world

import (
	"fmt"
	"github.com/juan-medina/goecs/pkg/system"
	"github.com/juan-medina/goecs/pkg/view"
	"reflect"
)

type World struct {
	*view.View
	systems map[string][]system.System
}

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

func (wld *World) AddSystemToGroup(sys system.System, group string) {
	if _, ok := wld.systems[group]; !ok {
		wld.systems[group] = make([]system.System, 0)
	}
	wld.systems[group] = append(wld.systems[group], sys)
}

func (wld *World) AddSystem(sys system.System) {
	wld.AddSystemToGroup(sys, system.DefaultGroup)
}

func (wld *World) UpdateGroup(group string) {
	for _, s := range wld.systems[group] {
		s.Update(wld.View)
	}
}

func (wld *World) Update() {
	wld.UpdateGroup(system.DefaultGroup)
}

func New() *World {
	return &World{
		View:    view.New(),
		systems: make(map[string][]system.System, 0),
	}
}
