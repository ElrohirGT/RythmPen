package rythmpen

import (
	"fmt"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type ScoreManager struct {
	currentScore float64
	positioner   AudioPositioner
	scoreMap     *SongMap
	currentIdx   int

	targetPen *Pen

	maxBeatDelta     time.Duration
	beatPoints       float64
	isLeft           bool
	prevFramePressed bool

	// Draw only
	penPressed bool
	diff       float64
	precision  float64
}

func NewScoreManger(
	audioManager AudioPositioner,
	scoreMap *SongMap,
	maxBeatDelta time.Duration,
	beatPoints float64,
	targetPen *Pen,
	isLeft bool,
) *ScoreManager {
	return &ScoreManager{
		currentScore: 0,
		positioner:   audioManager,
		scoreMap:     scoreMap,
		currentIdx:   0,
		maxBeatDelta: maxBeatDelta,
		beatPoints:   beatPoints,
		targetPen:    targetPen,
		isLeft:       isLeft,
	}
}

func (sm *ScoreManager) Score() float64 {
	return sm.currentScore
}

func (sm *ScoreManager) Update() {
	defer func() {
		sm.prevFramePressed = sm.penPressed
	}()
	sm.penPressed = false
	if sm.currentIdx >= len(sm.scoreMap.beats) {
		return
	}
	current := sm.scoreMap.beats[sm.currentIdx]

	isLeft := current.LeftSide == PressStatusEnum.PRESSED
	isRight := current.RightSide == PressStatusEnum.PRESSED
	sm.penPressed = sm.targetPen.State == PenStateEnum.DOWN

	shouldRegisterNewPress := sm.penPressed != sm.prevFramePressed
	if shouldRegisterNewPress && sm.penPressed { // Only activate if we press it! NOT on release
		if (isLeft && sm.isLeft && sm.penPressed) ||
			(isRight && !sm.isLeft && sm.penPressed) {
			fmt.Println("Pressed correct!")
			sm.UpdateScore()
			sm.currentIdx++
			return
		} else {
			fmt.Println("Pressed INCORRECT!")
		}
	}
	// We're already past the current beat so skip it.
	pos := sm.positioner.Position() + sm.maxBeatDelta
	if current.Position < pos {
		sm.currentIdx++
	}
}

func (sm *ScoreManager) Draw(parent *ebiten.Image, opt *ebiten.DrawImageOptions) {
	// if sm.currentIdx >= len(sm.scoreMap.beats) {
	// 	return
	// }
	//
	// current := sm.scoreMap.beats[sm.currentIdx]
	//
	// ebitenutil.DebugPrint(parent, fmt.Sprintf(
	// 	"Score: %.2f\nDiff: %.2f\nPrecision: %.2f\nPos: %.2f\nCurrent: %.2f\nPressed: %t\nLeft: %d, Rigth: %d",
	// 	sm.currentScore,
	// 	sm.diff,
	// 	sm.precision,
	// 	float64(sm.positioner.Position().Microseconds()),
	// 	float64(current.Position.Microseconds()),
	// 	sm.pressed,
	// 	current.LeftSide,
	// 	current.RightSide,
	// ))
}

func (sm *ScoreManager) UpdateScore() {
	pos := sm.positioner.Position()
	current := sm.scoreMap.beats[sm.currentIdx]

	sm.diff = math.Abs(float64(current.Position.Microseconds()) - float64(pos.Microseconds()))
	sm.precision = sm.diff / float64(sm.maxBeatDelta.Microseconds())
	sm.currentScore += sm.beatPoints * sm.precision
}
