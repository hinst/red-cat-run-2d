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
	// Input parameter for update
	CameraX float64

	particles []FloatPoint
}

const DUST_MAN_DISTANCE_BETWEEN_PARTICLES = 50

func (me *DustMan) Initialize() {
	for i := float64(0); i < me.AreaWidth; i += DUST_MAN_DISTANCE_BETWEEN_PARTICLES {
		var y = rand.Intn(RoundFloat64ToInt(me.ViewHeight))
		var x = i - 5 + float64(rand.Intn(10))
		var particle = FloatPoint{X: x, Y: float64(y)}
		me.particles = append(me.particles, particle)
	}
}

func (me *DustMan) Update(deltaTime float64) {
}

func (me *DustMan) Draw(screen *ebiten.Image) {
	for _, particle := range me.particles {
		if me.checkParticleVisible(particle) {
			vector.DrawFilledRect(screen, float32(particle.X-me.CameraX), float32(particle.Y), 1, 1,
				color.NRGBA{R: 255, G: 255, B: 255, A: 255 / 3 * 2}, true)
		}
	}
}

func (me *DustMan) checkParticleVisible(particle FloatPoint) bool {
	return me.CameraX-1 <= particle.X && particle.X <= me.CameraX+me.ViewWidth+1
}
