package main

import (
	"fmt"
	"github.com/bvisness/SQLJam/node"
)

var nodes []*node.Node

func main() {
		// init nodes
	filmTable := node.NewTable("film", "cool_films")
	nodes = append(nodes, filmTable)

	filter := node.NewFilter([]string{"rating = 'PG'", "rental_rate < 3"})
	filter.Inputs[0] = filmTable
	nodes = append(nodes, filter)

	pick := node.NewPickColumns("test_alias")
	pick.Data.(*node.PickColumns).Cols = append(pick.Data.(*node.PickColumns).Cols, "title", "description", "release_year")
	pick.Inputs[0] = filter
	nodes = append(nodes, pick)

	pick2 := node.NewPickColumns("other_alias")
	pick2.Data.(*node.PickColumns).Cols = append(pick2.Data.(*node.PickColumns).Cols, "title")
	pick2.Inputs[0] = pick
	nodes = append(nodes, pick2)

	// Recursive generate the context tree
	ctxTree := node.NewRecursiveGenerated(pick2)
	// Turn it into SQL
	//fmt.Println(fmt.Sprintf("Pick2 SRC: %s", ctxTree.Source))
	fmt.Println(ctxTree.SourceToSql())
}