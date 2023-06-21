package main

import (
	"bytes"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func LoadImage(data []byte) *ebiten.Image {
	var image, _, imageError = image.Decode(bytes.NewReader(data))
	AssertError(imageError)
	return ebiten.NewImageFromImage(image)
}
