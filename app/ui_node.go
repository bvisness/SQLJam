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

	isHoverOutputPin := rl.CheckCollisionPointRec(raygui.GetMousePositionWorld(), getPinRect(n.OutputPinPos, true))
	pinHoverColor := Tint(pinColor, 0.5)

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

		color := pinColor
		if isHoverOutputPin {
			color = pinHoverColor
		}
		rl.DrawRectangleRounded(stackBackRect, RoundnessPx(stackBackRect, 10), 6, color)

		shadowRect.Width += stackBackJut
	}

	bgColor := n.Color
	if canShowSnapHighlight() && rl.CheckCollisionPointRec(raygui.GetMousePositionWorld(), n.SnapTargetRect) {
		bgColor = Tint(bgColor, 0.3)
	}
	rl.DrawRectangleRounded(bgRect, RoundnessPx(bgRect, 10), 6, bgColor) // main node

	titleHeight := float32(32)

	titleBarRect := rl.Rectangle{nodeRect.X, nodeRect.Y, nodeRect.Width, titleHeight}

	drawBasicText(n.Title, nodeRect.X+6, nodeRect.Y+3, titleHeight, Brightness(n.Color, 0.4))

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

		if source, ok := didDropWire(); isHoverPin && ok {
			n.Inputs[i] = source
			MarkInspectorDirty(n)
		} else if n.Inputs[i] != nil {
			if tryDragNewWire(n.Inputs[i], getPinRect(n.InputPinPos[i], false)) {
				n.Inputs[i] = nil
			}
		}
	}
	if !n.HasChildren {
		var jut float32 = pinJut
		if n.Snapped {
			jut += stackBackJut
		}

		pinColor := pinColor
		if isHoverOutputPin {
			pinColor = pinHoverColor
		}

		drawPin(n.OutputPinPos, jut, true, pinColor)
		tryDragNewWire(n, getPinRect(n.OutputPinPos, true))
	}

	if tryStartDrag(n, titleBarRect, n.Pos) {
		n.Sort = nodeSortTop()
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
	} else {
		titlebarHover := rl.CheckCollisionPointRec(raygui.GetMousePositionWorld(), titleBarRect)
		if titlebarHover && rl.IsMouseButtonReleased(rl.MouseLeftButton) && !dragging {
			MarkInspectorDirty(n)
		}
	}

	n.DoUI()
}

func canShowSnapHighlight() bool {
	n, isDraggingNode := dragThing.(*Node)
	if !isDraggingNode {
		return false
	}
	return n.CanSnap
}
