package app

import (
	"github.com/bvisness/SQLJam/node"
	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const orderDirectionWidth = 40

func doOrderUpdate(n *node.Node, o *node.Order) {
	uiHeight := 0
	for range o.Cols {
		uiHeight += UIFieldHeight
		uiHeight += UIFieldSpacing
	}
	uiHeight += UIFieldHeight // for buttons

	n.UISize = rl.Vector2{300, float32(uiHeight)}

	opts := columnNameDropdownOpts(n.Inputs[0])
	for _, col := range o.Cols {
		col.ColDropdown.SetOptions(opts...)
	}
}

func doOrderUI(n *node.Node, o *node.Order) {
	openDropdown, isOpen := raygui.GetOpenDropdown(o.ColDropdowns())
	if isOpen {
		raygui.Disable()
		defer raygui.Enable()
	}

	// Render bottom to top to avoid overlap issues with dropdowns

	fieldY := n.UIRect.Y + n.UIRect.Height - UIFieldHeight
	if raygui.Button(rl.Rectangle{
		n.UIRect.X,
		fieldY,
		n.UIRect.Width/2 - UIFieldSpacing/2,
		UIFieldHeight,
	}, "+") {
		o.Cols = append(o.Cols, &node.OrderColumn{})
	}
	if raygui.Button(rl.Rectangle{
		n.UIRect.X + n.UIRect.Width/2 + UIFieldSpacing/2,
		fieldY,
		n.UIRect.Width/2 - UIFieldSpacing/2,
		UIFieldHeight,
	}, "-") {
		if len(o.Cols) > 1 {
			o.Cols = o.Cols[:len(o.Cols)-1]
		}
	}

	for i := len(o.Cols) - 1; i >= 0; i-- {
		func() {
			fieldY -= UIFieldSpacing + UIFieldHeight

			col := o.Cols[i]
			dropdown := &col.ColDropdown

			if openDropdown == &col.ColDropdown {
				raygui.Enable()
				defer raygui.Disable()
			}

			colName := dropdown.Do(rl.Rectangle{
				n.UIRect.X,
				fieldY,
				n.UIRect.Width - orderDirectionWidth - UIFieldSpacing,
				UIFieldHeight,
			})
			col.Col, _ = colName.(string)

			directionStr := "A-Z"
			if col.Descending {
				directionStr = "Z-A"
			}
			col.Descending = raygui.Toggle(rl.Rectangle{
				n.UIRect.X + n.UIRect.Width - orderDirectionWidth,
				fieldY,
				orderDirectionWidth,
				UIFieldHeight,
			}, directionStr, col.Descending)
		}()
	}
}
