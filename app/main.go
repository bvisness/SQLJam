package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	_ "github.com/mattn/go-sqlite3"
)

const screenWidth = 1920
const screenHeight = 1080

var nodes []*Node
var font rl.Font

func Main() {
	rl.InitWindow(screenWidth, screenHeight, "SQL Jam")
	defer rl.CloseWindow()

	// much fps or not you decide
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())))

	font = rl.LoadFont("JetBrainsMono-Regular.ttf")
	//rl.GenTextureMipmaps(&font.Texture) // kinda muddy? need second opinion
	rl.SetTextureFilter(font.Texture, rl.FilterBilinear) // FILTER_TRILINEAR requires generated mipmaps

	close := openDB()
	defer close()

	// main frame loop
	rl.SetExitKey(0)
	for !rl.WindowShouldClose() {
		doFrame()
	}
}

var latestResult *queryResult

var cam = rl.Camera2D{
	Offset: rl.Vector2{screenWidth / 2, screenHeight / 2},
	Target: rl.Vector2{screenWidth / 2, screenHeight / 2},
	Zoom:   1,
}
var panMouseStart rl.Vector2
var panCamStart rl.Vector2

const pinSize = 12
const wireThickness = 3

func drawBasicText(text string, x float32, y float32, size float32, color rl.Color) {
	rl.DrawTextEx(font, text, rl.Vector2{X: x, Y: y}, size, 2, color)
}

func doFrame() {
	raygui.Set2DCamera(nil)

	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(rl.RayWhite)

	updateDrag()

	// Pan/zoom camera
	{
		const minZoom = 0.2
		const maxZoom = 4

		zoomBefore := cam.Zoom
		zoomFactor := float32(rl.GetMouseWheelMove()) / 10
		cam.Zoom *= 1 + zoomFactor
		if cam.Zoom < minZoom {
			cam.Zoom = minZoom
		}
		if cam.Zoom > maxZoom {
			cam.Zoom = maxZoom
		}
		zoomAfter := cam.Zoom
		actualZoomFactor := zoomAfter/zoomBefore - 1

		mouseWorld := rl.GetScreenToWorld2D(raygui.GetMousePositionWorld(), cam)
		cam2mouse := rl.Vector2Subtract(mouseWorld, cam.Target)
		cam.Target = rl.Vector2Add(cam.Target, rl.Vector2Multiply(cam2mouse, rl.Vector2{actualZoomFactor, actualZoomFactor}))

		// TODO: Find a nice way of returning us to exactly 100% zoom.
		// But also supporting smooth trackpad zoom...?

		if rl.IsMouseButtonDown(rl.MouseMiddleButton) {
			mousePos := raygui.GetMousePositionWorld()
			if rl.IsMouseButtonPressed(rl.MouseMiddleButton) {
				panMouseStart = mousePos
				panCamStart = cam.Target
			}

			mouseDelta := rl.Vector2DivideV(rl.Vector2Subtract(mousePos, panMouseStart), rl.Vector2{cam.Zoom, cam.Zoom})
			cam.Target = rl.Vector2Subtract(panCamStart, mouseDelta) // camera moves opposite of mouse
		}
	}

	// update nodes
	for _, n := range nodes {
		n.UISize = rl.Vector2{}
		n.InputPinHeights = nil
		n.Update()
	}

	doLayout()

	raygui.Set2DCamera(&cam)
	rl.BeginMode2D(cam)
	{
		sort.SliceStable(nodes, func(i, j int) bool {
			/*
				A node should be less than another in the draw order if it
				should be drawn first. Thus, a stacked node should be "less
				than" its parent, but parents should be sorted according to
				their Sort values, and stacked nodes along with them.
			*/

			a := nodes[i]
			b := nodes[j]

			// Check if a is a stacked child of b
			if isSnappedUnder(a, b) {
				return true
			}

			// They're not stacked, but find their parents and sort by that.
			if SnapRoot(a).Sort < SnapRoot(b).Sort {
				return true
			}

			return false
		})

		// draw lines
		for _, n := range nodes {
			if n.Snapped {
				continue
			}
			for i, input := range n.Inputs {
				if input == nil {
					continue
				}
				rl.DrawLineBezier(input.OutputPinPos, n.InputPinPos[i], wireThickness, rl.Black)
			}
		}

		if draggingWire() {
			rl.DrawLineBezier(wireDragOutputNode.OutputPinPos, raygui.GetMousePositionWorld(), wireThickness, rl.Black)
		}

		getPinRect := func(pos rl.Vector2, right bool) rl.Rectangle {
			var x float32
			if right {
				x = pos.X
			} else {
				x = pos.X - pinSize
			}

			return rl.Rectangle{
				x,
				pos.Y - pinSize/2,
				pinSize,
				pinSize,
			}
		}

		// draw nodes
		for _, n := range nodes {
			nodeRect := n.Rect()
			bgRect := nodeRect
			if n.Snapped {
				const snappedOverlap = 20
				bgRect.Y -= snappedOverlap
				bgRect.Height += snappedOverlap
			}

			rl.DrawRectangleRounded(bgRect, RoundnessPx(bgRect, 10), 6, n.Color)
			//rl.DrawRectangleRoundedLines(bgRect, 0.16, 6, 2, rl.Black)

			titleBarRect := rl.Rectangle{nodeRect.X, nodeRect.Y, nodeRect.Width - 24, 24}
			previewRect := rl.Rectangle{nodeRect.X + nodeRect.Width - 24, nodeRect.Y, 24, 24}

			drawBasicText(n.Title, nodeRect.X+6, nodeRect.Y+4, 24, rl.Black)
			drawBasicText("P", previewRect.X+4, previewRect.Y+10, 14, rl.Black)

			for i, pinPos := range n.InputPinPos {
				if n.Snapped && i == 0 {
					continue
				}

				isHoverPin := rl.CheckCollisionPointRec(raygui.GetMousePositionWorld(), getPinRect(pinPos, false))

				pinColor := rl.Black
				if isHoverPin {
					pinColor = rl.Blue
				}
				rl.DrawRectangleRec(getPinRect(pinPos, false), pinColor)

				if isHoverPin {
					if source, ok := didDropWire(); ok {
						n.Inputs[i] = source
					} else if rl.IsMouseButtonPressed(rl.MouseLeftButton) && n.Inputs[i] != nil {
						tryDragNewWire(n.Inputs[i])
						n.Inputs[i] = nil
					}
				}
			}
			if !n.HasChildren {
				rl.DrawRectangleRec(getPinRect(n.OutputPinPos, true), rl.Black)
				if rl.CheckCollisionPointRec(raygui.GetMousePositionWorld(), getPinRect(n.OutputPinPos, true)) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
					tryDragNewWire(n)
				}
			}

			titleHover := rl.CheckCollisionPointRec(raygui.GetMousePositionWorld(), titleBarRect)
			if titleHover {
				toDraw := n.GenerateSql()
				charSize := 25
				lineHeight := charSize + 11
				numLines := strings.Count(toDraw, "\n") + 1
				drawBasicText(n.GenerateSql(), titleBarRect.X, SnapRoot(n).Pos.Y-float32(numLines*lineHeight), float32(charSize), rl.Black)
			}
			if titleHover && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				if tryStartDrag(n, n.Pos) {
					n.Sort = nodeSortTop()
				}
			}

			if draggingThis, done, canceled := dragState(n); draggingThis {
				n.Snapped = false
				n.Pos = dragNewPosition()
				if done {
					if canceled {
						n.Pos = dragObjStart
					} else {
						trySnapNode(n)
					}
				}
			}

			previewHover := rl.CheckCollisionPointRec(raygui.GetMousePositionWorld(), previewRect)
			if previewHover && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				latestResult = doQuery(n.GenerateSql())
			}

			n.DoUI()
		}
	}
	rl.EndMode2D()
	raygui.Set2DCamera(nil)

	drawToolbar()

	drawLatestResults()
}

var blerp rl.Vector2

func drawToolbar() {
	toolbarWidth := int32(rl.GetScreenWidth())
	toolbarHeight := int32(64)
	rl.DrawRectangle(0, 0, toolbarWidth, toolbarHeight, rl.ColorAlpha(rl.Black, 0.5))
	rl.DrawLine(0, toolbarHeight, toolbarWidth, toolbarHeight, rl.Black)
	rl.DrawLineEx(
		rl.Vector2{Y: float32(toolbarHeight)},
		rl.Vector2{X: float32(toolbarWidth), Y: float32(toolbarHeight)},
		5,
		rl.Black,
	)

	buttHeight := 40 // thicc

	if raygui.Button(rl.Rectangle{
		X:      20,
		Y:      float32(toolbarHeight/2) - float32(buttHeight/2),
		Width:  100,
		Height: float32(buttHeight),
	}, "Add Table") {
		table := NewTable()
		table.Pos = rl.Vector2{400, 400}
		nodes = append(nodes, table)
	}

	if raygui.Button(rl.Rectangle{
		X:      140,
		Y:      float32(toolbarHeight/2) - float32(buttHeight/2),
		Width:  100,
		Height: float32(buttHeight),
	}, "Add Filter") {
		filter := NewFilter()
		filter.Pos = rl.Vector2{400, 400}
		nodes = append(nodes, filter)
	}

	if raygui.Button(rl.Rectangle{
		X:      260,
		Y:      float32(toolbarHeight/2) - float32(buttHeight/2),
		Width:  180,
		Height: float32(buttHeight),
	}, "Add Pick Columns") {
		pc := NewPickColumns()
		pc.Pos = rl.Vector2{400, 400}
		nodes = append(nodes, pc)
	}

	if raygui.Button(rl.Rectangle{
		X:      460,
		Y:      float32(toolbarHeight/2) - float32(buttHeight/2),
		Width:  180,
		Height: float32(buttHeight),
	}, "Add Combine Rows") {
		cr := NewCombineRows(Union)
		cr.Pos = rl.Vector2{400, 400}
		nodes = append(nodes, cr)
	}

	if raygui.Button(rl.Rectangle{
		X:      660,
		Y:      float32(toolbarHeight/2) - float32(buttHeight/2),
		Width:  100,
		Height: float32(buttHeight),
	}, "Add Order") {
		pc := NewOrder()
		pc.Pos = rl.Vector2{400, 400}
		nodes = append(nodes, pc)
	}

	if raygui.Button(rl.Rectangle{
		X:      780,
		Y:      float32(toolbarHeight/2) - float32(buttHeight/2),
		Width:  100,
		Height: float32(buttHeight),
	}, "Add Join") {
		pc := NewJoin()
		pc.Pos = rl.Vector2{400, 400}
		nodes = append(nodes, pc)
	}

	if raygui.Button(rl.Rectangle{
		X:      900,
		Y:      float32(toolbarHeight/2) - float32(buttHeight/2),
		Width:  100,
		Height: float32(buttHeight),
	}, "Add Aggregate") {
		pc := NewAggregate()
		pc.Pos = rl.Vector2{400, 400}
		nodes = append(nodes, pc)
	}
}

const latestResultsHeight = 500

var latestResultsPanel raygui.ScrollPanelEx

func drawLatestResults() {
	// display query results (temporary)
	if latestResult != nil {
		const charWidth = 11
		const charHeight = 20
		const cellPadding = 5
		const rowHeight = cellPadding + charHeight + cellPadding

		var rows [][]string
		for i := -1; i < len(latestResult.Rows); i++ {
			if i < 0 {
				// headers
				rows = append(rows, latestResult.Columns)
			} else {
				row := latestResult.Rows[i]
				valStrings := make([]string, len(row))
				for i, v := range row {
					valStrings[i] = fmt.Sprintf("%v", v)
				}
				rows = append(rows, valStrings)
			}
		}

		colWidths := make([]int, len(rows[0]))
		for r := 0; r < len(rows); r++ {
			for c := 0; c < len(rows[0]); c++ {
				thisColWidth := cellPadding + len(rows[r][c])*charWidth + cellPadding
				if colWidths[c] < thisColWidth {
					colWidths[c] = thisColWidth
				}
			}
		}

		totalWidth := 0
		for _, w := range colWidths {
			totalWidth += w
		}

		panelBounds := rl.Rectangle{0, screenHeight - latestResultsHeight, screenWidth, latestResultsHeight}
		panelContents := rl.Rectangle{0, 0, float32(totalWidth), float32(len(rows) * rowHeight)}
		latestResultsPanel.Do(panelBounds, panelContents, func(start rl.Vector2, view rl.Rectangle) {
			cellPos := start
			for _, row := range rows {
				cellPos.X = 0
				for i, cell := range row {
					drawBasicText(cell, cellPos.X+cellPadding, cellPos.Y+cellPadding, charHeight, rl.Black)
					cellPos.X += float32(colWidths[i])
				}
				rl.DrawLine(0, int32(cellPos.Y+rowHeight), int32(view.X+view.Width), int32(cellPos.Y+rowHeight), rl.LightGray)
				cellPos.Y += rowHeight
			}

			gridX := 0
			for _, width := range colWidths {
				gridX += width
				rl.DrawLine(int32(gridX), int32(view.Y), int32(gridX), int32(view.Y+view.Height), rl.LightGray)
			}
		})
	}
}

func RoundnessPx(rect rl.Rectangle, radiusPx float32) float32 {
	minDimension := rect.Width
	if rect.Height < minDimension {
		minDimension = rect.Height
	}
	if minDimension == 0 {
		return 0
	}
	return radiusPx / minDimension
}

func doLayout() {
	/*
		Layout algo is as follows:

		- Calculate heights, widths, and input pins of all unsnapped nodes
		- Calculate heights, widths, and input pins of all snapped nodes
		- Do a pass across all nodes making them wider if necessary (yay snapping!)
		- Calculate output pins and final collisions of all nodes
	*/

	const titleBarHeight = 24
	const uiPadding = 10
	const snapRectHeight = 30
	const pinStartHeight = titleBarHeight + uiPadding

	const pinDefaultSpacing = 36 // used if the node does not specify pin heights in update

	basicLayout := func(n *Node) {
		n.Size = rl.Vector2{
			float32(n.UISize.X + 2*uiPadding),
			float32(titleBarHeight + uiPadding + int(n.UISize.Y) + uiPadding),
		}

		// use default input pin positions if not provided in update
		if len(n.InputPinHeights) < len(n.Inputs) {
			n.InputPinHeights = make([]int, len(n.Inputs))
			for i := range n.Inputs {
				n.InputPinHeights[i] = i * pinDefaultSpacing
			}
		}

		// init InputPinPos if necessary
		if len(n.InputPinPos) != len(n.Inputs) {
			n.InputPinPos = make([]rl.Vector2, len(n.Inputs))
		}

		for i := range n.Inputs {
			if n.Snapped && i == 0 {
				continue
			}
			n.InputPinPos[i] = rl.Vector2{n.Pos.X, n.Pos.Y + pinStartHeight + float32(n.InputPinHeights[i]) + pinSize/2}
		}
	}

	// sort nodes to ensure processing order
	sort.SliceStable(nodes, func(i, j int) bool {
		/*
			Here a node should be "less than" another if it should have its
			layout computed first. So a parent node should be "less than"
			its child.
		*/

		a := nodes[i]
		b := nodes[j]

		if isSnappedUnder(b, a) {
			return true
		}

		return false
	})

	// global setup
	for _, n := range nodes {
		n.HasChildren = false
	}

	// unsnapped
	for _, n := range nodes {
		if !n.Snapped {
			basicLayout(n)
		}
	}

	// snapped
	for _, n := range nodes {
		if n.Snapped {
			basicLayout(n)
			parent := n.Inputs[0]
			n.Pos = rl.Vector2{parent.Pos.X, parent.Pos.Y + parent.Size.Y}
		}
	}

	// fix sizing
	for _, n := range nodes {
		maxWidth := n.Size.X

		current := n
		for {
			if current.Size.X > maxWidth {
				maxWidth = current.Size.X
			}
			n.Size.X = maxWidth
			current.Size.X = maxWidth

			if current.Snapped && len(current.Inputs) > 0 {
				current = current.Inputs[0]
				continue
			}
			break
		}
	}

	// output pin positions (unsnapped)
	for _, n := range nodes {
		if !n.Snapped {
			n.OutputPinPos = rl.Vector2{n.Pos.X + n.Size.X, n.Pos.Y + float32(pinStartHeight) + pinSize/2}
		}
	}

	// final calculations
	for _, n := range nodes {
		if n.Snapped {
			current := n
			for {
				n.OutputPinPos = current.OutputPinPos
				if current != n {
					current.HasChildren = true
				}
				if current.Snapped && len(current.Inputs) > 0 {
					current = current.Inputs[0]
					continue
				}
				break
			}
		}
		n.UIRect = rl.Rectangle{
			n.Pos.X + uiPadding,
			n.Pos.Y + titleBarHeight + uiPadding,
			n.Size.X - 2*uiPadding,
			n.Size.Y - titleBarHeight - 2*uiPadding,
		}
		n.SnapTargetRect = rl.Rectangle{n.Pos.X, n.Pos.Y + n.Size.Y - snapRectHeight, n.Size.X, snapRectHeight}
	}
}

func trySnapNode(n *Node) {
	if !n.CanSnap {
		return
	}

	for _, other := range nodes {
		if n == other {
			continue
		}

		if rl.CheckCollisionPointRec(raygui.GetMousePositionWorld(), other.SnapTargetRect) {
			// See snapping.png.
			// INVARIANT: Nodes must always be pointing at the leaves of stacks.

			oldLeaf := SnapLeaf(other)
			newRoot := SnapRoot(other)
			newLeaf := SnapLeaf(n)

			// make nodes pointing at oldLeaf point to newLeaf
			for _, other := range nodes {
				for i := range other.Inputs {
					if other.Inputs[i] == oldLeaf {
						other.Inputs[i] = newLeaf
					}
				}
			}

			// break cycles - if new root points at new leaf, set it to nil
			for i := range newRoot.Inputs {
				if newRoot.Inputs[i] == newLeaf {
					newRoot.Inputs[i] = nil
				}
			}

			// Snap! ^_^
			n.Inputs[0] = SnapLeaf(other)
			n.Snapped = true

			break
		}
	}
}

var topSortValue = 0

func nodeSortTop() int {
	topSortValue++
	return topSortValue
}

func SnapRoot(n *Node) *Node {
	root, _ := SnapRootAndDistance(n)
	return root
}

func SnapRootAndDistance(n *Node) (*Node, int) {
	distance := 0
	root := n
	for {
		if root.Snapped && len(root.Inputs) > 0 && root.Inputs[0] != nil {
			root = root.Inputs[0]
			distance++
			continue
		}
		break
	}

	return root, distance
}

// this is not efficient, who cares
func SnapLeaf(n *Node) *Node {
	root := SnapRoot(n)

	// The leaf is the node farthest from the snap root
	leaf := n
	maxDistToRoot := 0
	for _, other := range nodes {
		otherRoot, distance := SnapRootAndDistance(other)
		if otherRoot == root && distance > maxDistToRoot {
			maxDistToRoot = distance
			leaf = other
		}
	}

	return leaf
}

func isSnappedUnder(a, b *Node) bool {
	current := a
	for current != nil {
		if current == b {
			return true
		}

		if current.Snapped && len(current.Inputs) > 0 {
			current = current.Inputs[0]
		} else {
			break
		}
	}

	return false
}
