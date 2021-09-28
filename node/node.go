package node

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Node struct {
	Data   NodeData
	Inputs []*Node

	// UI data
	Pos     rl.Vector2
	Size    rl.Vector2
	Title   string
	CanSnap bool // can snap to another node for its primary input
	Snapped bool

	// calculated fields
	InputPinPos    []rl.Vector2
	OutputPinPos   rl.Vector2
	SnapTargetRect rl.Rectangle

	Color rl.Color
}

func NewTable(table string, alias string) *Node {
	return &Node{
		Title:   "Table",
		CanSnap: false,
		Color: rl.NewColor(242, 201, 76, 255),
		Data: &Table{
			Table: table,
			Alias: alias,
		},
	}
}

func NewPickColumns(alias string) *Node {
	return &Node{
		Title:   "Pick Columns",
		CanSnap: true,
		Color: rl.NewColor(244, 143, 177, 255),
		Inputs:  make([]*Node, 1),
		Data: &PickColumns{
			Alias: alias,
		},
	}
}

func NewFilter(conditions []string) *Node {
	return &Node{
		Title:   "Filter",
		CanSnap: true,
		Color: rl.NewColor(111, 207, 151, 255),
		Inputs:  make([]*Node, 1),
		Data: &Filter{
			Conditions: conditions,
		},
	}
}

func NewCombineRows(combineType CombineType) *Node {
	return &Node{
		Title:   "Combine Rows",
		CanSnap: false,
		Inputs:  make([]*Node, 1),
		Data: &CombineRows{
			CombinationType: combineType,
		},
	}
}

func NewNodeGenContext() *NodeGenContext {
	return &NodeGenContext{}
}

func NewRecursiveGenerated(n *Node) *NodeGenContext {
	return NewNodeGenContext().RecursiveGenerate(n)
}

// SourceToSql Turns a context tree into an SQL statement string
func (ctx *NodeGenContext) SourceToSql() string {
	sql := "SELECT "

	if len(ctx.Cols) == 0 {
		sql += "*"
	} else {
		sql += strings.Join(ctx.Cols, ", ")
	}

	if ctx.Source == nil {
		panic("We can't generate SQL if our gen context has no source!")
	} else {
		fmt.Println(fmt.Sprintf("child element: %s ||| %s", reflect.TypeOf(ctx.Source), ctx.Source.SourceToSql()))
		fmt.Println(fmt.Sprintf("it's alias is: %s ### %s", ctx.Source.SourceAlias(), ctx.Alias))
		sql += fmt.Sprintf(" FROM (%s)", ctx.Source.SourceToSql())

		// Currently only shows alias if it's not empty
		if ctx.Source.SourceAlias() != "" {
			sql += fmt.Sprintf(" AS %s", ctx.Alias)
		}
	}

	if len(ctx.FilterConditions) > 0 {
		sql += " WHERE "
		sql += strings.Join(ctx.FilterConditions, " AND ")
	}

	// handle combining here

	return sql
}

// RecursiveGenerate Turns a node into a recursive context tree for SQL generation
func (ctx *NodeGenContext) RecursiveGenerate(n *Node) *NodeGenContext {
	switch d := n.Data.(type) {
	case *Table:
		ctx.Source = d
		ctx.Alias = d.Alias
		ctx.RecursiveGenerateAllChildren(n)
	case *PickColumns:
		if len(ctx.Cols) > 0 {
			ctx.Source = NewRecursiveGenerated(n)
			ctx.Alias = d.Alias
		} else {
			ctx.Cols = d.Cols
			ctx.RecursiveGenerateAllChildren(n)
		}
	case *Filter:
		ctx.FilterConditions = append(ctx.FilterConditions, d.Conditions...)
		ctx.RecursiveGenerateAllChildren(n)
	}
	return ctx
}

func (ctx *NodeGenContext) RecursiveGenerateAllChildren(n *Node) *NodeGenContext {
	for _, value := range n.Inputs {
		ctx.RecursiveGenerate(value)
	}
	return ctx
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
		var cols = d.Cols
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

func (n *Node) Rect() rl.Rectangle {
	return rl.Rectangle{n.Pos.X, n.Pos.Y, n.Size.X, n.Size.Y}
}
