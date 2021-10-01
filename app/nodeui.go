package app

import (
	"github.com/bvisness/SQLJam/node"
	"github.com/bvisness/SQLJam/raygui"
	"log"
)

// Before drawing. Sort out your data, set your layout info, etc.
func doNodeUpdate(n *node.Node) {
	switch d := n.Data.(type) {
	case *node.Table:
		doTableUpdate(n, d)
	case *node.Filter:
		doFilterUpdate(n, d)
	case *node.PickColumns:
		doPickColumnsUpdate(n, d)
	case *node.CombineRows:
		doCombineRowsUpdate(n, d)
	case *node.Order:
		doOrderUpdate(n, d)
	case *node.Join:
		doJoinUpdate(n, d)
	case *node.Aggregate:
		doAggregateUpdate(n, d)
	}
}



// Drawing and user input.
func doNodeUI(n *node.Node) {
	switch d := n.Data.(type) {
	case *node.Table:
		doTableUI(n, d)
	case *node.Filter:
		doFilterUI(n, d)
	case *node.PickColumns:
		doPickColumnsUI(n, d)
	case *node.CombineRows:
		doCombineRowsUI(n, d)
	case *node.Order:
		doOrderUI(n, d)
	case *node.Join:
		doJoinUI(n, d)
	case *node.Aggregate:
		doAggregateUI(n, d)
	}
}

const UIFieldHeight = 24
const UIFieldSpacing = 4
const orderDirectionWidth = 70


func getSchema(n *node.Node) ([]string, error) {
	rows, err := db.Query(n.GenerateSql() + " LIMIT 0") // TODO: The limit should be part of SQL generation, yeah?
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rows.Columns()
}

var errorOpts = []raygui.DropdownExOption{{"ERROR", "ERROR"}}

// Gets dropdown options for the table produced by the given node.
// Returns default options if no schema can be found.
func columnNameDropdownOpts(inputNode *node.Node) []raygui.DropdownExOption {
	if inputNode == nil {
		return errorOpts
	}

	var opts []raygui.DropdownExOption
	schemaCols, err := getSchema(inputNode)
	if err == nil {
		for _, col := range schemaCols {
			opts = append(opts, raygui.DropdownExOption{
				Name:  col,
				Value: col,
			})
		}
	} else {
		log.Print(err)
		return errorOpts
	}

	return opts
}
