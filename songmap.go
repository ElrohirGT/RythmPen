package rythmpen

import (
	"fmt"
	"io"
	"time"
)

type PressStatus int

var PressStatusEnum = struct {
	IGNORE  PressStatus
	LIFTED  PressStatus
	PRESSED PressStatus
}{
	IGNORE:  0,
	LIFTED:  1,
	PRESSED: 2,
}

type SongBeat struct {
	Position  time.Duration
	LeftSide  PressStatus
	RightSide PressStatus
}

type AudioPositioner interface {
	Position() time.Duration
}

type DummyAudioPositioner struct {
	D_Position func() time.Duration
}

func (d DummyAudioPositioner) Position() time.Duration {
	return d.D_Position()
}

type SongMap struct {
	beats        []SongBeat
	audioManager AudioPositioner
}

func NewSongMap(audioManager AudioPositioner) *SongMap {
	return &SongMap{
		beats:        make([]SongBeat, 0, 50),
		audioManager: audioManager,
	}
}

func (s *SongMap) LeftBeat() {
	s.beats = append(s.beats, SongBeat{
		Position: s.audioManager.Position(),
		LeftSide: PressStatusEnum.PRESSED,
	})
}

func (s *SongMap) RightBeat() {
	s.beats = append(s.beats, SongBeat{
		Position:  s.audioManager.Position(),
		RightSide: PressStatusEnum.PRESSED,
	})
}

func (s *SongMap) BothBeat() {
	s.beats = append(s.beats, SongBeat{
		Position:  s.audioManager.Position(),
		RightSide: PressStatusEnum.PRESSED,
		LeftSide:  PressStatusEnum.PRESSED,
	})
}

func (s *SongMap) WriteToFile(dst io.Writer) error {
	_, err := io.WriteString(dst, "position,left,right\n")
	if err != nil {
		return err
	}

	for _, b := range s.beats {
		_, err := fmt.Fprintf(dst, "%d,%d,%d\n", b.Position.Microseconds(), b.LeftSide, b.RightSide)
		if err != nil {
			return err
		}
	}

	return nil
}
