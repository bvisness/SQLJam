package node

import "github.com/bvisness/SQLJam/raygui"

type NodeData interface {
	_implNodeData() // delete if we ever actually put a meaningful method here
}

type SqlSource interface {
	SourceToSql() string
	SourceAlias() string
}

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

type PickColumns struct {
	NodeData
	SqlSource
	Alias string
	Cols  []string
}

func (pc *PickColumns) SourceAlias() string {
	return pc.Alias
}

type CombineType int

const (
	Union CombineType = iota + 1
	Intersect
	Except
	UnionAll
	IntersectAll
	ExceptAll
)

type CombineRows struct {
	NodeData
	CombinationType CombineType
}

type Join struct {
	NodeData
	Conditions []string
}

type Filter struct {
	NodeData
	Conditions []string // TODO: whatever data we actually need for our filter UI
	// for now it always does AND because I am testing
}

// A context for node generation recursion.
// At certain times, we will need to push

type NodeGenContext struct {
	SqlSource
	Alias            string
	Cols             []string
	Source           SqlSource // or NodeGenContext
	Combines         []NodeGenContext
	Joins            []NodeGenContext
	JoinConditions   []string
	FilterConditions []string
}

func (ctx *NodeGenContext) SourceAlias() string {
	return ctx.Alias
}
