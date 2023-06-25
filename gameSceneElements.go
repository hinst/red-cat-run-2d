package main

import (
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var SHAFT_COLOR_MULTIPLIER = float64(127) / 255
var SHAFT_COLOR = color.RGBA{
	R: uint8(50 * SHAFT_COLOR_MULTIPLIER),
	G: uint8(35 * SHAFT_COLOR_MULTIPLIER),
	B: uint8(27 * SHAFT_COLOR_MULTIPLIER),
	A: 255}
var TORCH_LIGHT_COLOR = color.NRGBA{R: 255, G: 244, B: 188, A: 10}

func DrawTorchLight(screen *ebiten.Image, x float32, y float32) {
	vector.DrawFilledCircle(screen, x, y, 16, TORCH_LIGHT_COLOR, false)
}

func DrawTorch(screen *ebiten.Image, torchImage *ebiten.Image, x float64, y float64) {
	var torchScale = 0.5
	var imageWidth = float64(torchImage.Bounds().Dx())
	var imageHeight = float64(torchImage.Bounds().Dy())
	var xScaleMultiplier float64 = 1
	if time.Now().Nanosecond() < 1000_000_000/2 {
		xScaleMultiplier = -1
	}
	var drawOptions = ebiten.DrawImageOptions{}
	ScaleCentered(&drawOptions, imageWidth, imageHeight, xScaleMultiplier, 1)
	drawOptions.GeoM.Scale(torchScale, torchScale)
	DrawTorchLight(screen, float32(x), float32(y))
	drawOptions.GeoM.Translate(
		math.Round(x-imageWidth/2*torchScale),
		math.Round(y-imageHeight/2*torchScale),
	)
	screen.DrawImage(torchImage, &drawOptions)
}
