package main

import (
	"bytes"
	"errors"
	"image"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	updateTime      time.Time
	justPressedKeys []ebiten.Key
	pressedKeys     []ebiten.Key
	menu            MenuUserInterface
	gameScene       GameScene
	isExiting       bool
	mode            int
	viewWidth       float64
	viewHeight      float64

	ebitengineReverseImage *ebiten.Image
	catWalkImage           *ebiten.Image
	catRunFrame            float64
}

func (me *Game) Initialize() {
	var ebitengineReverseImage, _, ebitengineReverseImageError = image.Decode(bytes.NewReader(ebitengineReverse))
	AssertError(ebitengineReverseImageError)
	me.ebitengineReverseImage = ebiten.NewImageFromImage(ebitengineReverseImage)
	me.viewWidth = 420
	me.viewHeight = 240
	me.updateTime = time.Now()
	me.menu = MenuUserInterface{
		Items: []MenuUserInterfaceItem{
			{
				Title: "New Game",
				Id:    1,
			},
			{
				Title: "Toggle Full Screen",
				Id:    2,
			},
			{
				Title: "Exit",
				Id:    3,
			},
		},
	}
	me.gameScene = GameScene{}
	me.gameScene.ViewWidth = me.viewWidth
	me.gameScene.ViewHeight = me.viewHeight
	me.gameScene.Initialize()

	var catWalkImage, _, catImageError = image.Decode(bytes.NewReader(catRun))
	AssertError(catImageError)
	me.catWalkImage = ebiten.NewImageFromImage(catWalkImage)
}

func (me *Game) Update() error {
	me.justPressedKeys = me.justPressedKeys[:0]
	me.justPressedKeys = inpututil.AppendJustPressedKeys(me.justPressedKeys)
	me.pressedKeys = me.pressedKeys[:0]
	me.pressedKeys = inpututil.AppendPressedKeys(me.pressedKeys)
	if me.isExiting {
		return errors.New("exiting")
	}
	return nil
}

func (me *Game) Draw(screen *ebiten.Image) {
	var updateTime = time.Now()
	me.update(updateTime.Sub(me.updateTime).Seconds())
	me.updateTime = updateTime
	me.justPressedKeys = me.justPressedKeys[:0]
	me.draw(screen)
}

func (me *Game) update(deltaTime float64) {
	if me.mode == me.GetModeMenu() {
		me.updateMenu(deltaTime)
	} else if me.mode == me.GetModeGame() {
		me.updateGameScene(deltaTime)
	}
}

func (me *Game) updateMenu(deltaTime float64) {
	me.menu.JustPressedKeys = me.justPressedKeys
	me.menu.Update(deltaTime)
	if me.menu.PressedItemId == 1 {
		me.mode = me.GetModeGame()
	} else if me.menu.PressedItemId == 2 {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	} else if me.menu.PressedItemId == 3 {
		me.isExiting = true
	}
	me.catRunFrame += deltaTime * CAT_RUN_ANIMATION_FRAME_PER_SECOND
	for me.catRunFrame >= CAT_RUN_ANIMATION_FRAME_COUNT {
		me.catRunFrame -= CAT_RUN_ANIMATION_FRAME_COUNT
	}
}

func (me *Game) updateGameScene(deltaTime float64) {
	me.gameScene.PressedKeys = me.pressedKeys
	me.gameScene.JustPressedKeys = me.justPressedKeys
	me.gameScene.Update(deltaTime)
}

func (me *Game) draw(screen *ebiten.Image) {
	if me.mode == me.GetModeMenu() {
		me.drawEbitenReverse(screen)
		me.drawCatAnimationTop(screen)
		me.menu.Draw(screen)
	} else if me.mode == me.GetModeGame() {
		me.gameScene.Draw(screen)
	}
}

func (me *Game) drawEbitenReverse(screen *ebiten.Image) {
	var drawOptions = ebiten.DrawImageOptions{}
	drawOptions.GeoM.Scale(0.25, 0.25)
	drawOptions.GeoM.Translate(200, 50)
	screen.DrawImage(me.ebitengineReverseImage, &drawOptions)
}

func (me *Game) drawCatAnimationBottom(screen *ebiten.Image) {
	var drawOptions = ebiten.DrawImageOptions{}
	drawOptions.GeoM.Scale(-1, 1)
	drawOptions.GeoM.Rotate(math.Pi)
	drawOptions.GeoM.Translate(285, 229)

	var spriteShiftX = float64(int(me.catRunFrame)) * CAT_RUN_ANIMATION_FRAME_WIDTH
	var rect = GetShiftedRectangle(spriteShiftX, CAT_RUN_ANIMATION_FRAME_WIDTH)
	screen.DrawImage(me.catWalkImage.SubImage(rect).(*ebiten.Image), &drawOptions)
}

func (me *Game) drawCatAnimationTop(screen *ebiten.Image) {
	var drawOptions = ebiten.DrawImageOptions{}
	drawOptions.GeoM.Scale(-1, 1)
	drawOptions.GeoM.Translate(320, 25)

	var spriteShiftX = float64(int(me.catRunFrame)) * CAT_RUN_ANIMATION_FRAME_WIDTH
	var rect = GetShiftedRectangle(spriteShiftX, CAT_RUN_ANIMATION_FRAME_WIDTH)
	screen.DrawImage(me.catWalkImage.SubImage(rect).(*ebiten.Image), &drawOptions)
}

func (me *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(me.viewWidth), int(me.viewHeight)
}

func (me *Game) GetModeMenu() int {
	return 0
}

func (me *Game) GetModeGame() int {
	return 1
}
