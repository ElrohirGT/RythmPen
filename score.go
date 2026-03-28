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
	}
}

func (sm *ScoreManager) Reset() {
	sm.currentIdx = 0
	sm.currentScore = 0
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
	clampedDifference := Float64Clamp(0, 1, math.Abs(scaledDifference))

	sm.precision = 1.0 - clampedDifference
	log.Printf(
		"[DIFF: %2dms] \t 1 - clamp(0, 1, abs(abs(%d - %d) / %d)) = %.2f\n",
		(current.Position - pos).Round(time.Millisecond).Milliseconds(),
		current.Position.Microseconds(),
		pos.Microseconds(),
		sm.maxBeatDelta.Microseconds(),
		sm.precision,
	)
	sm.currentScore += sm.beatPoints * sm.precision
}
