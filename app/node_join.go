package app

import (
	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Join struct {
	NodeData
	Conditions []*JoinCondition
}

type JoinCondition struct {
	Condition string
	Left      bool
	Right     bool
	TextBox   raygui.TextBoxEx
}

type JoinType int

const (
	LeftJoin JoinType = iota + 1
	RightJoin
	InnerJoin
	OuterJoin
)

func NewJoin() *Node {
	return &Node{
		Title:   "Join",
		CanSnap: false,
		Color:   rl.NewColor(102, 187, 106, 255),
		Inputs:  make([]*Node, 2),
		Data: &Join{
			Conditions: []*JoinCondition{{}},
		},
	}
}

func (jc *JoinCondition) Type() JoinType {
	if jc.Left && jc.Right {
		return OuterJoin
	} else if jc.Left {
		return LeftJoin
	} else if jc.Right {
		return RightJoin
	} else {
		return InnerJoin
	}
}

func (jt JoinType) String() string {
	switch jt {
	case LeftJoin:
		return "LEFT JOIN"
	case RightJoin:
		return "RIGHT JOIN"
	case InnerJoin:
		return "JOIN"
	case OuterJoin:
		return "OUTER JOIN"
	default:
		return "BAD JOIN"
	}
}

func (d *Join) Update(n *Node) {
	n.InputPinHeights = make([]int, len(n.Inputs))

	uiHeight := UIFieldHeight // blank space for first table input
	n.InputPinHeights[0] = 0

	for i := range n.Inputs[1:] {
		uiHeight += UIFieldSpacing
		n.InputPinHeights[i+1] = uiHeight
		uiHeight += UIFieldHeight
	}

	uiHeight += UIFieldSpacing + UIFieldHeight // +/- buttons

	n.UISize = rl.Vector2{300, float32(uiHeight)}
}

func (d *Join) DoUI(n *Node) {
	fieldY := n.UIRect.Y + UIFieldHeight + UIFieldSpacing

	uiRight := n.UIRect.X + n.UIRect.Width
	boxWidth := n.UIRect.Width - (UIFieldSpacing+UIFieldHeight)*2

	for _, condition := range d.Conditions {
		boxRect := rl.Rectangle{
			n.UIRect.X,
			float32(fieldY),
			boxWidth,
			UIFieldHeight,
		}
		rl.DrawRectangleRec(boxRect, rl.White)
		condition.Condition = condition.TextBox.Do(boxRect, condition.Condition, 100)
		condition.Left = raygui.Toggle(rl.Rectangle{
			uiRight - (UIFieldHeight + UIFieldSpacing + UIFieldHeight),
			float32(fieldY),
			UIFieldHeight,
			UIFieldHeight,
		}, "L", condition.Left)
		raygui.Disable() // lol thank you sqlite
		{
			condition.Right = raygui.Toggle(rl.Rectangle{
				uiRight - UIFieldHeight,
				float32(fieldY),
				UIFieldHeight,
				UIFieldHeight,
			}, "R", condition.Right)
		}
		raygui.Enable()

		fieldY += UIFieldHeight + UIFieldSpacing
	}

	if raygui.Button(rl.Rectangle{
		n.UIRect.X,
		fieldY,
		n.UIRect.Width/2 - UIFieldSpacing/2,
		UIFieldHeight,
	}, "+") {
		n.Inputs = append(n.Inputs, nil)
		d.Conditions = append(d.Conditions, &JoinCondition{})
	}
	if raygui.Button(rl.Rectangle{
		n.UIRect.X + n.UIRect.Width/2 + UIFieldSpacing/2,
		fieldY,
		n.UIRect.Width/2 - UIFieldSpacing/2,
		UIFieldHeight,
	}, "-") {
		if len(d.Conditions) > 1 {
			n.Inputs = n.Inputs[:len(n.Inputs)-1]
			d.Conditions = d.Conditions[:len(d.Conditions)-1]
		}
	}
}
