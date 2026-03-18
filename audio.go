package rythmpen

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
)

type AudioManager struct {
	context *audio.Context
	player  *audio.Player
}

func NewAudioManager(sampleRate int) *AudioManager {
	ctx := audio.NewContext(sampleRate)
	return &AudioManager{
		context: ctx,
		// TODO: Create audio player
		// player: ctx.NewPlayerF32("")
	}
}
