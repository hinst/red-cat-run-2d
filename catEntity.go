package main

import "github.com/hajimehoshi/ebiten/v2"

type CatEntity struct {
	// Input parameter for initialization
	ViewWidth float64
	// Input parameter for initialization
	ViewHeight float64
	// Input parameter for initialization
	FloorY float64
	// Input parameter for initialization
	CeilingY float64
	// Input parameter for update
	JustPressedKeys []ebiten.Key
	// Input parameter for update
	PressedKeys []ebiten.Key
	// Input parameter for every draw
	CameraX float64
	// Input parameter for every draw
	CameraY float64

	X      float64
	Y      float64
	Width  float64
	Height float64
}
