package main

import (
	"bytes"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type CatEntity struct {
	X          float64
	Y          float64
	Speed      float64
	Width      float64
	FrameWidth float64
	Height     float64
	Location   int
	Status     int
	// Input parameter for initialization
	ViewWidth float64
	// Input parameter for initialization
	ViewHeight float64
	// Input parameter for every draw
	CameraX float64
	// Input parameter for every draw
	CameraY          float64
	DebugModeEnabled bool

	runImage          *ebiten.Image
	runFrame          float64
	runFramePerSecond float64
	runFrameCount     float64

	dieImage          *ebiten.Image
	dieFrame          float64
	dieFramePerSecond float64
	dieFrameCount     float64
}

func (me *CatEntity) Initialize() {
	var catWalkImage, _, catWalkImageError = image.Decode(bytes.NewReader(catRun))
	AssertError(catWalkImageError)
	me.runImage = ebiten.NewImageFromImage(catWalkImage)
	me.runFramePerSecond = 6
	me.runFrameCount = 6

	var catDieImage, _, catDieImageError = image.Decode(bytes.NewReader(catDie))
	AssertError(catDieImageError)
	me.dieImage = ebiten.NewImageFromImage(catDieImage)
	me.dieFramePerSecond = 6
	me.dieFrameCount = 4

	me.Width = 40
	me.FrameWidth = 48
	me.Height = 25
	me.Speed = 40
}

func (me *CatEntity) Update(deltaTime float64) {
	if me.Status == me.GetStatusRun() {
		me.runFrame += deltaTime * me.runFramePerSecond
		if me.runFrame >= me.runFrameCount {
			me.runFrame -= me.runFrameCount
		}
	} else if me.Status == me.GetStatusDead() {
		me.dieFrame += deltaTime * me.dieFramePerSecond
		if me.dieFrame >= me.dieFrameCount {
			me.dieFrame = me.dieFrameCount - 1
		}
		if me.Location == me.GetLocationFloor() {
			if me.Y < me.ViewHeight {
				me.Y += deltaTime * me.Speed;
			}
		}
	}
	me.X += deltaTime * me.Speed
}

func (me *CatEntity) Draw(screen *ebiten.Image) {
	if me.DebugModeEnabled {
		vector.DrawFilledRect(screen, float32(me.X-me.CameraX), float32(me.Y-me.CameraY),
			float32(me.Width), float32(me.Height), color.RGBA{R: 100, G: 100, B: 100}, false)
	}
	var drawOptions = ebiten.DrawImageOptions{}
	drawOptions.GeoM.Translate(me.X, me.Y)
	drawOptions.GeoM.Translate(-me.CameraX, -me.CameraY)
	if me.Status == me.GetStatusRun() {
		var spriteShiftX = float64(int(me.runFrame)) * me.FrameWidth
		var rect = image.Rect(
			RoundFloat64ToInt(spriteShiftX), 0,
			RoundFloat64ToInt(spriteShiftX+me.FrameWidth), RoundFloat64ToInt(me.FrameWidth),
		)
		screen.DrawImage(me.runImage.SubImage(rect).(*ebiten.Image), &drawOptions)
	} else if me.Status == me.GetStatusDead() {
		var spriteShiftX = float64(int(me.dieFrame)) * me.FrameWidth
		var rect = image.Rect(
			RoundFloat64ToInt(spriteShiftX), 0,
			RoundFloat64ToInt(spriteShiftX+me.FrameWidth), RoundFloat64ToInt(me.FrameWidth),
		)
		screen.DrawImage(me.dieImage.SubImage(rect).(*ebiten.Image), &drawOptions)
	}
}

func (me *CatEntity) GetStatusRun() int {
	return 0
}

func (me *CatEntity) GetStatusJump() int {
	return 1
}

func (me *CatEntity) GetStatusDead() int {
	return 2
}

func (me *CatEntity) GetLocationFloor() int {
	return 0
}

func (me *CatEntity) GetLocationCeiling() int {
	return 1
}
