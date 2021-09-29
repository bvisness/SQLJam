package node

import (
	"fmt"
	"strings"

	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Node struct {
	Data   NodeData
	Inputs []*Node

	// UI data
	Pos     rl.Vector2
	Size    rl.Vector2
	Title   string
	Color   rl.Color
	CanSnap bool // can snap to another node for its primary input
	Snapped bool
	Sort    int

	// Set in node update, affects layout
	UISize rl.Vector2

	// calculated fields
	InputPinPos    []rl.Vector2
	OutputPinPos   rl.Vector2
	HasChildren    bool
	SnapTargetRect rl.Rectangle
	UIRect         rl.Rectangle // roughly where the UI should fit
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

func NewPickColumns() *Node {
	return &Node{
		Title:   "Pick Columns",
		CanSnap: true,
		Color:   rl.NewColor(244, 143, 177, 255),
		Inputs:  make([]*Node, 1),
		Data: &PickColumns{
			Cols: make([]string, 1),
		},
	}
}

func NewFilter() *Node {
	return &Node{
		Title:   "Filter",
		CanSnap: true,
		Color:   rl.NewColor(111, 207, 151, 255),
		Inputs:  make([]*Node, 1),
		Data:    &Filter{},
	}
}

func NewCombineRows(combineType CombineType) *Node {
	return &Node{
		Title:   "Combine Rows",
		CanSnap: false,
		Color: rl.NewColor(178, 223, 219, 255),
		Inputs:  make([]*Node, 2),
		Data: &CombineRows{
			CombinationType: combineType,
		},
	}
}

func (n *Node) OldSqlGen(hasParent bool) string {
	if n == nil {
		return ""
	}

	// TODO: Optimizations :P

	switch d := n.Data.(type) {
	case *Table:
		ourQuery := ""
		if hasParent {
			ourQuery += d.Table
		} else {
			ourQuery += fmt.Sprintf("SELECT * FROM (%s)", d.Table)
		}
		if d.Alias != "" {
			ourQuery += fmt.Sprintf(" AS %s", d.Alias)
		}
		return ourQuery
	case *PickColumns:
		var cols = d.Cols
		colsJoined := strings.Join(cols, ", ")

		if len(n.Inputs) == 0 {
			// TODO: Return some kind of nice compile error
			return "ERROR"
		} else if len(n.Inputs) == 1 {

			return fmt.Sprintf("SELECT %s FROM (%s)", colsJoined, n.Inputs[0].OldSqlGen(true))
		} else {
			panic("Pick Columns node had more than one input")
		}
	case *Filter:
		wrappedConditions := fmt.Sprintf("(%s)", d.Conditions)

		if len(n.Inputs) == 0 {
			// TODO: Return some kind of nice compile error
			return "ERROR"
		} else if len(n.Inputs) == 1 {
			return fmt.Sprintf("SELECT * FROM (%s) WHERE %s", n.Inputs[0].OldSqlGen(true), wrappedConditions)
		} else {
			panic("Pick Columns node had more than one input")
		}
	case *CombineRows:
		if len(n.Inputs) == 2 {
			used := ""
			switch d.CombinationType {
			case Union:
				used = "UNION"
			case Intersect:
				used = "INTERSECT"
			case Except:
				used = "EXCEPT"
			case UnionAll:
				used = "UNION ALL"
			case IntersectAll:
				used = "INTERSECT ALL"
			case ExceptAll:
				used = "EXCEPT ALL"
			}
			return fmt.Sprintf("%s %s %s", n.Inputs[0], used, n.Inputs[1])
		} else {
			panic("Combine rows did not have two inputs")
		}
	default:
		return "SELECT NULL LIMIT 0" // empty result set
	}
}

func (n *Node) Rect() rl.Rectangle {
	return rl.Rectangle{n.Pos.X, n.Pos.Y, n.Size.X, n.Size.Y}
}
