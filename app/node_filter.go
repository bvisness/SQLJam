package app

import (
	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Filter struct {
	NodeData
	Conditions string // TODO: whatever data we actually need for our filter UI

	// UI data
	TextBox raygui.TextBoxEx
}

func NewFilter() *Node {
	return &Node{
		Title:   "Filter",
		CanSnap: true,
		Color:   rl.NewColor(111, 207, 151, 255),
		Inputs:  make([]*Node, 1),
		Data:    &Filter{},
	}
}

func (d *Filter) Update(n *Node) {
	n.UISize = rl.Vector2{360, 24}
}

func (d *Filter) DoUI(n *Node) {
	rl.DrawRectangleRec(n.UIRect, rl.White)
	d.Conditions = d.TextBox.Do(n.UIRect, d.Conditions, 100)
}
