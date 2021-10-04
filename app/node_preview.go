package app

import (
	"fmt"

	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const previewMinWidth = 300
const previewMinHeight = 200

type Preview struct {
	Panel     QueryResultPanel
	Size      rl.Vector2 // this will be applied to UISize which will determine the node Size. Make sense???
	StartSize rl.Vector2
}

var _ NodeData = &Preview{}

func NewPreview() *Node {
	return &Node{
		Title:   "Preview",
		CanSnap: true,
		Color:   rl.Gray,
		Inputs:  make([]*Node, 1),
		Data: &Preview{
			Size: rl.Vector2{600, 400},
		},
	}
}

func (d *Preview) Update(n *Node) {
	if n.Schema == nil {
		d.Panel.Update(doQuery(n.GenerateSql(true)))
	}

	n.UISize = d.Size
}

func (d *Preview) DoUI(n *Node) {
	LoadStyleMain()

	d.Panel.Draw(n.UIRect)
	if rl.CheckCollisionPointRec(raygui.GetMousePositionWorld(), n.UIRect) {
		didCaptureScrollThisFrame = true
	}

	bottomRight := rl.Vector2{n.Pos.X + n.Size.X, n.Pos.Y + n.Size.Y}
	resizeRect := rl.Rectangle{bottomRight.X - 20, bottomRight.Y - 20, 20, 20}

	// resize handle
	rl.DrawLineV(
		rl.Vector2{bottomRight.X - 10, bottomRight.Y - 2},
		rl.Vector2{bottomRight.X - 2, bottomRight.Y - 10},
		Brightness(n.Color, 0.5),
	)
	rl.DrawLineV(
		rl.Vector2{bottomRight.X - 6, bottomRight.Y - 2},
		rl.Vector2{bottomRight.X - 2, bottomRight.Y - 6},
		Brightness(n.Color, 0.5),
	)

	resizeDragKey := fmt.Sprintf("resize: %p", d)
	if tryStartDrag(resizeDragKey, resizeRect, rl.Vector2{}) {
		d.StartSize = d.Size
	}

	if resizingThis, _, canceled := dragState(resizeDragKey); resizingThis {
		if canceled {
			d.Size = d.StartSize
		} else {
			newSize := rl.Vector2Add(d.StartSize, dragOffset())
			if newSize.X < previewMinWidth {
				newSize.X = previewMinWidth
			}
			if newSize.Y < previewMinHeight {
				newSize.Y = previewMinHeight
			}
			d.Size = newSize
		}
	}
}

func (d *Preview) Serialize() (string, bool) {
	return "", false
}
