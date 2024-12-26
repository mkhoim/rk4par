package gravitation

import (
	"math"
	"proj3/rk4"
)

type Gravitation struct {
	Masses             []float64
	Initial_positions  [][]float64
	Initial_velocities [][]float64
	G                  float64
	Num_bodies         int
}

func flatten(arr [][]float64) []float64 {
	flat := make([]float64, 0)
	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr[i]); j++ {
			flat = append(flat, arr[i][j])
		}
	}
	return flat
}

func concat(a []float64, b []float64) []float64 {
	concat := make([]float64, len(a)+len(b))
	for i := 0; i < len(a); i++ {
		concat[i] = a[i]
	}
	for i := 0; i < len(b); i++ {
		concat[len(a)+i] = b[i]
	}
	return concat
}

func distance(a []float64, b []float64) float64 {
	distance := 0.0
	for i := 0; i < len(a); i++ {
		distance += (a[i] - b[i]) * (a[i] - b[i])
	}
	return math.Sqrt(distance)
}

func (g *Gravitation) Acceleration(positions [][]float64) [][]float64 {
	acc := make([][]float64, len(positions))
	for i := 0; i < len(positions); i++ {
		acc[i] = make([]float64, len(positions[i]))
		for j := 0; j < len(positions); j++ {
			if i != j {
				r := distance(positions[i], positions[j])
				diff_x := positions[j][0] - positions[i][0]
				diff_y := positions[j][1] - positions[i][1]
				diff_z := positions[j][2] - positions[i][2]
				if r > 1e-5 {
					acc[i][0] += g.G * g.Masses[j] * diff_x / (r * r)
					acc[i][1] += g.G * g.Masses[j] * diff_y / (r * r)
					acc[i][2] += g.G * g.Masses[j] * diff_z / (r * r)
				}
			}
		}
	}
	return acc
}

func (g *Gravitation) Dynamics(t float64, state []float64) []float64 {
	positions := state[:g.Num_bodies*3]
	velocities := state[g.Num_bodies*3:]

	reshaped_positions := make([][]float64, g.Num_bodies)
	reshaped_velocities := make([][]float64, g.Num_bodies)

	for i := 0; i < g.Num_bodies; i++ {
		reshaped_positions[i] = positions[i*3 : i*3+3]
		reshaped_velocities[i] = velocities[i*3 : i*3+3]
	}

	acc := g.Acceleration(reshaped_positions)

	new_state := concat(flatten(reshaped_velocities), flatten(acc))

	return new_state
}

func (g *Gravitation) Simulate(t0 float64, tf float64, step_size float64, impl string, numThreads int) ([]float64, [][][]float64, [][][]float64) {
	initial_state := concat(flatten(g.Initial_positions), flatten(g.Initial_velocities))

	numerical, _ := rk4.New(step_size, t0, tf)

	times := make([]float64, 0)
	states := make([][]float64, 0)

	if impl == "seq" {
		times, states = numerical.SeqRK4(g.Dynamics, initial_state)
	} else if impl == "par" {
		times, states = numerical.PipelinedRK4(g.Dynamics, initial_state, numThreads)
	} else {
		times, states = numerical.WorkStealRK4(g.Dynamics, initial_state, numThreads)
	}

	len_states := int((tf-t0)/step_size) + 1
	positions := make([][][]float64, len_states)
	velocities := make([][][]float64, len_states)

	for i := 0; i < len_states; i++ {
		positions[i] = make([][]float64, g.Num_bodies)
		velocities[i] = make([][]float64, g.Num_bodies)
		for j := 0; j < g.Num_bodies; j++ {
			positions[i][j] = states[i][j*3 : j*3+3]
			velocities[i][j] = states[i][g.Num_bodies*3+j*3 : g.Num_bodies*3+j*3+3]
		}
	}

	return times, positions, velocities
}
