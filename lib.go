package rythmpen

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Vec2 struct {
	X float64
	Y float64
}

func NewVec2(x, y float64) Vec2 {
	return Vec2{X: x, Y: y}
}

func Vec2Lerp(start, end Vec2, t float64) Vec2 {
	return NewVec2(
		Float64Lerp(start.X, end.X, t),
		Float64Lerp(start.Y, end.Y, t),
	)
}

func (v Vec2) Vec2TranslateGeom(opt *ebiten.DrawImageOptions) {
	opt.GeoM.Translate(v.X, v.Y)
}

func Float64Lerp(start, end, t float64) float64 {
	return start + t*(end-start)
}

func Float64Aproximately(a, b float64) bool {
	return a-b <= math.SmallestNonzeroFloat64 || b-a <= math.SmallestNonzeroFloat64
}

func Float64Clamp(maximum, minimum, value float64) float64 {
	if value < minimum {
		return minimum
	}

	if value > maximum {
		return maximum
	}

	return value
}

func SlicesRemoveWithoutOrder[T any](slice []T, idx int) []T {
	if len(slice) > 0 {
		slice[idx] = slice[len(slice)-1]
		return slice[:len(slice)-1]
	} else {
		return slice
	}
}

type Optional[T any] struct {
	value   T
	isValid bool
}

func NewOptionalValue[T any](obj T) Optional[T] {
	return Optional[T]{
		value:   obj,
		isValid: true,
	}
}

func (op Optional[T]) IsValid() bool {
	return op.isValid
}

func (op Optional[T]) Value() T {
	if !op.IsValid() {
		log.Panicf("Can't get value of invalid optional!")
	}

	return op.value
}
