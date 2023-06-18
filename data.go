package main

import (
	_ "embed"
)

var (
	//go:embed data/catWalk.png
	catRun []byte
	//go:embed data/catDie.png
	catDie []byte
	//go:embed data/brickBlock.png
	brickBlock []byte
	//go:embed data/ebitengine-reverse.png
	ebitengineReverse []byte
)
