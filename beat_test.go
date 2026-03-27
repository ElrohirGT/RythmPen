package rythmpen

import (
	"testing"
	"time"

	"github.com/assertgo/assert"
)

func Test_AddBeat(t *testing.T) {
	ass := assert.New(t)

	manager := NewBeatManager(BeatConfig{}, BeatConfig{})
	b1 := &Beat{}
	manager.AddBeat(b1)
	manager.AddBeat(b1)
	manager.AddBeat(b1)

	ass.ThatInt(len(manager.beats)).
		IsEqualTo(3)

	manager.AddLeftBeat(2 * time.Second)
	manager.AddRightBeat(2 * time.Second)

	ass.ThatInt(len(manager.beats)).
		IsEqualTo(5)
}
