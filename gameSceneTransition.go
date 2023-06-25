package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameSceneTransition struct {
	// Input parameter for initialization
	ViewWidth  float64
	ViewHeight float64
	// Input parameter for initialization
	CatViewX float64
	// Input parameter for initialization
	CatSpeedX float64
	// Input parameter for initialization
	FloorY float64
	// Input parameter for start of the scene
	CatRunFrame float64
	// Input parameter for initialization
	PaddingWidth float64
	// Input parameter for initialization
	TorchGapY float64
	// Input parameter for initialization
	TorchSpeedY        float64
	catRunImage        *ebiten.Image
	angle              float64
	catSpeedY          float64
	catX               float64
	catY               float64
	torchY             float64
	torchImage         *ebiten.Image
	secondPhaseEnabled bool
	Complete           bool
}

func (me *GameSceneTransition) Initialize() {
	me.catRunImage = LoadImage(CAT_RUN_IMAGE_BYTES)
	me.catX = me.ViewWidth - CAT_RUN_ANIMATION_FRAME_WIDTH - me.CatViewX
	me.catY = me.FloorY - float64(me.catRunImage.Bounds().Dy())
	me.torchImage = LoadImage(TORCH_IMAGE_BYTES)
	me.torchY = me.ViewHeight
}

func (me *GameSceneTransition) Update(deltaTime float64) {
	me.CatRunFrame += deltaTime * CAT_RUN_ANIMATION_FRAME_PER_SECOND
	for me.CatRunFrame >= CAT_RUN_ANIMATION_FRAME_COUNT {
		me.CatRunFrame -= CAT_RUN_ANIMATION_FRAME_COUNT
	}
	if !me.secondPhaseEnabled {
		if me.angle > -math.Pi/2 {
			me.angle -= deltaTime
		}
		if me.catY < me.ViewHeight {
			me.catX -= me.CatSpeedX * deltaTime
			me.catSpeedY += deltaTime * 9.8 * 10
			me.catY += me.catSpeedY * deltaTime
		} else {
			me.secondPhaseEnabled = true
			me.catX = me.ViewWidth / 2
			me.catSpeedY = -100
		}
	} else {
		me.torchY -= me.TorchSpeedY * deltaTime
		me.torchY += me.catSpeedY * deltaTime
		me.catY += me.catSpeedY * deltaTime
	}
}

func (me *GameSceneTransition) Draw(screen *ebiten.Image) {
	if !me.secondPhaseEnabled {
		me.drawCatRun(screen)
	} else {
		me.drawCatRun(screen)
		me.drawTorches(screen)
	}
}

func (me *GameSceneTransition) drawCatRun(screen *ebiten.Image) {
	var drawOptions ebiten.DrawImageOptions
	ScaleCentered(&drawOptions, CAT_RUN_ANIMATION_FRAME_WIDTH, float64(me.catRunImage.Bounds().Dy()), -1, 1)
	RotateCentered(&drawOptions, CAT_RUN_ANIMATION_FRAME_WIDTH, float64(me.catRunImage.Bounds().Dy()), me.angle)
	drawOptions.GeoM.Translate(
		me.catX,
		me.catY)
	var spriteShiftX = float64(int(me.CatRunFrame)) * CAT_RUN_ANIMATION_FRAME_WIDTH
	var rectangle = GetShiftedRectangle(spriteShiftX,
		CAT_RUN_ANIMATION_FRAME_WIDTH, float64(me.catRunImage.Bounds().Dy()))
	screen.DrawImage(me.catRunImage.SubImage(rectangle).(*ebiten.Image), &drawOptions)
}

func (me *GameSceneTransition) drawTorches(screen *ebiten.Image) {
	for y := me.torchY; y < me.ViewHeight+float64(me.torchImage.Bounds().Dy()); y += me.PaddingWidth {
		var x = me.PaddingWidth / 2
		DrawTorch(screen, me.torchImage, x, y)
		x = me.ViewWidth - me.PaddingWidth/2
		DrawTorch(screen, me.torchImage, x, y)
	}
}
