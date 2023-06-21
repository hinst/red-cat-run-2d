package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type GameScene struct {
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
	cameraX                 float64
	cameraY                 float64
	transitionTimeRemaining float64
}

const GAME_SCENE_TRANSITION_TIME = 2

func (me *GameScene) Initialize() {
	me.terrainMan.ViewHeight = me.ViewHeight
	me.terrainMan.ViewWidth = me.ViewWidth
	me.terrainMan.AreaWidth = 100
	me.terrainMan.FloorY = me.GetFloorY()
	me.terrainMan.CeilingY = me.GetCeilingY()
	me.terrainMan.Initialize()

	me.catEntity.ViewHeight = me.ViewHeight
	me.catEntity.ViewWidth = me.ViewWidth
	me.catEntity.X = me.GetCatViewX()
	me.catEntity.FloorY = me.GetFloorY()
	me.catEntity.CeilingY = me.GetCeilingY()
	me.catEntity.Initialize()
}

func (me *GameScene) Update(deltaTime float64) {
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
			me.transitionTimeRemaining = GAME_SCENE_TRANSITION_TIME
			me.catEntity.Direction = DIRECTION_LEFT
		}
	} else {
		me.transitionTimeRemaining -= deltaTime
		if me.transitionTimeRemaining <= 0 {
			me.transitionTimeRemaining = 0
		}
		me.cameraX = me.getCameraXGoingLeft() + me.transitionTimeRemaining*(me.getCameraXGoingRight()-me.getCameraXGoingLeft())/2
	}
}

func (me *GameScene) Draw(screen *ebiten.Image) {
	me.catEntity.CameraX = me.cameraX
	me.catEntity.CameraY = me.cameraY
	me.catEntity.Draw(screen)

	me.terrainMan.CameraX = me.cameraX
	me.terrainMan.CameraY = me.cameraY
	me.terrainMan.Draw(screen)
}

func (me *GameScene) GetFloorY() float64 {
	return 200
}

func (me *GameScene) GetCeilingY() float64 {
	return 40
}

// The distance from left view border to the cat
func (me *GameScene) GetCatViewX() float64 {
	return 10
}

func (me *GameScene) CheckCatHold() bool {
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

func (me *GameScene) CheckCatAtRightEndOfTerrain() bool {
	var catRight = me.catEntity.X + me.catEntity.Width
	var terrainRight = float64(me.terrainMan.AreaWidth) * float64(me.terrainMan.GetTileWidth())
	return catRight >= terrainRight
}

func (me *GameScene) getCameraXGoingRight() float64 {
	return me.catEntity.X - me.GetCatViewX()
}

func (me *GameScene) getCameraXGoingLeft() float64 {
	return me.catEntity.X + me.catEntity.Width - me.ViewWidth + me.GetCatViewX()
}
