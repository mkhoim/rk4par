package rk4

import (
	"errors"
)

type Function func(t float64, y []float64) []float64

type Numerical struct {
	Step  float64
	Start float64
	End   float64
}

type Task struct {
	StartIndex, EndIndex   int
	TStart, TEnd, StepSize float64
	YStart                 []float64
}

func New(Step float64, Start float64, End float64) (*Numerical, error) {
	if Step <= 0 {
		return nil, errors.New("step must be greater than 0")
	}
	if Start >= End {
		return nil, errors.New("start must be less than end")
	}
	return &Numerical{Step, Start, End}, nil
}

func Step(f Function, t float64, y []float64, h float64) []float64 {
	numVar := len(y)
	yNext := make([]float64, numVar)
	y_copy := make([]float64, numVar)
	copy(y_copy, y)
	fs := make([]float64, numVar*4)
	k1 := fs[0*numVar : 1*numVar]
	k2 := fs[1*numVar : 2*numVar]
	k3 := fs[2*numVar : 3*numVar]
	k4 := fs[3*numVar : 4*numVar]

	k1 = f(t, y_copy)

	for i := 0; i < numVar; i++ {
		yNext[i] = y_copy[i] + 0.5*h*k1[i]
	}
	k2 = f(t+0.5*h, yNext)

	for i := 0; i < numVar; i++ {
		yNext[i] = y_copy[i] + 0.5*h*k2[i]
	}
	k3 = f(t+0.5*h, yNext)

	for i := 0; i < numVar; i++ {
		yNext[i] = y_copy[i] + h*k3[i]
	}
	k4 = f(t+h, yNext)

	for i := 0; i < numVar; i++ {
		y_copy[i] = y_copy[i] + (h/6)*(k1[i]+2*k2[i]+2*k3[i]+k4[i])
	}

	return y_copy
}
