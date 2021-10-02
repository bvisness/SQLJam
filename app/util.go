package app

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func Vector2Rotate(v rl.Vector2, rad float32) rl.Vector2 {
	return rl.Vector2{
		X: v.X*float32(math.Cos(float64(rad))) - v.Y*float32(math.Sin(float64(rad))),
		Y: v.X*float32(math.Sin(float64(rad))) + v.Y*float32(math.Cos(float64(rad))),
	}
}

func Lerp(a, b, t float32) float32 {
	return (1-t)*a + t*b
}

func EaseInOutCubic(x float32) float32 {
	if x < 0.5 {
		return 4 * x * x * x
	} else {
		return 1 - float32(math.Pow(-2*float64(x)+2, 3))/2
	}
}

func Clamp(x, min, max float32) float32 {
	if x < min {
		return min
	} else if x > max {
		return max
	} else {
		return x
	}
}
