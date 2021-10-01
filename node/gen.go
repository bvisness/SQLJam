package node

import (
	"fmt"
	"strings"
)

func NewQueryContext(source SqlSource) *QueryContext {
	return &QueryContext{
		Source: source,
	}
}

func NewQueryContextFromNode(n *Node) *QueryContext {
	return NewQueryContext(nil).CreateQuery(n)
}

// SourceToSql Turns a context tree into an SQL statement string
func (ctx *QueryContext) SourceToSql() string {
	var sql string

	if len(ctx.Combines) > 0 {
		sql += ctx.Source.SourceToSql()
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
			sql = fmt.Sprintf("%s %s %s", sql, used, gc.Context.SourceToSql())
		}
	} else {
		sql += "SELECT "

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
			return "Error: No SQL Source"
		} else {
			//fmt.Println(fmt.Sprintf("child element: %s ||| %s", reflect.TypeOf(ctx.Source), ctx.Source.SourceToSql()))
			//fmt.Println(fmt.Sprintf("it's alias is: %s ### %s", ctx.Source.SourceAlias(), ctx.Alias))
			//fmt.Println(fmt.Sprintf("Doot %s", reflect.TypeOf(ctx.Source)))
			switch ctx.Source.(type) {
			case *Table:
				sql += fmt.Sprintf(" FROM %s", ctx.Source.SourceToSql())
			default:
				sql += fmt.Sprintf(" FROM (%s)", ctx.Source.SourceToSql())
			}

			// Currently only shows alias if it's not empty
			if ctx.Source.SourceAlias() != "" {
				sql += fmt.Sprintf(" AS %s", ctx.Alias)
			}
		}

		if len(ctx.FilterConditions) > 0 && ctx.FilterConditions[0] != "" {
			sql += " WHERE "
			sql += strings.Join(ctx.FilterConditions, " AND ")
		}
	}

	if len(ctx.Orders) > 0 {
		sql += " ORDER BY "
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
			ctx = NewQueryContext(ctx)
		}

		ctx.Cols = d.Cols()
		ctx.ColAliases = d.Aliases()
		// ctx.Alias = d.Alias // TODO: ???????
	case *Filter:
		ctx = ctx.CreateQuery(n.Inputs[0])
		if len(ctx.Cols) > 0 {
			ctx = NewQueryContext(ctx)
		}

		ctx.FilterConditions = append(ctx.FilterConditions, d.Conditions) // TODO: This should be split into multiple again? Right??
	case *Order:
		ctx = ctx.CreateQuery(n.Inputs[0])
		if len(ctx.Orders) > 0 {
			ctx = NewQueryContext(ctx)
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

		ctx = NewQueryContext(firstCtx)

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
		numNotNull := 0
		for _, input := range n.Inputs {
			if input != nil {
				numNotNull++
			}
		}

		// Top input gets recursive generated as normal in same context
		if n.Inputs[0] != nil {
			ctx.CreateQuery(n.Inputs[0])
			// Only do combines if first wire and at least one other wire is connected
			if numNotNull >= 2 {
				// All other inputs get thrown into a new recursive context
				for _, input := range n.Inputs[1:] {
					if input != nil {
						newCtx := NewQueryContextFromNode(input)
						ctx.Joins = append(ctx.Joins, newCtx)
					}
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
	return NewQueryContextFromNode(n).SourceToSql()
}
