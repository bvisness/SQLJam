package node

type NodeData interface {
	_implNodeData() // delete if we ever actually put a meaningful method here
}

type SqlSource interface {
	SourceToSql() string
	SourceAlias() string
}

type GenCombine struct {
	Context *QueryContext
	Type    CombineType
}

type GenOrder struct {
	Col        string
	Descending bool
}

type GenJoin struct {
	Source    SqlSource
	Condition string
	Type      JoinType
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
	Joins    []GenJoin

	FilterConditions []string

	Orders []GenOrder
}

var _ SqlSource = &QueryContext{}

func (ctx *QueryContext) SourceAlias() string {
	return ctx.Alias
}
