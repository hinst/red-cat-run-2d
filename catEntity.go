package main

import (
	"bytes"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type CatEntity struct {
	X     float64
	Y     float64
	Speed float64

	runImage          *ebiten.Image
	runFrame          float64
	runFramePerSecond float64
	runFrameCount     float64

	dieImage          *ebiten.Image
	dieFrame          float64
	dieFramePerSecond float64
	dieFrameCount     float64

	Width  float64
	Height float64
	Status int
	// Input parameter for every draw
	CameraX float64
	// Input parameter for every draw
	CameraY float64
}

func (me *CatEntity) Initialize() {
	var catWalkImage, _, catWalkImageError = image.Decode(bytes.NewReader(catRun))
	Assert(catWalkImageError)
	me.runImage = ebiten.NewImageFromImage(catWalkImage)
	me.runFramePerSecond = 6
	me.runFrameCount = 6

	var catDieImage, _, catDieImageError = image.Decode(bytes.NewReader(catDie))
	Assert(catDieImageError)
	me.dieImage = ebiten.NewImageFromImage(catDieImage)
	me.dieFramePerSecond = 6
	me.dieFrameCount = 4

	me.Width = 48
	me.Height = 25
	me.Speed = 40
}

func (me *CatEntity) Update(deltaTime float64) {
	if me.Status == me.GetStatusFloor() || me.Status == me.GetStatusCeiling() {
		me.runFrame += deltaTime * me.runFramePerSecond
		if me.runFrame >= me.runFrameCount {
			me.runFrame = 0
		}
	} else if me.Status == me.GetStatusDead() {
		me.dieFrame += deltaTime * me.dieFramePerSecond
		if me.dieFrame >= me.dieFrameCount {
			me.dieFrame -= me.dieFrameCount
		}
	}
	me.X += deltaTime * me.Speed
}

func (me *CatEntity) Draw(screen *ebiten.Image) {
	var drawOptions = ebiten.DrawImageOptions{}
	drawOptions.GeoM.Translate(me.X, me.Y)
	drawOptions.GeoM.Translate(-me.CameraX, -me.CameraY)
	if me.Status == me.GetStatusFloor() || me.Status == me.GetStatusCeiling() {
		var spriteShiftX = float64(int(me.runFrame)) * me.Width
		var rect = image.Rect(
			RoundFloat64ToInt(spriteShiftX), 0,
			RoundFloat64ToInt(spriteShiftX+me.Width), RoundFloat64ToInt(me.Width),
		)
		screen.DrawImage(me.runImage.SubImage(rect).(*ebiten.Image), &drawOptions)
	} else if me.Status == me.GetStatusDead() {
		var spriteShiftX = float64(int(me.dieFrame)) * me.Width
		var rect = image.Rect(
			RoundFloat64ToInt(spriteShiftX), 0,
			RoundFloat64ToInt(spriteShiftX+me.Width), RoundFloat64ToInt(me.Width),
		)
		screen.DrawImage(me.dieImage.SubImage(rect).(*ebiten.Image), &drawOptions)
	}
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
