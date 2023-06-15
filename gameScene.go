package main

import (
	"bytes"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameScene struct {
	catRunImage           *ebiten.Image
	catRunFrame           float64
	catRunFrameCount      float64
	catRunFrameSize       int
	catRunFramesPerSecond float64
}

func (me *GameScene) Initialize() {
	var catWalkImage, _, err = image.Decode(bytes.NewReader(catRun))
	me.catRunImage = ebiten.NewImageFromImage(catWalkImage)
	me.catRunFrameCount = 6
	me.catRunFrameSize = 48
	me.catRunFramesPerSecond = 6
	Assert(err)
}

func (me *GameScene) Update(deltaTime float64) {
	me.catRunFrame += deltaTime * me.catRunFramesPerSecond
	if me.catRunFrame >= me.catRunFrameCount {
		me.catRunFrame = 0
	}
}

func (me *GameScene) Draw(screen *ebiten.Image) {
	var drawOptions = ebiten.DrawImageOptions{}
	var shiftX = int(me.catRunFrame) * int(me.catRunFrameSize)
	var rect = image.Rect(shiftX, 0, shiftX+me.catRunFrameSize, me.catRunFrameSize)
	screen.DrawImage(me.catRunImage.SubImage(rect).(*ebiten.Image), &drawOptions)
}
