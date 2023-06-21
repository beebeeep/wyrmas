package main

import "math"

var (
	_dirE  = Direction{1, 0}
	_dirNE = Direction{1, -1}
	_dirN  = Direction{0, -1}
	_dirNW = Direction{-1, -1}
	_dirW  = Direction{-1, 0}
	_dirSW = Direction{-1, 1}
	_dirS  = Direction{0, 1}
	_dirSE = Direction{1, 1}
)

type Dist int8
type Direction [2]Dist

type activationFn func(s *State, w *Wyrm, n *Neuron) float64

type Wyrm struct {
	x, y       Dist
	age        int
	direction  Direction
	sensors    []Neuron
	innerLayer []Neuron
	actions    []Neuron
}

type Link struct {
	weight float64
	source *Neuron
}

type Neuron struct {
	potential      float64
	responsiveness float64
	activate       activationFn
	inputs         []Link
}

func (d Direction) normalize() float64 {
	var s float64 = 1.0 / 7.0
	switch d {
	case _dirE:
		return 0
	case _dirNE:
		return s
	case _dirN:
		return s * 2
	case _dirNW:
		return s * 3
	case _dirW:
		return s * 4
	case _dirSW:
		return s * 5
	case _dirS:
		return s * 6
	case _dirSE:
		return s * 7
	default:
		return 0
	}
}

func (n *Neuron) tanhActivate() {
	var sum float64
	for _, l := range n.inputs {
		sum += l.source.potential * l.weight
	}
	n.potential = math.Tanh(n.responsiveness * sum)

}
