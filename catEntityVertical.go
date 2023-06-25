package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type CatEntityVertical struct {
	CatEntity
	// Input parameter for every update
	CameraY float64
	// Input parameter for every update
	Collided              bool
	flyImage              *ebiten.Image
	flyAnimationDirection float64
	flyAnimationFrame     float64
	DebugModeEnabled      bool
	angle                 float64
	collidedY             float64
}

func (me *CatEntityVertical) Initialize() {
	me.Width = 19
	me.Height = 48
	me.flyImage = LoadImage(CAT_FLY_DOWN_IMAGE_BYTES)
	me.flyAnimationDirection = 1
	me.DebugModeEnabled = true
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
	if me.Collided {
		me.collidedY += deltaTime * me.getCollidedSpeedY()
		me.angle = UnwindAngle(me.angle + deltaTime*me.getCollidedRotationSpeed())
	} else {
		me.updateSteer(deltaTime)
		me.Y += me.GetSpeedY() * deltaTime
	}
}

func (me *CatEntityVertical) updateSteer(deltaTime float64) {
	for _, key := range me.PressedKeys {
		if key == ebiten.KeyLeft {
			me.X -= deltaTime * me.GetSteerSpeed()
			break
		} else if key == ebiten.KeyRight {
			me.X += deltaTime * me.GetSteerSpeed()
		}
	}
}

func (me *CatEntityVertical) Draw(screen *ebiten.Image) {
	var drawOptions ebiten.DrawImageOptions
	RotateCentered(&drawOptions, CAT_FLY_ANIMATION_FRAME_WIDTH, float64(me.flyImage.Bounds().Dy()), me.angle)
	drawOptions.GeoM.Translate(me.X, me.Y-me.CameraY+me.collidedY)
	var spriteShiftX = float64(int(me.flyAnimationFrame)) * CAT_FLY_ANIMATION_FRAME_WIDTH
	var rectangle = GetShiftedRectangle(spriteShiftX, me.Width, me.Height)
	if me.DebugModeEnabled {
		var box = me.GetHitBox()
		vector.DrawFilledRect(screen,
			float32(box.A.X), float32(box.A.Y-me.CameraY), float32(box.GetWidth()), float32(box.GetHeight()),
			color.NRGBA{R: 255, G: 255, B: 255, A: 127}, true)
	}
	screen.DrawImage(me.flyImage.SubImage(rectangle).(*ebiten.Image), &drawOptions)
}

// Measurement unit: pixels per second
func (me *CatEntityVertical) GetSpeedY() float64 {
	return 50
}

func (me *CatEntityVertical) GetSteerSpeed() float64 {
	return 80
}

func (me *CatEntityVertical) GetHitBox() Rectangle {
	var rect = Rectangle{
		A: FloatPoint{
			X: me.X,
			Y: me.Y,
		},
	}
	rect.B.X = rect.A.X + me.Width
	rect.B.Y = rect.A.Y + me.Height
	return rect.Shrink(1)
}

func (me *CatEntityVertical) getCollidedSpeedY() float64 {
	return -50
}

func (me *CatEntityVertical) getCollidedRotationSpeed() float64 {
	return 7
}
