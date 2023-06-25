package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type FallObstacleMan struct {
	// Input parameter for initialization. Measurement unit: pixels
	AreaWidth float64
	// Input parameter for initialization. Measurement unit: pixels
	AreaHeight float64
	// Input parameter for initialization. Measurement unit: pixels
	ViewWidth float64
	// Input parameter for initialization. Measurement unit: pixels
	ViewHeight float64
	// Input parameter for every update
	CameraY          float64
	ObstacleWidth    float64
	obstacles        []FloatPoint
	obstacleImage    *ebiten.Image
	animationAngle   float64
	DebugModeEnabled bool
}

func (me *FallObstacleMan) Initialize() {
	me.obstacleImage = LoadImage(OBSTACLE_IMAGE_BYTES)
	me.ObstacleWidth = float64(me.obstacleImage.Bounds().Dx()) * 2
	var previousX float64
	for y := me.ViewHeight; y < me.AreaHeight-me.ViewHeight; y += me.getDistanceBetweenObstacles() {
		for i := 0; i < 2; i++ {
			var width = me.AreaWidth - me.ObstacleWidth - me.getPadding()*2
			var findX = func() float64 {
				return me.getShaftLeft() + me.getPadding() +
					me.ObstacleWidth/2 + rand.Float64()*width
			}
			var x = findX()
			for i := 0; i < 4; i++ {
				if math.Abs(x-previousX) < me.ObstacleWidth*1.5 {
					x = findX()
				} else {
					break
				}
			}
			var obstacle = FloatPoint{
				X: x,
				Y: y + (rand.Float64()-0.5)*me.getFluctuationY(),
			}
			me.obstacles = append(me.obstacles, obstacle)
			previousX = x
		}
	}
	me.DebugModeEnabled = false
}

func (me *FallObstacleMan) Update(deltaTime float64) {
	me.animationAngle += deltaTime * me.getAnimationAngleSpeed()
	me.animationAngle = UnwindAngle(me.animationAngle)
}

func (me *FallObstacleMan) Draw(screen *ebiten.Image) {
	me.drawShaftArea(screen)
	for index, obstacle := range me.obstacles {
		me.drawObstacle(screen, index, obstacle)
	}
}

func (me *FallObstacleMan) getDistanceBetweenObstacles() float64 {
	return 90
}

func (me *FallObstacleMan) getFluctuationY() float64 {
	return 30
}

func (me *FallObstacleMan) getAnimationAngleSpeed() float64 {
	return 0.5
}

func (me *FallObstacleMan) getShaftLeft() float64 {
	return me.ViewWidth/2 - me.AreaWidth/2
}

func (me *FallObstacleMan) getPadding() float64 {
	return 10
}

func (me *FallObstacleMan) drawObstacle(screen *ebiten.Image, index int, obstacle FloatPoint) {
	var drawOptions ebiten.DrawImageOptions
	RotateCentered(&drawOptions,
		float64(me.obstacleImage.Bounds().Dx()), float64(me.obstacleImage.Bounds().Dy()),
		UnwindAngle(me.animationAngle+float64(index)))
	drawOptions.GeoM.Scale(
		me.ObstacleWidth/float64(me.obstacleImage.Bounds().Dx()),
		me.ObstacleWidth/float64(me.obstacleImage.Bounds().Dy()))
	drawOptions.GeoM.Translate(obstacle.X-me.ObstacleWidth/2, obstacle.Y-me.CameraY-me.ObstacleWidth/2)
	screen.DrawImage(me.obstacleImage, &drawOptions)
	if me.DebugModeEnabled {
		var r = me.GetCollisionRectangle(obstacle)
		vector.DrawFilledRect(screen,
			float32(r.A.X), float32(r.A.Y-me.CameraY), float32(r.GetWidth()), float32(r.GetHeight()),
			color.White, false)
	}
}

func (me *FallObstacleMan) drawShaftArea(screen *ebiten.Image) {
	vector.DrawFilledRect(screen,
		float32(me.getShaftLeft()), 0, float32(me.AreaWidth), float32(me.ViewHeight),
		color.NRGBA{R: 255, G: 255, B: 255, A: 1}, false)
}

func (me *FallObstacleMan) GetCollisionRectangle(obstacle FloatPoint) (result Rectangle) {
	result = Rectangle{
		A: FloatPoint{
			X: obstacle.X - me.ObstacleWidth/2,
			Y: obstacle.Y - me.ObstacleWidth/2,
		},
	}
	result.B.X = result.A.X + me.ObstacleWidth
	result.B.Y = result.A.Y + me.ObstacleWidth
	return result.Shrink(5)
}

func (me *FallObstacleMan) CheckCollided(rectangle Rectangle) bool {
	for _, obstacle := range me.obstacles {
		if me.GetCollisionRectangle(obstacle).CheckCollides(rectangle) {
			return true
		}
	}
	return false
}
