package main

import "github.com/hajimehoshi/ebiten/v2"

type GameSceneStatus int

const (
	GAME_SCENE_STATUS_HORIZONTAL GameSceneStatus = iota
	GAME_SCENE_STATUS_VERTICAL
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
	me.Status = GAME_SCENE_STATUS_HORIZONTAL
}

func (me *GameScene) Update(deltaTime float64) {
	switch me.Status {
	case GAME_SCENE_STATUS_HORIZONTAL:
		me.gameSceneHorizontal.JustPressedKeys = me.JustPressedKeys
		me.gameSceneHorizontal.PressedKeys = me.PressedKeys
		me.gameSceneHorizontal.Update(deltaTime)
	case GAME_SCENE_STATUS_VERTICAL:
		me.gameSceneVertical.JustPressedKeys = me.JustPressedKeys
		me.gameSceneVertical.PressedKeys = me.PressedKeys
		me.gameSceneVertical.Update(deltaTime)
	}
}

func (me *GameScene) Draw(screen *ebiten.Image) {
	switch me.Status {
	case GAME_SCENE_STATUS_HORIZONTAL:
		me.gameSceneHorizontal.Draw(screen)
	case GAME_SCENE_STATUS_VERTICAL:
		me.gameSceneVertical.Draw(screen)
	}
}
