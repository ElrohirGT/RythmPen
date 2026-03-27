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

// func init() {
// 	failedBeatImage = ebiten.NewImage(50, 20)
// 	grayRGB := 0
//
// 	failedBeatImage.Fill(color.RGBA{R: uint8(grayRGB), G: uint8(grayRGB), B: uint8(grayRGB), A: 255})
// }

type Beat struct {
	Index      int
	StartPos   Vec2
	EndPos     Vec2
	CurrentPos Vec2
	Expiration time.Time
	LifeSpan   time.Duration
	Image      *ebiten.Image
	Radians    float64
	IsPlucked  bool
}

func NewBeat(
	image *ebiten.Image,
	startPos, endPos Vec2,
	lifeSpan time.Duration,
	index int,
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
		Index:      index,
		StartPos:   startPos,
		EndPos:     endPos,
		LifeSpan:   lifeSpan,
		Image:      image,
		Expiration: time.Now().Add(lifeSpan),
	}
}

func (b *Beat) PluckWithPrecision(precision float64) {
	b.IsPlucked = true
	b.Expiration = b.Expiration.Add(1 * time.Second)
	var preffix string
	if precision >= 0.9 { // Excelent
		preffix = "PERFECT!"
		b.Image = perfectBeatImage
	} else if precision >= 0.5 { // Good
		preffix = "GOOD!"
		b.Image = goodBeatImage
	} else { // Bad
		preffix = "BAD!"
		b.Image = badBeatImage
	}

	log.Printf("%s: %f\n", preffix, precision)
}

func (b *Beat) FailedPluck() { // The pluck was failed
	b.Expiration = b.Expiration.Add(1 * time.Second)
	log.Println("MISSED!")
	b.IsPlucked = true
	b.Image = failedBeatImage
}

func (b *Beat) Update() bool {
	if !b.IsPlucked {
		t := 1 - float64(time.Until(b.Expiration).Microseconds())/float64(b.LifeSpan.Microseconds())
		b.CurrentPos = Vec2Lerp(b.StartPos, b.EndPos, t)
	}
	return time.Until(b.Expiration).Microseconds() <= 0
}

func (b Beat) Draw(parent *ebiten.Image, opt *ebiten.DrawImageOptions) {
	opt.GeoM.Translate(b.CurrentPos.X, b.CurrentPos.Y)
	opt.GeoM.Rotate(b.Radians)
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
}

type BeatManager struct {
	beats           []*Beat
	leftBeatConfig  BeatConfig
	rightBeatConfig BeatConfig
	currentIdx      int
	maxDelta        time.Duration
}

func NewBeatManager(maxDelta time.Duration, leftBeatConfig, rightBeatConfig BeatConfig) *BeatManager {
	return &BeatManager{
		beats:           make([]*Beat, 0, 50),
		leftBeatConfig:  leftBeatConfig,
		rightBeatConfig: rightBeatConfig,
		maxDelta:        maxDelta,
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
	b := NewBeat(config.Image, startPos, config.End, lifeSpan, len(manager.beats))
	manager.AddBeat(b)
}

func (manager *BeatManager) AddLeftBeat(lifeSpan time.Duration) {
	config := manager.leftBeatConfig
	startPosX := config.End.X - config.PixelsPerMicro*float64(lifeSpan.Microseconds())
	// config.End.X = config.End.X + config.PixelsPerMicro*float64(manager.maxDelta.Microseconds())
	manager.AddBeatWithConfig(config, NewVec2(startPosX, config.End.Y), lifeSpan)
}

func (manager *BeatManager) AddRightBeat(lifeSpan time.Duration) {
	config := manager.rightBeatConfig
	startPosX := config.End.X + config.PixelsPerMicro*float64(lifeSpan.Microseconds())
	// config.End.X = config.End.X - config.PixelsPerMicro*float64(manager.maxDelta.Microseconds())
	manager.AddBeatWithConfig(config, NewVec2(startPosX, config.End.Y), lifeSpan)
}
