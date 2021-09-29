package raygui

import (
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type DropdownEx struct {
	options []DropdownExOption
	active  int
	open    bool

	str string
}

type DropdownExOption struct {
	Name  string
	Value interface{}
}

func NewDropdownEx(opts ...DropdownExOption) DropdownEx {
	d := DropdownEx{}
	d.SetOptions(opts...)

	return d
}

func (d *DropdownEx) Do(bounds rl.Rectangle) interface{} {
	if d.active >= len(d.options) {
		d.active = len(d.options) - 1
	}
	if d.active < 0 {
		d.active = 0
	}

	toggle := DropdownBox(bounds, d.str, &d.active, d.open)
	if toggle {
		d.open = !d.open
	}

	if len(d.options) == 0 {
		return nil
	} else {
		return d.options[d.active].Value
	}
}

func (d *DropdownEx) GetOptions() []DropdownExOption {
	return d.options
}

func (d *DropdownEx) SetOptions(opts ...DropdownExOption) {
	var names []string
	for _, opt := range opts {
		names = append(names, opt.Name)
	}

	d.options = opts
	d.str = strings.Join(names, ";")
}

type TextBoxEx struct {
	active bool
}

func NewTextBoxEx() TextBoxEx {
	return TextBoxEx{}
}

func (t *TextBoxEx) Do(bounds rl.Rectangle, text string, textSize int) string {
	newText, toggle := TextBox(bounds, text, textSize, t.active)
	if toggle {
		t.active = !t.active
	}
	return newText
}
