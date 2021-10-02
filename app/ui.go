package app

import (
	"fmt"
	"log"

	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const UIFieldHeight = 24
const UIFieldSpacing = 4

func getSchemaOfSqlSource(src SqlSource) ([]string, error) {
	srcToRun := src.SourceToSql(0)

	if src.IsPure() {
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

func getSchema(n *Node) ([]string, error) {
	ctx := NewQueryContextFromNode(n)
	var colsToShow = make([]string, 0)

	currentSourceCols, _ := getSchemaOfSqlSource(ctx.Source)

	for _, col := range currentSourceCols {
		colsToShow = append(colsToShow, ctx.SourceAlias() + "." +col)
	}

	for _, join := range ctx.Joins {
		joinCols, _ := getSchemaOfSqlSource(join.Source)
		for _, col := range joinCols {
			colsToShow = append(colsToShow, join.Source.SourceAlias() + "." + col)
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
