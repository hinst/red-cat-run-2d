package main

import (
	"bytes"
	"io"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

const SAMPLE_RATE = 44100

func StreamToBytes(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

func DecodeVorbis(data []byte) []byte {
	if audio.CurrentContext() == nil {
		audio.NewContext(SAMPLE_RATE)
	}
	var stream, streamError = vorbis.DecodeWithSampleRate(audio.CurrentContext().SampleRate(), bytes.NewReader(data))
	AssertError(streamError)
	return StreamToBytes(stream)
}

func InitializeSound() {
	ACHIEVEMENT_SOUND_BYTES = DecodeVorbis(ACHIEVEMENT_SOUND_BYTES)
	JUMP_SOUND_BYTES = DecodeVorbis(JUMP_SOUND_BYTES)
	HIT_SOUND_BYTES = DecodeVorbis(HIT_SOUND_BYTES)
	REVERSE_SOUND_BYTES = DecodeVorbis(REVERSE_SOUND_BYTES)
	EXPLOSION_SOUND_BYTES = DecodeVorbis(EXPLOSION_SOUND_BYTES)
	ASCENDED_SOUND_BYTES = DecodeVorbis(ASCENDED_SOUND_BYTES)
	VIBIN_SOUND_BYTES = DecodeVorbis(VIBIN_SOUND_BYTES)
}

func PlaySound(data []byte, volume float64) *audio.Player {
	var player = audio.CurrentContext().NewPlayerFromBytes(data)
	player.SetVolume(volume)
	player.Play()
	return player
}
