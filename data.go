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
	//go:embed data/smallFish.png
	FISH_IMAGE_BYTES []byte
	//go:embed data/catFlyDown.png
	CAT_FLY_DOWN_IMAGE_BYTES []byte
	//go:embed data/torch.png
	TORCH_IMAGE_BYTES []byte
	//go:embed data/obstacle.png
	OBSTACLE_IMAGE_BYTES []byte
	//go:embed data/270330__littlerobotsoundfactory__jingle_achievement_01.ogg
	ACHIEVEMENT_SOUND_BYTES []byte
	//go:embed data/270342__littlerobotsoundfactory__pickup_03.ogg
	JUMP_SOUND_BYTES []byte
	//go:embed data/270326__littlerobotsoundfactory__hit_01.ogg
	HIT_SOUND_BYTES []byte
	//go:embed data/270341__littlerobotsoundfactory__pickup_04.ogg
	REVERSE_SOUND_BYTES []byte
)
