package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var SHAFT_COLOR_MULTIPLIER = float64(100) / 255
var SHAFT_COLOR = color.RGBA{R: 0, G: uint8(43 * SHAFT_COLOR_MULTIPLIER), B: uint8(54 * SHAFT_COLOR_MULTIPLIER), A: 255}
var TORCH_LIGHT_COLOR = color.NRGBA{R: 255, G: 244, B: 188, A: 10}

func DrawTorchLight(screen *ebiten.Image, x float32, y float32) {
	vector.DrawFilledCircle(screen, x, y, 16, TORCH_LIGHT_COLOR, false)
}
