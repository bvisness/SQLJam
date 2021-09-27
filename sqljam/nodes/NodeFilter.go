package nodes

import "github.com/bvisness/SQLJam/sqljam"

type NodeFilter struct {
	sqljam.SqlNode
	sqljam.SqlNodeStacker
	TableSource string
	output *sqljam.SqlNodePin
}

