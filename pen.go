package main

import "github.com/hajimehoshi/ebiten/v2"

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
	if p.State == PenStateEnum.UP {
		opt.GeoM.Translate(p.UpPosition.X, p.UpPosition.Y)
	} else {
		opt.GeoM.Translate(p.DownPosition.X, p.DownPosition.Y)
	}
	parent.DrawImage(p.Image, opt)
}
