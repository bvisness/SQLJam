package node

import (
	"fmt"
	"strings"
)

func NewNodeGenContext() *NodeGenContext {
	return &NodeGenContext{}
}

func NewRecursiveContext(n *Node) *NodeGenContext {
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

	// handle combining here

	if len(ctx.Combines) > 0 {
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
			sql += fmt.Sprintf(" %s %s ", used, gc.Context.SourceToSql())
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

// RecursiveGenerate Turns a node into a recursive context tree for SQL generation
func (ctx *NodeGenContext) RecursiveGenerate(n *Node) *NodeGenContext {
	//fmt.Println(fmt.Sprintf("test1 %s", n))
	//fmt.Println(fmt.Sprintf("test2 %s", reflect.TypeOf(n)))
	switch d := n.Data.(type) {
	case *Table:
		ctx.Source = d
		ctx.Alias = d.Alias
		ctx.RecursiveGenerateChildren(n)
	case *PickColumns:
		if len(ctx.Cols) > 0 {
			ctx.Source = NewRecursiveContext(n)
			ctx.Alias = d.Alias
		} else {
			ctx.Cols = d.Cols
			ctx.RecursiveGenerateChildren(n)
		}
	case *Filter:
		ctx.FilterConditions = append(ctx.FilterConditions, d.Conditions) // TODO: This should be split into multiple again? Right??
		ctx.RecursiveGenerateChildren(n)
	case *Order:
		for _, col := range d.Cols {
			ctx.Orders = append(ctx.Orders, GenOrder{
				Col:        col.Col,
				Descending: col.Descending,
			})
		}
		ctx.RecursiveGenerateChildren(n)
	case *CombineRows:
		numNotNull := 0
		for _, input := range n.Inputs {
			if input != nil {
				numNotNull++
			}
		}

		// Top input gets recursive generated as normal in same context
		if n.Inputs[0] != nil {
			ctx.RecursiveGenerate(n.Inputs[0])
			// Only do combines if first wire and at least one other wire is connected
			if numNotNull >= 2 {
				// All other inputs get thrown into a new recursive context
				for _, input := range n.Inputs[1:] {
					if input != nil {
						ctx.Combines = append(ctx.Combines, GenCombine{
							Context: NewRecursiveContext(input),
							Type:    d.CombinationType,
						})
					}
				}
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
			ctx.RecursiveGenerate(n.Inputs[0])
			// Only do combines if first wire and at least one other wire is connected
			if numNotNull >= 2 {
				// All other inputs get thrown into a new recursive context
				for _, input := range n.Inputs[1:] {
					if input != nil {
						ctx.Joins = append(ctx.Joins, *NewRecursiveContext(input))
					}
				}
			}
		}
	}
	return ctx
}

func (ctx *NodeGenContext) RecursiveGenerateChildren(n *Node) *NodeGenContext {
	for _, value := range n.Inputs {
		if value != nil {
			ctx.RecursiveGenerate(value)
		}
	}
	return ctx
}

func (n *Node) GenerateSql() string {
	return NewRecursiveContext(n).SourceToSql()
}
