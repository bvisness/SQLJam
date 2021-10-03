package app

import (
	"fmt"

	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type CombineRows struct {
	NodeData
	CombinationType CombineType
	Dropdown        raygui.DropdownEx
}

type CombineType int

const (
	Union CombineType = iota + 1
	UnionAll
	Intersect
	Except
)

func NewCombineRows(combineType CombineType) *Node {
	return &Node{
		Title:   "Combine Rows",
		CanSnap: false,
		Color:   rl.NewColor(77,155,230, 255),
		Inputs:  make([]*Node, 2),
		Data: &CombineRows{
			CombinationType: combineType,
		},
	}
}

var combineRowsOpts = []raygui.DropdownExOption{
	{"UNION", Union},
	{"UNION ALL", UnionAll},
	{"INTERSECT", Intersect},
	{"EXCEPT", Except},
}

func (d *CombineRows) Update(n *Node) {
	n.UISize = rl.Vector2{X: 240, Y: UIFieldHeight * 1.25}
	d.Dropdown.SetOptions(combineRowsOpts...)
}

func (d *CombineRows) DoUI(n *Node) {
	chosen := d.Dropdown.Do(n.UIRect)
	d.CombinationType = chosen.(CombineType)
}

func (d *CombineRows) Serialize() string {
	return fmt.Sprintf("%v", d.CombinationType)
}
