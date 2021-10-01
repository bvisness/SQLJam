package node

import (
	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Order struct {
	NodeData

	Alias string
	Cols  []*OrderColumn
}

type OrderColumn struct {
	Col         string
	Descending  bool
	ColDropdown raygui.DropdownEx
}

func NewOrder() *Node {
	return &Node{
		Title:   "Order",
		CanSnap: true,
		Color:   rl.NewColor(255, 204, 128, 255),
		Inputs:  make([]*Node, 1),
		Data: &Order{
			Cols: []*OrderColumn{{}},
		},
	}
}

func (oc *Order) ColDropdowns() []*raygui.DropdownEx {
	res := make([]*raygui.DropdownEx, len(oc.Cols))
	for i := range res {
		res[i] = &oc.Cols[i].ColDropdown
	}
	return res
}
