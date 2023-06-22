package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func ScaleCentered(drawOptions *ebiten.DrawImageOptions, width float64, height float64, x float64, y float64) {
	var halfWidth = width / 2
	var halfHeight = height / 2
	drawOptions.GeoM.Translate(-halfWidth, -halfHeight)
	drawOptions.GeoM.Scale(x, y)
	drawOptions.GeoM.Translate(halfWidth, halfHeight)
}

func GetShiftedRectangle(shiftX float64, frameWidth float64, frameHeight float64) image.Rectangle {
	return image.Rect(
		RoundFloat64ToInt(shiftX), 0,
		RoundFloat64ToInt(shiftX+frameWidth), RoundFloat64ToInt(frameHeight),
	)
}
