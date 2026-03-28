package rythmpen

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type GameScreen int

const (
	GameScreen_Game GameScreen = iota
	GameScreen_End
)

type Game struct {
	Screen      GameScreen
	LeftPen     *Pen
	RightPen    *Pen
	BeatManager *BeatManager

	DebugManager *DebugImageManager
	AudioManager *AudioManager

	SongMap      *SongMap
	ScoreManager *ScoreManager
}

func (g *Game) Update() error {
	if g.Screen == GameScreen_Game {
		g.DebugManager.Update()

		g.BeatManager.Update()
		g.LeftPen.Update()
		g.RightPen.Update()

		g.ScoreManager.Update()
		if !g.AudioManager.IsPlaying() {
			g.Screen = GameScreen_End
		}
	} else if g.Screen == GameScreen_End {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.StartLevel()
		}
	}

	return nil
}

func (g *Game) StartLevel() {
	g.Screen = GameScreen_Game
	g.ScoreManager.Reset()

	g.BeatManager.Reset()
	g.BeatManager.SpawnBeats(g.SongMap)

	g.AudioManager.Reset()
	g.AudioManager.Play()
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.Screen == GameScreen_Game {

		op := &ebiten.DrawImageOptions{}
		g.DebugManager.Draw(screen, op)

		g.LeftPen.Draw(screen, op)
		g.RightPen.Draw(screen, op)

		g.BeatManager.Draw(screen, op)
		g.ScoreManager.Draw(screen, op)

		score := g.ScoreManager.Score()
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %.2f", score))
	} else if g.Screen == GameScreen_End {
		score := g.ScoreManager.Score()
		txtOpt := &text.DrawOptions{}
		txtOpt.GeoM.Translate(WindowWidth/2, WindowHeight/2)
		text.Draw(screen, fmt.Sprintf("%.2f", score), basicFace, txtOpt)
		txtOpt.GeoM.Translate(0, 20)
		text.Draw(screen, "Press R to Restart", basicFace, txtOpt)
	}
}

const WindowHeightWidthRatio float64 = 1080.0 / 1920.0
const WindowWidth float64 = 1500.0
const WindowHeight float64 = WindowWidth * WindowHeightWidthRatio

func ComputeDiscreteHeight(heightWidthRatio float64, width float64) int {
	ab := width * WindowHeightWidthRatio
	return int(ab)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, ComputeDiscreteHeight(WindowHeightWidthRatio, float64(outsideWidth))
}
