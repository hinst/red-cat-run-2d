package main

import "github.com/hajimehoshi/ebiten/v2"

type CatEntityVertical struct {
	CatEntity
	flyImage *ebiten.Image
}

func (me *CatEntityVertical) Initialize() {
	me.Width = 19
	me.Height = 48
	me.flyImage = LoadImage(CAT_FLY_DOWN_IMAGE_BYTES)
}

func (me *CatEntityVertical) Update(deltaTime float64) {
}

func (me *CatEntityVertical) Draw(screen *ebiten.Image) {
	var drawOptions ebiten.DrawImageOptions
	drawOptions.GeoM.Translate(me.X, me.Y)
	var rectangle = GetShiftedRectangle(0, me.Width, me.Height)
	screen.DrawImage(me.flyImage.SubImage(rectangle).(*ebiten.Image), &drawOptions)
}
