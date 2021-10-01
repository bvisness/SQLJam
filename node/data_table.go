package node

import (
	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Table struct {
	NodeData
	SqlSource
	Alias string
	Table string

	// UI data
	TableDropdown raygui.DropdownEx
}


func (t *Table) SourceToSql() string {
	return t.Table
}

func (t *Table) SourceAlias() string {
	return t.Alias
}


func NewTable() *Node {
	return &Node{
		Title:   "Table",
		CanSnap: false,
		Color:   rl.NewColor(242, 201, 76, 255),
		Data: &Table{
			TableDropdown: raygui.NewDropdownEx(),
		},
	}
}
