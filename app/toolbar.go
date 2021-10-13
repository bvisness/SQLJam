package app

import (
	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func drawToolbar() {
	toolbarWidth := int32(rl.GetScreenWidth())
	toolbarHeight := int32(64)
	rl.DrawRectangle(0, 0, toolbarWidth, toolbarHeight, rl.ColorAlpha(rl.Black, 0.5))
	rl.DrawLineEx(
		rl.Vector2{0, float32(toolbarHeight)},
		rl.Vector2{float32(toolbarWidth), float32(toolbarHeight)},
		dividerThickness,
		rl.Black,
	)

	const buttHeight = 40  // thicc
	const buttSpacing = 15 // extra thicc

	initNewNode := func(n *Node, defaultSize rl.Vector2) {
		n.Pos = rl.Vector2Subtract(cam.Target, rl.Vector2DivideV(defaultSize, rl.Vector2{2, 2}))
		n.Sort = nodeSortTop()
	}

	buttonRect := func(x, width float32) rl.Rectangle {
		return rl.Rectangle{
			X:      x,
			Y:      float32(toolbarHeight/2) - float32(buttHeight/2),
			Width:  width,
			Height: float32(buttHeight),
		}
	}

	doToolbarButton := func(text, description string, bounds rl.Rectangle, color rl.Color, makeNode func() *Node) (nextX float32) {
		LoadThemeForColor(color)
		if raygui.Button(bounds, text) {
			nodes = append(nodes, makeNode())
		}

		if rl.CheckCollisionPointRec(rl.GetMousePosition(), bounds) {
			drawBasicText(description, 20, float32(toolbarHeight)+20, 24, PaneFontColor)
		}

		return bounds.X + bounds.Width + buttSpacing
	}

	var nextX float32 = buttSpacing

	nextX = doToolbarButton(
		"Table", "Get the contents of a database table.",
		buttonRect(nextX, 100),
		TableColor,
		func() *Node {
			n := NewTable()
			initNewNode(n, rl.Vector2{300, 100})
			return n
		},
	)

	nextX = doToolbarButton(
		"Filter", "Pick only the rows matching the condition.",
		buttonRect(nextX, 120),
		FilterColor,
		func() *Node {
			n := NewFilter()
			initNewNode(n, rl.Vector2{400, 100})
			return n
		},
	)

	nextX = doToolbarButton(
		"Pick Columns", "Pick only the specified columns. Can also rename columns.",
		buttonRect(nextX, 200),
		PickColumnsColor,
		func() *Node {
			n := NewPickColumns()
			initNewNode(n, rl.Vector2{450, 200})
			return n
		},
	)

	nextX = doToolbarButton(
		"Sort", "Sort the table by one or more columns.",
		buttonRect(nextX, 100),
		SortColor,
		func() *Node {
			n := NewSort()
			initNewNode(n, rl.Vector2{350, 150})
			return n
		},
	)

	nextX = doToolbarButton(
		"Aggregate", "Aggregate all the rows of the table into a single result, optionally grouping by one or more columns.",
		buttonRect(nextX, 160),
		AggregateColor,
		func() *Node {
			n := NewAggregate()
			initNewNode(n, rl.Vector2{600, 200})
			return n
		},
	)

	nextX = doToolbarButton(
		"Join", "Combine multiple tables into one by pairing up rows that match a certain condition.",
		buttonRect(nextX, 120),
		JoinColor,
		func() *Node {
			n := NewJoin()
			initNewNode(n, rl.Vector2{600, 200})
			return n
		},
	)

	nextX = doToolbarButton(
		"Combine Rows", "Take two tables with the same schema and combine their rows using set operations.",
		buttonRect(nextX, 200),
		CombineRowsColor,
		func() *Node {
			n := NewCombineRows(Union)
			initNewNode(n, rl.Vector2{300, 150})
			return n
		},
	)

	doToolbarButton(
		"Preview", "View the results of a query as you work.",
		buttonRect(screenWidth-buttSpacing-160, 160),
		PreviewColor,
		func() *Node {
			n := NewPreview()
			initNewNode(n, rl.Vector2{600, 400})
			return n
		},
	)

	doToolbarButton(
		"Chart", "Plot and visualize data.",
		buttonRect(screenWidth-buttSpacing-160-buttSpacing-160, 160),
		ChartColor,
		func() *Node {
			n := NewChart()
			initNewNode(n, rl.Vector2{600, 400})
			return n
		},
	)

	LoadStyleMain()

	rl.DrawRectangle(0, 0, toolbarWidth, toolbarHeight, rl.ColorAlpha(rl.Black, 0.25))
}
