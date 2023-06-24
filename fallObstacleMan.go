package main

import (
	"image/color"
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
	CameraY        float64
	ObstacleWidth  float64
	obstacles      []FloatPoint
	obstacleImage  *ebiten.Image
	animationAngle float64
}

func (me *FallObstacleMan) Initialize() {
	me.obstacleImage = LoadImage(OBSTACLE_IMAGE_BYTES)
	me.ObstacleWidth = float64(me.obstacleImage.Bounds().Dx())
	for y := me.ViewHeight; y < me.AreaHeight-me.ViewHeight; y += me.getDistanceBetweenObstacles() {
		var obstacle = FloatPoint{
			X: me.getShaftLeft() + me.ObstacleWidth/2 + rand.Float64()*(me.AreaWidth-me.ObstacleWidth),
			Y: y + (rand.Float64()-0.5)*me.getFluctuationY(),
		}
		me.obstacles = append(me.obstacles, obstacle)
	}
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
	return 100
}

func (me *FallObstacleMan) getFluctuationY() float64 {
	return 40
}

func (me *FallObstacleMan) getAnimationAngleSpeed() float64 {
	return 0.5
}

func (me *FallObstacleMan) getShaftLeft() float64 {
	return me.ViewWidth/2 - me.AreaWidth/2
}

func (me *FallObstacleMan) getShaftRight() float64 {
	return me.ViewWidth/2 + me.AreaWidth/2
}

func (me *FallObstacleMan) drawObstacle(screen *ebiten.Image, index int, obstacle FloatPoint) {
	var drawOptions ebiten.DrawImageOptions
	RotateCentered(&drawOptions,
		float64(me.obstacleImage.Bounds().Dx()), float64(me.obstacleImage.Bounds().Dy()),
		UnwindAngle(me.animationAngle+float64(index)))
	drawOptions.GeoM.Translate(obstacle.X-me.ObstacleWidth/2, obstacle.Y-me.CameraY-me.ObstacleWidth/2)
	screen.DrawImage(me.obstacleImage, &drawOptions)
}

func (me *FallObstacleMan) drawShaftArea(screen *ebiten.Image) {
	vector.DrawFilledRect(screen,
		float32(me.getShaftLeft()), 0, float32(me.AreaWidth), float32(me.ViewHeight),
		color.NRGBA{R: 255, G: 255, B: 255, A: 1}, false)
}
