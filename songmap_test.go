package rythmpen

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func Test_LeftBeat(t *testing.T) {
	past := time.Now()
	positioner := DummyAudioPositioner{
		D_Position: func() time.Duration {
			return time.Since(past)
		},
	}

	songMapper := NewSongMap(positioner)
	songMapper.LeftBeat()
	time.Sleep(4 * time.Millisecond)

	songMapper.RightBeat()
	time.Sleep(4 * time.Millisecond)

	songMapper.BothBeat()

	out := bytes.Buffer{}
	err := songMapper.WriteToFile(&out)
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")

	expectedHeader := "position,left,right"
	if string(lines[0]) != expectedHeader {
		t.Fatalf("header don't match!\n%s!=%s\n", expectedHeader, lines[0])
	}

	if len(lines) != 4 {
		t.Fatalf("Line length don't match!\n%#v", lines)
	}
}
