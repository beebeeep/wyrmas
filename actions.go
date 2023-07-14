package main

import (
	"math"
	"math/rand"
)

var (
	_eps = 1e-3
)

func aKill(s *Simulation, w *Wyrm, n *Neuron) float64 {
	return 0 // thou shalt not kill
}

// set neuron responsiveness
func aResp(s *Simulation, w *Wyrm, n *Neuron) float64 {
	p := tanhActivation(s, w, n)
	if sgn, ok := activate(p); ok {
		n.responsiveness *= 0.1 * sgn
	}
	return p
}

// move forward/backward (in direction wyrm is facing)
func aMoveF(s *Simulation, w *Wyrm, n *Neuron) float64 {
	p := tanhActivation(s, w, n)
	if sgn, ok := activate(p); ok {
		ds := Dist(sgn)
		move(s, w, Direction{w.direction[0] * ds, w.direction[1] * ds})
	}
	return p
}

// move east/west
func aMoveEW(s *Simulation, w *Wyrm, n *Neuron) float64 {
	p := tanhActivation(s, w, n)
	if sgn, ok := activate(p); ok {
		w.direction[0] = Dist(sgn)
		move(s, w, w.direction)
	}
	return p
}

// move north/south
func aMoveNS(s *Simulation, w *Wyrm, n *Neuron) float64 {
	p := tanhActivation(s, w, n)
	if sgn, ok := activate(p); ok {
		w.direction[1] = Dist(sgn)
		move(s, w, w.direction)
	}
	return p
}

func activate(p float64) (float64, bool) {
	ok := rand.Float64() <= math.Abs(p)
	if p >= _eps {
		return 1, ok
	}
	return -1, ok

}

func move(s *Simulation, w *Wyrm, d Direction) {
	x1 := w.x + d[0]
	y1 := w.y + d[1]
	if x1 < 0 || y1 < 0 || x1 >= s.sizeX || y1 >= s.sizeY {
		return
	}
	if s.world[x1][y1] != nil {
		return
	}

	s.world[w.x][w.y] = nil
	w.x = x1
	w.y = y1
	s.world[w.x][w.y] = w
}
