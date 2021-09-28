package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/bvisness/SQLJam/node"
	rl "github.com/gen2brain/raylib-go/raylib"
	_ "github.com/mattn/go-sqlite3"
)

const screenWidth = 1920
const screenHeight = 1080

var db *sql.DB
var nodes []*node.Node

func main() {
	rl.InitWindow(screenWidth, screenHeight, "SQL Jam")
	defer rl.CloseWindow()

	rl.SetTargetFPS(120) // wew

	var err error
	db, err = sql.Open("sqlite3", "./sakila.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// init nodes
	filmTable := node.NewTable("film")
	filmTable.Pos = rl.Vector2{60, 100}
	nodes = append(nodes, filmTable)

	filter := node.NewFilter([]string{"rating = 'PG'", "rental_rate < 3"})
	filter.Pos = rl.Vector2{160, 100}
	filter.Inputs = append(filter.Inputs, filmTable)
	nodes = append(nodes, filter)

	pick := node.NewPickColumns()
	pick.Pos = rl.Vector2{260, 100}
	pick.Data.(*node.PickColumns).Cols["title"] = true
	pick.Inputs = append(pick.Inputs, filter)
	nodes = append(nodes, pick)

	// main frame loop
	for !rl.WindowShouldClose() {
		doFrame()
	}
}

var latestResult *queryResult

func doFrame() {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(rl.RayWhite)

	for _, node := range nodes {
		nodeRect := rl.Rectangle{node.Pos.X, node.Pos.Y, 80, 60}
		rl.DrawRectangleLinesEx(nodeRect, 2, rl.Black)

		isHover := rl.CheckCollisionPointRec(rl.GetMousePosition(), nodeRect)
		isClick := isHover && rl.IsMouseButtonPressed(rl.MouseLeftButton) // TODO: better clicking (on release)
		if isHover {
			rl.DrawText(node.SQL(), int32(nodeRect.X), int32(nodeRect.Y)-22, 20, rl.Black)
		}
		if isClick {
			latestResult = doQuery(node.SQL())
		}
	}

	if latestResult != nil {
		rowPos := rl.Vector2{60, 400}
		for i := -1; i < len(latestResult.Rows); i++ {
			if i < 0 {
				// print headers
				rl.DrawText(strings.Join(latestResult.Columns, "    "), int32(rowPos.X), int32(rowPos.Y), 20, rl.Black)
			} else {
				row := latestResult.Rows[i]
				valStrings := make([]string, len(row))
				for i, v := range row {
					valStrings[i] = fmt.Sprintf("%v", v)
				}
				rl.DrawText(strings.Join(valStrings, "    "), int32(rowPos.X), int32(rowPos.Y), 20, rl.Black)
			}

			rowPos.Y += 24
		}
	}
}

func makeDropdownOptions(opts ...string) string {
	return strings.Join(opts, ";")
}

// TODO: Surely this is pretty temporary. I just need to display boring query output.
type queryResult struct {
	Columns []string
	Rows    [][]interface{}
}

func doQuery(q string) *queryResult {
	rows, err := db.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var res queryResult

	res.Columns, err = rows.Columns()
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		row := make([]interface{}, len(res.Columns))
		rowPointers := make([]interface{}, len(row))
		for i := range row {
			rowPointers[i] = &row[i]
		}

		err = rows.Scan(rowPointers...)
		if err != nil {
			panic(err)
		}
		res.Rows = append(res.Rows, row)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return &res
}
