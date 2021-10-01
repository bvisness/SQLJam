package app

import (
	"fmt"
	"github.com/bvisness/SQLJam/node"
	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func doJoinUpdate(n *node.Node, j *node.Join) {
	n.UISize = rl.Vector2{X: 200, Y: float32(48)}

	// TODO probably put this somewhere else?
	for i := range n.Inputs[1:] {
		fmt.Println(i)
		if len(j.Conditions) < i+1 {
			j.Conditions = append(j.Conditions, &node.JoinCondition{
				Type:      node.InnerJoin,
				Condition: "???",
				TextBox:   &raygui.TextBoxEx{},
			})
		}
	}
}

func doJoinUI(n *node.Node, j *node.Join) {
	for _, condition := range j.Conditions {
		condition.Condition = condition.TextBox.Do(n.UIRect, condition.Condition, 100)
	}
}
