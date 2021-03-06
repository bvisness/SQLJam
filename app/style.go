package app

import (
	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func ToggleTheme() {
	//d := Clamp(5 + 255 * float32(5), 0, 255)
}

func ToHexNum(c rl.Color) uint {
	col := uint(0)
	col += uint(c.R) << 24
	col += uint(c.G) << 16
	col += uint(c.B) << 8
	col += uint(c.A) << 0
	return col
}

// ok I think I botched these but um this gets us there for now

func Tint(c rl.Color, amt float32) rl.Color {
	base := float32(255)
	return rl.NewColor(
		AffectColor(c.R, amt, base),
		AffectColor(c.G, amt, base),
		AffectColor(c.B, amt, base),
		AffectColor(c.A, amt, base),
	)
}

func Brightness(c rl.Color, amt float32) rl.Color {
	hsv := rl.ColorToHSV(c)
	return rl.ColorFromHSV(hsv.X, hsv.Y, Clamp(amt*hsv.Z, 0, 1))
}

func AffectColor(cv uint8, amt float32, base float32) uint8 {
	return uint8(Clamp(float32(cv)+(base-float32(cv))*amt, 0, 255))
}

func SetStyleColor(control raygui.Control, property raygui.ControlProperty, color rl.Color) {
	raygui.SetStyle(control, property, ToHexNum(color))
}

func LoadThemeForNode(n *Node) {
	LoadThemeForColor(n.Color)
}

func LoadThemeForColor(color rl.Color) {
	dark1 := Brightness(color, 0.4)
	light1 := Tint(color, 0.3)

	SetStyleColor(raygui.Default, raygui.BaseColorNormalProp, color)
	SetStyleColor(raygui.Default, raygui.BaseColorFocusedProp, dark1)
	SetStyleColor(raygui.Default, raygui.BaseColorPressedProp, light1)

	SetStyleColor(raygui.Default, raygui.TextColorNormalProp, dark1)
	SetStyleColor(raygui.Default, raygui.TextColorFocusedProp, color)
	SetStyleColor(raygui.Default, raygui.TextColorPressedProp, dark1)

	SetStyleColor(raygui.Default, raygui.BorderColorNormalProp, dark1)
	SetStyleColor(raygui.Default, raygui.BorderColorFocusedProp, dark1)
	SetStyleColor(raygui.Default, raygui.BorderColorPressedProp, dark1)

	SetStyleColor(raygui.Default, raygui.LineColorProp, dark1)
	SetStyleColor(raygui.Default, raygui.BackgroundColorProp, color)

	SetStyleColor(raygui.TextBoxControl, raygui.BaseColorPressedProp, light1)
	SetStyleColor(raygui.TextBoxControl, raygui.BorderColorPressedProp, dark1)

	// Disabling stuff
	SetStyleColor(raygui.Default, raygui.BaseColorDisabledProp, light1)
	SetStyleColor(raygui.Default, raygui.BorderColorDisabledProp, dark1)
	SetStyleColor(raygui.Default, raygui.TextColorDisabledProp, dark1)

	//SetStyleColor(raygui.DropdownBoxControl, raygui.BackgroundColorProp, rl.Red)
}
