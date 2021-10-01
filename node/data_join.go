package node

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
