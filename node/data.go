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
	UnionAll
	Intersect
	Except
)

type CombineRows struct {
	NodeData
	CombinationType CombineType
	Dropdown        raygui.DropdownEx
}

type JoinType int

const (
	LeftJoin JoinType = iota + 1
	RightJoin
	InnerJoin
	OuterJoin
)

type JoinCondition struct {
	Type      JoinType
	Condition string
	TextBox   *raygui.TextBoxEx
}

type Join struct {
	NodeData
	Conditions []*JoinCondition
}

type Filter struct {
	NodeData
	Conditions string // TODO: whatever data we actually need for our filter UI

	// UI data
	TextBox raygui.TextBoxEx
}

type Order struct {
	NodeData

	Alias string
	Cols  []OrderColumn

	ColDropdowns []raygui.DropdownEx
}

type OrderColumn struct {
	Col        string
	Descending bool
}

type GenCombine struct {
	Context *NodeGenContext
	Type    CombineType
}

type GenOrder struct {
	Col        string
	Descending bool
}

// A context for node generation recursion.
// Eventually, we can no longer add onto this query. Thus,
// we continue recursive generation with a new Source context object.
// Thus this is basically a recursive tree

type NodeGenContext struct {
	SqlSource
	Alias  string
	Cols   []string
	Source SqlSource // or NodeGenContext

	Combines []GenCombine

	Joins          []NodeGenContext
	JoinConditions []string

	FilterConditions []string

	Orders []GenOrder
}

func (ctx *NodeGenContext) SourceAlias() string {
	return ctx.Alias
}
