package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type FallObstacleMan struct {
	// Input parameter for initialization. Measurement unit: pixels
	AreaHeight float64
	// Input parameter for initialization. Measurement unit: pixels
	ViewWidth float64
	// Input parameter for initialization. Measurement unit: pixels
	ViewHeight float64
	// Input parameter for drawing
	CameraY       float64
	obstacles     []FloatPoint
	obstacleImage *ebiten.Image
}

func (me *FallObstacleMan) Initialize() {
	me.obstacleImage = LoadImage(OBSTACLE_IMAGE_BYTES)
	for y := me.ViewHeight; y < me.AreaHeight-me.ViewHeight; y += me.getDistanceBetweenObstacles() {
		var obstacle = FloatPoint{
			X: rand.Float64() * me.ViewWidth,
			Y: y + (rand.Float64()-0.5)*me.getFluctuationY(),
		}
		me.obstacles = append(me.obstacles, obstacle)
	}
}

func (me *FallObstacleMan) Update(deltaTime float64) {

}

func (me *FallObstacleMan) Draw(screen *ebiten.Image) {
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

func (me *FallObstacleMan) drawObstacle(screen *ebiten.Image, index int, obstacle FloatPoint) {
	var drawOptions ebiten.DrawImageOptions
	drawOptions.GeoM.Translate(obstacle.X, obstacle.Y-me.CameraY)
	screen.DrawImage(me.obstacleImage, &drawOptions)
}
