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
	me.CreateObstacles()
	me.DebugModeEnabled = false
}

func (me *FallObstacleMan) CreateObstacles() {
	me.obstacles = me.obstacles[:0]
	var previousX float64
	var previousType int
	for y := me.ViewHeight; y < me.AreaHeight-me.ViewHeight; y += me.getDistanceBetweenObstacles() {
		for xIndex := 0; xIndex < 2; xIndex++ {
			var x float64
			if xIndex == 0 {
				var placementType = rand.Intn(3)
				for i := 0; i < 4 && placementType == previousType; i++ {
					placementType = rand.Intn(3)
				}
				switch placementType {
				case 0:
					x = me.getShaftLeft() + me.ObstacleWidth*0.77
				case 1:
					x = me.ViewWidth / 2
					if previousType == 0 {
						x += -1 - rand.Float64()
					} else if previousType == 2 {
						x += +1 + rand.Float64()
					} else {
						x += 0.5 + rand.Float64()
					}
				case 2:
					x = me.getShaftRight() - me.ObstacleWidth*0.77
				}
				previousType = placementType
			} else {
				if previousX < me.ViewWidth/2 {
					x = previousX + me.ObstacleWidth*1.5
				} else {
					x = previousX - me.ObstacleWidth*1.5
				}
			}
			var obstacle = FloatPoint{
				X: x,
				Y: y,
			}
			if xIndex == 1 {
				obstacle.Y += me.getFluctuationY()*0.5 + rand.Float64()*me.getFluctuationY()
			}
			if xIndex == 0 || rand.Intn(4) != 0 {
				me.obstacles = append(me.obstacles, obstacle)
			}
			previousX = x
		}
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
	return 144
}

func (me *FallObstacleMan) getFluctuationY() float64 {
	return 0
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
	var x = obstacle.X - me.ObstacleWidth/2
	var y = obstacle.Y - me.CameraY - me.ObstacleWidth/2
	var isVisible = -me.ObstacleWidth*1.6 < y && y < me.ViewHeight+me.ObstacleWidth
	if !isVisible {
		return
	}
	var drawOptions ebiten.DrawImageOptions
	RotateCentered(&drawOptions,
		float64(me.obstacleImage.Bounds().Dx()), float64(me.obstacleImage.Bounds().Dy()),
		UnwindAngle(me.animationAngle+float64(index)))
	drawOptions.GeoM.Scale(
		me.ObstacleWidth/float64(me.obstacleImage.Bounds().Dx()),
		me.ObstacleWidth/float64(me.obstacleImage.Bounds().Dy()))
	drawOptions.GeoM.Translate(x, y)
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
			Y: obstacle.Y - me.ObstacleWidth/4,
		},
	}
	result.B.X = result.A.X + me.ObstacleWidth
	result.B.Y = result.A.Y + me.ObstacleWidth/2
	return result.Shrink(1)
}

func (me *FallObstacleMan) CheckCollided(rectangle Rectangle) bool {
	for _, obstacle := range me.obstacles {
		if me.GetCollisionRectangle(obstacle).CheckCollides(rectangle) {
			return true
		}
	}
	return false
}
