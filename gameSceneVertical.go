package main

import (
	"log"

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
	TorchY          float64
	torchImage      *ebiten.Image
	brickImage      *ebiten.Image
	dirtImage       *ebiten.Image
	wallAlpha       float64
}

func (me *GameSceneVertical) Initialize() {
	if me.ViewHeight == 0 || me.ViewWidth == 0 {
		log.Println("Warning: view size is missing")
	}
	me.catEntity.Initialize()
	me.cameraY = me.catEntity.Y - 10
	me.catEntity.CameraY = me.cameraY
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

func (me *GameSceneVertical) Update(deltaTime float64) {
	me.catEntity.PressedKeys = me.PressedKeys
	me.catEntity.JustPressedKeys = me.JustPressedKeys
	me.catEntity.CameraY = me.cameraY
	me.catEntity.Update(deltaTime)
	if me.catEntity.X < me.GetPaddingWidth() {
		me.catEntity.X = me.GetPaddingWidth()
	}
	if me.catEntity.X >= me.ViewWidth-me.GetPaddingWidth()-me.catEntity.Width {
		me.catEntity.X = me.ViewWidth - me.GetPaddingWidth() - me.catEntity.Width
	}
	me.fallObstacleMan.CameraY = me.cameraY
	me.fallObstacleMan.Update(deltaTime)
	me.cameraY = me.catEntity.Y - me.GetCatViewY()
	me.TorchY -= deltaTime * me.GetTorchSpeedY()
	for me.TorchY < -me.GetTorchGapY() {
		me.TorchY += me.GetTorchGapY()
	}
	if me.wallAlpha < 1 {
		me.wallAlpha += deltaTime * me.getWallAlphaSpeed()
		if me.wallAlpha >= 1 {
			me.wallAlpha = 1
		}
	}
}

func (me *GameSceneVertical) Draw(screen *ebiten.Image) {
	me.drawDecorations(screen)
	me.catEntity.Draw(screen)
	me.fallObstacleMan.Draw(screen)
}

func (me *GameSceneVertical) GetAreaWidth() float64 {
	return 220
}

func (me *GameSceneVertical) GetTorchSpeedY() float64 {
	return 100
}

func (me *GameSceneVertical) GetTorchGapY() float64 {
	return 200
}

func (me *GameSceneVertical) GetPaddingWidth() float64 {
	return (me.ViewWidth - me.GetAreaWidth()) / 2
}

func (me *GameSceneVertical) getTorchScale() float64 {
	return 0.5
}

func (me *GameSceneVertical) GetCatViewY() float64 {
	return 10
}

func (me *GameSceneVertical) drawDecorations(screen *ebiten.Image) {
	me.drawShaftBackground(screen)
	for y := me.TorchY - me.GetTorchGapY(); y < me.ViewHeight+me.GetTorchGapY(); y += me.GetTorchGapY() {
		me.drawTorchPair(screen, y)
		me.drawFloors(screen, y)
	}
}

func (me *GameSceneVertical) drawTorchPair(screen *ebiten.Image, y float64) {
	var x = me.GetPaddingWidth() / 2
	me.drawTorch(screen, x, y)
	x = me.ViewWidth - me.GetPaddingWidth()/2
	me.drawTorch(screen, x, y)
}

func (me *GameSceneVertical) drawTorch(screen *ebiten.Image, x float64, y float64) {
	DrawTorch(screen, me.torchImage, x, y)
}

func (me *GameSceneVertical) drawFloors(screen *ebiten.Image, y float64) {
	var brickImageWidth = float64(me.brickImage.Bounds().Dx())
	for x := float64(0); x <= me.GetPaddingWidth()-brickImageWidth; x += brickImageWidth {
		me.drawFloorPart(screen, x, y)
		me.drawFloorPart(screen, me.ViewWidth-x-brickImageWidth, y)
	}
}

func (me *GameSceneVertical) drawFloorPart(screen *ebiten.Image, baseX float64, baseY float64) {
	var y = baseY + float64(me.brickImage.Bounds().Dy())*3 + float64(me.torchImage.Bounds().Dy())*me.getTorchScale()
	var drawOptions ebiten.DrawImageOptions
	drawOptions.GeoM.Translate(baseX, y)
	drawOptions.ColorScale.Scale(float32(me.wallAlpha), float32(me.wallAlpha), float32(me.wallAlpha), float32(me.wallAlpha))
	screen.DrawImage(me.brickImage, &drawOptions)
	for dirtIndex := 0; dirtIndex < 10; dirtIndex++ {
		y += float64(me.brickImage.Bounds().Dy())
		var drawOptions ebiten.DrawImageOptions
		drawOptions.GeoM.Translate(baseX, y)
		drawOptions.ColorScale.Scale(float32(me.wallAlpha), float32(me.wallAlpha), float32(me.wallAlpha), float32(me.wallAlpha))
		screen.DrawImage(me.dirtImage, &drawOptions)
	}
}

func (me *GameSceneVertical) drawShaftBackground(screen *ebiten.Image) {
	var color = MultiplyColor(SHAFT_COLOR, me.wallAlpha)
	var width = me.GetPaddingWidth()
	vector.DrawFilledRect(screen, 0, 0, float32(width), float32(me.ViewHeight), color, false)
	vector.DrawFilledRect(screen, float32(me.ViewWidth)-float32(width), 0, float32(width), float32(me.ViewHeight), color, false)
}

func (me *GameSceneVertical) getWallAlphaSpeed() float64 {
	return 0.3
}
