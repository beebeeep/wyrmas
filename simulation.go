package main

import "math"

type State struct {
	sizeX, sizeY Dist
	maxDist      Dist // max distance wyrm can see
	oscPeriod    int
	oscValue     float64
	maxAge       int
	tick         int // simulation tick
	wyrmas       []Wyrm
	world        [][]*Wyrm // x, y -> wyrm
}

func (s *State) simulate() {
	s.tick++
	s.oscValue = 0.5 + math.Cos(
		2.0*math.Pi*float64(s.tick%s.oscPeriod)/float64(s.oscPeriod),
	)/2.0
}
