package rythmpen

import (
	"io"
	"time"

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
	}
}

func (m *AudioManager) NewAudioPlayer(src io.Reader) error {
	var err error
	m.player, err = m.context.NewPlayer(src)
	if err != nil {
		return err
	}

	return nil
}

func (m *AudioManager) Play() {
	// if m.context.IsReady() {
	m.player.Play()
	// } else {
	// 	log.Panic("Failed to play audio because context is not ready!")
	// }
}

func (m *AudioManager) Position() time.Duration {
	return m.player.Position()
}
