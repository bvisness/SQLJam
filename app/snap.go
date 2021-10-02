package app

import (
	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

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
