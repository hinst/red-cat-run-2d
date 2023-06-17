package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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

	terrainMan TerrainMan
	catEntity  CatEntity
	cameraX    float64
	cameraY    float64
}

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
	if me.catEntity.Status == me.catEntity.GetStatusRun() && !me.CheckCatHold() {
		me.catEntity.Status = me.catEntity.GetStatusDead()
	}
	me.catEntity.JustPressedKeys = me.JustPressedKeys
	me.catEntity.Update(deltaTime)
	me.cameraX = me.catEntity.X - me.GetCatViewX()
}

func (me *GameScene) Draw(screen *ebiten.Image) {
	me.terrainMan.CameraX = me.cameraX
	me.terrainMan.CameraY = me.cameraY
	me.terrainMan.Draw(screen)

	if me.catEntity.Status == me.catEntity.GetStatusRun() {
		if me.catEntity.Location == me.catEntity.GetLocationFloor() {
			for _, key := range me.PressedKeys {
				if key == ebiten.KeyUp {
					me.drawAimLine(screen, true)
				}
			}
		} else if me.catEntity.Location == me.catEntity.GetLocationCeiling() {
			for _, key := range me.PressedKeys {
				if key == ebiten.KeyDown {
					me.drawAimLine(screen, false)
				}
			}
		}
	}

	me.catEntity.CameraX = me.cameraX
	me.catEntity.CameraY = me.cameraY
	me.catEntity.Draw(screen)
}

func (me *GameScene) drawAimLine(screen *ebiten.Image, up bool) {
	var y1 = me.catEntity.Y+me.catEntity.Height/2-me.cameraY
	if up {
		y1 -= me.catEntity.JumpSpeed
	} else {
		y1 += me.catEntity.JumpSpeed
	}
	vector.StrokeLine(screen,
		float32(me.catEntity.X+me.catEntity.Width/2-me.cameraX),
		float32(me.catEntity.Y+me.catEntity.Height/2-me.cameraY),
		float32(me.catEntity.X+me.catEntity.Width/2+me.catEntity.Speed-me.cameraX),
		float32(y1),
		1,
		color.RGBA{R: 150, G: 100, B: 100, A: 255},
		false,
	)
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
