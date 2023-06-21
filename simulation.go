package main

import "math"

type State struct {
	sizeX, sizeY Dist
	maxDist      Dist // max distance wyrm can see
	oscPeriod    int
	oscValue     float64
	maxAge       int
	tick         int // simulation tick
	wyrmas       []*Wyrm
	world        [][]*Wyrm // x, y -> wyrm
}

func (s *State) simulationStep() {
	s.tick++
	s.oscValue = 0.5 + math.Cos(
		2.0*math.Pi*float64(s.tick%s.oscPeriod)/float64(s.oscPeriod),
	)/2.0
	for _, w := range s.wyrmas {
		for _, n := range w.sensorLayer {
			n.potential = n.activate(s, w, nil)
		}
		for _, n := range w.innerLayer {
			n.potential = n.activate(nil, nil, n)
		}
		for _, n := range w.actionLayer {
			n.potential = n.activate(s, w, n)
		}
	}
}
