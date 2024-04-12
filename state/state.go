package state

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Code int

const (
	Nil = iota
	Change
	Add
	Delete
)

type Result struct {
	Code      Code
	NextState State
}

type State interface {
	Init(int)
	Update(int) Result
	Draw(*ebiten.Image, int)
}

type Machine struct {
	StateStack []State

	LayoutWidth, LayoutHeight int
}

func (g *Machine) Update() error {

	if g.IsEmpty() {
		return nil
	}

	res := g.StateStack[len(g.StateStack)-1].Update(len(g.StateStack) - 1)

	switch res.Code {
	case Change:
		g.StateChange(res.NextState)
	case Add:
		g.StateAdd(res.NextState)
	case Delete:
		g.StateStack = g.StateStack[0 : len(g.StateStack)-1]
	}

	return nil
}

func (g *Machine) Draw(screen *ebiten.Image) {
	for d, e := range g.StateStack {
		e.Draw(screen, d)
	}
}

func (g *Machine) IsEmpty() bool {
	return len(g.StateStack) == 0
}

func (g *Machine) StateAdd(v State) {
	g.StateStack = append(g.StateStack, v)
	v.Init(len(g.StateStack) - 1)
}

func (g *Machine) StateChange(v State) {
	g.StateStack[len(g.StateStack)-1] = v
	v.Init(len(g.StateStack) - 1)
}

func (g *Machine) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.LayoutWidth, g.LayoutHeight
}

func (g *Machine) Init() {
}
