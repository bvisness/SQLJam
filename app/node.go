package app

import rl "github.com/gen2brain/raylib-go/raylib"

type Node struct {
	Data   NodeData
	Inputs []*Node

	// Static properties, should only be set when created
	CanSnap bool // can snap to another node for its primary input

	// Fields to set in the node update pass. Used for layout.
	UISize          rl.Vector2
	InputPinHeights []int

	// UI data. Can be changed during the node UI pass.
	Pos     rl.Vector2
	Title   string
	Color   rl.Color
	Snapped bool
	Sort    int

	// calculated fields
	InputPinPos    []rl.Vector2
	OutputPinPos   rl.Vector2
	HasChildren    bool
	SnapTargetRect rl.Rectangle
	Size           rl.Vector2   // size of the entire node - calculated based on UISize
	UIRect         rl.Rectangle // the UI content area
}

type NodeData interface {
	Update(n *Node) // Set up data for later, set UI and layout stuff.
	DoUI(n *Node)   // Run all of this node's specific UI.
}

func (n *Node) Rect() rl.Rectangle {
	return rl.Rectangle{n.Pos.X, n.Pos.Y, n.Size.X, n.Size.Y}
}

func (n *Node) Update() {
	n.Data.Update(n)
}

func (n *Node) DoUI() {
	n.Data.DoUI(n)
}
