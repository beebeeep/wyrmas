package main

import "fmt"

const (
	_maxAge    float32 = 300
	_oscPeriod         = 10
)

const (
	_dirE = iota
	_dirNE
	_dirN
	_dirNW
	_dirW
	_dirSW
	_dirS
	_dirSE
)

type actFn func(s *State, w *Wyrm, n *Neuron) float32

type Wyrm struct {
	x, y       int
	age        int
	direction  byte
	sensors    []Neuron
	innerLayer []Neuron
	actions    []Neuron
}

type Link struct {
	weight float32
	Dest   *Neuron
}

type Neuron struct {
	val    float32
	fun    actFn
	inputs []Link
}

type State struct {
	tick         int
	sizeX, sizeY int
	wyrmas       []Wyrm
	world        [][]*Wyrm // x, y -> wyrm
}

func main() {
	fmt.Println("wyrmaas")

}
