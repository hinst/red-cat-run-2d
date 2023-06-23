package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

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
	if me.ViewHeight == 0 || me.ViewWidth == 0 {
		log.Println("Warning: view size is missing")
	}
	me.catEntity.Initialize()
	me.catEntity.X = me.ViewWidth/2 - me.catEntity.Width/2
}

func (me *GameSceneVertical) Draw(screen *ebiten.Image) {
	me.catEntity.Draw(screen)
}

func (me *GameSceneVertical) Update(deltaTime float64) {
	me.catEntity.Update(deltaTime)
}
