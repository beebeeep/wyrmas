package main

import (
	"math/rand"
)

// wyrm age
func sAge(s *State, w *Wyrm, _ *Neuron) float64 {
	return float64(w.age) / float64(s.maxAge)
}

// random
func sRand(s *State, w *Wyrm, _ *Neuron) float64 {
	return rand.Float64()
}

// population density nearby. 1 is max density
func sPop(s *State, w *Wyrm, _ *Neuron) float64 {
	c := 0
	for x := Dist(-1); x <= 1; x++ {
		x0 := w.x + x
		if x0 < 0 || x0 >= s.sizeX {
			continue
		}
		for y := Dist(-1); y <= 1; y++ {
			y0 := w.y + y
			if (x == 0 && y == 0) || y0 < 0 || y0 >= s.sizeX {
				continue
			}
			if s.world[x0][y0] != nil {
				c++
			}
		}
	}
	return float64(c) / 8.0
}

// distance to nearest wyrm
func sDistN(s *State, w *Wyrm, _ *Neuron) float64 {
	d, _ := findNearest(s, w)
	return float64(d) / float64(s.maxDist)
}

// direction of nearest wyrm.
func sDirN(s *State, w *Wyrm, _ *Neuron) float64 {
	_, d := findNearest(s, w)
	return d.normalize()
}

// distance to nearest wyrm in forward direction. 0 - no wyrm, 1 - wyrm in next cell
func sDistF(s *State, w *Wyrm, _ *Neuron) float64 {
	for t := Dist(1); t <= s.maxDist; t++ {
		x := w.x + t*w.direction[0]
		y := w.y + t*w.direction[1]
		if x >= s.sizeX || x < 0 || y >= s.sizeY || y < 0 {
			return 0
		}
		if s.world[x][y] != nil {
			return float64(t) / float64(s.maxDist)
		}
	}
	return 0
}

// oscillator
func sOsc(s *State, w *Wyrm, _ *Neuron) float64 {
	return s.oscValue
}

// latitude (0 is north, 1 is south)
func sLat(s *State, w *Wyrm, _ *Neuron) float64 {
	return float64(w.y) / float64(s.sizeY-1)
}

// longitude (0 is west, 1 is east)
func sLon(s *State, w *Wyrm, _ *Neuron) float64 {
	return float64(w.x) / float64(s.sizeX-1)
}

func findNearest(s *State, w *Wyrm) (Dist, Direction) {
	for t := Dist(1); t <= s.maxDist; t++ {
		for dx := Dist(-1); dx <= 1; dx++ {
			x := w.x + t*dx
			if x >= s.sizeX || x < 0 {
				continue
			}
			for dy := Dist(-1); dx <= 1; dx++ {
				y := w.y + t*dy
				if (dx == 0 && dy == 0) || y >= s.sizeY || y < 0 {
					continue
				}
				if s.world[x][y] != nil {
					return t, Direction{dx, dy}
				}
			}
		}
	}
	return 0, _dirE
}
