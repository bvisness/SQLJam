package app

import (
	"fmt"

	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Join struct {
	FirstAlias        string
	FirstAliasTextbox raygui.TextBoxEx
	Conditions        []*JoinCondition
}

type JoinCondition struct {
	Alias            string
	Condition        string
	Left             bool
	Right            bool
	AliasTextBox     raygui.TextBoxEx
	ConditionTextBox raygui.TextBoxEx
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
		Color:   rl.NewColor(113, 170, 52, 255),
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

	n.InputPinHeights[0] = 0
	uiHeight := UIFieldHeight + UIFieldSpacing*2 // first input (no condition)

	if d.FirstAlias == "" && !d.FirstAliasTextbox.Active {
		d.FirstAlias = "a"
	}

	for i := range n.Inputs[1:] {
		n.InputPinHeights[i+1] = uiHeight
		uiHeight += UIFieldHeight + UIFieldSpacing   // alias field
		uiHeight += UIFieldHeight + 2*UIFieldSpacing // condition field

		cond := d.Conditions[i]
		if cond.Alias == "" && !cond.AliasTextBox.Active {
			cond.Alias = string(rune('b' + i))
		}
	}

	uiHeight += UIFieldHeight // +/- buttons

	n.UISize = rl.Vector2{480, float32(uiHeight)}
}

func (d *Join) DoUI(n *Node) {
	fieldY := n.UIRect.Y

	// first alias
	{
		aliasRect := rl.Rectangle{
			n.UIRect.X,
			float32(fieldY),
			n.UIRect.Width,
			UIFieldHeight,
		}
		d.FirstAlias, _ = d.FirstAliasTextbox.Do(aliasRect, d.FirstAlias, 100)
	}

	fieldY += UIFieldHeight + 2*UIFieldSpacing

	uiRight := n.UIRect.X + n.UIRect.Width
	boxWidth := n.UIRect.Width - (UIFieldSpacing+UIFieldHeight)*2

	for _, condition := range d.Conditions {
		// alias
		aliasRect := rl.Rectangle{
			n.UIRect.X,
			float32(fieldY),
			n.UIRect.Width,
			UIFieldHeight,
		}
		condition.Alias, _ = condition.AliasTextBox.Do(aliasRect, condition.Alias, 100)

		fieldY += UIFieldHeight + UIFieldSpacing

		// condition
		conditionRect := rl.Rectangle{
			n.UIRect.X,
			float32(fieldY),
			boxWidth,
			UIFieldHeight,
		}
		condition.Condition, _ = condition.ConditionTextBox.Do(conditionRect, condition.Condition, 100)
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

		fieldY += UIFieldHeight + 2*UIFieldSpacing
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

func (j *Join) Serialize() (res string, active bool) {
	res += j.FirstAlias
	for _, cond := range j.Conditions {
		res += cond.Alias
		res += cond.Condition
		res += fmt.Sprintf("%v", cond.Left)
		res += fmt.Sprintf("%v", cond.Right)
		if cond.AliasTextBox.Active || cond.ConditionTextBox.Active {
			active = true
		}
	}
	return
}
