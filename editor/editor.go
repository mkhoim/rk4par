package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"proj3/gravitation"
	"proj3/rk4"
	"strconv"
	"strings"
	"time"
)

func readFile(file string) (*gravitation.Gravitation, *rk4.Numerical, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("failed to scan file: %v", err)
	}

	if len(lines) != 10 {
		return nil, nil, fmt.Errorf("invalid file format")
	}

	massStrings := strings.Split(lines[0], ",")
	masses := make([]float64, len(massStrings))
	for i, val := range massStrings {
		mass, err := strconv.ParseFloat(strings.TrimSpace(val), 64)
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing mass on line 1: %v", err)
		}
		masses[i] = mass
	}

	initialPos := make([][]float64, len(masses))
	initialVel := make([][]float64, len(masses))
	for i := 0; i < len(masses); i++ {
		initialPos[i] = make([]float64, 3)
		initialVel[i] = make([]float64, 3)
	}

	for coord, line := range lines[1:4] {
		coords := strings.Split(line, ",")
		for i, val := range coords {
			pos, err := strconv.ParseFloat(strings.TrimSpace(val), 64)
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing initial position on line %d: %v", coord+2, err)
			}
			initialPos[i][coord] = pos
		}
	}

	for coord, line := range lines[4:7] {
		coords := strings.Split(line, ",")
		for i, val := range coords {
			vel, err := strconv.ParseFloat(strings.TrimSpace(val), 64)
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing initial velocity on line %d: %v", coord+5, err)
			}
			initialVel[i][coord] = vel
		}
	}

	g_val, err := strconv.ParseFloat(strings.TrimSpace(lines[7]), 64)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing G value on line 8: %v", err)
	}

	interval_slice := make([]float64, 2)
	for i, val := range strings.Split(lines[8], ",") {
		interval_slice[i], _ = strconv.ParseFloat(strings.TrimSpace(val), 64)
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing interval on line 9: %v", err)
		}
	}

	stepSize, err := strconv.ParseFloat(strings.TrimSpace(lines[9]), 64)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing step size on line 10: %v", err)
	}

	g := &gravitation.Gravitation{
		Masses:             masses,
		Initial_positions:  initialPos,
		Initial_velocities: initialVel,
		G:                  g_val,
		Num_bodies:         len(masses),
	}

	n, err := rk4.New(stepSize, interval_slice[0], interval_slice[1])
	if err != nil {
		return nil, nil, fmt.Errorf("error creating numerical: %v", err)
	}

	return g, n, nil
}

func write2JSON(filename string, times []float64, positions [][][]float64, velocities [][][]float64) error {
	if len(times) != len(positions) || len(times) != len(velocities) {
		return fmt.Errorf("times, positions, and velocities must have the same length")
	}

	numBodies := len(positions[0])
	bodies := make([]map[string]interface{}, numBodies)

	for i := 0; i < numBodies; i++ {
		positionData := []map[string]interface{}{}
		velocityData := []map[string]interface{}{}

		for j, time := range times {
			positionData = append(positionData, map[string]interface{}{
				"time":     time,
				"position": [3]float64{positions[j][i][0], positions[j][i][1], positions[j][i][2]},
			})
			velocityData = append(velocityData, map[string]interface{}{
				"time":     time,
				"velocity": [3]float64{velocities[j][i][0], velocities[j][i][1], velocities[j][i][2]},
			})
		}

		bodies[i] = map[string]interface{}{
			"body":       i + 1,
			"positions":  positionData,
			"velocities": velocityData,
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print JSON
	if err := encoder.Encode(bodies); err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}

	return nil
}

func filler(times []float64, positions [][][]float64, velocities [][][]float64) {
	return
}

const usage = "Usage: go run editor.go <input_file> <seq/par/ws> <num_threads>\n" +
	"  <input_file> is the path to the input file\n" +
	"  <seq/par/ws> is the mode to run the simulation in\n" +
	"  <num_threads> is the number of threads to use in the parallel implementation, ignored in sequential mode"

func main() {
	args := os.Args[1:]
	if len(args) != 3 {
		fmt.Println(usage)
		return
	}

	g, n, err := readFile(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	if args[1] != "seq" && args[1] != "par" && args[1] != "ws" {
		fmt.Println(usage)
		return
	}

	numThreads, _ := strconv.Atoi(args[2])

	start := time.Now()
	times, positions, velocities := g.Simulate(n.Start, n.End, n.Step, args[1], numThreads)
	elapsed := time.Since(start)
	fmt.Printf("%f", elapsed.Seconds())
	// write2JSON("benchmark/in_out/output.json", times, positions, velocities)
	filler(times, positions, velocities)
}
