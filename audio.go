package main

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
)

type AudioManager struct {
	context *audio.Context
}

func NewAudioManager(sampleRate int) *AudioManager {
	return &AudioManager{
		context: audio.NewContext(sampleRate),
	}
}
