package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type GameScene struct {
	terrainMan TerrainMan
	catEntity  CatEntity
	// Initialization input parameter
	ViewHeight float64
	// Initialization input parameter
	ViewWidth float64
	CameraX   float64
	CameraY   float64
}

func (me *GameScene) Initialize() {
	me.terrainMan.ViewHeight = me.ViewHeight
	me.terrainMan.ViewWidth = me.ViewWidth
	me.terrainMan.AreaWidth = 100
	me.terrainMan.Initialize()
	me.catEntity.Initialize()
	me.catEntity.X = me.GetCatViewX()
	me.catEntity.Y = me.GetCatY() - float64(me.catEntity.Height)
}

func (me *GameScene) Update(deltaTime float64) {
	me.catEntity.Update(deltaTime)
	me.CameraX = me.catEntity.X - me.GetCatViewX()
}

func (me *GameScene) Draw(screen *ebiten.Image) {
	me.terrainMan.CameraX = me.CameraX
	me.terrainMan.CameraY = me.CameraY
	me.terrainMan.Draw(screen)
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

func (me *GameScene) CheckCatFall() {
}
