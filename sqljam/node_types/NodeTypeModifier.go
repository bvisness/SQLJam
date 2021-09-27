package node_types

import "github.com/bvisness/SQLJam/sqljam"

type NodeTypeModifier struct {
	input *sqljam.SqlNodePin
	output *sqljam.SqlNodePin
}
