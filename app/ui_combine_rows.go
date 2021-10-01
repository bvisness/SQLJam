package app

import (
	"github.com/bvisness/SQLJam/node"
	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func doCombineRowsUpdate(n *node.Node, c *node.CombineRows) {
	n.UISize = rl.Vector2{X: 200, Y: float32(48)}
	c.Dropdown.SetOptions(combineRowsOpts...)
}

var combineRowsOpts = []raygui.DropdownExOption{
	{"UNION", node.Union},
	{"UNION ALL", node.UnionAll},
	{"INTERSECT", node.Intersect},
	{"EXCEPT", node.Except},
}

func doCombineRowsUI(n *node.Node, d *node.CombineRows) {
	chosen := d.Dropdown.Do(n.UIRect)
	d.CombinationType = chosen.(node.CombineType)
}
