package main

import (
	"math"
	"math/rand"
)

var (
	_eps = 1e-3
)

func aKill(s *State, w *Wyrm, n *Neuron) float64 {
	return 0 // thou shalt not kill
}

// set neuron responsiveness
func aResp(s *State, w *Wyrm, n *Neuron) float64 {
	n.tanhActivate()
	if sgn, ok := activate(n); ok {
		n.responsiveness *= 0.1 * sgn
	}
	return 0
}

// move forward/backward (in direction wyrm is facing)
func aMoveF(s *State, w *Wyrm, n *Neuron) float64 {
	n.tanhActivate()
	if sgn, ok := activate(n); ok {
		ds := Dist(sgn)
		move(s, w, Direction{w.direction[0] * ds, w.direction[1] * ds})
	}
	return 0
}

// move east/west
func aMoveEW(s *State, w *Wyrm, n *Neuron) float64 {
	n.tanhActivate()
	if sgn, ok := activate(n); ok {
		w.direction[0] = Dist(sgn)
		move(s, w, w.direction)
	}
	return 0
}

// move north/south
func aMoveNS(s *State, w *Wyrm, n *Neuron) float64 {
	n.tanhActivate()
	if sgn, ok := activate(n); ok {
		w.direction[1] = Dist(sgn)
		move(s, w, w.direction)
	}
	return 0
}

func activate(n *Neuron) (float64, bool) {
	ok := rand.Float64() <= math.Abs(n.potential)
	if n.potential >= _eps {
		return 1, ok
	}
	return -1, ok

}

func move(s *State, w *Wyrm, d Direction) {
	w.x += d[0]
	w.y += d[1]
	if w.x < 0 {
		w.x = 0
	}
	if w.x >= s.sizeX {
		w.x = s.sizeX - 1
	}
	if w.y < 0 {
		w.y = 0
	}
	if w.y >= s.sizeY {
		w.y = s.sizeY - 1
	}
}
