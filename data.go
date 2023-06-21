package main

import (
	_ "embed"
)

var (
	//go:embed data/title.png
	TITLE_IMAGE_BYTES []byte
	//go:embed data/catWalk.png
	CAT_RUN_IMAGE_BYTES []byte
	//go:embed data/catDie.png
	CAT_DIE_IMAGE_BYTES []byte
	//go:embed data/brickBlock.png
	BRICK_BLOCK_IMAGE_BYTES []byte
	//go:embed data/dirtBlock.png
	DIRT_BLOCK_IMAGE_BYTES []byte
	//go:embed data/ebitengine-reverse.png
	EBITENGINE_REVERSE_IMAGE_BYTES []byte
	//go:embed data/waterBlockTop.png
	WATER_BLOCK_TOP_IMAGE_BYTES []byte
	//go:embed data/waterBlock.png
	WATER_BLOCK_IMAGE_BYTES []byte
)
