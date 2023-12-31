package main

import (
	"math/rand"
)

// wyrm age
func sAge(s *Simulation, w *Wyrm, _ *Neuron) float64 {
	return float64(w.age) / float64(s.maxAge)
}

// random
func sRand(s *Simulation, w *Wyrm, _ *Neuron) float64 {
	return rand.Float64()
}

// population density nearby. 1 is max density
func sPop(s *Simulation, w *Wyrm, _ *Neuron) float64 {
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
func sDistN(s *Simulation, w *Wyrm, _ *Neuron) float64 {
	d, _ := findNearest(s, w)
	return 1.0 - float64(d)/float64(s.maxDist)
}

// direction of nearest wyrm.
func sDirN(s *Simulation, w *Wyrm, _ *Neuron) float64 {
	_, d := findNearest(s, w)
	return d.normalize()
}

// distance to nearest wyrm in forward direction. 0 - no wyrm, 1 - wyrm in next cell
func sDistF(s *Simulation, w *Wyrm, _ *Neuron) float64 {
	for t := Dist(1); t <= s.maxDist; t++ {
		x := w.x + t*w.direction[0]
		y := w.y + t*w.direction[1]
		if x >= s.sizeX || x < 0 || y >= s.sizeY || y < 0 {
			return 0
		}
		if s.world[x][y] != nil {
			return 1.0 - float64(t)/float64(s.maxDist)
		}
	}
	return 0
}

// oscillator
func sOsc(s *Simulation, w *Wyrm, _ *Neuron) float64 {
	return s.oscValue
}

// latitude (0 is north, 1 is south)
func sLat(s *Simulation, w *Wyrm, _ *Neuron) float64 {
	return float64(w.y) / float64(s.sizeY-1)
}

// longitude (0 is west, 1 is east)
func sLon(s *Simulation, w *Wyrm, _ *Neuron) float64 {
	return float64(w.x) / float64(s.sizeX-1)
}

// count of good places in forward direction
func sGoodCF(s *Simulation, w *Wyrm, _ *Neuron) float64 {
	count := 0
	for t := Dist(0); t <= s.maxDist; t++ {
		x := w.x + t*w.direction[0]
		y := w.y + t*w.direction[1]
		if x >= s.sizeX || x < 0 || y >= s.sizeY || y < 0 {
			return 0
		}
		if s.selectionArea[x][y] {
			count++
		}
	}
	return float64(count) / float64(s.maxDist)
}

// count of good places around
func sGoodCA(s *Simulation, w *Wyrm, _ *Neuron) float64 {
	count := 0
	dist := s.maxDist / 3 // see less far
	for t := Dist(0); t <= dist; t++ {
		x := w.x + t*w.direction[0]
		y := w.y + t*w.direction[1]
		if x >= s.sizeX || x < 0 || y >= s.sizeY || y < 0 {
			return 0
		}
		if s.selectionArea[x][y] {
			count++
		}
	}
	return float64(count) / float64(dist*dist)
	//return 0
}

// distance to good place in forward direction
func sGoodD(s *Simulation, w *Wyrm, _ *Neuron) float64 {
	if s.selectionArea[w.x][w.y] {
		return 1
	}
	for t := Dist(1); t <= s.maxDist; t++ {
		x := w.x + t*w.direction[0]
		y := w.y + t*w.direction[1]
		if x >= s.sizeX || x < 0 || y >= s.sizeY || y < 0 {
			return 0
		}
		if s.selectionArea[x][y] {
			return 1 - float64(t)/float64(s.maxDist)
		}
	}
	return 0
}

func findNearest(s *Simulation, w *Wyrm) (Dist, Direction) {
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
