package rythmpen

import (
	"io"
	"log"
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

func (m *AudioManager) Reset() {
	if err := m.player.Rewind(); err != nil {
		log.Panicf("Failed to rewind song! %v\n", err)
	}
}

func (m *AudioManager) Play() {
	m.player.Play()
}

func (m *AudioManager) IsPlaying() bool {
	return m.player.IsPlaying()
}

func (m *AudioManager) Position() time.Duration {
	return m.player.Position()
}
