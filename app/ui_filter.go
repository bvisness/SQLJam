package app

import (
	"github.com/bvisness/SQLJam/node"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func doFilterUpdate(n *node.Node, f *node.Filter) {
	n.UISize = rl.Vector2{360, 24}
}

func doFilterUI(n *node.Node, f *node.Filter) {
	rl.DrawRectangleRec(n.UIRect, rl.White)
	f.Conditions = f.TextBox.Do(n.UIRect, f.Conditions, 100)
}