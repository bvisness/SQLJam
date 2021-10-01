package node

import (
	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type PickColumns struct {
	NodeData
	SqlSource
	Alias string

	Entries []*PickColumnsEntry
}

type PickColumnsEntry struct {
	Col          string
	ColDropdown  raygui.DropdownEx
	Alias        string
	AliasTextbox raygui.TextBoxEx
}

func (pc *PickColumns) SourceAlias() string {
	return pc.Alias
}

func (pc *PickColumns) Cols() []string {
	res := make([]string, len(pc.Entries))
	for i := range res {
		res[i] = pc.Entries[i].Col
	}
	return res
}

func (pc *PickColumns) Aliases() []string {
	res := make([]string, len(pc.Entries))
	for i := range res {
		res[i] = pc.Entries[i].Alias
	}
	return res
}

func (pc *PickColumns) ColDropdowns() []*raygui.DropdownEx {
	res := make([]*raygui.DropdownEx, len(pc.Entries))
	for i := range res {
		res[i] = &pc.Entries[i].ColDropdown
	}
	return res
}

func NewPickColumns() *Node {
	return &Node{
		Title:   "Pick Columns",
		CanSnap: true,
		Color:   rl.NewColor(244, 143, 177, 255),
		Inputs:  make([]*Node, 1),
		Data: &PickColumns{
			Entries: []*PickColumnsEntry{{}},
		},
	}
}
