package main

import (
	"bytes"
	"errors"
	"image"
	"math"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/exp/slices"
)

type GameMode int

const (
	GAME_MODE_MENU = iota
	GAME_MODE_GAME
	GAME_MODE_INFORMATION
)

type Game struct {
	updateTime                       time.Time
	justPressedKeys                  []ebiten.Key
	pressedKeys                      []ebiten.Key
	menu                             MenuUserInterface
	gameScene                        *GameScene
	gameInfoScene                    GameInfoScene
	isExiting                        bool
	mode                             GameMode
	viewWidth                        float64
	viewHeight                       float64
	updatesToSkip                    int
	initialized                      bool
	titleImage                       *ebiten.Image
	ebitengineReverseImage           *ebiten.Image
	catWalkImage                     *ebiten.Image
	catRunFrame                      float64
	initialInformationAcknowledged   bool
	fpsCounterEnabled                bool
	isFirstUpdateAfterInitialization bool
	isNewGamePlusUnlocked            bool
}

const (
	GAME_MENU_ITEM_ID_NEW_GAME = iota
	GAME_MENU_ITEM_ID_NEW_GAME_PLUS
	GAME_MENU_ITEM_ID_TOGGLE_FULL_SCREEN
	GAME_MENU_ITEM_ID_GENERAL_INFORMATION
	GAME_MENU_ITEM_ID_TOGGLE_VSYNC
	GAME_MENU_ITEM_ID_TOGGLE_FPS_COUNTER
	GAME_MENU_ITEM_ID_EXIT
)

const GAME_TEXT_CONTROLS = "press [up + space] to jump up\n" +
	"press [right + space] to jump right\n" +
	"\n" +
	"press any key to start"

func (me *Game) Initialize() {
	me.viewWidth = 420
	me.viewHeight = 240
	me.updatesToSkip = 4
}

func (me *Game) initializeInternal() {
	var titleImage, _, titleImageError = image.Decode(bytes.NewReader(TITLE_IMAGE_BYTES))
	AssertError(titleImageError)
	me.titleImage = ebiten.NewImageFromImage(titleImage)
	var ebitengineReverseImage, _, ebitengineReverseImageError = image.Decode(bytes.NewReader(EBITENGINE_REVERSE_IMAGE_BYTES))
	AssertError(ebitengineReverseImageError)
	me.ebitengineReverseImage = ebiten.NewImageFromImage(ebitengineReverseImage)
	me.updateTime = time.Now()
	me.menu = MenuUserInterface{Items: me.createMenuItems()}
	var catWalkImage, _, catImageError = image.Decode(bytes.NewReader(CAT_RUN_IMAGE_BYTES))
	AssertError(catImageError)
	me.catWalkImage = ebiten.NewImageFromImage(catWalkImage)
	InitializeSound()
	me.initialized = true
	me.isFirstUpdateAfterInitialization = true
}

func (me *Game) initializeGameScene() {
	me.gameScene = &GameScene{}
	me.gameScene.ViewWidth = me.viewWidth
	me.gameScene.ViewHeight = me.viewHeight
	me.gameScene.Initialize()
}

func (me *Game) Update() error {
	if me.updatesToSkip > 0 {
		me.updatesToSkip--
		return nil
	}
	if !me.initialized {
		me.initializeInternal()
		return nil
	}
	if me.isFirstUpdateAfterInitialization {
		// Ignore keys pressed during the initialization screen by the impatient user
		me.isFirstUpdateAfterInitialization = false
		return nil
	}
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
	if !me.initialized {
		ebitenutil.DebugPrint(screen, "Waiting for initialization...")
		return
	}
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
		if me.initialInformationAcknowledged {
			me.updateGameScene(deltaTime)
		} else if len(me.justPressedKeys) > 0 {
			me.initialInformationAcknowledged = true
		}
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
		if !me.initialInformationAcknowledged {
			ebitenutil.DebugPrintAt(screen, GAME_TEXT_CONTROLS, 80, 80)
		}
	} else if me.mode == GAME_MODE_INFORMATION {
		me.gameInfoScene.Draw(screen)
	}
	if me.fpsCounterEnabled {
		ebitenutil.DebugPrint(screen, strconv.FormatFloat(ebiten.ActualFPS(), 'f', 0, 64))
	}
}

func (me *Game) updateMenu(deltaTime float64) {
	me.menu.JustPressedKeys = me.justPressedKeys
	me.menu.Update(deltaTime)
	switch me.menu.PressedItemId {
	case GAME_MENU_ITEM_ID_NEW_GAME:
		me.mode = GAME_MODE_GAME
		me.initializeGameScene()
	case GAME_MENU_ITEM_ID_TOGGLE_FULL_SCREEN:
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	case GAME_MENU_ITEM_ID_GENERAL_INFORMATION:
		me.gameInfoScene.Text = GAME_INFO_SCENE_TEXT_GENERAL
		me.mode = GAME_MODE_INFORMATION
	case GAME_MENU_ITEM_ID_TOGGLE_VSYNC:
		ebiten.SetVsyncEnabled(!ebiten.IsVsyncEnabled())
		for i := 0; i < len(me.menu.Items); i++ {
			if me.menu.Items[i].Id == GAME_MENU_ITEM_ID_TOGGLE_VSYNC {
				me.menu.Items[i].Title = "Toggle Vsync: " + strconv.FormatBool(ebiten.IsVsyncEnabled())
			}
		}
	case GAME_MENU_ITEM_ID_TOGGLE_FPS_COUNTER:
		me.fpsCounterEnabled = !me.fpsCounterEnabled
	case GAME_MENU_ITEM_ID_EXIT:
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
	if me.gameScene.Completed {
		me.gameScene.Close()
		me.mode = GAME_MODE_MENU
	}
}

func (me *Game) drawEbitenReverse(screen *ebiten.Image) {
	var drawOptions = ebiten.DrawImageOptions{}
	drawOptions.GeoM.Scale(0.25, 0.25)
	drawOptions.GeoM.Translate(200, 50)
	screen.DrawImage(me.ebitengineReverseImage, &drawOptions)
}

func (me *Game) drawCatAnimationTop(screen *ebiten.Image) {
	var drawOptions = ebiten.DrawImageOptions{}
	drawOptions.GeoM.Scale(-1, 1)
	drawOptions.GeoM.Translate(320, 25)

	var spriteShiftX = float64(int(me.catRunFrame)) * CAT_RUN_ANIMATION_FRAME_WIDTH
	var rect = GetShiftedRectangle(spriteShiftX, CAT_RUN_ANIMATION_FRAME_WIDTH, float64(me.catWalkImage.Bounds().Dy()))
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

func (me *Game) createMenuItems() (items []MenuUserInterfaceItem) {
	items = []MenuUserInterfaceItem{
		{
			Title: "New Game",
			Id:    GAME_MENU_ITEM_ID_NEW_GAME,
		},
		{
			Title: "Information",
			Id:    GAME_MENU_ITEM_ID_GENERAL_INFORMATION,
		},
		{
			Title: "Toggle Full Screen",
			Id:    GAME_MENU_ITEM_ID_TOGGLE_FULL_SCREEN,
		},
		{
			Title: "Toggle Vsync: " + strconv.FormatBool(ebiten.IsVsyncEnabled()),
			Id:    GAME_MENU_ITEM_ID_TOGGLE_VSYNC,
		},
		{
			Title: "Toggle FPS counter",
			Id:    GAME_MENU_ITEM_ID_TOGGLE_FPS_COUNTER,
		},
		{
			Title: "Exit",
			Id:    GAME_MENU_ITEM_ID_EXIT,
		},
	}
	if me.isNewGamePlusUnlocked {
		slices.Insert(items, 1, MenuUserInterfaceItem{
			Title: "New Game Plus",
			Id:    GAME_MENU_ITEM_ID_NEW_GAME_PLUS,
		})
	}
	return items
}
