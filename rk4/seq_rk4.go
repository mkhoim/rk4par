package rk4

func (n *Numerical) SeqRK4(f Function, y0 []float64) ([]float64, [][]float64) {
	t0 := n.Start
	tf := n.End
	stepSize := n.Step
	nSteps := int((tf - t0) / n.Step)
	numVar := len(y0)
	ts := make([]float64, nSteps+1)

	y := make([]float64, len(y0))
	t := t0
	ts[0] = t
	copy(y, y0)

	ys := make([][]float64, nSteps+1)
	ys[0] = make([]float64, numVar)
	copy(ys[0], y)

	for i := 0; i < nSteps; i++ {
		y = Step(f, t, y, stepSize)
		t += stepSize
		ys[i+1] = make([]float64, numVar)
		copy(ys[i+1], y)
		ts[i+1] = t
	}

	return ts, ys
}
