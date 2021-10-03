package app

import rl "github.com/gen2brain/raylib-go/raylib"

type Preview struct {
	Panel QueryResultPanel
}

var _ NodeData = &Preview{}

func NewPreview() *Node {
	return &Node{
		Title:   "Preview",
		CanSnap: true,
		Color:   rl.Gray,
		Data:    &Preview{},
	}
}

func (d *Preview) Update(n *Node) {
	if n.Schema == nil {
		d.Panel.Update(doQuery(n.GenerateSql()))
	}

	n.UISize = rl.Vector2{600, 400}
}

func (d *Preview) DoUI(n *Node) {
	d.Panel.Draw(n.UIRect)
}

func (d *Preview) Serialize() (string, bool) {
	return "", false
}
