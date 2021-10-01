package app

import (
	"github.com/bvisness/SQLJam/node"
	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func doOrderUpdate(n *node.Node, o *node.Order) {
	uiHeight := 0
	for range o.Cols {
		uiHeight += UIFieldHeight
		uiHeight += UIFieldSpacing
	}
	uiHeight += UIFieldHeight // for buttons

	n.UISize = rl.Vector2{300, float32(uiHeight)}

	// This will obliterate existing selections on resize,
	// but this shouldn't happen anyway if we're resizing correctly.
	if len(o.Cols) != len(o.ColDropdowns) {
		o.ColDropdowns = make([]*raygui.DropdownEx, len(o.Cols))
	}

	opts := columnNameDropdownOpts(n.Inputs[0])
	for _, dropdown := range o.ColDropdowns {
		dropdown.SetOptions(opts...)
	}
}

func doOrderUI(n *node.Node, o *node.Order) {
	openDropdown, isOpen := raygui.GetOpenDropdown(o.ColDropdowns)
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
		o.Cols = append(o.Cols, node.OrderColumn{})
		o.ColDropdowns = append(o.ColDropdowns, &raygui.DropdownEx{})
	}
	if raygui.Button(rl.Rectangle{
		n.UIRect.X + n.UIRect.Width/2 + UIFieldSpacing/2,
		fieldY,
		n.UIRect.Width/2 - UIFieldSpacing/2,
		UIFieldHeight,
	}, "-") {
		if len(o.Cols) > 1 {
			o.Cols = o.Cols[:len(o.Cols)-1]
			o.ColDropdowns = o.ColDropdowns[:len(o.ColDropdowns)-1]
		}
	}

	for i := len(o.ColDropdowns) - 1; i >= 0; i-- {
		fieldY -= UIFieldSpacing + UIFieldHeight
		func() {
			col := &o.Cols[i]
			dropdown := o.ColDropdowns[i]
			if openDropdown == dropdown {
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

			activeSort := 0
			if col.Descending {
				activeSort = 1
			}
			newSort := raygui.ComboBox(rl.Rectangle{
				n.UIRect.X + n.UIRect.Width - orderDirectionWidth,
				fieldY,
				orderDirectionWidth,
				UIFieldHeight,
			}, "A-Z;Z-A", activeSort)
			switch newSort {
			case 1:
				col.Descending = true
			default:
				col.Descending = false
			}
		}()
	}
}
