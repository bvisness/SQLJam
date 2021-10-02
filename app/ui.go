package app

import (
	"fmt"
	"log"

	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const UIFieldHeight = 24
const UIFieldSpacing = 4

func getSchemaOfQueryContext(ctx *QueryContext) ([]string, error) {
	srcToRun := ctx.SourceToSql(0)

	rows, err := db.Query(srcToRun + " LIMIT 0")
	if err != nil {
		//fmt.Println(src.SourceToSql(0))
		fmt.Println(err.Error())
		return nil, err
	}
	return rows.Columns()
}

func getSchema(n *Node) ([]string, error) {
	ctx := NewQueryContextFromNode(n)
	var colsToShow = make([]string, 0)

	// ### Get CURRENT source rows ###

	var currentSourceRows []string

	// TODO figure out why we're grabbing 2x as many current ctx columns as needed

	if len(ctx.Cols) == 0 {
		currentSourceRows, _ = getSchemaOfQueryContext(ctx)
	} else {
		currentSourceRows = ctx.Cols
	}

	for _, row := range currentSourceRows {
		colsToShow = append(colsToShow, ctx.SourceAlias()+"."+row)
	}

	// ### Get JOIN source rows ###

	for i, join := range ctx.Joins {
		input := n.Inputs[i+1]
		joinCols, _ := getSchemaOfQueryContext(NewQueryContextFromNode(input))
		for _, col := range joinCols {
			colsToShow = append(colsToShow, join.Alias+"."+col)
		}
	}

	return colsToShow, nil
}

var errorOpts = []raygui.DropdownExOption{{"ERROR", "ERROR"}}

// Gets dropdown options for the table produced by the given node.
// Returns default options if no schema can be found.
func columnNameDropdownOpts(inputNode *Node) []raygui.DropdownExOption {
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

func drawBasicText(text string, x float32, y float32, size float32, color rl.Color) {
	rl.DrawTextEx(font, text, rl.Vector2{X: x, Y: y}, size, 2, color)
}

func measureBasicText(text string, size float32) rl.Vector2 {
	return rl.MeasureTextEx(font, text, size, 2)
}
