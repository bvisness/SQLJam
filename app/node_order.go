package app

import (
	"fmt"

	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var SortColor = rl.NewColor(255, 204, 128, 255)

const sortDirectionWidth = 50

type Sort struct {
	Cols []*SortColumn
}

type SortColumn struct {
	Col         string
	Descending  bool
	ColDropdown raygui.DropdownEx
}

func NewSort() *Node {
	return &Node{
		Title:   "Sort",
		CanSnap: true,
		Color:   rl.NewColor(255, 204, 128, 255),
		Inputs:  make([]*Node, 1),
		Data: &Sort{
			Cols: []*SortColumn{{}},
		},
	}
}

func (oc *Sort) ColDropdowns() []*raygui.DropdownEx {
	res := make([]*raygui.DropdownEx, len(oc.Cols))
	for i := range res {
		res[i] = &oc.Cols[i].ColDropdown
	}
	return res
}

func (d *Sort) Update(n *Node) {
	uiHeight := 0
	for range d.Cols {
		uiHeight += UIFieldHeight
		uiHeight += UIFieldSpacing
	}
	uiHeight += UIFieldHeight // for buttons

	n.UISize = rl.Vector2{300, float32(uiHeight)}

	opts := columnNameDropdownOpts(n.Inputs[0])
	for _, col := range d.Cols {
		col.ColDropdown.SetOptions(opts...)
	}
}

func (d *Sort) DoUI(n *Node) {
	openDropdown, isOpen := raygui.GetOpenDropdown(d.ColDropdowns())
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
		d.Cols = append(d.Cols, &SortColumn{})
	}
	if raygui.Button(rl.Rectangle{
		n.UIRect.X + n.UIRect.Width/2 + UIFieldSpacing/2,
		fieldY,
		n.UIRect.Width/2 - UIFieldSpacing/2,
		UIFieldHeight,
	}, "-") {
		if len(d.Cols) > 1 {
			d.Cols = d.Cols[:len(d.Cols)-1]
		}
	}

	for i := len(d.Cols) - 1; i >= 0; i-- {
		func() {
			fieldY -= UIFieldSpacing + UIFieldHeight

			col := d.Cols[i]
			dropdown := &col.ColDropdown

			if openDropdown == &col.ColDropdown {
				raygui.Enable()
				defer raygui.Disable()
			}

			colName := dropdown.Do(rl.Rectangle{
				n.UIRect.X,
				fieldY,
				n.UIRect.Width - sortDirectionWidth - UIFieldSpacing,
				UIFieldHeight,
			})
			col.Col, _ = colName.(string)

			directionStr := "A-Z"
			if col.Descending {
				directionStr = "Z-A"
			}
			col.Descending = raygui.Toggle(rl.Rectangle{
				n.UIRect.X + n.UIRect.Width - sortDirectionWidth,
				fieldY,
				sortDirectionWidth,
				UIFieldHeight,
			}, directionStr, col.Descending)
		}()
	}
}

func (d *Sort) Serialize() (res string, active bool) {
	for _, col := range d.Cols {
		res += col.Col
		res += fmt.Sprintf("%v", col.Descending)
	}
	return res, false
}
