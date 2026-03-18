package main

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	frameCount int64
	leftPen    *Pen
	rightPen   *Pen
	beats      []*Beat
}

func (g *Game) Update() error {
	g.frameCount++

	g.leftPen.Update()
	g.rightPen.Update()

	for i, b := range g.beats {
		if shouldRemove := b.Update(); shouldRemove {
			g.beats = SlicesRemoveWithoutOrder(g.beats, i)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	g.leftPen.Draw(screen, op)
	op.GeoM.Reset()
	g.rightPen.Draw(screen, op)
	op.GeoM.Reset()

	for _, b := range g.beats {
		b.Draw(screen, op)
		op.GeoM.Reset()
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, ComputeDiscreteHeight(WindowHeightWidthRatio, float64(outsideWidth))
}

const WindowHeightWidthRatio float64 = 1080.0 / 1920.0
const WindowWidth float64 = 1500.0
const WindowHeight float64 = WindowWidth * WindowHeightWidthRatio

func ComputeDiscreteHeight(heightWidthRatio float64, width float64) int {
	ab := width * WindowHeightWidthRatio
	return int(ab)
}

func main() {
	ebiten.SetWindowSize(int(WindowWidth), ComputeDiscreteHeight(WindowHeightWidthRatio, WindowWidth))
	ebiten.SetWindowTitle("RythmPen")

	leftPenImg := ebiten.NewImage(50, 100)
	leftPenImg.Fill(color.RGBA{R: 255, G: 0, B: 0, A: 255})
	rightPenImg := ebiten.NewImage(50, 100)
	rightPenImg.Fill(color.RGBA{R: 0, G: 0, B: 255, A: 255})

	yCenter := WindowHeight / 2
	yPenCenterDelta := WindowHeight / 4
	yPenStart := yCenter - yPenCenterDelta
	yPenEnd := yCenter + yPenCenterDelta

	xCenter := WindowWidth / 2
	xPenCenterDelta := WindowWidth / 8
	leftX := float64(xCenter - xPenCenterDelta)
	leftPen := NewPen(
		leftPenImg,
		NewVec2(leftX, yPenStart),
		NewVec2(leftX, yPenEnd),
		ebiten.KeyF,
	)

	rightX := float64(xCenter + xPenCenterDelta)
	rightPen := NewPen(
		rightPenImg,
		NewVec2(rightX, yPenStart),
		NewVec2(rightX, yPenEnd),
		ebiten.KeyJ,
	)

	beatImage := ebiten.NewImage(50, 20)
	beatImage.Fill(color.RGBA{R: 255, G: 0, B: 255, A: 255})
	beat := NewBeat(
		beatImage,
		NewVec2(WindowWidth, yPenEnd+100),
		NewVec2(leftX, yPenEnd+100),
		5*time.Second,
	)

	game := &Game{
		leftPen:  leftPen,
		rightPen: rightPen,
		beats:    []*Beat{beat},
	}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
