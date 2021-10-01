package raygui

import (
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var cam *rl.Camera2D

func Set2DCamera(camera *rl.Camera2D) {
	cam = camera
}

// Gets mouse position, respecting the 2D camera's transform.
func GetMousePositionWorld() rl.Vector2 {
	if cam == nil {
		return rl.GetMousePosition()
	}
	return rl.GetScreenToWorld2D(rl.GetMousePosition(), *cam)
}

type DropdownEx struct {
	Open bool

	options []DropdownExOption
	active  int
	str     string
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

func MakeDropdownExList(n int, opts ...DropdownExOption) []*DropdownEx {
	list := make([]*DropdownEx, n)
	for i := range list {
		d := NewDropdownEx(opts...)
		list[i] = &d
	}
	return list
}

func (d *DropdownEx) Do(bounds rl.Rectangle) interface{} {
	d.fixupActive()

	toggle := DropdownBox(bounds, d.str, &d.active, d.Open)
	if toggle {
		d.Open = !d.Open
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

func (d *DropdownEx) fixupActive() {
	if d.active >= len(d.options) {
		d.active = len(d.options) - 1
	}
	if d.active < 0 {
		d.active = 0
	}
}

func GetOpenDropdown(dropdowns []*DropdownEx) (*DropdownEx, bool) {
	for _, other := range dropdowns {
		if other.Open {
			return other, true
		}
	}

	return nil, false
}

type TextBoxEx struct {
	active bool
}

func NewTextBoxEx() TextBoxEx {
	return TextBoxEx{}
}

func MakeTextBoxExList(n int) []*TextBoxEx {
	list := make([]*TextBoxEx, n)
	for i := range list {
		t := NewTextBoxEx()
		list[i] = &t
	}
	return list
}

func (t *TextBoxEx) Do(bounds rl.Rectangle, text string, textSize int) string {
	newText, toggle := TextBox(bounds, text, textSize, t.active)
	if toggle {
		t.active = !t.active
	}
	return newText
}
