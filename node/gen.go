package node

import (
	"fmt"
	"strings"
)

// NewQueryContext Creates an empty query context.
func NewQueryContext() *QueryContext {
	return &QueryContext{}
}

// NewQueryContextFromNode Creates a query context and fills it with data from a given node.
// For example, calling this with a Table/Filter will result in this
// structure:
//
//	context:
//	  source: table
//    filters: filters
func NewQueryContextFromNode(n *Node) *QueryContext {
	return NewQueryContext().CreateQuery(n)
}

// WrapQueryContext Wraps the given query context in a new, empty context. Use when you
// want to break something...into a subquery.
func WrapQueryContext(ctx *QueryContext) *QueryContext {
	return &QueryContext{
		Source: ctx,
	}
}

func Indented(s string, amount int) string {
	return strings.Repeat("\t", amount) + s
}

// SourceToSql Turns a context tree into an SQL statement string
func (ctx *QueryContext) SourceToSql(indent int) string {
	var sql string

	if len(ctx.Combines) > 0 {
		sql += ctx.Source.SourceToSql(indent)
		for _, gc := range ctx.Combines {
			used := ""
			switch gc.Type {
			case Union:
				used = "UNION"
			case Intersect:
				used = "INTERSECT"
			case Except:
				used = "EXCEPT"
			case UnionAll:
				used = "UNION ALL"
			}
			sql += "\n" + Indented(used, indent)
			sql += "\n" + gc.Context.SourceToSql(indent) + "\n"
		}
	} else {
		sql += Indented("SELECT ", indent)

		if len(ctx.Cols) == 0 {
			sql += "*"
		} else {
			colStrings := make([]string, len(ctx.Cols))
			for i := range ctx.Cols {
				aliasStr := ""
				if ctx.ColAliases[i] != "" {
					aliasStr = fmt.Sprintf(" AS %s", ctx.ColAliases[i])
				}
				colStrings[i] = fmt.Sprintf("%s%s", ctx.Cols[i], aliasStr)
			}
			sql += strings.Join(colStrings, ", ")
		}

		if ctx.Source == nil {
			return Indented("Error: No SQL Source", indent)
		} else {
			sql += "\n" + Indented("", indent)
			switch ctx.Source.(type) {
			case *Table:
				sql += fmt.Sprintf("FROM %s", ctx.Source.SourceToSql(0))
			default:
				sql += fmt.Sprintf("FROM (\n%s", ctx.Source.SourceToSql(indent + 1))
				sql += "\n" + Indented(")", indent)
			}

			sql += " AS a"
		}

		for _, join := range ctx.Joins {
			sql += "\n" + Indented(join.Type.String(), indent) + " "
			switch join.Source.(type) {
			case *Table:
				sql += join.Source.SourceToSql(indent)
			default:
				sql += "(\n" + join.Source.SourceToSql(indent + 1)
				sql += "\n" + Indented(")", indent)
			}
			sql += fmt.Sprintf(" AS %s", join.Source.SourceAlias())
			if join.Condition != "" {
				sql += fmt.Sprintf(" ON %s", join.Condition)
			}
		}

		if len(ctx.FilterConditions) > 0 && ctx.FilterConditions[0] != "" {
			sql += "\n" + Indented("WHERE ", indent)
			sql += strings.Join(ctx.FilterConditions, " AND ")
		}
	}

	if len(ctx.Orders) > 0 {
		sql += Indented("\nORDER BY ", indent)
		var orderStrings []string
		for _, order := range ctx.Orders {
			direction := ""
			if order.Descending {
				direction = " DESC"
			}
			orderStrings = append(orderStrings, fmt.Sprintf("%s%s", order.Col, direction))
		}
		sql += strings.Join(orderStrings, ", ")
	}

	return sql
}

// CreateQuery Turns a node into a recursive context tree for SQL generation
func (ctx *QueryContext) CreateQuery(n *Node) *QueryContext {
	if n == nil {
		return ctx
	}

	switch d := n.Data.(type) {
	case *Table:
		ctx.Source = d
		ctx.Alias = d.Alias
	case *PickColumns:
		ctx = ctx.CreateQuery(n.Inputs[0])
		if len(ctx.Cols) > 0 {
			ctx = WrapQueryContext(ctx)
		}

		ctx.Cols = d.Cols()
		ctx.ColAliases = d.Aliases()
		// ctx.Alias = d.Alias // TODO: ???????
	case *Filter:
		ctx = ctx.CreateQuery(n.Inputs[0])
		if len(ctx.Cols) > 0 {
			ctx = WrapQueryContext(ctx)
		}

		ctx.FilterConditions = append(ctx.FilterConditions, d.Conditions) // TODO: This should be split into multiple again? Right??
	case *Order:
		ctx = ctx.CreateQuery(n.Inputs[0])
		if len(ctx.Orders) > 0 {
			ctx = WrapQueryContext(ctx)
		}

		for _, col := range d.Cols {
			ctx.Orders = append(ctx.Orders, GenOrder{
				Col:        col.Col,
				Descending: col.Descending,
			})
		}
	case *CombineRows:
		firstCtx := NewQueryContextFromNode(n.Inputs[0])
		firstCtx.Orders = nil // anything involved in Combine Rows can't use ORDER BY

		ctx = WrapQueryContext(firstCtx)

		// All other inputs get thrown into a new recursive context
		for _, input := range n.Inputs[1:] {
			if input != nil {
				newCtx := NewQueryContextFromNode(input)
				newCtx.Orders = nil
				ctx.Combines = append(ctx.Combines, GenCombine{
					Context: newCtx,
					Type:    d.CombinationType,
				})
			}
		}
	case *Join:
		if n.Inputs[0] != nil {
			if table, ok := n.Inputs[0].Data.(*Table); ok {
				ctx = NewQueryContext()
				ctx.Source = table
			} else {
				firstCtx := NewQueryContextFromNode(n.Inputs[0])
				ctx = WrapQueryContext(firstCtx)
			}

			// All other inputs get thrown into a new recursive context
			for i, input := range n.Inputs[1:] {
				if input != nil {
					aliasChar := string(rune('b' + i))
					cond := d.Conditions[i]
					var source SqlSource
					if table, ok := input.Data.(*Table); ok {
						table.Alias = aliasChar
						source = table
					} else {
						newQuery := NewQueryContextFromNode(input)
						newQuery.Alias = aliasChar
						source = newQuery
					}
					ctx.Joins = append(ctx.Joins, GenJoin{
						Source:    source,
						Condition: cond.Condition,
						Type:      cond.Type(),
					})
				}
			}
		}
	}

	return ctx
}

func (ctx *QueryContext) RecursiveGenerateInputs(n *Node) *QueryContext {
	for _, value := range n.Inputs {
		if value != nil {
			ctx.CreateQuery(value)
		}
	}
	return ctx
}

func (n *Node) GenerateSql() string {
	return NewQueryContextFromNode(n).SourceToSql(0)
}
