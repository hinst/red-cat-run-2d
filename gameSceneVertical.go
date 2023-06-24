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

	catEntity       CatEntityVertical
	fallObstacleMan FallObstacleMan
	cameraY         float64
}

func (me *GameSceneVertical) Initialize() {
	if me.ViewHeight == 0 || me.ViewWidth == 0 {
		log.Println("Warning: view size is missing")
	}
	me.catEntity.Initialize()
	me.cameraY = me.catEntity.Y - 10
	me.catEntity.X = me.ViewWidth/2 - me.catEntity.Width/2
	me.fallObstacleMan.AreaWidth = me.GetAreaWidth()
	me.fallObstacleMan.AreaHeight = me.ViewHeight * 10
	me.fallObstacleMan.ViewWidth = me.ViewWidth
	me.fallObstacleMan.ViewHeight = me.ViewHeight
	me.fallObstacleMan.Initialize()
}

func (me *GameSceneVertical) Draw(screen *ebiten.Image) {
	me.catEntity.Draw(screen)
	me.fallObstacleMan.Draw(screen)
}

func (me *GameSceneVertical) Update(deltaTime float64) {
	me.catEntity.CameraY = me.cameraY
	me.catEntity.Update(deltaTime)
	me.fallObstacleMan.CameraY = me.cameraY
	me.fallObstacleMan.Update(deltaTime)
	me.cameraY = me.catEntity.Y - 10
}

func (me *GameSceneVertical) GetAreaWidth() float64 {
	return 220
}
