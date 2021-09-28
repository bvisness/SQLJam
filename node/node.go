package node

import (
	"fmt"
	"sort"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Node struct {
	Data   NodeData
	Inputs []*Node

	// UI data
	Pos     rl.Vector2
	CanSnap bool // can snap to another node for its primary input
	Snapped bool

	// calculated fields
	InputPinPos  []rl.Vector2
	OutputPinPos []rl.Vector2
}

func NewTable(table string) *Node {
	return &Node{
		CanSnap: false,
		Data: &Table{
			Table: table,
		},
	}
}

func NewPickColumns() *Node {
	return &Node{
		CanSnap: true,
		Inputs:  make([]*Node, 1),
		Data: &PickColumns{
			Cols: make(map[string]bool),
		},
	}
}

func NewFilter(conditions []string) *Node {
	return &Node{
		CanSnap: true,
		Inputs:  make([]*Node, 1),
		Data: &Filter{
			Conditions: conditions,
		},
	}
}

func (n *Node) SQL() string {
	// TODO: Optimizations :P

	switch d := n.Data.(type) {
	case *Table:
		return fmt.Sprintf("SELECT * FROM %s", d.Table)
	case *PickColumns:
		// TODO: Someday allow custom order of columns
		var cols []string
		for col, yep := range d.Cols {
			if yep {
				cols = append(cols, col)
			}
		}
		sort.Strings(cols)
		colsJoined := strings.Join(cols, ", ")

		if len(n.Inputs) == 0 {
			// TODO: Return some kind of nice compile error
			return "ERROR"
		} else if len(n.Inputs) == 1 {
			return fmt.Sprintf("SELECT %s FROM (%s)", colsJoined, n.Inputs[0].SQL())
		} else {
			panic("Pick Columns node had more than one input")
		}
	case *Filter:
		wrappedConditions := make([]string, len(d.Conditions))
		for i, cond := range d.Conditions {
			wrappedConditions[i] = fmt.Sprintf("(%s)", cond)
		}
		joinedConditions := strings.Join(wrappedConditions, " AND ")

		if len(n.Inputs) == 0 {
			// TODO: Return some kind of nice compile error
			return "ERROR"
		} else if len(n.Inputs) == 1 {
			return fmt.Sprintf("SELECT * FROM (%s) WHERE %s", n.Inputs[0].SQL(), joinedConditions)
		} else {
			panic("Pick Columns node had more than one input")
		}
	default:
		return "SELECT NULL LIMIT 0" // empty result set
	}
}
