package test

import (
	"math"
	"proj3/rk4"
	"testing"
)

func linear_ode(t float64, state []float64) []float64 {
	return []float64{1}
}

func linear1(t float64) float64 {
	return t
}

func linear2(t float64) float64 {
	return t + 1
}

func TestLinear1(t *testing.T) {
	y0 := []float64{0}
	n, _ := rk4.New(0.1, 0, 10)

	times, res := n.SeqRK4(linear_ode, y0)

	for i := 0; i < len(times); i++ {
		if math.Abs(linear1(times[i])-res[i][0]) > 1e-5 {
			t.Errorf("expected %f, got %f", linear1(times[i]), res[i][0])
		}
	}
}

func TestLinear2(t *testing.T) {
	y0 := []float64{2}
	n, _ := rk4.New(0.1, 1, 1000)

	times, res := n.SeqRK4(linear_ode, y0)

	for i := 0; i < len(times); i++ {
		if math.Abs(linear2(times[i])-res[i][0]) > 1e-5 {
			t.Errorf("expected %f, got %f", linear2(times[i]), res[i][0])
		}
	}
}

func TestLinearPar1(t *testing.T) {
	y0 := []float64{0}
	n, _ := rk4.New(0.05, 0, 100)

	times, res := n.PipelinedRK4(linear_ode, y0, 10)

	for i := 0; i < len(times); i++ {
		if math.Abs(linear1(times[i])-res[i][0]) >= 0.5 {
			t.Errorf("expected %f, got %f", linear1(times[i]), res[i][0])
			break
		}
	}
}

func TestLinearPar2(t *testing.T) {
	y0 := []float64{2}
	n, _ := rk4.New(0.1, 1, 1000)

	times, res := n.PipelinedRK4(linear_ode, y0, 2)

	for i := 0; i < len(times); i++ {
		if math.Abs(linear2(times[i])-res[i][0]) > 1 {
			t.Errorf("expected %f, got %f", linear2(times[i]), res[i][0])
		}
	}
}

func TestLinearWS1(t *testing.T) {
	y0 := []float64{0}
	n, _ := rk4.New(0.1, 0, 10)

	times, res := n.WorkStealRK4(linear_ode, y0, 10)

	for i := 0; i < len(times); i++ {
		if math.Abs(linear1(times[i])-res[i][0]) >= 0.7 {
			t.Errorf("expected %f, got %f", linear1(times[i]), res[i][0])
		}
	}
}

func TestLinearWS2(t *testing.T) {
	y0 := []float64{2}
	n, _ := rk4.New(0.1, 1, 1000)

	times, res := n.WorkStealRK4(linear_ode, y0, 2)

	for i := 0; i < len(times); i++ {
		if math.Abs(linear2(times[i])-res[i][0]) > 2 {
			t.Errorf("expected %f, got %f", linear2(times[i]), res[i][0])
		}
	}
}
