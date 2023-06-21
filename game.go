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

type GameMode int

const (
	GAME_MODE_MENU = iota
	GAME_MODE_GAME
	GAME_MODE_INFORMATION
)

type Game struct {
	updateTime      time.Time
	justPressedKeys []ebiten.Key
	pressedKeys     []ebiten.Key
	menu            MenuUserInterface
	gameScene       GameSceneHorizontal
	gameInfoScene   GameInfoScene
	isExiting       bool
	mode            GameMode
	viewWidth       float64
	viewHeight      float64

	titleImage             *ebiten.Image
	ebitengineReverseImage *ebiten.Image
	catWalkImage           *ebiten.Image
	catRunFrame            float64
}

const (
	GAME_MENU_ITEM_ID_NEW_GAME = iota
	GAME_MENU_ITEM_ID_TOGGLE_FULL_SCREEN
	GAME_MENU_ITEM_ID_INFORMATION
	GAME_MENU_ITEM_ID_EXIT
)

func (me *Game) Initialize() {
	var titleImage, _, titleImageError = image.Decode(bytes.NewReader(TITLE_IMAGE_BYTES))
	AssertError(titleImageError)
	me.titleImage = ebiten.NewImageFromImage(titleImage)
	var ebitengineReverseImage, _, ebitengineReverseImageError = image.Decode(bytes.NewReader(EBITENGINE_REVERSE_IMAGE_BYTES))
	AssertError(ebitengineReverseImageError)
	me.ebitengineReverseImage = ebiten.NewImageFromImage(ebitengineReverseImage)
	me.viewWidth = 420
	me.viewHeight = 240
	me.updateTime = time.Now()
	me.menu = MenuUserInterface{
		Items: []MenuUserInterfaceItem{
			{
				Title: "New Game",
				Id:    GAME_MENU_ITEM_ID_NEW_GAME,
			},
			{
				Title: "Toggle Full Screen",
				Id:    GAME_MENU_ITEM_ID_TOGGLE_FULL_SCREEN,
			},
			{
				Title: "Information",
				Id:    GAME_MENU_ITEM_ID_INFORMATION,
			},
			{
				Title: "Exit",
				Id:    GAME_MENU_ITEM_ID_EXIT,
			},
		},
	}
	me.gameScene = GameSceneHorizontal{}
	me.gameScene.ViewWidth = me.viewWidth
	me.gameScene.ViewHeight = me.viewHeight
	me.gameScene.Initialize()

	var catWalkImage, _, catImageError = image.Decode(bytes.NewReader(CAT_RUN_IMAGE_BYTES))
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
	var deltaTime = math.Min(1, updateTime.Sub(me.updateTime).Seconds())
	me.update(deltaTime)
	me.updateTime = updateTime
	me.justPressedKeys = me.justPressedKeys[:0]
	me.draw(screen)
}

func (me *Game) update(deltaTime float64) {
	if me.mode == GAME_MODE_MENU {
		me.updateMenu(deltaTime)
	} else if me.mode == GAME_MODE_GAME {
		me.updateGameScene(deltaTime)
	} else if me.mode == GAME_MODE_INFORMATION {
		if len(me.justPressedKeys) > 0 {
			me.mode = GAME_MODE_MENU
		}
	}
}

func (me *Game) draw(screen *ebiten.Image) {
	if me.mode == GAME_MODE_MENU {
		if false {
			me.drawTitle(screen)
		}
		me.drawEbitenReverse(screen)
		me.drawCatAnimationTop(screen)
		me.menu.Draw(screen)
	} else if me.mode == GAME_MODE_GAME {
		me.gameScene.Draw(screen)
	} else if me.mode == GAME_MODE_INFORMATION {
		me.gameInfoScene.Draw(screen)
	}
}

func (me *Game) updateMenu(deltaTime float64) {
	me.menu.JustPressedKeys = me.justPressedKeys
	me.menu.Update(deltaTime)
	if me.menu.PressedItemId == GAME_MENU_ITEM_ID_NEW_GAME {
		me.mode = GAME_MODE_GAME
	} else if me.menu.PressedItemId == GAME_MENU_ITEM_ID_TOGGLE_FULL_SCREEN {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	} else if me.menu.PressedItemId == GAME_MENU_ITEM_ID_INFORMATION {
		me.mode = GAME_MODE_INFORMATION
	} else if me.menu.PressedItemId == GAME_MENU_ITEM_ID_EXIT {
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

func (me *Game) drawTitle(screen *ebiten.Image) {
	var drawOptions = ebiten.DrawImageOptions{}
	drawOptions.GeoM.Scale(0.6, 0.6)
	drawOptions.GeoM.Translate(230, 210)
	drawOptions.ColorScale.Scale(0.6, 0.6, 0.6, 1)
	drawOptions.Filter = ebiten.FilterLinear
	screen.DrawImage(me.titleImage, &drawOptions)
}
