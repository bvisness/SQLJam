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

func NewTable(table string, alias string) *Node {
	return &Node{
		CanSnap: false,
		Data: &Table{
			Table: table,
			NodeAlias: NodeAlias{
				Alias: alias,
			},
		},
	}
}

func NewPickColumns(alias string) *Node {
	return &Node{
		CanSnap: true,
		Inputs:  make([]*Node, 1),
		Data: &PickColumns{
			Cols: make(map[string]bool),
			NodeAlias: NodeAlias{
				Alias: alias,
			},
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

func NewCombineRows(combineType CombineType) *Node {
	return &Node {
		CanSnap: false,
		Data: &CombineRows{
			CombinationType: combineType,
		},
	}
}

func (n *Node) SQL(hasParent bool) string {
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

			return fmt.Sprintf("SELECT %s FROM (%s)", colsJoined, n.Inputs[0].SQL(true))
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
			return fmt.Sprintf("SELECT * FROM (%s) WHERE %s", n.Inputs[0].SQL(true), joinedConditions)
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
