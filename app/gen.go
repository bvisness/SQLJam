package app

import (
	"fmt"
	"strings"
)

type SqlSource interface {
	SourceToSql(indent int) string
	SourceTableName() string
	IsTable() bool
}

// QueryContext A context for node generation recursion.
// Eventually, we can no longer add onto this query. Thus,
// we continue recursive generation with a new Source context object.
// Thus this is basically a recursive tree
type QueryContext struct {
	Source SqlSource // or NodeGenContext

	// Picked columns and aggregates are mutually exclusive.
	Cols      []GenColumn
	Aggregate *GenAggregate

	Combines         []GenCombine
	JoinSourceAlias  string
	Joins            []GenJoin
	WhereConditions  []string
	HavingConditions []string
	Orders           []GenOrder
}

var _ SqlSource = &QueryContext{}

func (ctx *QueryContext) IsTable() bool {
	return false
}

func (ctx *QueryContext) SourceTableName() string {
	return ctx.Source.SourceTableName()
}

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

type GenColumn struct {
	Col   string
	Alias string
}

type GenAggregate struct {
	GroupByCols []string
	Aggs        []GenAggregateEntry
}

type GenAggregateEntry struct {
	Type  AggregateType
	Col   string
	Alias string
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
	Type      JoinType
	Source    SqlSource
	Alias     string
	Condition string
}

func indented(s string, amount int) string {
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
			sql += "\n" + indented(used, indent)
			sql += "\n" + gc.Context.SourceToSql(indent) + "\n"
		}
	} else {
		sql += indented("SELECT ", indent)

		if len(ctx.Cols) > 0 && ctx.Aggregate != nil {
			return indented("Error: Pick columns and aggregate on the same context", indent)
		}

		if ctx.Aggregate != nil {
			var colStrings []string
			for _, gbCol := range ctx.Aggregate.GroupByCols {
				colStrings = append(colStrings, gbCol)
			}
			for _, agg := range ctx.Aggregate.Aggs {
				op := "ERROR("
				switch agg.Type {
				case Avg:
					op = "AVG("
				case Max:
					op = "MAX("
				case Min:
					op = "MIN("
				case Sum:
					op = "SUM("
				case Count:
					op = "COUNT("
				case CountDistinct:
					op = "COUNT(DISTINCT "
				}

				aliasStr := ""
				if agg.Alias != "" {
					aliasStr = fmt.Sprintf(" AS %s", agg.Alias)
				}

				colStrings = append(colStrings, fmt.Sprintf("%s%s)%s", op, agg.Col, aliasStr))
			}
			sql += strings.Join(colStrings, ", ")
		} else if len(ctx.Cols) == 0 {
			sql += "*"
		} else {
			colStrings := make([]string, len(ctx.Cols))
			for i, col := range ctx.Cols {
				aliasStr := ""
				if col.Alias != "" {
					aliasStr = fmt.Sprintf(" AS %s", col.Alias)
				}
				colStrings[i] = fmt.Sprintf("%s%s", col.Col, aliasStr)
			}
			sql += strings.Join(colStrings, ", ")
		}

		if ctx.Source == nil {
			return indented("Error: No SQL Source", indent)
		} else {
			sql += "\n" + indented("", indent)
			if ctx.Source.IsTable() {
				sql += fmt.Sprintf("FROM %s", ctx.Source.SourceToSql(0))
			} else {
				sql += fmt.Sprintf("FROM (\n%s", ctx.Source.SourceToSql(indent+1))
				sql += "\n" + indented(")", indent)
			}
		}

		// JOIN
		if ctx.JoinSourceAlias != "" {
			sql += fmt.Sprintf(" AS %s", ctx.JoinSourceAlias)
		}
		for _, join := range ctx.Joins {
			sql += "\n" + indented(join.Type.String(), indent) + " "
			if join.Source.IsTable() {
				sql += join.Source.(*Table).Table
			} else {
				sql += "(\n" + join.Source.SourceToSql(indent+1)
				sql += "\n" + indented(")", indent)
			}
			if join.Alias != "" {
				sql += fmt.Sprintf(" AS %s", join.Alias)
			}
			if join.Condition != "" {
				sql += fmt.Sprintf(" ON %s", join.Condition)
			}
		}

		if len(ctx.WhereConditions) > 0 && ctx.WhereConditions[0] != "" {
			sql += "\n" + indented("WHERE ", indent)
			sql += strings.Join(ctx.WhereConditions, " AND ")
		}
	}

	if ctx.Aggregate != nil && len(ctx.Aggregate.GroupByCols) > 0 {
		sql += "\n" + indented("GROUP BY ", indent)
		sql += strings.Join(ctx.Aggregate.GroupByCols, ", ")
	}

	if len(ctx.HavingConditions) > 0 {
		sql += "\n" + indented("HAVING ", indent)
		sql += strings.Join(ctx.HavingConditions, " AND ")
	}

	if len(ctx.Orders) > 0 {
		sql += indented("\nORDER BY ", indent)
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
	case *PickColumns:
		ctx = ctx.CreateQuery(n.Inputs[0])
		if len(ctx.Cols) > 0 || ctx.Aggregate != nil {
			ctx = WrapQueryContext(ctx)
		}

		for _, entry := range d.Entries {
			ctx.Cols = append(ctx.Cols, GenColumn{
				Col:   entry.Col,
				Alias: entry.Alias,
			})
		}
	case *Filter:
		ctx = ctx.CreateQuery(n.Inputs[0])
		if len(ctx.Cols) > 0 {
			ctx = WrapQueryContext(ctx)
		}

		if ctx.Aggregate != nil {
			ctx.HavingConditions = append(ctx.HavingConditions, d.Conditions)
		} else {
			ctx.WhereConditions = append(ctx.WhereConditions, d.Conditions)
		}
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
		if n.Inputs[0] == nil {
			break
		}

		/*
			We get our list of column names by getting the schemas of all our
			inputs and concatenating the lists together. We then have to
			resolve duplicate column names - if there are no duplicates, we
			can do SELECT * as usual, but if there are, we need to generate a
			list of column names like "a.film_id AS a_film_id, b.film_id AS
			b_film_id, ...".
		*/
		type inputSchema struct {
			Alias       string
			ColumnNames []string
		}
		var inputSchemas []inputSchema

		if table, ok := n.Inputs[0].Data.(*Table); ok {
			ctx = NewQueryContext()
			ctx.Source = table
		} else {
			firstCtx := NewQueryContextFromNode(n.Inputs[0])
			ctx = WrapQueryContext(firstCtx)
		}

		ctx.JoinSourceAlias = d.FirstAlias
		inputSchemas = append(inputSchemas, inputSchema{
			Alias:       d.FirstAlias,
			ColumnNames: getSchema(n.Inputs[0]),
		})

		// All other inputs get thrown into a new recursive context
		for i, input := range n.Inputs[1:] {
			if input == nil {
				continue
			}

			alias := d.Conditions[i].Alias

			inputSchemas = append(inputSchemas, inputSchema{
				Alias:       alias,
				ColumnNames: getSchema(input),
			})

			var source SqlSource
			if table, ok := input.Data.(*Table); ok {
				source = table
			} else {
				source = NewQueryContextFromNode(input)
			}

			cond := d.Conditions[i]
			ctx.Joins = append(ctx.Joins, GenJoin{
				Source:    source,
				Condition: cond.Condition,
				Type:      cond.Type(),
				Alias:     alias,
			})
		}

		colCounts := map[string]int{}
		for _, schema := range inputSchemas {
			for _, col := range schema.ColumnNames {
				current, _ := colCounts[col]
				colCounts[col] = current + 1
			}
		}

		anyDuplicates := false
		for _, count := range colCounts {
			if count > 1 {
				anyDuplicates = true
				break
			}
		}

		if anyDuplicates {
			for _, schema := range inputSchemas {
				for _, col := range schema.ColumnNames {
					specificCol := fmt.Sprintf("%s.%s", schema.Alias, col)

					alias := ""
					if colCounts[col] > 1 {
						alias = fmt.Sprintf("%s_%s", schema.Alias, col)
					}

					ctx.Cols = append(ctx.Cols, GenColumn{
						Col:   specificCol,
						Alias: alias,
					})
				}
			}
		}
	case *Aggregate:
		ctx = ctx.CreateQuery(n.Inputs[0])
		if len(ctx.Cols) > 0 || ctx.Aggregate != nil {
			ctx = WrapQueryContext(ctx)
		}

		groupByCols := make([]string, len(d.GroupBys))
		for i, gb := range d.GroupBys {
			groupByCols[i] = gb.Col
		}

		aggs := make([]GenAggregateEntry, len(d.Aggregates))
		for i, agg := range d.Aggregates {
			aggs[i] = GenAggregateEntry{
				Type:  agg.Type,
				Col:   agg.Col,
				Alias: agg.Alias,
			}
		}

		ctx.Aggregate = &GenAggregate{
			GroupByCols: groupByCols,
			Aggs:        aggs,
		}
	case *Preview:
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

func (n *Node) GenerateSql(limit bool) string {
	sql := NewQueryContextFromNode(n).SourceToSql(0)
	if limit {
		sql += " LIMIT 1000"
	}
	return sql
}
