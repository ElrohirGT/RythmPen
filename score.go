package rythmpen

import (
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type ScoreManager struct {
	currentScore float64
	positioner   AudioPositioner
	scoreMap     *SongMap
	currentIdx   int
	beatManager  *BeatManager

	leftPen  *Pen
	rightPen *Pen

	maxBeatDelta     time.Duration
	beatPoints       float64
	isLeft           bool
	prevFramePressed bool

	// Draw only
	anyPenIsPressed bool
	diff            float64
	precision       float64
}

func NewScoreManger(
	audioManager AudioPositioner,
	beatManager *BeatManager,
	scoreMap *SongMap,
	maxBeatDelta time.Duration,
	beatPoints float64,
	leftPen *Pen,
	rightPen *Pen,
	isLeft bool,
) *ScoreManager {
	return &ScoreManager{
		currentScore: 0,
		beatManager:  beatManager,
		positioner:   audioManager,
		scoreMap:     scoreMap,
		currentIdx:   0,
		maxBeatDelta: maxBeatDelta,
		beatPoints:   beatPoints,
		leftPen:      leftPen,
		rightPen:     rightPen,
		isLeft:       isLeft,
	}
}

func (sm *ScoreManager) Score() float64 {
	return sm.currentScore
}

func (sm *ScoreManager) Update() {
	defer func() {
		sm.prevFramePressed = sm.anyPenIsPressed
	}()
	sm.anyPenIsPressed = false
	if sm.currentIdx >= len(sm.scoreMap.beats) {
		return
	}
	current := sm.scoreMap.beats[sm.currentIdx]
	currentBeat := sm.beatManager.Beat(sm.currentIdx)

	isLeftBeat := current.LeftSide == PressStatusEnum.PRESSED
	isRightBeat := current.RightSide == PressStatusEnum.PRESSED
	isBothBeat := isLeftBeat && isRightBeat

	leftPenIsPressed := sm.leftPen.State == PenStateEnum.DOWN
	rightPenIsPressed := sm.rightPen.State == PenStateEnum.DOWN
	bothPenPressed := leftPenIsPressed && rightPenIsPressed
	sm.anyPenIsPressed = leftPenIsPressed || rightPenIsPressed

	shouldRegisterNewPress := sm.anyPenIsPressed != sm.prevFramePressed
	if shouldRegisterNewPress && sm.anyPenIsPressed { // Only activate if we press it! NOT on release
		if (isBothBeat && bothPenPressed) ||
			(isLeftBeat && leftPenIsPressed) ||
			(isRightBeat && rightPenIsPressed) {
			sm.UpdateScore()
			currentBeat.PluckWithPrecision(sm.precision)
			sm.currentIdx++
			return
		} else {
			currentBeat.FailedPluck()
		}
	}
	pos := sm.positioner.Position() - sm.maxBeatDelta
	// We're already past the current beat so skip it.
	if current.Position < pos {
		sm.currentIdx++
	}
}

func (sm *ScoreManager) Draw(parent *ebiten.Image, opt *ebiten.DrawImageOptions) {
}

func (sm *ScoreManager) UpdateScore() {
	pos := sm.positioner.Position()
	current := sm.scoreMap.beats[sm.currentIdx]

	sm.diff = math.Abs(float64(current.Position.Microseconds()) - float64(pos.Microseconds()))
	scaledDifference := sm.diff / float64(sm.maxBeatDelta.Microseconds())
	clampedDifference := Float64Clamp(1.0, 0.0, math.Abs(scaledDifference))
	log.Println("Difference:", current.Position-pos, "Scaled:", scaledDifference, "Abs and Clamped:", clampedDifference)

	sm.precision = 1.0 - clampedDifference
	sm.currentScore += sm.beatPoints * sm.precision
}
