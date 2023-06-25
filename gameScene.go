package main

import "github.com/hajimehoshi/ebiten/v2"

type GameSceneStatus int

const (
	GAME_SCENE_STATUS_HORIZONTAL GameSceneStatus = iota
	GAME_SCENE_STATUS_TRANSITION
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

	Status          GameSceneStatus
	sceneHorizontal GameSceneHorizontal
	sceneTransition GameSceneTransition
	sceneVertical   GameSceneVertical
}

func (me *GameScene) Initialize() {
	me.sceneHorizontal.ViewWidth = me.ViewWidth
	me.sceneHorizontal.ViewHeight = me.ViewHeight
	me.sceneHorizontal.Initialize()
	me.sceneTransition.ViewWidth = me.ViewWidth
	me.sceneTransition.CatViewX = me.sceneHorizontal.GetCatViewX()
	me.sceneTransition.FloorY = me.sceneHorizontal.GetFloorY()
	me.sceneTransition.Initialize()
	me.sceneVertical.ViewWidth = me.ViewWidth
	me.sceneVertical.ViewHeight = me.ViewHeight
	me.sceneVertical.Initialize()
	me.Status = GAME_SCENE_STATUS_HORIZONTAL
}

func (me *GameScene) Update(deltaTime float64) {
	switch me.Status {
	case GAME_SCENE_STATUS_HORIZONTAL:
		me.sceneHorizontal.JustPressedKeys = me.JustPressedKeys
		me.sceneHorizontal.PressedKeys = me.PressedKeys
		me.sceneHorizontal.Update(deltaTime)
		if me.sceneHorizontal.Completed {
			me.Status = GAME_SCENE_STATUS_TRANSITION
			me.sceneTransition.CatRunFrame = me.sceneHorizontal.catEntity.runFrame
		}
	case GAME_SCENE_STATUS_TRANSITION:
		me.sceneTransition.Update(deltaTime)
	case GAME_SCENE_STATUS_VERTICAL:
		me.sceneVertical.JustPressedKeys = me.JustPressedKeys
		me.sceneVertical.PressedKeys = me.PressedKeys
		me.sceneVertical.Update(deltaTime)
	}
}

func (me *GameScene) Draw(screen *ebiten.Image) {
	switch me.Status {
	case GAME_SCENE_STATUS_HORIZONTAL:
		me.sceneHorizontal.Draw(screen)
	case GAME_SCENE_STATUS_TRANSITION:
		me.sceneTransition.Draw(screen)
	case GAME_SCENE_STATUS_VERTICAL:
		me.sceneVertical.Draw(screen)
	}
}
