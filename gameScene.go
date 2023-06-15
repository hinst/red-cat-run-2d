package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type GameScene struct {
	catEntity CatEntity
}

func (me *GameScene) Initialize() {
	me.catEntity.Initialize()
	me.catEntity.X = 10
	me.catEntity.Y = 200 - float64(me.catEntity.Height)
}

func (me *GameScene) Update(deltaTime float64) {
	me.catEntity.Update(deltaTime)
}

func (me *GameScene) Draw(screen *ebiten.Image) {
	me.catEntity.Draw(screen)
}