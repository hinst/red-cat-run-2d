package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type GameScene struct {
	terrainMan TerrainMan
	catEntity  CatEntity
}

func (me *GameScene) Initialize() {
	me.terrainMan.ViewHeight = 240
	me.terrainMan.ViewWidth = 320
	me.terrainMan.AreaWidth = 100
	me.terrainMan.Initialize()
	me.catEntity.Initialize()
	me.catEntity.X = 10
	me.catEntity.Y = 200 - float64(me.catEntity.Height)
}

func (me *GameScene) Update(deltaTime float64) {
	me.catEntity.Update(deltaTime)
}

func (me *GameScene) Draw(screen *ebiten.Image) {
	me.terrainMan.Draw(screen)
	me.catEntity.Draw(screen)
}
