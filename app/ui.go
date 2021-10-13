package app

import (
	"fmt"

	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const UIFieldHeight = 36
const UIFieldSpacing = 4

func getSchemaOfSqlSource(src SqlSource) ([]string, error) {
	srcToRun := src.SourceToSql(0)

	if src.IsTable() {
		srcToRun = fmt.Sprintf("SELECT * FROM %s", srcToRun)
	}

	rows, err := db.Query(srcToRun + " LIMIT 0")
	if err != nil {
		//fmt.Println(src.SourceToSql(0))
		fmt.Println(err.Error())
		return nil, err
	}
	return rows.Columns()
}

func getSchema(n *Node) []string {
	if n.Schema != nil {
		return n.Schema
	}

	ctx := NewQueryContextFromNode(n)
	var colsToShow []string

	currentSourceCols, _ := getSchemaOfSqlSource(ctx)
	for _, col := range currentSourceCols {
		colsToShow = append(colsToShow, col)
	}

	n.Schema = colsToShow

	return colsToShow
}

var errorOpts = []raygui.DropdownExOption{{"ERROR", "ERROR"}}

// Gets dropdown options for the table produced by the given node.
// Returns default options if no schema can be found.
func columnNameDropdownOpts(inputNode *Node) []raygui.DropdownExOption {
	if inputNode == nil {
		return errorOpts
	}

	var opts []raygui.DropdownExOption
	schemaCols := getSchema(inputNode)
	for _, col := range schemaCols {
		opts = append(opts, raygui.DropdownExOption{
			Name:  col,
			Value: col,
		})
	}

	return opts
}

const basicTextSpacingRatio = 2 / 24

func drawBasicText(text string, x float32, y float32, size float32, color rl.Color) {
	rl.DrawTextEx(font, text, rl.Vector2{X: x, Y: y}, size, basicTextSpacingRatio*size, color)
}

func measureBasicText(text string, size float32) rl.Vector2 {
	return rl.MeasureTextEx(font, text, size, basicTextSpacingRatio*size)
}

func drawResizeHandle(bottomRight rl.Vector2, nodeColor rl.Color) {
	rl.DrawLineV(
		rl.Vector2{bottomRight.X - 10, bottomRight.Y - 2},
		rl.Vector2{bottomRight.X - 2, bottomRight.Y - 10},
		Brightness(nodeColor, 0.5),
	)
	rl.DrawLineV(
		rl.Vector2{bottomRight.X - 6, bottomRight.Y - 2},
		rl.Vector2{bottomRight.X - 2, bottomRight.Y - 6},
		Brightness(nodeColor, 0.5),
	)
}
