package app

import (
	"log"

	"github.com/bvisness/SQLJam/raygui"
)

const UIFieldHeight = 24
const UIFieldSpacing = 4

func getSchema(n *Node) ([]string, error) {
	rows, err := db.Query(n.GenerateSql() + " LIMIT 0") // TODO: The limit should be part of SQL generation, yeah?
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rows.Columns()
}

var errorOpts = []raygui.DropdownExOption{{"ERROR", "ERROR"}}

// Gets dropdown options for the table produced by the given
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
