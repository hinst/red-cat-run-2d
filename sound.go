package main

import (
	"bytes"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

const SAMPLE_RATE = 44000

func PlaySound(data []byte, volume float64) {
	go func() {
		if audio.CurrentContext() == nil {
			audio.NewContext(SAMPLE_RATE)
		}
		var stream, streamError = vorbis.DecodeWithSampleRate(SAMPLE_RATE, bytes.NewReader(data))
		AssertError(streamError)
		var player, playerError = audio.CurrentContext().NewPlayer(stream)
		player.SetVolume(volume)
		AssertError(playerError)
		player.Play()
	}()
}
