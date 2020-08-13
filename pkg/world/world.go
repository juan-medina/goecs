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

func (w World) String() string {
	result := ""

	result += fmt.Sprintf("World[view: %v, systems: [", w.View)

	for g,_ := range w.systems {
		result += fmt.Sprintf("%s:[", g)
		for _, s := range w.systems[g] {
			result += fmt.Sprintf("%s,", reflect.TypeOf(s).String())
		}
		result += "],"
	}
	result += "]"

	return result
}

func (w *World) AddSystemToGroup(s system.System, g string) {
	if _, ok := w.systems[g]; !ok {
		w.systems[g] = make([]system.System, 0)
	}
	w.systems[g] = append(w.systems[g], s)
}

func (w *World) AddSystem(s system.System) {
	w.AddSystemToGroup(s, system.DefaultGroup)
}

func (w *World) UpdateGroup(g string) {
	for _, s := range w.systems[g] {
		s.Update(w.View)
	}
}

func (w *World) Update() {
	w.UpdateGroup(system.DefaultGroup)
}

func New() *World {
	return &World{
		View:    view.New(),
		systems: make(map[string][]system.System, 0),
	}
}
