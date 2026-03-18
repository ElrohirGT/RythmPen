package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	leftPen  *ebiten.Image
	rightPen *ebiten.Image
)

type Game struct {
	frameCount int64
}

func (g *Game) Update() error {
	g.frameCount++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	const limit = int64(WindowWidth)
	rrow := float64(g.frameCount / limit)
	ccol := float64(g.frameCount)
	if g.frameCount >= limit {
		ccol = float64(g.frameCount % limit)
	}
	op.GeoM.Translate(ccol, rrow)
	screen.DrawImage(leftPen, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, ComputeDiscreteHeight(WindowHeightWidthRatio, float32(outsideWidth))
}

const WindowHeightWidthRatio float32 = 1080.0 / 1920.0
const WindowWidth float32 = 1500.0

func ComputeDiscreteHeight(heightWidthRatio float32, width float32) int {
	ab := width * WindowHeightWidthRatio
	return int(ab)
}

func main() {
	ebiten.SetWindowSize(int(WindowWidth), ComputeDiscreteHeight(WindowHeightWidthRatio, WindowWidth))
	ebiten.SetWindowTitle("Test")

	leftPen = ebiten.NewImage(50, 100)
	leftPen.Fill(color.RGBA{R: 255, G: 0, B: 0, A: 255})
	rightPen = ebiten.NewImage(50, 100)
	rightPen.Fill(color.RGBA{R: 0, G: 0, B: 255, A: 255})

	game := &Game{}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
