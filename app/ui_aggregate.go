package app

import (
	"github.com/bvisness/SQLJam/node"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func doAggregateUpdate(n *node.Node, d *node.Aggregate) {
	n.UISize = rl.Vector2{X: 200, Y: float32(48)}
}

func doAggregateUI(n *node.Node, d *node.Aggregate) {

}
