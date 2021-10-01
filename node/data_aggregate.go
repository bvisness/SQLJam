package node

import (
	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type AggregateType int

const (
	Min AggregateType = iota + 1
	Max
	Sum
)

type AggregateColumn struct {
	Type AggregateType
	Col string
}

type Aggregate struct {
	NodeData

	Alias string
	Aggregates []AggregateColumn
	Groups []string

	ColAggregateFields []*raygui.DropdownEx
}

func NewAggregate() *Node {
	return &Node{
		Title: "Aggregate",
		CanSnap: true,
		Color: rl.NewColor(78, 186, 170, 255),
		Inputs: make([]*Node, 1),
		Data: &Aggregate{},
	}
}
