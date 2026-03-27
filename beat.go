package rythmpen

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type BeatState int

var failedBeatImage *ebiten.Image
var badBeatImage *ebiten.Image
var goodBeatImage *ebiten.Image
var perfectBeatImage *ebiten.Image

type Beat struct {
	/* Rendering */

	Index     int
	Image     *ebiten.Image
	IsPlucked bool

	/* Movement */

	StartPos     Vec2
	EndPos       Vec2
	CurrentPos   Vec2
	Positioner   AudioPositioner
	LifeSpan     time.Duration
	MaxBeatDelta time.Duration
}

func NewBeat(
	image *ebiten.Image,
	startPos, endPos Vec2,
	lifeSpan time.Duration,
	index int,
	positioner AudioPositioner,
	maxBeatDelta time.Duration,
) *Beat {
	if failedBeatImage == nil {
		failedBeatImage = ebiten.NewImage(image.Bounds().Dx(), image.Bounds().Dy())
		failedBeatImage.Fill(color.Black)

		badBeatImage = ebiten.NewImage(image.Bounds().Dx(), image.Bounds().Dy())
		badBeatImage.Fill(color.RGBA{R: 100, G: 100, B: 100, A: 255})

		goodBeatImage = ebiten.NewImage(image.Bounds().Dx(), image.Bounds().Dy())
		goodBeatImage.Fill(color.RGBA{R: 200, G: 200, B: 10, A: 255})

		perfectBeatImage = ebiten.NewImage(image.Bounds().Dx(), image.Bounds().Dy())
		perfectBeatImage.Fill(color.RGBA{R: 10, G: 200, B: 10, A: 255})
	}

	return &Beat{
		Index:        index,
		StartPos:     startPos,
		EndPos:       endPos,
		LifeSpan:     lifeSpan,
		Image:        image,
		Positioner:   positioner,
		MaxBeatDelta: maxBeatDelta,
	}
}

func (b *Beat) PluckWithPrecision(precision float64) {
	b.IsPlucked = true
	var preffix string
	if precision >= 0.8 { // Excellent
		preffix = "PERFECT!"
		b.Image = perfectBeatImage
	} else if precision >= 0.5 { // Good
		preffix = "GOOD!"
		b.Image = goodBeatImage
	} else { // Bad
		preffix = "BAD!"
		b.Image = badBeatImage
	}

	log.Printf("%s: %.2f\n", preffix, precision)
}

func (b *Beat) FailedPluck() { // Failed to pluck the beat
	log.Println("MISSED!")
	b.IsPlucked = true
	b.Image = failedBeatImage
}

func (b *Beat) Update() bool {
	currentPosition := b.Positioner.Position()
	t := float64(currentPosition) / float64(b.LifeSpan)
	if !b.IsPlucked {
		b.CurrentPos = Vec2Lerp(b.StartPos, b.EndPos, t)
	}
	return currentPosition > (b.LifeSpan + b.MaxBeatDelta)
}

func (b Beat) Draw(parent *ebiten.Image, opt *ebiten.DrawImageOptions) {
	opt.GeoM.Translate(b.CurrentPos.X, b.CurrentPos.Y)
	parent.DrawImage(b.Image, opt)

	txtOpt := &text.DrawOptions{}
	txtOpt.GeoM.Translate(b.CurrentPos.X, b.CurrentPos.Y)
	text.Draw(parent, fmt.Sprintf("%d", b.Index), basicFace, txtOpt)
	// ebitenutil.DebugPrint(parent, fmt.Sprintf("%f %f", b.CurrentPos.X, b.CurrentPos.Y))
}

type BeatConfig struct {
	Image          *ebiten.Image
	PixelsPerMicro float64
	End            Vec2
	Positioner     AudioPositioner
	MaxDelta       time.Duration
}

type BeatManager struct {
	beats           []*Beat
	leftBeatConfig  BeatConfig
	rightBeatConfig BeatConfig
	currentIdx      int
}

func NewBeatManager(leftBeatConfig, rightBeatConfig BeatConfig) *BeatManager {
	return &BeatManager{
		beats:           make([]*Beat, 0, 50),
		leftBeatConfig:  leftBeatConfig,
		rightBeatConfig: rightBeatConfig,
	}
}

func (manager *BeatManager) Beat(idx int) *Beat {
	return manager.beats[idx]
}

func (manager *BeatManager) Beats() []*Beat {
	return manager.beats
}

func (manager *BeatManager) Update() {
	for i := manager.currentIdx; i < len(manager.beats); i++ {
		b := manager.beats[i]
		if shouldRemove := b.Update(); shouldRemove {
			manager.currentIdx = i + 1
		}
	}
}

func (manager *BeatManager) Draw(parent *ebiten.Image, opt *ebiten.DrawImageOptions) {
	for i := manager.currentIdx; i < len(manager.beats); i++ {
		beat := manager.beats[i]
		beat.Draw(parent, opt)
		opt.GeoM.Reset()
	}
}

func (manager *BeatManager) AddBeat(beat *Beat) {
	manager.beats = append(manager.beats, beat)
}
func (manager *BeatManager) AddBeatWithConfig(config BeatConfig, startPos Vec2, lifeSpan time.Duration) {
	b := NewBeat(
		config.Image,
		startPos,
		config.End,
		lifeSpan,
		len(manager.beats),
		config.Positioner,
		config.MaxDelta,
	)
	manager.AddBeat(b)
}

func (manager *BeatManager) AddLeftBeat(lifeSpan time.Duration) {
	config := manager.leftBeatConfig
	startPosX := config.End.X - config.PixelsPerMicro*float64(lifeSpan.Microseconds())
	manager.AddBeatWithConfig(config, NewVec2(startPosX, config.End.Y), lifeSpan)
}

func (manager *BeatManager) AddRightBeat(lifeSpan time.Duration) {
	config := manager.rightBeatConfig
	startPosX := config.End.X + config.PixelsPerMicro*float64(lifeSpan.Microseconds())
	manager.AddBeatWithConfig(config, NewVec2(startPosX, config.End.Y), lifeSpan)
}
