package main

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type GameSceneVertical struct {
	// Initialization input parameter
	ViewHeight float64
	// Initialization input parameter
	ViewWidth float64
	// Input parameter for every update
	JustPressedKeys []ebiten.Key
	// Input parameter for every update
	PressedKeys []ebiten.Key

	catEntity       CatEntityVertical
	fallObstacleMan FallObstacleMan
	cameraY         float64
	torchY          float64
	torchImage      *ebiten.Image
	brickImage      *ebiten.Image
	dirtImage       *ebiten.Image
}

func (me *GameSceneVertical) Initialize() {
	if me.ViewHeight == 0 || me.ViewWidth == 0 {
		log.Println("Warning: view size is missing")
	}
	me.catEntity.Initialize()
	me.cameraY = me.catEntity.Y - 10
	me.catEntity.X = me.ViewWidth/2 - me.catEntity.Width/2
	me.fallObstacleMan.AreaWidth = me.GetAreaWidth()
	me.fallObstacleMan.AreaHeight = me.ViewHeight * 10
	me.fallObstacleMan.ViewWidth = me.ViewWidth
	me.fallObstacleMan.ViewHeight = me.ViewHeight
	me.fallObstacleMan.Initialize()
	me.torchImage = LoadImage(TORCH_IMAGE_BYTES)
	me.brickImage = LoadImage(BRICK_BLOCK_IMAGE_BYTES)
	me.dirtImage = LoadImage(DIRT_BLOCK_IMAGE_BYTES)
}

func (me *GameSceneVertical) Draw(screen *ebiten.Image) {
	me.drawDecorations(screen)
	me.catEntity.Draw(screen)
	me.fallObstacleMan.Draw(screen)
}

func (me *GameSceneVertical) Update(deltaTime float64) {
	me.catEntity.CameraY = me.cameraY
	me.catEntity.Update(deltaTime)
	me.fallObstacleMan.CameraY = me.cameraY
	me.fallObstacleMan.Update(deltaTime)
	me.cameraY = me.catEntity.Y - 10
	me.torchY -= math.Round(deltaTime * me.getTorchSpeed())
	for me.torchY < -me.getTorchGapY() {
		me.torchY += me.getTorchGapY()
	}
}

func (me *GameSceneVertical) GetAreaWidth() float64 {
	return 220
}

func (me *GameSceneVertical) getTorchSpeed() float64 {
	return 100
}

func (me *GameSceneVertical) getTorchGapY() float64 {
	return 200
}

func (me *GameSceneVertical) getPaddingWidth() float64 {
	return (me.ViewWidth - me.GetAreaWidth()) / 2
}

func (me *GameSceneVertical) getTorchScale() float64 {
	return 0.5
}

func (me *GameSceneVertical) drawDecorations(screen *ebiten.Image) {
	me.drawShaftBackground(screen)
	for y := me.torchY - me.getTorchGapY(); y < me.ViewHeight+me.getTorchGapY(); y += me.getTorchGapY() {
		me.drawTorch(screen, y)
		me.drawFloors(screen, y)
	}
}

func (me *GameSceneVertical) drawTorch(screen *ebiten.Image, y float64) {
	var torchScale = me.getTorchScale()
	{
		var drawOptions = ebiten.DrawImageOptions{}
		drawOptions.GeoM.Scale(torchScale, torchScale)
		var x = me.getPaddingWidth() / 2
		DrawTorchLight(screen, float32(x), float32(y))
		drawOptions.GeoM.Translate(
			x-float64(me.torchImage.Bounds().Dx())/2*torchScale,
			y-float64(me.torchImage.Bounds().Dy())/2*torchScale,
		)
		screen.DrawImage(me.torchImage, &drawOptions)
	}
	{
		var drawOptions = ebiten.DrawImageOptions{}
		drawOptions.GeoM.Scale(torchScale, torchScale)
		var x = me.ViewWidth - me.getPaddingWidth()/2
		DrawTorchLight(screen, float32(x), float32(y))
		drawOptions.GeoM.Translate(
			x-float64(me.torchImage.Bounds().Dx())/2*torchScale,
			y-float64(me.torchImage.Bounds().Dy())/2*torchScale,
		)
		screen.DrawImage(me.torchImage, &drawOptions)
	}
}

func (me *GameSceneVertical) drawFloors(screen *ebiten.Image, y float64) {
	var brickImageWidth = float64(me.brickImage.Bounds().Dx())
	for x := float64(0); x < me.getPaddingWidth()-brickImageWidth; x += brickImageWidth {
		me.drawFloorPart(screen, x, y)
		me.drawFloorPart(screen, me.ViewWidth-x-brickImageWidth, y)
	}
}

func (me *GameSceneVertical) drawFloorPart(screen *ebiten.Image, baseX float64, baseY float64) {
	var y = baseY + float64(me.brickImage.Bounds().Dy()) + float64(me.torchImage.Bounds().Dy())*me.getTorchScale()
	var drawOptions ebiten.DrawImageOptions
	drawOptions.GeoM.Translate(baseX, y)
	screen.DrawImage(me.brickImage, &drawOptions)
	for dirtIndex := 0; dirtIndex < 10; dirtIndex++ {
		y += float64(me.brickImage.Bounds().Dy())
		var drawOptions ebiten.DrawImageOptions
		drawOptions.GeoM.Translate(baseX, y)
		screen.DrawImage(me.dirtImage, &drawOptions)
	}
}

func (me *GameSceneVertical) drawShaftBackground(screen *ebiten.Image) {
	var width = (int(me.getPaddingWidth())/me.brickImage.Bounds().Dx() - 1) * me.brickImage.Bounds().Dx()
	vector.DrawFilledRect(screen, 0, 0, float32(width), float32(me.ViewHeight), SHAFT_COLOR, false)
	vector.DrawFilledRect(screen, float32(me.ViewWidth)-float32(width), 0, float32(width), float32(me.ViewHeight), SHAFT_COLOR, false)
}
