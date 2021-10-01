package node

import (
	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type CombineRows struct {
	NodeData
	CombinationType CombineType
	Dropdown        raygui.DropdownEx
}

type CombineType int

const (
	Union CombineType = iota + 1
	UnionAll
	Intersect
	Except
)

func NewCombineRows(combineType CombineType) *Node {
	return &Node{
		Title:   "Combine Rows",
		CanSnap: false,
		Color:   rl.NewColor(178, 223, 219, 255),
		Inputs:  make([]*Node, 2),
		Data: &CombineRows{
			CombinationType: combineType,
		},
	}
}
