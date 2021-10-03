package app

import (
	"sort"
	"strings"

	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	_ "github.com/mattn/go-sqlite3"
)

var screenWidth = float32(1920)
var screenHeight = float32(1080)

const resultsMaxHeight = 400
const currentSQLWidth = 640

const dividerThickness = 4
const pinRadius = 8
const pinClickSize = 30 // size of clickable/droppable area, square
const pinJut = 2
const wireThickness = 3

var nodes []*Node
var font rl.Font

var dark = true
var mainColorLight = rl.RayWhite
var mainColorDark = rl.NewColor(46, 34, 47, 255)
var pinColor = rl.NewColor(123, 107, 127, 255)

var PaneFontColor = rl.NewColor(253, 203, 176, 255)
var PaneLineColor = rl.NewColor(139, 94, 91, 255)

func MainColor() rl.Color {
	if dark {
		return mainColorDark
	} else {
		return mainColorLight
	}
}

// LoadStyleMain Per frame custom style settings
func LoadStyleMain() {
	raygui.SetFont(font)
	raygui.SetStyle(raygui.Default, raygui.BorderWidthProp, 2)

	raygui.SetStyle(raygui.ScrollBarControl, raygui.ArrowsVisible, 1)
	raygui.SetStyle(raygui.DropdownBoxControl, raygui.DropdownItemsPadding, 0)
	raygui.SetStyle(raygui.TextBoxControl, raygui.BorderWidthProp, 2)

	raygui.SetStyle(raygui.Default, raygui.BaseColorNormalProp, 0x3E3546FF)
	raygui.SetStyle(raygui.Default, raygui.BorderColorNormalProp, 0x3E3546FF)
	raygui.SetStyle(raygui.Default, raygui.TextColorNormalProp, 0x625565FF)

	raygui.SetStyle(raygui.Default, raygui.TextSizeProp, 28)

	raygui.SetStyle(raygui.Default, raygui.BaseColorFocusedProp, 0x625565FF)
	raygui.SetStyle(raygui.Default, raygui.BorderColorFocusedProp, 0x3E3546FF)
	raygui.SetStyle(raygui.Default, raygui.TextColorFocusedProp, ToHexNum(MainColor()))

	SetStyleColor(raygui.ScrollBarControl, raygui.BorderColorPressedProp, Tint(MainColor(), 0.3))

	SetStyleColor(raygui.SliderControl, raygui.BaseColorNormalProp, Tint(MainColor(), 0.3))

	raygui.SetStyle(raygui.ListViewControl, raygui.ScrollBarWidth, 20)

}

func MarkInspectorDirtyCurrent() {
	MarkInspectorDirty(selectedNode)
}

func MarkInspectorDirty(n *Node) {
	selectedNode = n
	inspectorDirty = true
}

func UpdateInspectorIfNeeded() {
	if inspectorDirty && selectedNode != nil {
		sql := selectedNode.GenerateSql()
		currentSQL = sql
		resultsOpen = true
		latestResults.Update(doQuery(sql))
	}
	inspectorDirty = false
}

func Main() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(int32(screenWidth), int32(screenHeight), "SQL Jam")
	defer rl.CloseWindow()

	monWidth := float32(rl.GetMonitorWidth(rl.GetCurrentMonitor()))
	monHeight := float32(rl.GetMonitorHeight(rl.GetCurrentMonitor()))
	rl.SetWindowSize(int(monWidth*0.8), int(monHeight*0.8))
	rl.SetWindowPosition(int(monWidth*0.1), int(monHeight*0.1))

	// much fps or not you decide
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())))

	font = rl.LoadFont("JetBrainsMono-Regular.ttf")
	//rl.GenTextureMipmaps(&font.Texture) // kinda muddy? need second opinion
	rl.SetTextureFilter(font.Texture, rl.FilterBilinear) // FILTER_TRILINEAR requires generated mipmaps

	LoadStyleMain()

	close := openDB()
	defer close()

	// main frame loop
	rl.SetExitKey(0)
	for !rl.WindowShouldClose() {
		doFrame()
	}
}

var cam = rl.Camera2D{
	Offset: rl.Vector2{screenWidth / 2, screenHeight / 2},
	Target: rl.Vector2{screenWidth / 2, screenHeight / 2},
	Zoom:   1,
}
var panning bool
var panMouseStart rl.Vector2
var panCamStart rl.Vector2

var currentSQL string
var selectedNode *Node
var inspectorDirty bool

const minZoom = 0.25
const maxZoom = 4
const zoomSnapRadius = 0.15 // percent deviation from snap point - e.g. radius of 0.2 means snap point +/- 20%

var zoom float32 = 1
var zoomSnapPoints = []float32{0.25, 0.5, 1, 2, 3, 4}

var didCaptureScrollThisFrame bool

func doFrame() {
	screenWidth = float32(rl.GetScreenWidth())
	screenHeight = float32(rl.GetScreenHeight())

	raygui.Set2DCamera(nil)

	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(MainColor())

	didCaptureScrollThisFrame = false

	DoPane(rl.Rectangle{0, 0, screenWidth, screenHeight - resultsCurrentHeight}, func(p Pane) {
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
			updateDrag()

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
					rl.DrawLineBezier(input.OutputPinPos, n.InputPinPos[i], wireThickness, pinColor)
				}
			}

			if draggingWire() {
				rl.DrawLineBezier(wireDragOutputNode.OutputPinPos, raygui.GetMousePositionWorld(), wireThickness, pinColor)
			}

			// draw nodes
			for _, n := range nodes {
				drawNode(n)
			}

			// Reset to default style
			raygui.LoadStyleDefault()
			LoadStyleMain()
		}
		rl.EndMode2D()
		raygui.Set2DCamera(nil)

		UpdateInspectorIfNeeded()

		drawToolbar()

		// Pan/zoom camera
		{
			zoomBefore := cam.Zoom
			zoomFactor := float32(rl.GetMouseWheelMove()) / 10
			if !p.MouseInPane() || didCaptureScrollThisFrame {
				zoomFactor = 0
			}
			zoom = zoom * (1 + zoomFactor) // actual zoom does not snap
			if zoom < minZoom {
				zoom = minZoom
			}
			if zoom > maxZoom {
				zoom = maxZoom
			}
			snappedZoom := zoom
			for _, snapPoint := range zoomSnapPoints {
				if Abs(snappedZoom-snapPoint) <= snapPoint*zoomSnapRadius {
					snappedZoom = snapPoint
					break
				}
			}
			cam.Zoom = snappedZoom
			zoomAfter := cam.Zoom
			actualZoomFactor := zoomAfter/zoomBefore - 1

			mouseWorld := rl.GetScreenToWorld2D(raygui.GetMousePositionWorld(), cam)
			cam2mouse := rl.Vector2Subtract(mouseWorld, cam.Target)
			cam.Target = rl.Vector2Add(cam.Target, rl.Vector2Multiply(cam2mouse, rl.Vector2{actualZoomFactor, actualZoomFactor}))

			// TODO: Find a nice way of returning us to exactly 100% zoom.
			// But also supporting smooth trackpad zoom...?

			if rl.IsMouseButtonDown(rl.MouseRightButton) {
				if rl.IsMouseButtonPressed(rl.MouseRightButton) && p.MouseInPane() {
					panning = true
					panMouseStart = raygui.GetMousePositionWorld()
					panCamStart = cam.Target
				}
			} else {
				panning = false
			}

			if panning {
				mousePos := raygui.GetMousePositionWorld()
				mouseDelta := rl.Vector2DivideV(rl.Vector2Subtract(mousePos, panMouseStart), rl.Vector2{cam.Zoom, cam.Zoom})
				cam.Target = rl.Vector2Subtract(panCamStart, mouseDelta) // camera moves opposite of mouse
			}
		}
	})

	drawLatestResults()
	drawCurrentSQL()

}

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

	buttHeight := 40 // thicc

	initNewNode := func(n *Node, defaultSize rl.Vector2) {
		n.Pos = rl.Vector2Subtract(cam.Target, rl.Vector2DivideV(defaultSize, rl.Vector2{2, 2}))
		n.Sort = nodeSortTop()
	}

	if raygui.Button(rl.Rectangle{
		X:      20,
		Y:      float32(toolbarHeight/2) - float32(buttHeight/2),
		Width:  160,
		Height: float32(buttHeight),
	}, "Add Table") {
		n := NewTable()
		initNewNode(n, rl.Vector2{300, 100})
		nodes = append(nodes, n)
	}

	if raygui.Button(rl.Rectangle{
		X:      200,
		Y:      float32(toolbarHeight/2) - float32(buttHeight/2),
		Width:  180,
		Height: float32(buttHeight),
	}, "Add Filter") {
		n := NewFilter()
		initNewNode(n, rl.Vector2{400, 100})
		nodes = append(nodes, n)
	}

	if raygui.Button(rl.Rectangle{
		X:      400,
		Y:      float32(toolbarHeight/2) - float32(buttHeight/2),
		Width:  260,
		Height: float32(buttHeight),
	}, "Add Pick Columns") {
		n := NewPickColumns()
		initNewNode(n, rl.Vector2{450, 200})
		nodes = append(nodes, n)
	}

	if raygui.Button(rl.Rectangle{
		X:      680,
		Y:      float32(toolbarHeight/2) - float32(buttHeight/2),
		Width:  260,
		Height: float32(buttHeight),
	}, "Add Combine Rows") {
		n := NewCombineRows(Union)
		initNewNode(n, rl.Vector2{300, 150})
		nodes = append(nodes, n)
	}

	if raygui.Button(rl.Rectangle{
		X:      960,
		Y:      float32(toolbarHeight/2) - float32(buttHeight/2),
		Width:  160,
		Height: float32(buttHeight),
	}, "Add Order") {
		n := NewOrder()
		initNewNode(n, rl.Vector2{350, 150})
		nodes = append(nodes, n)
	}

	if raygui.Button(rl.Rectangle{
		X:      1140,
		Y:      float32(toolbarHeight/2) - float32(buttHeight/2),
		Width:  140,
		Height: float32(buttHeight),
	}, "Add Join") {
		n := NewJoin()
		initNewNode(n, rl.Vector2{600, 200})
		nodes = append(nodes, n)
	}

	if raygui.Button(rl.Rectangle{
		X:      1300,
		Y:      float32(toolbarHeight/2) - float32(buttHeight/2),
		Width:  220,
		Height: float32(buttHeight),
	}, "Add Aggregate") {
		n := NewAggregate()
		initNewNode(n, rl.Vector2{600, 200})
		nodes = append(nodes, n)
	}

	if raygui.Button(rl.Rectangle{
		X:      screenWidth - 220 - 20,
		Y:      float32(toolbarHeight/2) - float32(buttHeight/2),
		Width:  220,
		Height: float32(buttHeight),
	}, "Add Preview") {
		n := NewPreview()
		initNewNode(n, rl.Vector2{600, 400})
		nodes = append(nodes, n)
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

var topSortValue = 0

func nodeSortTop() int {
	topSortValue++
	return topSortValue
}

var currentSQLPanel raygui.ScrollPanelEx

func drawCurrentSQL() {
	var lineX float32 = screenWidth - currentSQLWidth + dividerThickness/2
	rl.DrawLineEx(
		rl.Vector2{lineX, screenHeight - resultsCurrentHeight},
		rl.Vector2{lineX, screenHeight},
		dividerThickness, rl.Black,
	)
	DoPane(rl.Rectangle{screenWidth - currentSQLWidth + dividerThickness, screenHeight - resultsCurrentHeight, currentSQLWidth - dividerThickness, resultsMaxHeight}, func(p Pane) {
		const headerHeight = 40
		const bottomButtonHeights = 40
		const padding = 6
		const fontSize = 20
		const lineHeight = 24

		rl.DrawRectangleRec(p.Bounds, MainColor())

		topBounds := p.Bounds
		topBounds.Height = headerHeight

		if selectedNode != nil {
			rl.DrawRectangleRec(topBounds, selectedNode.Color)
			centerOffset := topBounds.Width/2 - float32(raygui.GetTextWidth(selectedNode.Title))/2
			drawBasicText(selectedNode.Title, topBounds.X+centerOffset, topBounds.Y+5, 32, Brightness(selectedNode.Color, 0.45))
		}

		lines := strings.Split(currentSQL, "\n")

		var maxLineLength float32
		for _, line := range lines {
			lineWidth := measureBasicText(line, fontSize)
			if lineWidth.X > maxLineLength {
				maxLineLength = lineWidth.X
			}
		}

		scrollPanelBounds := p.Bounds
		scrollPanelBounds.Height = p.Bounds.Height - bottomButtonHeights - headerHeight
		scrollPanelBounds.Y += headerHeight
		currentSQLPanel.Do(
			scrollPanelBounds,
			rl.Rectangle{
				Width:  padding + maxLineLength + padding,
				Height: padding + float32(len(lines)*lineHeight) + padding,
			},
			func(scroll raygui.ScrollContext) {

				for i, line := range lines {
					drawBasicText(line, scroll.Start.X+padding, scroll.Start.Y+padding+float32(i)*lineHeight, fontSize, PaneFontColor)
				}
			},
		)

		if raygui.Button(rl.Rectangle{p.Bounds.X, p.Bounds.Y + p.Bounds.Height - bottomButtonHeights, p.Bounds.Width / 2, bottomButtonHeights}, "Copy Text") {
			rl.SetClipboardText(currentSQL)
		}
		if raygui.Button(rl.Rectangle{p.Bounds.X + p.Bounds.Width/2, p.Bounds.Y + p.Bounds.Height - bottomButtonHeights, p.Bounds.Width / 2, bottomButtonHeights}, "Delete Node") {
			for i, scanNode := range nodes {

				for k, input := range scanNode.Inputs {
					if input == selectedNode {
						if len(selectedNode.Inputs) == 1 {
							scanNode.Inputs[k] = selectedNode.Inputs[0]
						} else {
							scanNode.Inputs[k] = nil
						}
						if scanNode.Inputs[k] == nil {
							scanNode.Snapped = false
						}
					}
				}

				if scanNode == selectedNode {
					nodes = append(nodes[:i], nodes[i+1:]...)
				}

			}

			selectedNode = nil
			resultsOpen = false
		}
	})
}
