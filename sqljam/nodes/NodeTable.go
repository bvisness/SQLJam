package nodes

import "github.com/bvisness/SQLJam/sqljam"

type NodeTable struct {
	sqljam.SqlNode
	sqljam.SqlNodeStacker
	TableSource string
	output *sqljam.SqlNodePin
}

