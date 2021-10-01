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

	ColDropdowns []*raygui.DropdownEx
}

type OrderColumn struct {
	Col        string
	Descending bool
}

type GenCombine struct {
	Context *QueryContext
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

type QueryContext struct {
	Alias      string
	Cols       []string
	ColAliases []string
	Source     SqlSource // or NodeGenContext

	Combines []GenCombine

	Joins          []*QueryContext
	JoinConditions []string

	FilterConditions []string

	Orders []GenOrder
}

var _ SqlSource = &QueryContext{}

func (ctx *QueryContext) SourceAlias() string {
	return ctx.Alias
}
