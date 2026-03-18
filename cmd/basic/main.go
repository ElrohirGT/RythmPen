package main

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var LeftColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
var RightColor = color.RGBA{R: 0, G: 0, B: 255, A: 255}

type Game struct {
	leftPen     *Pen
	rightPen    *Pen
	beatManager *BeatManager

	debugManager *DebugImageManager
	audioManager *AudioManager
}

func (g *Game) Update() error {
	g.leftPen.Update()
	g.rightPen.Update()

	g.beatManager.Update()
	g.debugManager.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	g.leftPen.Draw(screen, op)
	g.rightPen.Draw(screen, op)

	g.beatManager.Draw(screen, op)
	g.debugManager.Draw(screen, op)
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

	debugManager := NewDebugImageManager(ebiten.KeyB)

	leftPenImg := ebiten.NewImage(50, 100)
	leftPenImg.Fill(LeftColor)
	rightPenImg := ebiten.NewImage(50, 100)
	rightPenImg.Fill(RightColor)

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

	beatWidth := 50
	beatHeight := 20

	rightBeatStart := NewVec2(WindowWidth, yPenEnd+100)
	rightBeatEnd := NewVec2(rightX, yPenEnd+100)

	leftBeatStart := NewVec2(-float64(beatWidth), yPenEnd+100)
	leftBeatEnd := NewVec2(leftX, yPenEnd+100)

	rightBeatImage := ebiten.NewImage(beatWidth, beatHeight)
	rightBeatImage.Fill(RightColor)
	rightBeat := NewBeat(
		rightBeatImage,
		rightBeatStart,
		rightBeatEnd,
		2*time.Second,
	)

	leftBeatImage := ebiten.NewImage(beatWidth, beatHeight)
	leftBeatImage.Fill(LeftColor)
	leftBeat := NewBeat(
		leftBeatImage,
		leftBeatStart,
		leftBeatEnd,
		2*time.Second,
	)

	debugManager.Add(NewDebugImage(rightBeatStart))
	debugManager.Add(NewDebugImage(rightBeatEnd))

	debugManager.Add(NewDebugImage(leftBeatStart))
	debugManager.Add(NewDebugImage(leftBeatEnd))

	beatManager := NewBeatManager(
		BeatConfig{
			Image:    leftBeatImage,
			Start:    leftBeatStart,
			End:      leftBeatEnd,
			LifeSpan: 3 * time.Second,
		},
		BeatConfig{
			Image:    rightBeatImage,
			Start:    rightBeatStart,
			End:      rightBeatEnd,
			LifeSpan: 3 * time.Second,
		},
	)
	beatManager.AddBeat(leftBeat)
	beatManager.AddBeat(rightBeat)

	beatManager.AddRightBeat()
	beatManager.AddRightBeat()

	audioManager := NewAudioManager(44100)

	game := &Game{
		leftPen:      leftPen,
		rightPen:     rightPen,
		beatManager:  beatManager,
		debugManager: debugManager,
		audioManager: audioManager,
	}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
