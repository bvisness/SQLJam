package app

import rl "github.com/gen2brain/raylib-go/raylib"

func DoPane(bounds rl.Rectangle, draw func(p Pane)) {
	rl.BeginScissorMode(int32(bounds.X), int32(bounds.Y), int32(bounds.Width), int32(bounds.Height))
	draw(Pane{
		Bounds: bounds,
	})
	rl.EndScissorMode()
}

type Pane struct {
	Bounds rl.Rectangle
}

func (p *Pane) MouseInPane() bool {
	return rl.CheckCollisionPointRec(rl.GetMousePosition(), p.Bounds)
}
