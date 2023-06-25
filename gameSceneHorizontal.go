package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type GameSceneHorizontal struct {
	// Initialization input parameter
	ViewWidth float64
	// Initialization input parameter
	ViewHeight float64
	// Input parameter for every update
	JustPressedKeys []ebiten.Key
	// Input parameter for every update
	PressedKeys []ebiten.Key

	terrainMan              TerrainMan
	CatEntity               CatEntityHorizontal
	dustMan                 DustMan
	cameraX                 float64
	cameraY                 float64
	transitionTimeRemaining float64
	fishImage               *ebiten.Image

	// Output parameter after each update
	Completed bool
}

const GAME_SCENE_TRANSITION_TIME = 2

func (me *GameSceneHorizontal) Initialize() {
	me.terrainMan.ViewWidth = me.ViewWidth
	me.terrainMan.ViewHeight = me.ViewHeight
	me.terrainMan.AreaWidth = 10
	me.terrainMan.FloorY = me.GetFloorY()
	me.terrainMan.CeilingY = me.GetCeilingY()
	me.terrainMan.Initialize()

	me.CatEntity.ViewWidth = me.ViewWidth
	me.CatEntity.ViewHeight = me.ViewHeight
	me.CatEntity.X = me.GetCatViewX()
	me.CatEntity.FloorY = me.GetFloorY()
	me.CatEntity.CeilingY = me.GetCeilingY()
	me.CatEntity.Initialize()

	me.dustMan.ViewWidth = me.ViewWidth
	me.dustMan.ViewHeight = me.ViewHeight
	me.dustMan.AreaWidth = me.GetAreaWidth()
	me.dustMan.Initialize()

	me.fishImage = LoadImage(FISH_IMAGE_BYTES)
}

func (me *GameSceneHorizontal) Update(deltaTime float64) {
	me.terrainMan.Update(deltaTime)
	if me.transitionTimeRemaining == 0 {
		if me.CatEntity.Status == CAT_ENTITY_STATUS_RUN && !me.CheckCatHold() {
			me.CatEntity.Status = CAT_ENTITY_STATUS_DEAD
		}
		me.CatEntity.JustPressedKeys = me.JustPressedKeys
		me.CatEntity.PressedKeys = me.PressedKeys
		me.CatEntity.Update(deltaTime)
		if me.CatEntity.Status != CAT_ENTITY_STATUS_DEAD {
			if me.CatEntity.Direction == DIRECTION_RIGHT {
				me.cameraX = me.getCameraXGoingRight()
			} else {
				me.cameraX = me.getCameraXGoingLeft()
			}
		}
		if me.CheckCatAtRightEndOfTerrain() && me.CatEntity.Direction == DIRECTION_RIGHT {
			me.switchDirection()
			PlaySound(REVERSE_SOUND_BYTES, 0.20)
		}
		if me.CheckCatAtLeftEndOfTerrain() && me.CatEntity.Direction == DIRECTION_LEFT {
			me.Completed = true
			PlaySound(REVERSE_SOUND_BYTES, 0.20)
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
	me.drawFish(screen)

	me.drawShaftBackground(screen)
	me.terrainMan.CameraX = me.cameraX
	me.terrainMan.CameraY = me.cameraY
	me.terrainMan.DrawLowerLayer(screen)

	me.CatEntity.CameraX = me.cameraX
	me.CatEntity.CameraY = me.cameraY
	me.CatEntity.Draw(screen)

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
		var isFittingBlock = me.CatEntity.Location == block.Location
		if isFittingBlock {
			if CheckDualIntersect(
				me.CatEntity.X,
				me.CatEntity.X+me.CatEntity.Width,
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
	var catRight = me.CatEntity.X + me.CatEntity.Width
	return catRight >= me.GetAreaWidth()
}

func (me *GameSceneHorizontal) CheckCatAtLeftEndOfTerrain() bool {
	return me.CatEntity.X <= 0
}

// Measurement unit: pixels
func (me *GameSceneHorizontal) GetAreaWidth() float64 {
	return float64(me.terrainMan.AreaWidth) * float64(me.terrainMan.GetTileWidth())
}

func (me *GameSceneHorizontal) getCameraXGoingRight() float64 {
	return me.CatEntity.X - me.GetCatViewX()
}

func (me *GameSceneHorizontal) getCameraXGoingLeft() float64 {
	return me.CatEntity.X + me.CatEntity.Width - me.ViewWidth + me.GetCatViewX()
}

func (me *GameSceneHorizontal) switchDirection() {
	me.transitionTimeRemaining = GAME_SCENE_TRANSITION_TIME
	me.CatEntity.Direction = DIRECTION_LEFT
	me.terrainMan.Shuffle()
}

func (me *GameSceneHorizontal) drawFish(screen *ebiten.Image) {
	if me.CatEntity.Direction == DIRECTION_RIGHT {
		var drawOptions ebiten.DrawImageOptions
		var y float64
		if me.terrainMan.GetLastBlock().Location == TERRAIN_LOCATION_FLOOR {
			y = me.GetFloorY() - float64(me.fishImage.Bounds().Dy()) - 1
		} else if me.terrainMan.GetLastBlock().Location == TERRAIN_LOCATION_CEILING {
			y = me.GetCeilingY() + 1
		}
		var x = me.GetAreaWidth() - float64(me.fishImage.Bounds().Dx())
		drawOptions.GeoM.Translate(x-me.cameraX, y)
		screen.DrawImage(me.fishImage, &drawOptions)
	} else if me.CatEntity.Direction == DIRECTION_LEFT {
		var drawOptions ebiten.DrawImageOptions
		var y = me.GetFloorY() - float64(me.fishImage.Bounds().Dy()) - 1
		var x float64 = 0
		drawOptions.GeoM.Translate(x-me.cameraX, y)
		screen.DrawImage(me.fishImage, &drawOptions)
	}
}

func (me *GameSceneHorizontal) drawShaftBackground(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, 0, float32(me.ViewWidth), float32(me.GetCeilingY()), SHAFT_COLOR, false)
	vector.DrawFilledRect(screen, 0, float32(me.GetFloorY()), float32(me.ViewWidth), float32(me.ViewHeight), SHAFT_COLOR, false)
}
