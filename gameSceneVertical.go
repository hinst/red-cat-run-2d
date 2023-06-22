package main

import "github.com/hajimehoshi/ebiten/v2"

type GameSceneVertical struct {
	// Initialization input parameter
	ViewHeight float64
	// Initialization input parameter
	ViewWidth float64
	// Input parameter for every update
	JustPressedKeys []ebiten.Key
	// Input parameter for every update
	PressedKeys []ebiten.Key

	catEntity CatEntityVertical
}

func (me *GameSceneVertical) Initialize() {
}
