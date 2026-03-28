package main

import (
	"embed"
	"flag"
	"image/color"
	"log"
	"time"

	rythmpen "github.com/ElrohirGT/RythmPen"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

//go:embed song.mp3
//go:embed song.map
var f embed.FS

var LeftColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
var RightColor = color.RGBA{R: 0, G: 0, B: 255, A: 255}

type Parameters struct {
	AudioSrc      string
	MapSrc        string
	AudioDuration time.Duration
}

var Params Parameters

func ParseParams() {
	flag.StringVar(&Params.AudioSrc, "audio", "song.mp3", "The path for the audio file")
	flag.StringVar(&Params.MapSrc, "map", "song.map", "The path for the map file")
	flag.DurationVar(&Params.AudioDuration, "duration", 7*time.Second, "The duration of the provided song")
	flag.Parse()
}

func main() {
	ParseParams()

	ebiten.SetWindowSize(int(rythmpen.WindowWidth), rythmpen.ComputeDiscreteHeight(rythmpen.WindowHeightWidthRatio, rythmpen.WindowWidth))
	ebiten.SetWindowTitle("RythmPen")
	const SampleRate = 48000

	debugManager := rythmpen.NewDebugImageManager(ebiten.KeyB)
	audioManager := rythmpen.NewAudioManager(SampleRate)

	leftPenImg := ebiten.NewImage(50, 100)
	leftPenImg.Fill(LeftColor)
	rightPenImg := ebiten.NewImage(50, 100)
	rightPenImg.Fill(RightColor)

	yCenter := rythmpen.WindowHeight / 2
	yPenCenterDelta := rythmpen.WindowHeight / 4
	yPenStart := yCenter - yPenCenterDelta
	yPenEnd := yCenter + yPenCenterDelta

	xCenter := rythmpen.WindowWidth / 2
	xPenCenterDelta := rythmpen.WindowWidth / 8
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

	rightBeatStart := rythmpen.NewVec2(rythmpen.WindowWidth, yPenEnd+100)
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

	audioSrc, err := f.Open(Params.AudioSrc)
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

	mapSrc, err := f.Open(Params.MapSrc)
	if err != nil {
		log.Panicf("%s\nFailed to read map source!\n", err)
	}

	songMap := rythmpen.SongMapReadFromFile(mapSrc, audioManager)
	songMapBeatCount := len(songMap.Beats())

	scoreManager := rythmpen.NewScoreManger(
		audioManager,
		beatManager,
		songMap,
		rythmpen.DefaultMaxBeatDelta,
		5.0,
		leftPen,
		rightPen,
	)

	game := &rythmpen.Game{
		LeftPen:      leftPen,
		RightPen:     rightPen,
		BeatManager:  beatManager,
		DebugManager: debugManager,
		AudioManager: audioManager,
		SongMap:      songMap,
		ScoreManager: scoreManager,
		SongDuration: Params.AudioDuration,
	}
	game.StartLevel()

	beatManagerBeatCount := len(beatManager.Beats())
	if songMapBeatCount != beatManagerBeatCount {
		log.Panicf("SongMap beats (%d) != (%d) BeatManager beats!\n", songMapBeatCount, beatManagerBeatCount)
	}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
