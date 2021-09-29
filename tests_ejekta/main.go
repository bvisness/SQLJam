package main

import (
	"fmt"
	"github.com/bvisness/SQLJam/node"
)

var nodes []*node.Node

func main() {
		// init nodes
	filmTable := node.NewTable()
	nodes = append(nodes, filmTable)

	filter := node.NewFilter()
	filter.Inputs[0] = filmTable
	nodes = append(nodes, filter)

	pick := node.NewPickColumns()
	pick.Data.(*node.PickColumns).Cols = append(pick.Data.(*node.PickColumns).Cols, "title", "description", "release_year")
	pick.Inputs[0] = filter
	nodes = append(nodes, pick)

	pick2 := node.NewPickColumns()
	pick2.Data.(*node.PickColumns).Cols = append(pick2.Data.(*node.PickColumns).Cols, "title")
	pick2.Inputs[0] = pick
	nodes = append(nodes, pick2)

	// Recursive generate the context tree
	ctxTree := node.NewRecursiveContext(pick2)
	// Turn it into SQL
	//fmt.Println(fmt.Sprintf("Pick2 SRC: %s", ctxTree.Source))
	fmt.Println(ctxTree.SourceToSql())
}