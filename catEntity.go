package main

import (
	"bytes"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type CatEntity struct {
	X float64
	Y float64

	runImage          *ebiten.Image
	runFrame          float64
	runFramePerSecond float64
	runFrameCount     float64

	dieImage      *ebiten.Image
	dieFrameCount float64

	Width  int
	Height int
	Status int
}

func (me *CatEntity) Initialize() {
	var catWalkImage, _, catWalkImageError = image.Decode(bytes.NewReader(catRun))
	Assert(catWalkImageError)
	me.runImage = ebiten.NewImageFromImage(catWalkImage)
	me.runFrameCount = 6

	var catDieImage, _, catDieImageError = image.Decode(bytes.NewReader(catDie))
	Assert(catDieImageError)
	me.dieImage = ebiten.NewImageFromImage(catDieImage)
	me.dieFrameCount = 4

	me.Width = 48
	me.Height = 25
	me.runFramePerSecond = 6
}

func (me *CatEntity) Update(deltaTime float64) {
	me.runFrame += deltaTime * me.runFramePerSecond
	if me.runFrame >= me.runFrameCount {
		me.runFrame = 0
	}
}

func (me *CatEntity) Draw(screen *ebiten.Image) {
	var drawOptions = ebiten.DrawImageOptions{}
	drawOptions.GeoM.Translate(me.X, me.Y)
	var shiftX = int(me.runFrame) * int(me.Width)
	var rect = image.Rect(shiftX, 0, shiftX+me.Width, me.Width)
	screen.DrawImage(me.runImage.SubImage(rect).(*ebiten.Image), &drawOptions)
}

func (me *CatEntity) GetStatusFloor() int {
	return 0
}

func (me *CatEntity) GetStatusCeiling() int {
	return 1
}

func (me *CatEntity) GetStatusFly() int {
	return 2
}

func (me *CatEntity) GetStatusDead() int {
	return 3
}
