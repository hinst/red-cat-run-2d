package main

import (
	"bytes"
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type TerrainMan struct {
	brickBlockImage *ebiten.Image
	blocks          []*TerrainBlock
	// Initialization parameter, pixels
	ViewWidth int
	// Initialization parameter, pixels
	ViewHeight int
	// Initialization parameter, tiles
	AreaWidth int
	CameraX   float64
}

func (me *TerrainMan) GetMinBlockWidth() int {
	return 8
}

func (me *TerrainMan) GetMaxBlockWidth() int {
	return 16
}

func (me *TerrainMan) GetMinGapWidth() int {
	return 3
}

func (me *TerrainMan) GetMaxGapWidth() int {
	return 6
}

func (me *TerrainMan) GetTileWidth() int {
	return 10
}

func (me *TerrainMan) GetTileHeight() int {
	return 10
}

func (me *TerrainMan) Initialize() {
	var brickBlockImage, _, brickBlockImageError = image.Decode(bytes.NewReader(brickBlock))
	Assert(brickBlockImageError)
	me.brickBlockImage = ebiten.NewImageFromImage(brickBlockImage)
	for me.GetLastBlock() == nil || me.GetLastBlock().X+me.GetLastBlock().Width < me.AreaWidth {
		var block = &TerrainBlock{}
		if me.GetLastBlock() == nil {
			block.Type = block.GetTypeFloor()
			block.X = 0
			block.Width = me.GetMaxBlockWidth()
		} else {
			block.Type = rand.Intn(2)
			var gap = GetRandomNumberBetween(me.GetMinGapWidth(), me.GetMaxGapWidth())
			block.X = me.GetLastBlock().X + me.GetLastBlock().Width + gap
			block.Width = GetRandomNumberBetween(me.GetMinBlockWidth(), me.GetMaxBlockWidth())
		}
		me.blocks = append(me.blocks, block)
	}
}

func (me *TerrainMan) GetLastBlock() *TerrainBlock {
	if len(me.blocks) > 0 {
		return me.blocks[len(me.blocks)-1]
	} else {
		return nil
	}
}

func (me *TerrainMan) Draw(screen *ebiten.Image) {
	for _, block := range me.blocks {
		var drawOptions ebiten.DrawImageOptions
		drawOptions.GeoM.Translate(float64(me.GetTileWidth())*float64(block.X), 0)
		if block.Type == block.GetTypeFloor() {
			drawOptions.GeoM.Translate(0, 200)
		} else {
			drawOptions.GeoM.Translate(0, 30)
		}
		for i := 0; i < block.Width; i++ {
			screen.DrawImage(me.brickBlockImage, &drawOptions)
			drawOptions.GeoM.Translate(float64(me.GetTileWidth()), 0)
		}
	}
}