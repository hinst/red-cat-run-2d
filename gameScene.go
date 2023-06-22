package main

import "github.com/hajimehoshi/ebiten/v2"

type GameSceneStatus int

const (
	GAME_SCENE_STATUS_HORIZONTAL GameSceneStatus = iota
)

type GameScene struct {
	// Initialization input parameter. Measurement unit: pixels
	ViewWidth float64
	// Initialization input parameter. Measurement unit: pixels
	ViewHeight float64
	// Input parameter for every update
	JustPressedKeys []ebiten.Key
	// Input parameter for every update
	PressedKeys []ebiten.Key

	Status              GameSceneStatus
	gameSceneHorizontal GameSceneHorizontal
	gameSceneVertical   GameSceneVertical
}

func (me *GameScene) Initialize() {
	me.gameSceneHorizontal.ViewWidth = me.ViewWidth
	me.gameSceneHorizontal.ViewHeight = me.ViewHeight
	me.gameSceneHorizontal.Initialize()
	me.gameSceneVertical.ViewWidth = me.ViewWidth
	me.gameSceneVertical.ViewHeight = me.ViewHeight
	me.gameSceneVertical.Initialize()
}

func (me *GameScene) Update(deltaTime float64) {
	if me.Status == GAME_SCENE_STATUS_HORIZONTAL {
		me.gameSceneHorizontal.JustPressedKeys = me.JustPressedKeys
		me.gameSceneHorizontal.PressedKeys = me.PressedKeys
		me.gameSceneHorizontal.Update(deltaTime)
	}
}

func (me *GameScene) Draw(screen *ebiten.Image) {
	if me.Status == GAME_SCENE_STATUS_HORIZONTAL {
		me.gameSceneHorizontal.Draw(screen)
	}
}
