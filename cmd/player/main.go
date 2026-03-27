package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	rythmpen "github.com/ElrohirGT/RythmPen"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var LeftColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
var RightColor = color.RGBA{R: 0, G: 0, B: 255, A: 255}

type Parameters struct {
	AudioSrc string
	MapSrc   string
}

var Params Parameters

func ParseParams() {
	flag.StringVar(&Params.AudioSrc, "audio", "song.mp3", "The path for the audio file")
	flag.StringVar(&Params.MapSrc, "map", "song.map", "The path for the map file")
	flag.Parse()
}

type Game struct {
	leftPen     *rythmpen.Pen
	rightPen    *rythmpen.Pen
	beatManager *rythmpen.BeatManager

	debugManager *rythmpen.DebugImageManager
	audioManager *rythmpen.AudioManager

	songMap      *rythmpen.SongMap
	scoreManager *rythmpen.ScoreManager
}

func (g *Game) Update() error {
	g.debugManager.Update()

	g.beatManager.Update()
	g.leftPen.Update()
	g.rightPen.Update()

	g.scoreManager.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	g.debugManager.Draw(screen, op)

	g.leftPen.Draw(screen, op)
	g.rightPen.Draw(screen, op)

	g.beatManager.Draw(screen, op)
	g.scoreManager.Draw(screen, op)

	score := g.scoreManager.Score()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %.2f", score))
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

	const SampleRate = 44100
	audioManager := rythmpen.NewAudioManager(SampleRate)
	pixelsPerMicro := 0.2 / float64(time.Microsecond)
	beatManager := rythmpen.NewBeatManager(
		rythmpen.BeatConfig{
			Image:          leftBeatImage,
			End:            leftBeatEnd,
			PixelsPerMicro: pixelsPerMicro,
			Positioner:     audioManager,
			MaxDelta:       rythmpen.DefaultMaxBeatDelta,
		},
		rythmpen.BeatConfig{
			Image:          rightBeatImage,
			End:            rightBeatEnd,
			PixelsPerMicro: pixelsPerMicro,
			Positioner:     audioManager,
			MaxDelta:       rythmpen.DefaultMaxBeatDelta,
		},
	)

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

	mapSrc, err := os.Open(Params.MapSrc)
	if err != nil {
		log.Panicf("%s\nFailed to read map source!\n", err)
	}

	songMap := rythmpen.SongMapReadFromFile(mapSrc, audioManager)
	songMapBeatCount := len(songMap.Beats())
	for _, b := range songMap.Beats() {
		lifeSpan := b.Position
		log.Printf("Beat: %#v\n", b)

		if b.LeftSide == rythmpen.PressStatusEnum.PRESSED {
			log.Printf("Left!")
			beatManager.AddLeftBeat(lifeSpan)
		}

		if b.RightSide == rythmpen.PressStatusEnum.PRESSED {
			log.Printf("Right!")
			beatManager.AddRightBeat(lifeSpan)
		}
	}
	beatManagerBeatCount := len(beatManager.Beats())
	fmt.Println("Added", beatManagerBeatCount, "beats")
	if songMapBeatCount != beatManagerBeatCount {
		log.Panicf("SongMap beats (%d) != (%d) BeatManager beats!\n", songMapBeatCount, beatManagerBeatCount)
	}

	scoreManager := rythmpen.NewScoreManger(
		audioManager,
		beatManager,
		songMap,
		rythmpen.DefaultMaxBeatDelta,
		5.0,
		leftPen,
		rightPen,
		true,
	)

	game := &Game{
		leftPen:      leftPen,
		rightPen:     rightPen,
		beatManager:  beatManager,
		debugManager: debugManager,
		audioManager: audioManager,
		songMap:      songMap,
		scoreManager: scoreManager,
	}
	audioManager.Play()
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
