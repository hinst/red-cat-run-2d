package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type CatEntityVertical struct {
	CatEntity
	flyImage              *ebiten.Image
	flyAnimationDirection float64
	flyAnimationFrame     float64
	// Input parameter for every update
	CameraY float64
}

func (me *CatEntityVertical) Initialize() {
	me.Width = 19
	me.Height = 48
	me.flyImage = LoadImage(CAT_FLY_DOWN_IMAGE_BYTES)
	me.flyAnimationDirection = 1
}

func (me *CatEntityVertical) Update(deltaTime float64) {
	me.flyAnimationFrame += deltaTime * CAT_FLY_ANIMATION_FRAME_PER_SECOND * me.flyAnimationDirection
	if CAT_FLY_ANIMATION_FRAME_COUNT <= me.flyAnimationFrame {
		me.flyAnimationFrame = CAT_FLY_ANIMATION_FRAME_COUNT - 1
		me.flyAnimationDirection = -1
	}
	if me.flyAnimationFrame <= 0 {
		me.flyAnimationFrame = 1
		me.flyAnimationDirection = 1
	}
	me.Y += me.GetSpeedY() * deltaTime
	me.updateSteer(deltaTime)
}

func (me *CatEntityVertical) updateSteer(deltaTime float64) {
	for _, key := range me.PressedKeys {
		if key == ebiten.KeyLeft {
			me.X -= deltaTime * 50
			break
		} else if key == ebiten.KeyRight {
			me.X += deltaTime * 50
		}
	}
}

func (me *CatEntityVertical) Draw(screen *ebiten.Image) {
	var drawOptions ebiten.DrawImageOptions
	drawOptions.GeoM.Translate(me.X, me.Y-me.CameraY)
	var spriteShiftX = float64(int(me.flyAnimationFrame)) * CAT_FLY_ANIMATION_FRAME_WIDTH
	var rectangle = GetShiftedRectangle(spriteShiftX, me.Width, me.Height)
	screen.DrawImage(me.flyImage.SubImage(rectangle).(*ebiten.Image), &drawOptions)
}

// Measurement unit: pixels per second
func (me *CatEntityVertical) GetSpeedY() float64 {
	return 25
}
