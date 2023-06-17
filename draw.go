package main

import "github.com/hajimehoshi/ebiten/v2"

func ScaleCentered(drawOptions *ebiten.DrawImageOptions, width float64, height float64, x float64, y float64) {
	var halfWidth = width / 2
	var halfHeight = height / 2
	drawOptions.GeoM.Translate(-halfWidth, -halfHeight)
	drawOptions.GeoM.Scale(x, y)
	drawOptions.GeoM.Translate(halfWidth, halfHeight)
}
