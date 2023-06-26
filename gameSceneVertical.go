package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	// Output parameter
	Completed bool

	catEntity        CatEntityVertical
	obstacleMan      FallObstacleMan
	dustMan          DustMan
	cameraY          float64
	TorchY           float64
	torchImage       *ebiten.Image
	brickImage       *ebiten.Image
	dirtImage        *ebiten.Image
	wallAlpha        float64
	dead             bool
	deadMessageDelay float64
	fishImage        *ebiten.Image
	Ascended         bool
	ascendedImage    *ebiten.Image
	ascendedAngle    float64
	ascendedPulse    float64
}

func (me *GameSceneVertical) Initialize() {
	if me.ViewHeight == 0 || me.ViewWidth == 0 {
		log.Println("Warning: view size is missing")
	}
	me.catEntity.Initialize()
	me.cameraY = me.catEntity.Y - 10
	me.catEntity.CameraY = me.cameraY
	me.catEntity.X = me.ViewWidth/2 - me.catEntity.Width/2
	me.obstacleMan.AreaWidth = me.GetAreaWidth()
	me.obstacleMan.AreaHeight = me.GetAreaHeight()
	me.obstacleMan.ViewWidth = me.ViewWidth
	me.obstacleMan.ViewHeight = me.ViewHeight
	me.obstacleMan.Initialize()
	me.torchImage = LoadImage(TORCH_IMAGE_BYTES)
	me.brickImage = LoadImage(BRICK_BLOCK_IMAGE_BYTES)
	me.dirtImage = LoadImage(DIRT_BLOCK_IMAGE_BYTES)
	me.fishImage = LoadImage(FISH_IMAGE_BYTES)
	me.ascendedImage = LoadImage(ASCENDED_IMAGE_BYTES)
	me.dustMan.ViewWidth = me.ViewWidth
	me.dustMan.ViewHeight = me.ViewHeight
	me.dustMan.AreaWidth = me.GetAreaWidth()
	me.dustMan.AreaHeight = me.GetAreaHeight()
	me.dustMan.Direction = DIRECTION_BOTTOM
	me.dustMan.Initialize()
}

func (me *GameSceneVertical) Update(deltaTime float64) {
	me.updateCatEntity(deltaTime)
	if me.catEntity.Direction == DIRECTION_BOTTOM {
		me.cameraY = me.catEntity.Y - me.GetCatViewY()
	} else if me.catEntity.Direction == DIRECTION_TOP {
		me.cameraY = me.catEntity.Y - me.ViewHeight + me.catEntity.Height + me.GetCatViewY()
	}
	if me.cameraY < me.GetAreaHeight()-me.ViewHeight {
		me.TorchY -= deltaTime * me.GetTorchSpeedY()
	} else {
		me.cameraY = me.GetAreaHeight() - me.ViewHeight
	}
	me.obstacleMan.CameraY = me.cameraY
	me.obstacleMan.Update(deltaTime)
	for me.TorchY < -me.GetTorchGapY() {
		me.TorchY += me.GetTorchGapY()
	}
	for me.TorchY > me.GetTorchGapY() {
		me.TorchY -= me.GetTorchGapY()
	}
	if me.wallAlpha < 1 {
		me.wallAlpha += deltaTime * me.getWallAlphaSpeed()
		if me.wallAlpha >= 1 {
			me.wallAlpha = 1
		}
	}
	if me.dead && len(me.JustPressedKeys) > 0 && me.deadMessageDelay <= 0 {
		me.Completed = true
	}
	if me.Ascended {
		for _, key := range me.JustPressedKeys {
			if key == ebiten.KeyEnter {
				me.Completed = true
				break
			}
		}
	}
	if me.deadMessageDelay > 0 {
		me.deadMessageDelay -= deltaTime
	}
	if !me.dead {
		if me.checkCollided() {
			me.dead = true
			me.deadMessageDelay = 2
			me.catEntity.Collided = true
			PlaySound(EXPLOSION_SOUND_BYTES, 0.20)
		}
	}
	if me.catEntity.Direction == DIRECTION_BOTTOM && me.checkBottomReached() {
		if me.checkFishReached() {
			PlaySound(REVERSE_SOUND_BYTES, 0.20)
			me.catEntity.Direction = DIRECTION_TOP
			me.obstacleMan.CreateObstacles()
		} else if !me.dead {
			me.dead = true
			me.deadMessageDelay = 2
		}
	}
	if me.catEntity.Direction == DIRECTION_TOP && me.catEntity.Y < 0 && !me.Ascended {
		me.Ascended = true
		PlaySound(ASCENDED_SOUND_BYTES, 0.20)
	}
	if me.Ascended {
		me.ascendedAngle += deltaTime
		me.ascendedAngle = UnwindAngle(me.ascendedAngle)
		me.ascendedPulse += deltaTime
		if me.ascendedPulse >= math.Pi {
			me.ascendedPulse = 0
		}
	}
	me.dustMan.CameraY = me.cameraY / 2
}

func (me *GameSceneVertical) updateCatEntity(deltaTime float64) {
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
}

func (me *GameSceneVertical) Draw(screen *ebiten.Image) {
	me.dustMan.Draw(screen)
	me.drawDecorations(screen)
	me.catEntity.Draw(screen)
	me.obstacleMan.Draw(screen)
	if me.catEntity.Direction == DIRECTION_BOTTOM {
		me.drawFish(screen)
	}
	if me.dead && me.deadMessageDelay <= 0 {
		vector.DrawFilledRect(screen, 0, 0, float32(me.ViewWidth), float32(me.ViewHeight), color.NRGBA{R: 0, G: 0, B: 0, A: 128}, false)
		ebitenutil.DebugPrintAt(screen, "YOU DIED\n"+"press any key", 180, 100)
	}
	if me.Ascended {
		var drawOptions ebiten.DrawImageOptions
		RotateCentered(&drawOptions, float64(me.ascendedImage.Bounds().Dx()), float64(me.ascendedImage.Bounds().Dy()), me.ascendedAngle)
		ScaleCentered(&drawOptions, float64(me.ascendedImage.Bounds().Dx()), float64(me.ascendedImage.Bounds().Dy()), 16, 16)
		drawOptions.GeoM.Translate(me.ViewWidth/2-float64(me.ascendedImage.Bounds().Dx())/2, me.ViewHeight)
		drawOptions.ColorScale.Scale(float32(math.Sin(me.ascendedPulse)), float32(math.Sin(me.ascendedPulse)), float32(math.Sin(me.ascendedPulse)), float32(math.Sin(me.ascendedPulse)))
		screen.DrawImage(me.ascendedImage, &drawOptions)
		ebitenutil.DebugPrintAt(screen, "YOU HAVE ASCENDED\n"+"   press Enter", 160, 50)
	}
}

func (me *GameSceneVertical) GetAreaWidth() float64 {
	return 220
}

func (me *GameSceneVertical) GetTorchSpeedY() (result float64) {
	result = 120
	if me.catEntity.Direction == DIRECTION_TOP {
		result = -result
	}
	return
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
	return 0.5
}

func (me *GameSceneVertical) checkCollided() bool {
	return me.obstacleMan.CheckCollided(me.catEntity.GetHitBox())
}

func (me *GameSceneVertical) drawFish(screen *ebiten.Image) {
	var drawOptions = ebiten.DrawImageOptions{}
	drawOptions.GeoM.Translate(
		me.ViewWidth/2-float64(me.fishImage.Bounds().Dx())/2,
		me.obstacleMan.AreaHeight-float64(me.fishImage.Bounds().Dy())-me.cameraY)
	screen.DrawImage(me.fishImage, &drawOptions)
}

func (me *GameSceneVertical) GetAreaHeight() float64 {
	return me.ViewHeight * 9
}

func (me *GameSceneVertical) checkBottomReached() bool {
	return me.catEntity.Y+me.catEntity.Height >= me.GetAreaHeight()
}

func (me *GameSceneVertical) checkFishReached() bool {
	return math.Abs(me.catEntity.X+me.catEntity.Width/2-me.ViewWidth/2) < float64(me.fishImage.Bounds().Dx())/2+me.catEntity.Width/2
}
