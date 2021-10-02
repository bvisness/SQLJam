package app

import (
	"sort"

	rl "github.com/gen2brain/raylib-go/raylib"
)

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
