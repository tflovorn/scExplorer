package seriesaccel

import (
	"math"
	"testing"
)

func TestLevinZeta2(t *testing.T) {
	fn := func(i int) float64 {
		return 1.0 / (float64(i) * float64(i))
	}
	val, absErr := Levin_u(fn, 1, 20)
	expected := math.Pi * math.Pi / 6.0
	if math.Abs(val-expected) > absErr*10 {
		t.Errorf("val - expected = %e is much greater than error estimate %e", val-expected, absErr)
	}
}
