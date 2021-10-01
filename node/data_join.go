package node

import (
	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type JoinType int

const (
	LeftJoin JoinType = iota + 1
	RightJoin
	InnerJoin
	OuterJoin
)

type JoinCondition struct {
	Type      JoinType
	Condition string
	TextBox   *raygui.TextBoxEx
}

type Join struct {
	NodeData
	Conditions []*JoinCondition
}

func NewJoin() *Node {
	return &Node{
		Title:   "Join",
		CanSnap: false,
		Color:   rl.NewColor(102, 187, 106, 255),
		Inputs:  make([]*Node, 2),
		Data:    &Join{},
	}
}
