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
	CameraX    float64
	CameraY    float64
}

func (me *GameScene) Initialize() {
	me.terrainMan.ViewHeight = me.ViewHeight
	me.terrainMan.ViewWidth = me.ViewWidth
	me.terrainMan.AreaWidth = 100
	me.terrainMan.Initialize()
	me.catEntity.ViewHeight = me.ViewHeight
	me.catEntity.ViewWidth = me.ViewWidth
	me.catEntity.Initialize()
	me.catEntity.X = me.GetCatViewX()
	me.catEntity.Y = me.GetCatY() - float64(me.catEntity.Height)
}

func (me *GameScene) Update(deltaTime float64) {
	if !me.CheckCatHold() {
		me.catEntity.Status = me.catEntity.GetStatusDead()
	}
	me.catEntity.Update(deltaTime)
	me.CameraX = me.catEntity.X - me.GetCatViewX()
}

func (me *GameScene) Draw(screen *ebiten.Image) {
	me.terrainMan.CameraX = me.CameraX
	me.terrainMan.CameraY = me.CameraY
	me.terrainMan.Draw(screen)

	if me.catEntity.Status == me.catEntity.GetStatusRun() {
		for _, key := range me.PressedKeys {
			if key == ebiten.KeyUp {
				vector.StrokeLine(screen,
					float32(me.catEntity.X+me.catEntity.Width/2-me.CameraX),
					float32(me.catEntity.Y+me.catEntity.Height/2-me.CameraY),
					float32(me.catEntity.X+me.catEntity.Width/2+30-me.CameraX),
					float32(me.catEntity.Y+me.catEntity.Height/2-60-me.CameraY),
					1,
					color.RGBA{R: 150, G: 100, B: 100, A: 255},
					false,
				)
			}
		}
	}

	me.catEntity.CameraX = me.CameraX
	me.catEntity.CameraY = me.CameraY
	me.catEntity.Draw(screen)
}

func (me *GameScene) GetCatY() float64 {
	return 200
}

func (me *GameScene) GetCatViewX() float64 {
	return 10
}

func (me *GameScene) CheckCatHold() bool {
	for _, block := range me.terrainMan.GetBlocks() {
		if me.catEntity.Status == me.catEntity.GetStatusRun() &&
			me.catEntity.Location == me.catEntity.GetLocationFloor() &&
			block.Type == block.GetTypeFloor() {
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
