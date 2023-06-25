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
	TorchSpeedY float64
	// Input parameter for initialization
	CatViewY              float64
	catRunImage           *ebiten.Image
	catFlyImage           *ebiten.Image
	flyAnimationFrame     float64
	flyAnimationDirection float64
	angle                 float64
	catSpeedY             float64
	catX                  float64
	catY                  float64
	TorchY                float64
	torchImage            *ebiten.Image
	secondPhaseEnabled    bool
	Complete              bool
}

func (me *GameSceneTransition) Initialize() {
	me.catRunImage = LoadImage(CAT_RUN_IMAGE_BYTES)
	me.catFlyImage = LoadImage(CAT_FLY_DOWN_IMAGE_BYTES)
	me.catX = me.ViewWidth - CAT_RUN_ANIMATION_FRAME_WIDTH - me.CatViewX
	me.catY = me.FloorY - float64(me.catRunImage.Bounds().Dy())
	me.torchImage = LoadImage(TORCH_IMAGE_BYTES)
	me.TorchY = me.ViewHeight * me.getInitialTorchFallScaleY()
}

func (me *GameSceneTransition) Update(deltaTime float64) {
	if !me.secondPhaseEnabled {
		me.updateFirstPhase(deltaTime)
	} else {
		me.updateSecondPhase(deltaTime)
	}
}

func (me *GameSceneTransition) updateFirstPhase(deltaTime float64) {
	me.CatRunFrame += deltaTime * CAT_RUN_ANIMATION_FRAME_PER_SECOND
	for me.CatRunFrame >= CAT_RUN_ANIMATION_FRAME_COUNT {
		me.CatRunFrame -= CAT_RUN_ANIMATION_FRAME_COUNT
	}
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
		me.catY = me.ViewHeight * me.getInitialCatFallScaleY()
		me.catSpeedY = -100
	}
}

func (me *GameSceneTransition) updateSecondPhase(deltaTime float64) {
	me.flyAnimationFrame += deltaTime * CAT_FLY_ANIMATION_FRAME_PER_SECOND * me.flyAnimationDirection
	if CAT_FLY_ANIMATION_FRAME_COUNT <= me.flyAnimationFrame {
		me.flyAnimationFrame = CAT_FLY_ANIMATION_FRAME_COUNT - 1
		me.flyAnimationDirection = -1
	}
	if me.flyAnimationFrame <= 0 {
		me.flyAnimationFrame = 1
		me.flyAnimationDirection = 1
	}

	me.TorchY -= me.TorchSpeedY * deltaTime
	me.TorchY += me.catSpeedY * deltaTime
	me.catY += me.catSpeedY * deltaTime
	if me.catY <= me.CatViewY {
		me.Complete = true
		me.catY = me.CatViewY
	}
}

func (me *GameSceneTransition) Draw(screen *ebiten.Image) {
	if !me.secondPhaseEnabled {
		me.drawCatRun(screen)
	} else {
		me.drawCatFly(screen)
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

func (me *GameSceneTransition) drawCatFly(screen *ebiten.Image) {
	println(int(me.catX))
	var drawOptions ebiten.DrawImageOptions
	drawOptions.GeoM.Translate(
		me.catX-float64(CAT_FLY_ANIMATION_FRAME_WIDTH)/2,
		me.catY)
	var spriteShiftX = float64(int(me.flyAnimationFrame)) * CAT_FLY_ANIMATION_FRAME_WIDTH
	var rectangle = GetShiftedRectangle(spriteShiftX,
		CAT_FLY_ANIMATION_FRAME_WIDTH, float64(me.catFlyImage.Bounds().Dy()))
	screen.DrawImage(me.catFlyImage.SubImage(rectangle).(*ebiten.Image), &drawOptions)
}

func (me *GameSceneTransition) drawTorches(screen *ebiten.Image) {
	for y := me.TorchY; y < me.ViewHeight+float64(me.torchImage.Bounds().Dy()); y += me.TorchGapY {
		var x = me.PaddingWidth / 2
		DrawTorch(screen, me.torchImage, x, y)
		x = me.ViewWidth - me.PaddingWidth/2
		DrawTorch(screen, me.torchImage, x, y)
	}
}

func (me *GameSceneTransition) getInitialCatFallScaleY() float64 {
	return 2.33
}

func (me *GameSceneTransition) getInitialTorchFallScaleY() float64 {
	return 2
}
