package main

import "github.com/hajimehoshi/ebiten/v2"

type GameSceneTransition struct {
	// Input parameter for initialization
	ViewWidth float64
	// Input parameter for initialization
	CatViewX float64
	// Input parameter for initialization
	FloorY float64
	// Input parameter for start of the scene
	CatRunFrame float64
	catRunImage *ebiten.Image
}

func (me *GameSceneTransition) Initialize() {
	me.catRunImage = LoadImage(CAT_RUN_IMAGE_BYTES)
}

func (me *GameSceneTransition) Update(deltaTime float64) {
	me.CatRunFrame += deltaTime * CAT_RUN_ANIMATION_FRAME_PER_SECOND
	for me.CatRunFrame >= CAT_RUN_ANIMATION_FRAME_COUNT {
		me.CatRunFrame -= CAT_RUN_ANIMATION_FRAME_COUNT
	}
}

func (me *GameSceneTransition) Draw(screen *ebiten.Image) {
	var drawOptions ebiten.DrawImageOptions
	ScaleCentered(&drawOptions, CAT_RUN_ANIMATION_FRAME_WIDTH, float64(me.catRunImage.Bounds().Dy()), -1, 1)
	drawOptions.GeoM.Translate(
		me.ViewWidth-CAT_RUN_ANIMATION_FRAME_WIDTH-me.CatViewX,
		me.FloorY-float64(me.catRunImage.Bounds().Dy()))
	var spriteShiftX = float64(int(me.CatRunFrame)) * CAT_RUN_ANIMATION_FRAME_WIDTH
	var rectangle = GetShiftedRectangle(spriteShiftX,
		CAT_RUN_ANIMATION_FRAME_WIDTH, float64(me.catRunImage.Bounds().Dy()))
	screen.DrawImage(me.catRunImage.SubImage(rectangle).(*ebiten.Image), &drawOptions)
}
