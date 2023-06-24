package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func DrawTorchLight(screen *ebiten.Image, x float32, y float32) {
	vector.DrawFilledCircle(screen, x, y, 16, color.NRGBA{R: 255, G: 244, B: 188, A: 16}, false)
}
