package rk4

import (
	"sync"
)

func (n *Numerical) FirstStage(f Function, y0 []float64, largeStepSize float64, taskChan chan Task) {
	t, y := n.Start, y0
	tf := n.End
	index := 0
	smallStepSize := n.Step

	for t < tf {
		nextT := t + largeStepSize
		if nextT > tf {
			nextT = tf
		}

		y = Step(f, t, y, nextT-t)
		taskChan <- Task{
			StartIndex: index,
			EndIndex:   index + int((nextT-t)/smallStepSize),
			TStart:     t,
			TEnd:       nextT,
			YStart:     y,
			StepSize:   smallStepSize,
		}
		index += int((nextT - t) / smallStepSize)
		t = nextT
	}
	close(taskChan)
}

func Worker(f Function, ts []float64, results [][]float64, taskChan <-chan Task, wg *sync.WaitGroup, threadID int) {
	defer wg.Done()

	for task := range taskChan {
		t, y := task.TStart, task.YStart
		for i := task.StartIndex; i < task.EndIndex; i++ {
			y = Step(f, t, y, task.StepSize)
			t += task.StepSize
			ts[i+1] = t
			results[i+1] = y
		}
	}
}

func (n *Numerical) PipelinedRK4(f Function, y0 []float64, numThreads int) ([]float64, [][]float64) {
	t0 := n.Start
	tf := n.End
	nSteps := int((tf - t0) / n.Step)
	results := make([][]float64, nSteps+1)
	results[0] = y0
	ts := make([]float64, nSteps+1)
	ts[0] = t0

	h := n.Step
	largeStepSize := 5 * h

	taskChan := make(chan Task, numThreads)

	go n.FirstStage(f, y0, largeStepSize, taskChan)

	var wg sync.WaitGroup
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go Worker(f, ts, results, taskChan, &wg, i)
	}

	wg.Wait()
	return ts, results
}
