package app

import (
	"fmt"
	"log"

	"github.com/bvisness/SQLJam/node"
	"github.com/bvisness/SQLJam/raygui"
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

func getSchemaOfQueryContext(ctx *node.QueryContext) ([]string, error) {
	srcToRun := ctx.SourceToSql(0)

	rows, err := db.Query(srcToRun + " LIMIT 0")
	if err != nil {
		//fmt.Println(src.SourceToSql(0))
		fmt.Println(err.Error())
		return nil, err
	}
	return rows.Columns()
}

func getSchema(n *node.Node) ([]string, error) {

	ctx := node.NewQueryContextFromNode(n)

	//ctx.Joins[0].Source.SourceAlias()

	var colsToShow = make([]string, 0)

	// ### Get CURRENT source rows ###

	var currentSourceRows []string

	if len(ctx.Cols) == 0 {
		currentSourceRows, _ = getSchemaOfQueryContext(ctx)
		fmt.Println("Current schema alias is: " + ctx.SourceAlias())
	} else {
		currentSourceRows = ctx.Cols
	}

	for _, row := range currentSourceRows {
		colsToShow = append(colsToShow, ctx.SourceAlias() + "." + row)
	}

	// ### Get JOIN source rows ###

	for i, join := range ctx.Joins {
		input := n.Inputs[i+1]
		joinCols, _ := getSchemaOfQueryContext(node.NewQueryContextFromNode(input))
		for _, col := range joinCols {
			colsToShow = append(colsToShow, join.Alias + "." + col)
		}
	}

	return colsToShow, nil
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
