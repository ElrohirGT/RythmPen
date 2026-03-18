package rythmpen

import "testing"

func Test_Float64Lerp(t *testing.T) {

	start, end, progress := 0.0, 10.0, 0.5
	half := Float64Lerp(start, end, progress)
	if !Float64Aproximately(5.0, half) {
		t.Errorf("5.0 != %f", half)
	}
}

func Test_Vec2Lerp(t *testing.T) {
	start, end, progress := NewVec2(0.0, 10.0), NewVec2(10.0, 10.0), 0.5
	half := Vec2Lerp(start, end, progress)

	if !Float64Aproximately(5.0, half.X) {
		t.Errorf("5.0 != %f", half)
	}

	if !Float64Aproximately(10.0, half.Y) {
		t.Errorf("5.0 != %f", half)
	}
}
