package rythmpen

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type PenState int

var PenStateEnum = struct {
	UP   PenState
	DOWN PenState
}{
	UP:   0,
	DOWN: 1,
}

type Pen struct {
	Image *ebiten.Image
	State PenState

	UpPosition   Vec2
	DownPosition Vec2

	ActivationKey ebiten.Key
}

func NewPen(image *ebiten.Image, upPosition, downPosition Vec2, actKey ebiten.Key) *Pen {
	return &Pen{
		Image:         image,
		UpPosition:    upPosition,
		DownPosition:  downPosition,
		State:         PenStateEnum.UP,
		ActivationKey: actKey,
	}
}

func (p *Pen) Update() {
	if ebiten.IsKeyPressed(p.ActivationKey) {
		p.State = PenStateEnum.DOWN
	} else {
		p.State = PenStateEnum.UP
	}
}

func (p Pen) Draw(parent *ebiten.Image, opt *ebiten.DrawImageOptions) {
	x, y := p.DownPosition.X, p.DownPosition.Y
	if p.State == PenStateEnum.UP {
		x, y = p.UpPosition.X, p.UpPosition.Y
	}
	opt.GeoM.Translate(x, y)
	parent.DrawImage(p.Image, opt)

	txtOpt := &text.DrawOptions{}
	txtOpt.GeoM.Translate(x, y)
	key := p.ActivationKey.String()
	text.Draw(parent, key, basicFace, txtOpt)

	opt.GeoM.Reset()
}
