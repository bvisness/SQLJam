package app

import (
	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Aggregate struct {
	Alias      string
	Aggregates []AggregateColumn
	Groups     []string

	ColAggregateFields []*raygui.DropdownEx
}

type AggregateColumn struct {
	Type AggregateType
	Col  string
}

type AggregateType int

const (
	Min AggregateType = iota + 1
	Max
	Sum
)

func NewAggregate() *Node {
	return &Node{
		Title:   "Aggregate",
		CanSnap: true,
		Color:   rl.NewColor(78, 186, 170, 255),
		Inputs:  make([]*Node, 1),
		Data:    &Aggregate{},
	}
}

func (d *Aggregate) Update(n *Node) {
	n.UISize = rl.Vector2{X: 200, Y: float32(48)}
}

func (d *Aggregate) DoUI(n *Node) {

}
