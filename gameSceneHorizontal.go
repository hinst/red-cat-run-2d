package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type GameSceneHorizontal struct {
	// Initialization input parameter
	ViewHeight float64
	// Initialization input parameter
	ViewWidth float64
	// Input parameter for every update
	JustPressedKeys []ebiten.Key
	// Input parameter for every update
	PressedKeys []ebiten.Key

	terrainMan              TerrainMan
	catEntity               CatEntity
	dustMan                 DustMan
	cameraX                 float64
	cameraY                 float64
	transitionTimeRemaining float64
}

const GAME_SCENE_TRANSITION_TIME = 2

func (me *GameSceneHorizontal) Initialize() {
	me.terrainMan.ViewWidth = me.ViewWidth
	me.terrainMan.ViewHeight = me.ViewHeight
	me.terrainMan.AreaWidth = 100
	me.terrainMan.FloorY = me.GetFloorY()
	me.terrainMan.CeilingY = me.GetCeilingY()
	me.terrainMan.Initialize()

	me.catEntity.ViewWidth = me.ViewWidth
	me.catEntity.ViewHeight = me.ViewHeight
	me.catEntity.X = me.GetCatViewX()
	me.catEntity.FloorY = me.GetFloorY()
	me.catEntity.CeilingY = me.GetCeilingY()
	me.catEntity.Initialize()

	me.dustMan.ViewWidth = me.ViewWidth
	me.dustMan.ViewHeight = me.ViewHeight
	me.dustMan.AreaWidth = me.GetAreaWidth()
	me.dustMan.Initialize()
}

func (me *GameSceneHorizontal) Update(deltaTime float64) {
	me.terrainMan.Update(deltaTime)
	if me.transitionTimeRemaining == 0 {
		if me.catEntity.Status == CAT_ENTITY_STATUS_RUN && !me.CheckCatHold() {
			me.catEntity.Status = CAT_ENTITY_STATUS_DEAD
		}
		me.catEntity.JustPressedKeys = me.JustPressedKeys
		me.catEntity.PressedKeys = me.PressedKeys
		me.catEntity.Update(deltaTime)
		if me.catEntity.Status != CAT_ENTITY_STATUS_DEAD {
			if me.catEntity.Direction == DIRECTION_RIGHT {
				me.cameraX = me.getCameraXGoingRight()
			} else {
				me.cameraX = me.getCameraXGoingLeft()
			}
		}
		if me.CheckCatAtRightEndOfTerrain() && me.catEntity.Direction == DIRECTION_RIGHT {
			me.switchDirection()
		}
	} else {
		me.transitionTimeRemaining -= deltaTime
		if me.transitionTimeRemaining <= 0 {
			me.transitionTimeRemaining = 0
		}
		me.cameraX = me.getCameraXGoingLeft() + me.transitionTimeRemaining*(me.getCameraXGoingRight()-me.getCameraXGoingLeft())/2
	}
	me.dustMan.CameraX = me.cameraX
	me.dustMan.Update(deltaTime)
}

func (me *GameSceneHorizontal) Draw(screen *ebiten.Image) {
	me.dustMan.Draw(screen)

	me.catEntity.CameraX = me.cameraX
	me.catEntity.CameraY = me.cameraY
	me.catEntity.Draw(screen)

	me.terrainMan.CameraX = me.cameraX
	me.terrainMan.CameraY = me.cameraY
	me.terrainMan.Draw(screen)
}

func (me *GameSceneHorizontal) GetFloorY() float64 {
	return 200
}

func (me *GameSceneHorizontal) GetCeilingY() float64 {
	return 40
}

// The distance from left view border to the cat
func (me *GameSceneHorizontal) GetCatViewX() float64 {
	return 10
}

func (me *GameSceneHorizontal) CheckCatHold() bool {
	for _, block := range me.terrainMan.GetBlocks() {
		var isFittingBlock = me.catEntity.Location == block.Location
		if isFittingBlock {
			if CheckDualIntersect(
				me.catEntity.X,
				me.catEntity.X+me.catEntity.Width,
				float64(block.X)*float64(me.terrainMan.GetTileWidth()),
				float64(block.X+block.Width)*float64(me.terrainMan.GetTileWidth()),
			) {
				return true
			}
		}
	}
	return false
}

func (me *GameSceneHorizontal) CheckCatAtRightEndOfTerrain() bool {
	var catRight = me.catEntity.X + me.catEntity.Width
	return catRight >= me.GetAreaWidth()
}

// Measurement unit: pixels
func (me *GameSceneHorizontal) GetAreaWidth() float64 {
	return float64(me.terrainMan.AreaWidth) * float64(me.terrainMan.GetTileWidth())
}

func (me *GameSceneHorizontal) getCameraXGoingRight() float64 {
	return me.catEntity.X - me.GetCatViewX()
}

func (me *GameSceneHorizontal) getCameraXGoingLeft() float64 {
	return me.catEntity.X + me.catEntity.Width - me.ViewWidth + me.GetCatViewX()
}

func (me *GameSceneHorizontal) switchDirection() {
	me.transitionTimeRemaining = GAME_SCENE_TRANSITION_TIME
	me.catEntity.Direction = DIRECTION_LEFT
	me.terrainMan.Shuffle()
}
