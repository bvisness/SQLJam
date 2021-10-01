package node

import (
	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Order struct {
	NodeData

	Alias string
	Cols  []OrderColumn

	ColDropdowns []*raygui.DropdownEx
}

type OrderColumn struct {
	Col        string
	Descending bool
}

func NewOrder() *Node {
	return &Node{
		Title:   "Order",
		CanSnap: true,
		Color:   rl.NewColor(255, 204, 128, 255),
		Inputs:  make([]*Node, 1),
		Data: &Order{
			Cols:         make([]OrderColumn, 1),
			ColDropdowns: raygui.MakeDropdownExList(1),
		},
	}
}
