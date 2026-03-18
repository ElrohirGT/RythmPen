package rythmpen

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Beat struct {
	StartPos   Vec2
	EndPos     Vec2
	CurrentPos Vec2
	Expiration time.Time
	LifeSpan   time.Duration
	Image      *ebiten.Image
}

func NewBeat(
	image *ebiten.Image,
	startPos, endPos Vec2,
	lifeSpan time.Duration,
) *Beat {
	return &Beat{
		StartPos:   startPos,
		EndPos:     endPos,
		LifeSpan:   lifeSpan,
		Image:      image,
		Expiration: time.Now().Add(lifeSpan),
	}
}

func (b *Beat) Update() bool {
	t := 1 - float64(time.Until(b.Expiration).Microseconds())/float64(b.LifeSpan.Microseconds())
	b.CurrentPos = Vec2Lerp(b.StartPos, b.EndPos, t)
	return time.Until(b.Expiration).Microseconds() <= 0
}

func (b Beat) Draw(parent *ebiten.Image, opt *ebiten.DrawImageOptions) {
	opt.GeoM.Translate(b.CurrentPos.X, b.CurrentPos.Y)
	parent.DrawImage(b.Image, opt)
	// ebitenutil.DebugPrint(parent, fmt.Sprintf("%f %f", b.CurrentPos.X, b.CurrentPos.Y))
}

type BeatConfig struct {
	Image    *ebiten.Image
	Start    Vec2
	End      Vec2
	LifeSpan time.Duration
}

type BeatManager struct {
	beats           []*Beat
	leftBeatConfig  BeatConfig
	rightBeatConfig BeatConfig
}

func NewBeatManager(leftBeatConfig, rightBeatConfig BeatConfig) *BeatManager {
	return &BeatManager{
		beats:           make([]*Beat, 0, 50),
		leftBeatConfig:  leftBeatConfig,
		rightBeatConfig: rightBeatConfig,
	}
}

func (manager *BeatManager) Update() {
	for i := 0; i < len(manager.beats); i++ {
		b := manager.beats[i]
		if shouldRemove := b.Update(); shouldRemove {
			manager.beats = SlicesRemoveWithoutOrder(manager.beats, i)
			i--
		}
	}
}

func (manager *BeatManager) Draw(parent *ebiten.Image, opt *ebiten.DrawImageOptions) {
	for _, b := range manager.beats {
		b.Draw(parent, opt)
		opt.GeoM.Reset()
	}
}

func (manager *BeatManager) AddBeat(beat *Beat) {
	manager.beats = append(manager.beats, beat)
}
func (manager *BeatManager) AddBeatWithConfig(config BeatConfig) {
	b := NewBeat(config.Image, config.Start, config.End, config.LifeSpan)
	manager.AddBeat(b)
}

func (manager *BeatManager) AddLeftBeat() {
	manager.AddBeatWithConfig(manager.leftBeatConfig)
}

func (manager *BeatManager) AddRightBeat() {
	manager.AddBeatWithConfig(manager.rightBeatConfig)
}
