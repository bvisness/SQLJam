package app

import (
	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type PickColumns struct {
	NodeData
	Entries []*PickColumnsEntry
}

type PickColumnsEntry struct {
	Col          string
	ColDropdown  raygui.DropdownEx
	Alias        string
	AliasTextbox raygui.TextBoxEx
}

func NewPickColumns() *Node {
	return &Node{
		Title:   "Pick Columns",
		CanSnap: true,
		Color:   rl.NewColor(255, 122, 125, 255),
		Inputs:  make([]*Node, 1),
		Data: &PickColumns{
			Entries: []*PickColumnsEntry{{}},
		},
	}
}

func (pc *PickColumns) Cols() []string {
	res := make([]string, len(pc.Entries))
	for i := range res {
		res[i] = pc.Entries[i].Col
	}
	return res
}

func (pc *PickColumns) Aliases() []string {
	res := make([]string, len(pc.Entries))
	for i := range res {
		res[i] = pc.Entries[i].Alias
	}
	return res
}

func (pc *PickColumns) ColDropdowns() []*raygui.DropdownEx {
	res := make([]*raygui.DropdownEx, len(pc.Entries))
	for i := range res {
		res[i] = &pc.Entries[i].ColDropdown
	}
	return res
}

func (p *PickColumns) Update(n *Node) {
	uiHeight := 0
	for range p.Entries {
		uiHeight += UIFieldHeight
		uiHeight += UIFieldSpacing
	}
	uiHeight += UIFieldHeight // for buttons

	n.UISize = rl.Vector2{400, float32(uiHeight)}

	opts := columnNameDropdownOpts(n.Inputs[0])
	for _, entry := range p.Entries {
		if len(opts) == 0 {
			entry.ColDropdown.SetOptions(errorOpts...)
		} else {
			entry.ColDropdown.SetOptions(opts...)
		}
	}
}

func (p *PickColumns) DoUI(n *Node) {
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
		p.Entries = append(p.Entries, &PickColumnsEntry{})
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
			entry.Col = col.(string)

			aliasRect := rl.Rectangle{
				n.UIRect.X + n.UIRect.Width/2 + UIFieldSpacing/2,
				fieldY,
				n.UIRect.Width/2 - UIFieldSpacing/2,
				UIFieldHeight,
			}
			entry.Alias, _ = entry.AliasTextbox.Do(aliasRect, entry.Alias, 100)
		}()
	}
}

func (d *PickColumns) Serialize() string {
	res := ""
	for _, entry := range d.Entries {
		res += entry.Col
		res += entry.Alias
	}
	return res
}
