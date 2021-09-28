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
	filmTable := node.NewTable("film", "cool_films")
	filmTable.Pos = rl.Vector2{60, 100}
	nodes = append(nodes, filmTable)

	filter := node.NewFilter([]string{"rating = 'PG'", "rental_rate < 3"})
	filter.Pos = rl.Vector2{360, 100}
	filter.Inputs[0] = filmTable
	nodes = append(nodes, filter)

	pick := node.NewPickColumns("test_alias")
	pick.Pos = rl.Vector2{260, 100}
	pick.Data.(*node.PickColumns).Cols = append(pick.Data.(*node.PickColumns).Cols, "title")
	pick.Inputs[0] = filter
	nodes = append(nodes, pick)


	// main frame loop
	for !rl.WindowShouldClose() {
		doFrame()
	}


	ctxTree := node.NewRecursiveGenerated(pick) // try recursive gen
	fmt.Println(ctxTree.SourceToSql())
}

var latestResult *queryResult

var cam = rl.Camera2D{
	Offset: rl.Vector2{screenWidth / 2, screenHeight / 2},
	Target: rl.Vector2{screenWidth / 2, screenHeight / 2},
	Zoom:   1,
}
var panMouseStart rl.Vector2
var panCamStart rl.Vector2

func doFrame() {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(rl.RayWhite)

	// Move camera
	if rl.IsMouseButtonDown(rl.MouseMiddleButton) {
		mousePos := rl.GetMousePosition()
		if rl.IsMouseButtonPressed(rl.MouseMiddleButton) {
			panMouseStart = mousePos
			panCamStart = cam.Target
		}

		mouseDelta := rl.Vector2Subtract(mousePos, panMouseStart)
		cam.Target = rl.Vector2Subtract(panCamStart, mouseDelta) // camera moves opposite of mouse
	}

	rl.BeginMode2D(cam)
	defer rl.EndMode2D()

	CheckCollisionPointRec2D := func(point rl.Vector2, rec rl.Rectangle) bool {
		screenRec := rl.Rectangle{
			X:      rec.X - (cam.Target.X - cam.Offset.X),
			Y:      rec.Y - (cam.Target.Y - cam.Offset.Y),
			Width:  rec.Width,
			Height: rec.Height,
		}
		return rl.CheckCollisionPointRec(point, screenRec)
	}

	if rl.IsMouseButtonUp(rl.MouseLeftButton) {
		clearDrag()
	}

	doLayout()

	// draw lines
	for _, n := range nodes {
		for i, input := range n.Inputs {
			rl.DrawLineEx(input.OutputPinPos, n.InputPinPos[i], 2, rl.Black)
		}
	}

	// draw nodes
	for _, n := range nodes {
		nodeRect := n.Rect()
		rl.DrawRectangleRounded(nodeRect, 0.08, 6, rl.LightGray)
		rl.DrawRectangleRoundedLines(nodeRect, 0.08, 6, 2, rl.Black)

		titleBarRect := rl.Rectangle{nodeRect.X, nodeRect.Y, nodeRect.Width - 24, 24}
		previewRect := rl.Rectangle{nodeRect.X + nodeRect.Width - 24, nodeRect.Y, 24, 24}

		rl.DrawText(n.Title, int32(nodeRect.X)+6, int32(nodeRect.Y)+4, 20, rl.Black) // title bar
		rl.DrawText("P", int32(previewRect.X)+4, int32(previewRect.Y)+10, 10, rl.Black)

		for _, pinPos := range n.InputPinPos {
			rl.DrawCircle(int32(pinPos.X), int32(pinPos.Y), 6, rl.Black)
		}
		rl.DrawCircle(int32(n.OutputPinPos.X), int32(n.OutputPinPos.Y), 6, rl.Black)

		titleHover := CheckCollisionPointRec2D(rl.GetMousePosition(), titleBarRect)
		if titleHover {
			rl.DrawText(n.SQL(false), int32(titleBarRect.X), int32(titleBarRect.Y)-22, 20, rl.Black)
		}
		if titleHover && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			tryStartDrag(n, n.Pos)
		}

		if d, ok := dragging.(*node.Node); ok && d == n {
			n.Pos = dragNewPosition()
		}

		previewHover := CheckCollisionPointRec2D(rl.GetMousePosition(), previewRect)
		if previewHover && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			latestResult = doQuery(n.SQL(false))
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

var dragging interface{}
var dragMouseStart rl.Vector2
var dragObjStart rl.Vector2

func tryStartDrag(thing interface{}, objStart rl.Vector2) {
	if dragging != nil {
		return
	}

	dragging = thing
	dragMouseStart = rl.GetMousePosition()
	dragObjStart = objStart
}

func clearDrag() {
	dragging = nil
}

func dragOffset() rl.Vector2 {
	return rl.Vector2Subtract(rl.GetMousePosition(), dragMouseStart)
}

func dragNewPosition() rl.Vector2 {
	return rl.Vector2Add(dragObjStart, dragOffset())
}

func doLayout() {
	const titleBarHeight = 24
	const pinDefaultSpacing = 36

	for _, n := range nodes {
		// TODO: do different stuff for different node types

		width := 280
		pinStartHeight := 36
		inputHeight := titleBarHeight
		outputHeight := titleBarHeight

		if len(n.InputPinPos) != len(n.Inputs) {
			n.InputPinPos = make([]rl.Vector2, len(n.Inputs))
		}

		pinHeight := pinStartHeight
		for i := range n.Inputs {
			n.InputPinPos[i] = rl.Vector2{n.Pos.X - 1, n.Pos.Y + float32(pinHeight)}
			pinHeight += pinDefaultSpacing
			inputHeight += pinDefaultSpacing
		}
		inputHeight += 20

		if !n.Snapped {
			outputHeight += pinDefaultSpacing
		}

		height := inputHeight
		if outputHeight > height {
			height = outputHeight
		}

		n.Size = rl.Vector2{float32(width), float32(height)}

		n.OutputPinPos = rl.Vector2{n.Pos.X + n.Size.X + 1, n.Pos.Y + float32(pinStartHeight)}
	}
}
