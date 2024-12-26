package rk4

import (
	"math/rand"
	"sync"
)

type Node struct {
	Task Task
	Next *Node
}

type Queue struct {
	Head, Tail *Node
	Lock       sync.Mutex
}

func (d *Queue) Push(task Task) {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	node := &Node{Task: task}
	if d.Head == nil {
		d.Head, d.Tail = node, node
	} else {
		node.Next = d.Head
		d.Head = node
	}
}

func (d *Queue) Pop() (Task, bool) {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	if d.Head == nil {
		return Task{}, false
	}
	task := d.Head.Task
	d.Head = d.Head.Next
	if d.Head == nil {
		d.Tail = nil
	}
	return task, true
}

func (d *Queue) Steal() (Task, bool) {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	if d.Tail == nil {
		return Task{}, false
	}
	task := d.Tail.Task
	if d.Head == d.Tail {
		d.Head, d.Tail = nil, nil
	} else {
		current := d.Head
		for current.Next != d.Tail {
			current = current.Next
		}
		current.Next = nil
		d.Tail = current
	}
	return task, true
}

func (n *Numerical) FirstStage_WS(f Function, y0 []float64, largeStepSize float64, queues []*Queue) {
	t, y := n.Start, y0
	tf := n.End
	index := 0
	smallStepSize := n.Step
	count := 0
	probabilities := make([]float64, len(queues))

	// Initialize probabilities non-uniformly
	total := 0.0
	for i := range probabilities {
		probabilities[i] = float64(i + 1)
		total += probabilities[i]
	}
	for i := range probabilities {
		probabilities[i] /= total
	}

	for t < tf {
		nextT := t + largeStepSize
		if nextT > tf {
			nextT = tf
		}

		y = Step(f, t, y, nextT-t)
		task := Task{
			StartIndex: index,
			EndIndex:   index + int((nextT-t)/smallStepSize),
			TStart:     t,
			TEnd:       nextT,
			YStart:     y,
			StepSize:   smallStepSize,
		}
		index += int((nextT - t) / smallStepSize)
		t = nextT

		idx := 0
		r := rand.Float64()
		for i, p := range probabilities {
			if r < p {
				idx = i
				break
			}
			r -= p
		}

		queues[idx].Push(task)

		count += 1
	}
}

func Worker_WS(f Function, ts []float64, results [][]float64, queue *Queue, queues []*Queue, wg *sync.WaitGroup, threadID int) {
	defer wg.Done()

	for {
		task, ok := queue.Pop()

		if !ok {
			stolen := false
			for i, otherQueue := range queues {
				if i != threadID {
					if task, ok = otherQueue.Steal(); ok {
						stolen = true
						break
					}
				}
			}
			if !stolen {
				return
			}
		}

		t, y := task.TStart, task.YStart
		for i := task.StartIndex; i < task.EndIndex; i++ {
			y = Step(f, t, y, task.StepSize)
			t += task.StepSize
			ts[i+1] = t
			results[i+1] = y
		}
	}
}

func (n *Numerical) WorkStealRK4(f Function, y0 []float64, numThreads int) ([]float64, [][]float64) {
	t0 := n.Start
	tf := n.End
	nSteps := int((tf - t0) / n.Step)
	results := make([][]float64, nSteps+1)
	results[0] = y0
	ts := make([]float64, nSteps+1)
	ts[0] = t0

	h := n.Step
	largeStepSize := 5 * h

	queues := make([]*Queue, numThreads)
	for i := range queues {
		queues[i] = &Queue{}
	}

	n.FirstStage_WS(f, y0, largeStepSize, queues)

	var wg sync.WaitGroup
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go Worker_WS(f, ts, results, queues[i], queues, &wg, i)
	}

	wg.Wait()
	return ts, results
}
