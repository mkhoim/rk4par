package test

import (
	"math"
	"proj3/rk4"
	"testing"
)

func quadratic_ode(t float64, state []float64) []float64 {
	return []float64{2 * t}
}

func quadratic1(t float64) float64 {
	return t * t
}

func quadratic2(t float64) float64 {
	return t*t + 2
}

func TestQuadratic1(t *testing.T) {
	y0 := []float64{0}
	n, _ := rk4.New(0.1, 0, 10)

	times, res := n.SeqRK4(quadratic_ode, y0)

	for i := 0; i < len(times); i++ {
		if math.Abs(quadratic1(times[i])-res[i][0]) > 1e-5 {
			t.Errorf("expected %f, got %f", quadratic1(times[i]), res[i][0])
		}
	}
}

func TestQuadratic2(t *testing.T) {
	y0 := []float64{3}
	n, _ := rk4.New(0.1, 1, 1000)

	times, res := n.SeqRK4(quadratic_ode, y0)

	for i := 0; i < len(times); i++ {
		if math.Abs(quadratic2(times[i])-res[i][0]) > 1e-5 {
			t.Errorf("expected %f, got %f", quadratic2(times[i]), res[i][0])
		}
	}
}

func TestQuadraticPar1(t *testing.T) {
	y0 := []float64{0}
	n, _ := rk4.New(0.1, 0, 10)

	times, res := n.PipelinedRK4(quadratic_ode, y0, 10)

	for i := 0; i < len(times); i++ {
		if math.Abs(quadratic1(times[i])-res[i][0]) > 10 {
			t.Errorf("expected %f, got %f", quadratic1(times[i]), res[i][0])
		}
	}
}

func TestQuadraticPar2(t *testing.T) {
	y0 := []float64{3}
	n, _ := rk4.New(0.1, 1, 1000)

	times, res := n.PipelinedRK4(quadratic_ode, y0, 10)

	for i := 0; i < len(times); i++ {
		if math.Abs(quadratic2(times[i])-res[i][0]) > 1000 {
			t.Errorf("expected %f, got %f", quadratic2(times[i]), res[i][0])
		}
	}
}

func TestQuadraticWS1(t *testing.T) {
	y0 := []float64{0}
	n, _ := rk4.New(0.1, 0, 10)

	times, res := n.WorkStealRK4(quadratic_ode, y0, 10)

	for i := 0; i < len(times); i++ {
		if math.Abs(quadratic1(times[i])-res[i][0]) > 10 {
			t.Errorf("expected %f, got %f", quadratic1(times[i]), res[i][0])
		}
	}
}

func TestQuadraticWS2(t *testing.T) {
	y0 := []float64{3}
	n, _ := rk4.New(0.1, 1, 1000)

	times, res := n.WorkStealRK4(quadratic_ode, y0, 10)

	for i := 0; i < len(times); i++ {
		if math.Abs(quadratic2(times[i])-res[i][0]) > 1000 {
			t.Errorf("expected %f, got %f", quadratic2(times[i]), res[i][0])
		}
	}
}
