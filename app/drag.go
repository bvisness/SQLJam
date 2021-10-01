package app

import (
	"fmt"

	"github.com/bvisness/SQLJam/node"
	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var dragging bool
var dragCanceled bool
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
		dragCanceled = true
		dragKey = ""
		dragMouseStart = rl.Vector2{}
		dragObjStart = rl.Vector2{}
	}
}

func tryStartDrag(key interface{}, objStart rl.Vector2) bool {
	if dragging {
		return false
	}

	dragging = true
	dragCanceled = false
	dragKey = getDragKey(key)
	dragMouseStart = raygui.GetMousePositionWorld()
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

var wireDragOutputNode *node.Node

func tryDragNewWire(outputNode *node.Node) {
	if outputNode == nil {
		return
	}
	if tryStartDrag(wireDragKey, rl.Vector2{}) {
		wireDragOutputNode = outputNode
	}
}

func draggingWire() bool {
	return dragging && dragKey == wireDragKey
}

func didDropWire() (*node.Node, bool) {
	matchesKey, done, canceled := dragState(wireDragKey)
	if !matchesKey {
		return nil, false
	}

	if done && !canceled {
		return wireDragOutputNode, true
	}

	return nil, false
}
