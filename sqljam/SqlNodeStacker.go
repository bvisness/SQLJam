package sqljam

import "github.com/bvisness/SQLJam/sqljam/node_types"

type SqlNodeStacker struct {
	node_types.NodeTypeModifier
	StackChild *SqlNode
}

func (stacker *SqlNodeStacker) SetChild(child *SqlNode) {
	// if child does not have NodeTypeModifier then panic?
}
