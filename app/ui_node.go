package app

import (
	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func getPinRect(pos rl.Vector2, right bool) rl.Rectangle {
	var x float32
	if right {
		x = pos.X
	} else {
		x = pos.X - pinClickSize
	}

	return rl.Rectangle{
		x,
		pos.Y - pinClickSize/2,
		pinClickSize,
		pinClickSize,
	}
}

func drawPin(pos rl.Vector2, jut float32, right bool, color rl.Color) {
	var startAngle, endAngle float32
	var jutRec rl.Rectangle

	if right {
		startAngle = 0
		endAngle = 180
		jutRec = rl.Rectangle{pos.X, pos.Y - pinRadius, jut, pinRadius * 2}
		pos.X += jut
	} else {
		startAngle = 180
		endAngle = 360
		jutRec = rl.Rectangle{pos.X - jut, pos.Y - pinRadius, jut, pinRadius * 2}
		pos.X -= jut
	}

	rl.DrawRectangleRec(jutRec, color)
	rl.DrawCircleSector(pos, pinRadius, startAngle, endAngle, 36, color)
}

func drawNode(n *Node) {
	LoadThemeForNode(n)

	nodeRect := n.Rect()
	bgRect := nodeRect
	if n.Snapped {
		const snappedOverlap = 20
		bgRect.Y -= snappedOverlap
		bgRect.Height += snappedOverlap
	}

	const stackBackJut = 5
	doStackBack := !n.HasChildren && n.Snapped

	shadowRect := bgRect
	shadowRect.Y += 4
	if doStackBack {
		shadowRect.Width += stackBackJut
	}

	// draw drop shadow
	rl.DrawRectangleRounded(shadowRect, RoundnessPx(shadowRect, 10), 6, rl.NewColor(0, 0, 0, 50))

	// draw a stack behinder thinger
	if doStackBack {
		root := SnapRoot(n)
		stackBackRect := rl.Rectangle{
			X:      n.Pos.X,
			Y:      root.Pos.Y,
			Width:  n.Size.X + stackBackJut,
			Height: n.Pos.Y + n.Size.Y - root.Pos.Y,
		}
		rl.DrawRectangleRounded(stackBackRect, RoundnessPx(stackBackRect, 10), 6, pinColor)

		shadowRect.Width += stackBackJut
	}

	_, isDraggingNode := dragThing.(*Node)
	bgColor := n.Color
	if isDraggingNode && rl.CheckCollisionPointRec(raygui.GetMousePositionWorld(), n.SnapTargetRect) {
		bgColor = Tint(bgColor, 0.3)
	}
	rl.DrawRectangleRounded(bgRect, RoundnessPx(bgRect, 10), 6, bgColor) // main node

	titleHeight := float32(32)

	titleBarRect := rl.Rectangle{nodeRect.X, nodeRect.Y, nodeRect.Width - 24, titleHeight}
	previewRect := rl.Rectangle{nodeRect.X + nodeRect.Width - 24, nodeRect.Y, 24, titleHeight}

	drawBasicText(n.Title, nodeRect.X+6, nodeRect.Y+3, titleHeight, Shade(n.Color, 0.4))
	drawBasicText("P", previewRect.X+3, previewRect.Y+5, 28, Shade(n.Color, 0.4))

	pinHoverColor := Tint(pinColor, 0.5)

	for i, pinPos := range n.InputPinPos {
		if n.Snapped && i == 0 {
			continue
		}

		isHoverPin := rl.CheckCollisionPointRec(raygui.GetMousePositionWorld(), getPinRect(pinPos, false))

		pinColor := pinColor
		if isHoverPin {
			pinColor = pinHoverColor
		}
		drawPin(pinPos, pinJut, false, pinColor)

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
		var jut float32 = pinJut
		if n.Snapped {
			jut += stackBackJut
		}

		isHoverPin := rl.CheckCollisionPointRec(raygui.GetMousePositionWorld(), getPinRect(n.OutputPinPos, true))
		pinColor := pinColor
		if isHoverPin {
			pinColor = pinHoverColor
		}

		drawPin(n.OutputPinPos, jut, true, pinColor)
		if isHoverPin && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			tryDragNewWire(n)
		}
	}

	titleHover := rl.CheckCollisionPointRec(raygui.GetMousePositionWorld(), titleBarRect)
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
		sql := n.GenerateSql()
		currentSQL = sql
		resultsOpen = true
		setLatestResult(doQuery(sql))
	}

	n.DoUI()
}
