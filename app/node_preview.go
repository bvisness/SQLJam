package app

import (
	"fmt"

	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const previewMinWidth = 300
const previewMinHeight = 100

var PreviewColor = rl.NewColor(155, 171, 178, 255)

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
		Color:   PreviewColor,
		Inputs:  make([]*Node, 1),
		Data: &Preview{
			Size: rl.Vector2{600, 400},
		},
	}
}

func (d *Preview) Update(n *Node) {
	if n.Schema == nil {
		d.Panel.Update(doQuery(n.GenerateSql(true)))
		n.Schema = getSchema(n)
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

	drawResizeHandle(bottomRight, n.Color)

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
	} else {
		d.Size = rl.Vector2{n.UIRect.Width, n.UIRect.Height}
	}
}

func (d *Preview) Serialize() (string, bool) {
	return "", false
}
