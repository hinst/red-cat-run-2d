package main

import "github.com/hajimehoshi/ebiten/v2"

type CatEntityVertical struct {
	CatEntity
}

func (me *CatEntityVertical) Initialize() {
	me.Width = 20
	me.Height = 40
}

func (me *CatEntityVertical) Draw(screen *ebiten.Image) {

}
