package app

import (
	"github.com/bvisness/SQLJam/node"
	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func doPickColumnsUpdate(n *node.Node, p *node.PickColumns) {
	uiHeight := 0
	for range p.Entries {
		uiHeight += UIFieldHeight
		uiHeight += UIFieldSpacing
	}
	uiHeight += UIFieldHeight // for buttons

	n.UISize = rl.Vector2{300, float32(uiHeight)}

	opts := columnNameDropdownOpts(n.Inputs[0])
	for _, entry := range p.Entries {
		entry.ColDropdown.SetOptions(opts...)
	}
}

func doPickColumnsUI(n *node.Node, p *node.PickColumns) {
	openDropdown, isOpen := raygui.GetOpenDropdown(p.ColDropdowns())
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
		p.Entries = append(p.Entries, &node.PickColumnsEntry{})
	}
	if raygui.Button(rl.Rectangle{
		n.UIRect.X + n.UIRect.Width/2 + UIFieldSpacing/2,
		fieldY,
		n.UIRect.Width/2 - UIFieldSpacing/2,
		UIFieldHeight,
	}, "-") {
		if len(p.Entries) > 1 {
			p.Entries = p.Entries[:len(p.Entries)-1]
		}
	}

	for i := len(p.Entries) - 1; i >= 0; i-- {
		fieldY -= UIFieldSpacing + UIFieldHeight
		func() {
			entry := p.Entries[i]
			if openDropdown == &entry.ColDropdown {
				raygui.Enable()
				defer raygui.Disable()
			}

			col := entry.ColDropdown.Do(rl.Rectangle{
				n.UIRect.X,
				fieldY,
				n.UIRect.Width/2 - UIFieldSpacing/2,
				UIFieldHeight,
			})
			entry.Col, _ = col.(string)

			aliasRect := rl.Rectangle{
				n.UIRect.X + n.UIRect.Width/2 + UIFieldSpacing/2,
				fieldY,
				n.UIRect.Width/2 - UIFieldSpacing/2,
				UIFieldHeight,
			}
			rl.DrawRectangleRec(aliasRect, rl.White)
			entry.Alias = entry.AliasTextbox.Do(aliasRect, entry.Alias, 100)
		}()
	}
}
