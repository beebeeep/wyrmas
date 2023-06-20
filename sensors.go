package main

import (
	"math"
	"math/rand"
)

// wyrm age
func sAge(s *State, w *Wyrm, n *Neuron) float32 {
	return float32(w.age) / _maxAge
}

// random
func sRand(s *State, w *Wyrm, n *Neuron) float32 {
	return rand.Float32()
}

// population density nearby. 1 is max density
func sPop(s *State, w *Wyrm, n *Neuron) float32 {
	c := 0
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			x0 := w.x + x
			y0 := w.y + y
			if x == 0 || y == 0 || x0 < 0 || x0 >= s.sizeX || y0 < 0 || y0 >= s.sizeX {
				continue
			}
			if s.world[x0][y0] != nil {
				c++
			}
		}
	}
	return float32(c) / 8.0
}

// distance to nearest wyrm
func sDistN(s *State, w *Wyrm, n *Neuron) float32 {
	// implement me
}

// direction of nearest wyrm. 0 is east, 0.25 is north, 0.5 is west, 0.75 is south
func sDirN(s *State, w *Wyrm, n *Neuron) float32 {
	// implement me
}

// distance to nearest wyrm in forward direction. 0 - no wyrm, 1 - wyrm in next cell
func sDistF(s *State, w *Wyrm, n *Neuron) float32 {
	// implement me
}

// oscillator
func sOsc(s *State, w *Wyrm, n *Neuron) float32 {
	return 0.5 + float32(math.Cos(
		2.0*math.Pi*float64(s.tick%_oscPeriod)/float64(_oscPeriod),
	))/2.0
}

// latitude (0 is north, 1 is south)
func sLat(s *State, w *Wyrm, n *Neuron) float32 {
	return float32(w.y) / float32(s.sizeY-1)
}

// longitude (0 is west, 1 is east)
func sLon(s *State, w *Wyrm, n *Neuron) float32 {
	return float32(w.x) / float32(s.sizeX-1)
}
