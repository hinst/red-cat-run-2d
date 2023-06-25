package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type CatEntityHorizontalStatus int

const (
	CAT_ENTITY_STATUS_RUN CatEntityHorizontalStatus = iota
	CAT_ENTITY_STATUS_JUMP_SWITCH
	CAT_ENTITY_STATUS_JUMP_FORWARD
	CAT_ENTITY_STATUS_DEAD
)

const CAT_ENTITY_HORIZONTAL_JUMP_TIME = 2

type CatEntityHorizontal struct {
	CatEntity
	Location                    TerrainLocation
	Status                      CatEntityHorizontalStatus
	DebugModeEnabled            bool
	Direction                   Direction
	horizontalJumpTimeRemaining float64
	aimLineAnimationTime        float64

	runImage *ebiten.Image
	runFrame float64

	dieImage *ebiten.Image
	dieFrame float64
}

func (me *CatEntityHorizontal) Initialize() {
	me.runImage = LoadImage(CAT_RUN_IMAGE_BYTES)
	me.dieImage = LoadImage(CAT_DIE_IMAGE_BYTES)

	me.Width = 40
	me.Height = 25

	me.Y = me.FloorY - me.Height
	me.Direction = DIRECTION_RIGHT
}

func (me *CatEntityHorizontal) Update(deltaTime float64) {
	if me.Status == CAT_ENTITY_STATUS_RUN {
		me.updateRun(deltaTime)
	} else if me.Status == CAT_ENTITY_STATUS_JUMP_SWITCH {
		me.updateJumpSwitch(deltaTime)
	} else if me.Status == CAT_ENTITY_STATUS_JUMP_FORWARD {
		me.updateJumpForward(deltaTime)
	} else if me.Status == CAT_ENTITY_STATUS_DEAD {
		me.updateDead(deltaTime)
	}
	if -me.Height < me.Y && me.Y < me.ViewHeight {
		me.X += deltaTime * me.GetSpeedX() * me.getSpeedDirection()
	}
	me.aimLineAnimationTime += deltaTime * math.Pi
	for me.aimLineAnimationTime > math.Pi {
		me.aimLineAnimationTime -= math.Pi
	}
}

func (me *CatEntityHorizontal) getSpeedDirection() (speedDirection float64) {
	if me.Direction == DIRECTION_LEFT {
		speedDirection = -1
	} else if me.Direction == DIRECTION_RIGHT {
		speedDirection = 1
	} else {
		speedDirection = 0
	}
	return
}

func (me *CatEntityHorizontal) updateRun(deltaTime float64) {
	me.runFrame += deltaTime * me.GetRunFramePerSecond()
	if me.runFrame >= CAT_RUN_ANIMATION_FRAME_COUNT {
		me.runFrame -= CAT_RUN_ANIMATION_FRAME_COUNT
	}
	for _, key := range me.JustPressedKeys {
		if key == ebiten.KeySpace {
			for _, key := range me.PressedKeys {
				if key == ebiten.KeyUp && me.Status == CAT_ENTITY_STATUS_RUN && me.Location == TERRAIN_LOCATION_FLOOR {
					me.Status = CAT_ENTITY_STATUS_JUMP_SWITCH
					PlaySound(JUMP_SOUND_BYTES, 0.25)
				}
				if key == ebiten.KeyDown && me.Status == CAT_ENTITY_STATUS_RUN && me.Location == TERRAIN_LOCATION_CEILING {
					me.Status = CAT_ENTITY_STATUS_JUMP_SWITCH
					PlaySound(JUMP_SOUND_BYTES, 0.25)
				}
				var isJumpForward = me.Status == CAT_ENTITY_STATUS_RUN &&
					(key == ebiten.KeyRight && me.Direction == DIRECTION_RIGHT ||
						key == ebiten.KeyLeft && me.Direction == DIRECTION_LEFT)
				if isJumpForward {
					me.Status = CAT_ENTITY_STATUS_JUMP_FORWARD
					me.horizontalJumpTimeRemaining = CAT_ENTITY_HORIZONTAL_JUMP_TIME
					PlaySound(JUMP_SOUND_BYTES, 0.25)
				}
			}
		}
	}
}

func (me *CatEntityHorizontal) updateJumpSwitch(deltaTime float64) {
	me.runFrame += deltaTime * me.GetRunFramePerSecond() / 2
	if me.runFrame >= CAT_RUN_ANIMATION_FRAME_COUNT {
		me.runFrame -= CAT_RUN_ANIMATION_FRAME_COUNT
	}
	if me.Location == TERRAIN_LOCATION_FLOOR {
		me.Y -= deltaTime * me.GetSwitchJumpSpeedY()
		if me.Y <= me.CeilingY {
			me.Status = CAT_ENTITY_STATUS_RUN
			me.Location = TERRAIN_LOCATION_CEILING
			me.Y = me.CeilingY
		}
	} else if me.Location == TERRAIN_LOCATION_CEILING {
		me.Y += deltaTime * me.GetSwitchJumpSpeedY()
		if me.Y+me.Height >= me.FloorY {
			me.Status = CAT_ENTITY_STATUS_RUN
			me.Location = TERRAIN_LOCATION_FLOOR
			me.Y = me.FloorY - me.Height
		}
	}
}

func (me *CatEntityHorizontal) updateJumpForward(deltaTime float64) {
	me.runFrame += deltaTime * me.GetRunFramePerSecond() / 2
	if me.runFrame >= CAT_RUN_ANIMATION_FRAME_COUNT {
		me.runFrame -= CAT_RUN_ANIMATION_FRAME_COUNT
	}
	me.horizontalJumpTimeRemaining -= deltaTime
	var elevation float64
	if me.horizontalJumpTimeRemaining > CAT_ENTITY_HORIZONTAL_JUMP_TIME/2 {
		var horizontalJumpTimePassed = CAT_ENTITY_HORIZONTAL_JUMP_TIME - me.horizontalJumpTimeRemaining
		elevation = me.GetForwardJumpSpeedY() * horizontalJumpTimePassed
	} else {
		elevation = me.GetForwardJumpSpeedY() * me.horizontalJumpTimeRemaining
	}
	if me.Location == TERRAIN_LOCATION_FLOOR {
		me.Y = me.FloorY - me.Height - elevation
	} else if me.Location == TERRAIN_LOCATION_CEILING {
		me.Y = me.CeilingY + elevation
	}
	if me.horizontalJumpTimeRemaining <= 0 {
		me.Status = CAT_ENTITY_STATUS_RUN
		if me.Location == TERRAIN_LOCATION_FLOOR {
			me.Y = me.FloorY - me.Height
		} else {
			me.Y = me.CeilingY
		}
	}
}

func (me *CatEntityHorizontal) updateDead(deltaTime float64) {
	me.dieFrame += deltaTime * me.GetDieFramePerSecond()
	if me.dieFrame >= CAT_DIE_ANIMATION_FRAME_COUNT {
		me.dieFrame = CAT_DIE_ANIMATION_FRAME_COUNT - 1
	}
	if me.Location == TERRAIN_LOCATION_FLOOR {
		if me.Y < me.ViewHeight {
			me.Y += deltaTime * me.GetFallSpeed()
		}
	} else if me.Location == TERRAIN_LOCATION_CEILING {
		if me.Y > -me.Height {
			me.Y -= deltaTime * me.GetFallSpeed()
		}
	}
}

func (me *CatEntityHorizontal) Draw(screen *ebiten.Image) {
	if me.DebugModeEnabled {
		vector.DrawFilledRect(screen, float32(me.X-me.CameraX), float32(me.Y-me.CameraY),
			float32(me.Width), float32(me.Height), color.RGBA{R: 100, G: 100, B: 100}, false)
	}
	me.drawAimLine(screen)
	var drawOptions = ebiten.DrawImageOptions{}
	if me.Location == TERRAIN_LOCATION_CEILING {
		ScaleCentered(&drawOptions, me.Width, me.Height, 1, -1)
	}
	if me.Direction == DIRECTION_LEFT {
		ScaleCentered(&drawOptions, me.Width, me.Height, -1, 1)
	}
	drawOptions.GeoM.Translate(me.X, me.Y)
	drawOptions.GeoM.Translate(-me.CameraX, -me.CameraY)
	var isRunDrawMode = me.Status == CAT_ENTITY_STATUS_RUN ||
		me.Status == CAT_ENTITY_STATUS_JUMP_SWITCH ||
		me.Status == CAT_ENTITY_STATUS_JUMP_FORWARD
	if isRunDrawMode {
		var spriteShiftX = float64(int(me.runFrame)) * CAT_RUN_ANIMATION_FRAME_WIDTH
		var animationFrameRectangle = GetShiftedRectangle(spriteShiftX,
			CAT_RUN_ANIMATION_FRAME_WIDTH, float64(me.runImage.Bounds().Dy()))
		screen.DrawImage(me.runImage.SubImage(animationFrameRectangle).(*ebiten.Image), &drawOptions)
	} else if me.Status == CAT_ENTITY_STATUS_DEAD {
		var spriteShiftX = float64(int(me.dieFrame)) * CAT_RUN_ANIMATION_FRAME_WIDTH
		var animationFrameRectangle = GetShiftedRectangle(spriteShiftX,
			CAT_RUN_ANIMATION_FRAME_WIDTH, float64(me.dieImage.Bounds().Dy()))
		screen.DrawImage(me.dieImage.SubImage(animationFrameRectangle).(*ebiten.Image), &drawOptions)
	}
}

func (me *CatEntityHorizontal) drawAimLine(screen *ebiten.Image) {
	for _, key := range me.PressedKeys {
		if me.isAimUp(key) {
			me.drawVerticalAimLine(screen, true)
			break
		} else if me.isAimForwardFromFloor(key) {
			me.drawHorizontalAimLine(screen, true)
			break
		} else if me.isAimDown(key) {
			me.drawVerticalAimLine(screen, false)
			break
		} else if me.isAimForwardFromCeiling(key) {
			me.drawHorizontalAimLine(screen, false)
			break
		}
	}
}

func (me *CatEntityHorizontal) isAimUp(key ebiten.Key) bool {
	return key == ebiten.KeyUp && me.Status == CAT_ENTITY_STATUS_RUN && me.Location == TERRAIN_LOCATION_FLOOR
}

func (me *CatEntityHorizontal) isAimForward(key ebiten.Key) bool {
	return key == ebiten.KeyRight && me.Direction == DIRECTION_RIGHT ||
		key == ebiten.KeyLeft && me.Direction == DIRECTION_LEFT
}

func (me *CatEntityHorizontal) isAimForwardFromFloor(key ebiten.Key) bool {
	return me.isAimForward(key) && me.Status == CAT_ENTITY_STATUS_RUN && me.Location == TERRAIN_LOCATION_FLOOR
}

func (me *CatEntityHorizontal) isAimDown(key ebiten.Key) bool {
	return key == ebiten.KeyDown && me.Status == CAT_ENTITY_STATUS_RUN && me.Location == TERRAIN_LOCATION_CEILING
}

func (me *CatEntityHorizontal) isAimForwardFromCeiling(key ebiten.Key) bool {
	return me.isAimForward(key) && me.Status == CAT_ENTITY_STATUS_RUN && me.Location == TERRAIN_LOCATION_CEILING
}

func (me *CatEntityHorizontal) drawVerticalAimLine(screen *ebiten.Image, up bool) {
	const multiplier = 1.66
	var y1 = me.Y + me.Height/2 - me.CameraY
	if up {
		y1 -= me.GetSwitchJumpSpeedY() * multiplier
	} else {
		y1 += me.GetSwitchJumpSpeedY() * multiplier
	}
	vector.StrokeLine(screen,
		float32(me.X+me.Width/2-me.CameraX),
		float32(me.Y+me.Height/2-me.CameraY),
		float32(me.X+me.Width/2+me.GetSpeedX()*me.getSpeedDirection()*multiplier-me.CameraX),
		float32(y1),
		1,
		me.GetAimLineColor(),
		false,
	)
}

// isFloor == false means ceiling
func (me *CatEntityHorizontal) drawHorizontalAimLine(screen *ebiten.Image, isFloor bool) {
	var y1 = me.Y + me.Height/2 - me.CameraY
	if isFloor {
		y1 -= me.GetForwardJumpSpeedY()
	} else {
		y1 += me.GetForwardJumpSpeedY()
	}
	vector.StrokeLine(screen,
		float32(me.X+me.Width/2-me.CameraX),
		float32(me.Y+me.Height/2-me.CameraY),
		float32(me.X+me.Width/2+me.GetSpeedX()*me.getSpeedDirection()-me.CameraX),
		float32(y1),
		1,
		me.GetAimLineColor(),
		false,
	)
	vector.StrokeLine(screen,
		float32(me.X+me.Width/2+me.GetSpeedX()*me.getSpeedDirection()-me.CameraX),
		float32(y1),
		float32(me.X+me.Width/2+me.GetSpeedX()*me.getSpeedDirection()*2-me.CameraX),
		float32(me.Y+me.Height/2-me.CameraY),
		1,
		me.GetAimLineColor(),
		false,
	)
}

func (me *CatEntityHorizontal) GetSpeedX() float64 {
	return 50
}

func (me *CatEntityHorizontal) GetFallSpeed() float64 {
	return 50
}

func (me *CatEntityHorizontal) GetSwitchJumpSpeedY() float64 {
	return 70
}

func (me *CatEntityHorizontal) GetForwardJumpSpeedY() float64 {
	return 50
}

func (me *CatEntityHorizontal) GetAimLineColor() color.Color {
	return color.NRGBA{R: 168, G: 111, B: 50, A: 100 + uint8(155*math.Sin(me.aimLineAnimationTime))}
}

func (me *CatEntityHorizontal) GetRunFramePerSecond() float64 {
	return 6
}

func (me *CatEntityHorizontal) GetDieFramePerSecond() float64 {
	return 6
}
