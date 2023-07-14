package main

import (
	"math"
	"math/rand"
)

type Simulation struct {
	sizeX, sizeY        Dist
	maxDist             Dist // max distance wyrm can see
	oscPeriod           int
	oscValue            float64
	mutationProbability float64
	numInnerNeurons     int
	maxAge              int
	tick                int // simulation tick
	wyrmas              []Wyrm
	world               [][]*Wyrm // x, y -> wyrm
}

func (s *Simulation) simulationStep() {
	s.tick++
	s.oscValue = 0.5 + 0.5*math.Cos(
		2.0*math.Pi*float64(s.tick%s.oscPeriod)/float64(s.oscPeriod),
	)
	for iw, w := range s.wyrmas {
		for _, n := range w.sensorLayer {
			n.potential = n.activate(s, &s.wyrmas[iw], nil)
		}
		for _, n := range w.innerLayer {
			n.potential = n.activate(nil, nil, n)
		}
		for _, n := range w.actionLayer {
			n.potential = n.activate(s, &s.wyrmas[iw], n)
		}
	}
}

func (s *Simulation) selectionEastSide() int {
	// leave only those who ended up near the eastern border
	survivors := len(s.wyrmas)
	for x := Dist(0); x <= s.sizeX*3/4; x++ {
		for y := Dist(0); y < s.sizeY; y++ {
			if w := s.world[x][y]; w != nil {
				w.dead = true
			}
			survivors--
		}
	}
	return survivors
}

func (s *Simulation) randomizePopulation(targetPopulation, genomeLen int) {
	for x := range s.world {
		for y := range s.world[x] {
			s.world[x][y] = nil
		}
	}
	s.wyrmas = make([]Wyrm, targetPopulation)
	for i := range s.wyrmas {
		var x, y Dist
		for {
			x = Dist(rand.Intn(len(s.world)))
			y = Dist(rand.Intn(len(s.world[0])))
			if s.world[x][y] == nil {
				break
			}
		}
		genome := make([]Gene, genomeLen)
		for j := range genome {
			genome[j] = Gene(rand.Uint32())
		}
		s.wyrmas[i] = NewWyrm(x, y, s.numInnerNeurons, genome)
		s.world[x][y] = &s.wyrmas[i]
	}

}

func (s *Simulation) repopulate(targetPopulation int) {
	survivors := make([]*Wyrm, 0, len(s.wyrmas))
	for i := range s.wyrmas {
		if !s.wyrmas[i].dead {
			survivors = append(survivors, &s.wyrmas[i])
		}
	}
	genomes := s.breed(survivors, targetPopulation)
	// wipe parents
	for x := range s.world {
		for y := range s.world[x] {
			s.world[x][y] = nil
		}
	}
	s.wyrmas = s.wyrmas[:0]
	for _, genome := range genomes {
		var x, y Dist
		for {
			x = Dist(rand.Intn(int(s.sizeX)))
			y = Dist(rand.Intn(int(s.sizeY)))
			if s.world[x][y] == nil {
				break
			}
		}
		s.wyrmas = append(s.wyrmas, NewWyrm(x, y, s.numInnerNeurons, genome))
		s.world[x][y] = &s.wyrmas[len(s.wyrmas)-1]
	}
	s.tick = 0
}

func (s *Simulation) breed(wyrmas []*Wyrm, targetPopulation int) [][]Gene {
	genomes := make([][]Gene, 0, targetPopulation)
	pop := len(wyrmas)
	numChildren := targetPopulation / pop

	f := func(idx int) []Gene {
		genome := mixGenomes(wyrmas[idx].genome, wyrmas[(idx+1)%pop].genome)
		for i := range genome {
			if rand.Float64() <= s.mutationProbability {
				genome[i].mutate()
			}
		}
		return genome
	}

	// generate random pairs from whole population
	// each pair will have at least numChildren children
	for _, idx := range rand.Perm(len(wyrmas)) {
		for i := 0; i < numChildren; i++ {
			genomes = append(genomes, f(idx))
		}
	}
	// generate some more random pairs to top up to the target population
	for _, idx := range rand.Perm(pop)[:(targetPopulation % pop)] {
		genomes = append(genomes, f(idx))
	}
	return genomes
}
