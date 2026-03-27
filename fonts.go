package rythmpen

import (
	"bytes"
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/text/language"
)

//go:embed fonts/Tiny5-Regular.ttf
var pixelArtTTF []byte
var pixelArtSource *text.GoTextFaceSource
var basicFace *text.GoTextFace

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(pixelArtTTF))
	if err != nil {
		log.Fatal(err)
	}
	pixelArtSource = s

	basicFace = &text.GoTextFace{
		Source:    pixelArtSource,
		Direction: text.DirectionLeftToRight,
		Size:      24,
		Language:  language.Spanish,
	}
}
