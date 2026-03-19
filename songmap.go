package rythmpen

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
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

func (s *SongMap) Beats() []SongBeat {
	return s.beats
}

func SongMapReadFromFile(src io.Reader, audioManager AudioPositioner) *SongMap {
	m := NewSongMap(audioManager)
	rd := bufio.NewReader(src)
	_, _ = rd.ReadString('\n') // Ignore first line

	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Panicf("%s:\nFailed to read line of file!\n", err)
		}

		parts := strings.Split(strings.TrimSpace(line), ",")
		beatDuration, err := time.ParseDuration(parts[0])
		if err != nil {
			log.Panicf("%s:\nFailed to parse duration!\n", err)
		}

		var leftStatus PressStatus
		switch parts[1] {
		case "0":
			leftStatus = PressStatusEnum.IGNORE
		case "1":
			leftStatus = PressStatusEnum.LIFTED
		case "2":
			leftStatus = PressStatusEnum.PRESSED
		default:
			log.Panicf("Can't parse left status! %s", parts[1])
		}

		var rightStatus PressStatus
		switch parts[2] {
		case "0":
			rightStatus = PressStatusEnum.IGNORE
		case "1":
			rightStatus = PressStatusEnum.LIFTED
		case "2":
			rightStatus = PressStatusEnum.PRESSED
		default:
			log.Panicf("Can't parse left status! %s", parts[2])
		}

		m.beats = append(m.beats, SongBeat{
			Position:  beatDuration,
			LeftSide:  leftStatus,
			RightSide: rightStatus,
		})
	}

	return m
}

func (s *SongMap) WriteToFile(dst io.Writer) error {
	_, err := io.WriteString(dst, "position,left,right\n")
	if err != nil {
		return err
	}

	for _, b := range s.beats {
		_, err := fmt.Fprintf(dst, "%dus,%d,%d\n", b.Position.Microseconds(), b.LeftSide, b.RightSide)
		if err != nil {
			return err
		}
	}

	return nil
}
