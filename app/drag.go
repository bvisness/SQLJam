package app

import (
	"fmt"

	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var dragging bool
var dragPending bool
var dragCanceled bool
var dragThing interface{}
var dragKey string
var dragMouseStart rl.Vector2
var dragObjStart rl.Vector2

func getDragKey(key interface{}) string {
	switch kt := key.(type) {
	case string:
		return kt
	default:
		return fmt.Sprintf("%p", key)
	}
}

// Call once per frame at the start of the frame.
func updateDrag() {
	if rl.IsKeyPressed(rl.KeyEscape) {
		dragging = false
		dragCanceled = true
	} else if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		dragging = false
	} else if rl.IsMouseButtonUp(rl.MouseLeftButton) {
		dragging = false
		dragPending = false
		dragCanceled = true
		dragThing = nil
		dragKey = ""
		dragMouseStart = rl.Vector2{}
		dragObjStart = rl.Vector2{}
	} else if rl.IsMouseButtonDown(rl.MouseLeftButton) {
		if !dragging && !dragPending {
			dragPending = true
			dragMouseStart = raygui.GetMousePositionWorld()
		}
	}
}

func tryStartDrag(thing interface{}, dragRegion rl.Rectangle, objStart rl.Vector2) bool {
	if thing == nil {
		panic("you must provide a thing to drag")
	}

	if dragging {
		// can't start a new drag while one is in progress
		return false
	}

	if !dragPending {
		// can't start a new drag with this item unless we have a pending one
		return false
	}

	if rl.Vector2Length(rl.Vector2Subtract(raygui.GetMousePositionWorld(), dragMouseStart)) < 3 {
		// haven't dragged far enough
		return false
	}

	if !rl.CheckCollisionPointRec(dragMouseStart, dragRegion) {
		// not dragging from the right place
		return false
	}

	dragging = true
	dragPending = false
	dragCanceled = false
	dragThing = thing
	dragKey = getDragKey(thing)
	dragObjStart = objStart

	return true
}

func dragOffset() rl.Vector2 {
	if !dragging && dragKey == "" {
		return rl.Vector2{}
	}
	return rl.Vector2Subtract(raygui.GetMousePositionWorld(), dragMouseStart)
}

func dragNewPosition() rl.Vector2 {
	return rl.Vector2Add(dragObjStart, dragOffset())
}

// Pass in an key and it will tell you the relevant drag state for that thing.
// matchesKey will be true if that object is the one currently being dragged.
// done will be true if the drag is complete this frame.
// canceled will be true if the drag is done but was canceled.
func dragState(key interface{}) (matchesKey bool, done bool, canceled bool) {
	matchesKey = true
	if key != nil {
		matchesKey = dragKey == getDragKey(key)
	}

	if !dragging && dragKey != "" && matchesKey {
		return matchesKey, true, dragCanceled
	} else {
		return matchesKey, false, false
	}
}

// Specific stuff for wire dragging

const wireDragKey = "NEW_WIRE"

var wireDragOutputNode *Node

func tryDragNewWire(outputNode *Node, dragRegion rl.Rectangle) bool {
	if outputNode == nil {
		return false
	}
	if tryStartDrag(wireDragKey, dragRegion, rl.Vector2{}) {
		wireDragOutputNode = outputNode
		return true
	}

	return false
}

func draggingWire() bool {
	return dragging && dragKey == wireDragKey
}

func didDropWire() (*Node, bool) {
	matchesKey, done, canceled := dragState(wireDragKey)
	if !matchesKey {
		return nil, false
	}

	if done && !canceled {
		return wireDragOutputNode, true
	}

	return nil, false
}
