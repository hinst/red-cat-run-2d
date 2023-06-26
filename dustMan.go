package main

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type DustMan struct {
	// Input parameter for initialization. Measurement unit: pixels
	ViewWidth float64
	// Input parameter for initialization. Measurement unit: pixels
	ViewHeight float64
	// Input parameter for initialization. Measurement unit: pixels
	AreaWidth float64
	// Input parameter for initialization. Measurement unit: pixels
	AreaHeight float64
	// Input parameter for initialization
	Direction Direction
	// Input parameter for update. Measurement unit: pixels
	CameraX float64
	// Input parameter for update. Measurement unit: pixels
	CameraY float64

	particles []FloatPoint
}

const DUST_MAN_DISTANCE_BETWEEN_PARTICLES = 50

func (me *DustMan) Initialize() {
	if me.Direction == DIRECTION_RIGHT {
		for i := -me.ViewWidth; i < me.AreaWidth+me.ViewWidth; i += DUST_MAN_DISTANCE_BETWEEN_PARTICLES {
			var y = rand.Intn(RoundFloat64ToInt(me.ViewHeight))
			var x = i - 5 + float64(rand.Intn(10))
			var particle = FloatPoint{X: float64(x), Y: float64(y)}
			me.particles = append(me.particles, particle)
		}
	} else if me.Direction == DIRECTION_BOTTOM {
		for i := -me.ViewHeight; i < me.AreaHeight+me.ViewHeight; i += DUST_MAN_DISTANCE_BETWEEN_PARTICLES {
			var x = (me.ViewWidth-me.AreaWidth)/2 + float64(rand.Intn(RoundFloat64ToInt(me.AreaWidth)))
			var y = i - 5 + float64(rand.Intn(10))
			var particle = FloatPoint{X: float64(x), Y: float64(y)}
			me.particles = append(me.particles, particle)
		}
	}
}

func (me *DustMan) Update(deltaTime float64) {
}

func (me *DustMan) Draw(screen *ebiten.Image) {
	for _, particle := range me.particles {
		if me.checkParticleVisible(particle) {
			vector.DrawFilledRect(screen,
				float32(particle.X-me.CameraX), float32(particle.Y-me.CameraY), 1, 1,
				color.NRGBA{R: 255, G: 255, B: 255, A: 255 / 3 * 2}, true)
		}
	}
}

func (me *DustMan) checkParticleVisible(particle FloatPoint) bool {
	return me.CameraX-1 <= particle.X && particle.X <= me.CameraX+me.ViewWidth+1 &&
		me.CameraY-1 <= particle.Y && particle.Y <= me.CameraY+me.ViewHeight+1
}
