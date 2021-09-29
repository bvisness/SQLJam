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

	ColDropdowns []raygui.DropdownEx
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
	Dropdown raygui.DropdownEx
}

type Join struct {
	NodeData
	Conditions []string
}

type Filter struct {
	NodeData
	Conditions string // TODO: whatever data we actually need for our filter UI

	// UI data
	TextBox raygui.TextBoxEx
}

type Order struct {
	NodeData

	Alias      string
	Cols       []string
	Descending bool

	ColDropdowns []raygui.DropdownEx
}

// A context for node generation recursion.
// Eventually, we can no longer add onto this query. Thus,
// we continue recursive generation with a new Source context object.
// Thus this is basically a recursive tree

type GenCombine struct {
	Context *NodeGenContext
	Type CombineType
}

type NodeGenContext struct {
	SqlSource
	Alias            string
	Cols             []string
	Source           SqlSource // or NodeGenContext

	Combines         []GenCombine

	Joins            []NodeGenContext
	JoinConditions   []string

	FilterConditions []string
}

func (ctx *NodeGenContext) SourceAlias() string {
	return ctx.Alias
}
