package main

import (
	"flag"
	"image/color"
	"log"
	"os"
	"time"

	rythmpen "github.com/ElrohirGT/RythmPen"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var LeftColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
var RightColor = color.RGBA{R: 0, G: 0, B: 255, A: 255}

type Parameters struct {
	AudioSrc string
	MapDst   string
}

var Params Parameters

func ParseParams() {
	flag.StringVar(&Params.AudioSrc, "src", "source.mp3", "The audio source for the music")
	flag.StringVar(&Params.MapDst, "dst", "song.map", "The destination file for the recorded map")
	flag.Parse()
}

type Game struct {
	leftPen     *rythmpen.Pen
	rightPen    *rythmpen.Pen
	beatManager *rythmpen.BeatManager

	debugManager *rythmpen.DebugImageManager
	audioManager *rythmpen.AudioManager

	recorder *rythmpen.SongMap
}

func (g *Game) Update() error {
	g.leftPen.Update()
	g.rightPen.Update()

	g.beatManager.Update()
	g.debugManager.Update()

	leftJustPressed := inpututil.IsKeyJustPressed(g.leftPen.ActivationKey)
	rightJustPressed := inpututil.IsKeyJustPressed(g.rightPen.ActivationKey)

	if leftJustPressed && rightJustPressed {
		g.recorder.BothBeat()
	} else if leftJustPressed {
		g.recorder.LeftBeat()
	} else if rightJustPressed {
		g.recorder.RightBeat()
	}

	if ebiten.IsWindowBeingClosed() {
		SaveIntoFile(Params.MapDst, g.recorder)
	}

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
	ParseParams()

	ebiten.SetWindowSize(int(WindowWidth), ComputeDiscreteHeight(WindowHeightWidthRatio, WindowWidth))
	ebiten.SetWindowTitle("RythmPen")

	debugManager := rythmpen.NewDebugImageManager(ebiten.KeyB)

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
	leftPen := rythmpen.NewPen(
		leftPenImg,
		rythmpen.NewVec2(leftX, yPenStart),
		rythmpen.NewVec2(leftX, yPenEnd),
		ebiten.KeyF,
	)

	rightX := float64(xCenter + xPenCenterDelta)
	rightPen := rythmpen.NewPen(
		rightPenImg,
		rythmpen.NewVec2(rightX, yPenStart),
		rythmpen.NewVec2(rightX, yPenEnd),
		ebiten.KeyJ,
	)

	beatWidth := 50
	beatHeight := 20

	rightBeatStart := rythmpen.NewVec2(WindowWidth, yPenEnd+100)
	rightBeatEnd := rythmpen.NewVec2(rightX, yPenEnd+100)

	leftBeatStart := rythmpen.NewVec2(-float64(beatWidth), yPenEnd+100)
	leftBeatEnd := rythmpen.NewVec2(leftX, yPenEnd+100)

	rightBeatImage := ebiten.NewImage(beatWidth, beatHeight)
	rightBeatImage.Fill(RightColor)

	leftBeatImage := ebiten.NewImage(beatWidth, beatHeight)
	leftBeatImage.Fill(LeftColor)

	debugManager.Add(rythmpen.NewDebugImage(rightBeatStart))
	debugManager.Add(rythmpen.NewDebugImage(rightBeatEnd))

	debugManager.Add(rythmpen.NewDebugImage(leftBeatStart))
	debugManager.Add(rythmpen.NewDebugImage(leftBeatEnd))

	beatManager := rythmpen.NewBeatManager(
		rythmpen.BeatConfig{
			Image:    leftBeatImage,
			Start:    leftBeatStart,
			End:      leftBeatEnd,
			LifeSpan: 3 * time.Second,
		},
		rythmpen.BeatConfig{
			Image:    rightBeatImage,
			Start:    rightBeatStart,
			End:      rightBeatEnd,
			LifeSpan: 3 * time.Second,
		},
	)

	const SampleRate = 44100
	audioManager := rythmpen.NewAudioManager(SampleRate)
	audioSrc, err := os.Open(Params.AudioSrc)
	if err != nil {
		log.Panicf("%s\nFailed to create reader from file!\n", err)
	}
	defer audioSrc.Close()

	mp3Stream, err := mp3.DecodeWithSampleRate(SampleRate, audioSrc)
	if err != nil {
		log.Panicf("%s\nFailed to decode mp3!\n", err)
	}

	err = audioManager.NewAudioPlayer(mp3Stream)
	if err != nil {
		log.Panicf("%s\nFailed to create audio player!\n", err)
	}

	audioManager.Play()
	mapRecorder := rythmpen.NewSongMap(audioManager)

	game := &Game{
		leftPen:      leftPen,
		rightPen:     rightPen,
		beatManager:  beatManager,
		debugManager: debugManager,
		audioManager: audioManager,
		recorder:     mapRecorder,
	}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

func SaveIntoFile(filePath string, record *rythmpen.SongMap) {
	dst, err := os.Create(filePath)
	if err != nil {
		log.Panicf("%s\nFailed to open file to write song map!", err)
	}
	defer dst.Close()

	err = record.WriteToFile(dst)
	if err != nil {
		log.Panicf("%s\nFailed to write song map!", err)
	}

	log.Println("File saved!")
}
