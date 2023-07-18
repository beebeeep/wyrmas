package main

import (
	"math"
	"math/rand"
	"sort"
)

type Simulation struct {
	sizeX, sizeY    Dist
	maxDist         Dist // max distance wyrm can see
	oscPeriod       int
	oscValue        float64
	mutationRate    float64
	numInnerNeurons int
	maxAge          int
	tick            int // simulation tick
	wyrmas          []Wyrm
	world           [][]*Wyrm // x, y -> wyrm
	selectionArea   [][]bool
}

func NewSimulation(sizeX, sizeY, oscPeriod, numInnerNeurons, maxAge, maxDist, genomeLen, population int, mutationRate float64) Simulation {
	simulation := Simulation{
		sizeX: Dist(sizeX), sizeY: Dist(sizeY), oscPeriod: oscPeriod,
		mutationRate: mutationRate, numInnerNeurons: numInnerNeurons,
		maxAge: maxAge, maxDist: Dist(maxDist),
		world: make([][]*Wyrm, sizeX),
	}
	for x := range simulation.world {
		simulation.world[x] = make([]*Wyrm, sizeY)
	}
	simulation.selectionArea = make([][]bool, sizeX)
	for x := range simulation.selectionArea {
		simulation.selectionArea[x] = make([]bool, sizeY)
	}
	simulation.randomizePopulation(population, genomeLen)
	simulation.createSelectionArea()
	return simulation
}

func (s *Simulation) createSelectionArea() {

	/*	for x := range s.selectionArea {
			for y := range s.selectionArea[x] {
				s.selectionArea[x][y] = false
			}
		}
		for i := 0; i <= 80; i++ {
			x := rand.Intn(int(s.sizeX-4)) + 2
			y := rand.Intn(int(s.sizeY-4)) + 2
			for dx := -2; dx < 2; dx++ {
				for dy := -2; dy < 2; dy++ {
					s.selectionArea[x+dx][y+dy] = true
				}
			}
		}
	*/
	sx := int(s.sizeX)
	sy := int(s.sizeY)
	for x := range s.selectionArea {
		for y := range s.selectionArea[x] {
			// survive on border
			//if x >= sx/8 && x <= sx*7/8 && y >= sy/8 && y <= sy*7/8 {
			//	continue
			//}

			//	checker patter
			//if x%20 < 15 && y%20 < 15 {
			//	continue
			//}

			//survive in middle
			if !(x >= sx*3/8 && x < sx*5/8 && y >= sy*3/8 && y < sy*5/8) {
				continue
			}

			s.selectionArea[x][y] = true
		}
	}

}

func (s *Simulation) simulationStep() {
	s.tick++
	s.oscValue = 0.5 + 0.5*math.Cos(
		2.0*math.Pi*float64(s.tick%s.oscPeriod)/float64(s.oscPeriod),
	)
	for iw, w := range s.wyrmas {
		for _, n := range w.sensorLayer {
			if n.responsiveness == 0 {
				continue
			}
			n.potential = n.responsiveness * n.activate(s, &s.wyrmas[iw], nil)
		}
		for _, n := range w.innerLayer {
			if n.responsiveness == 0 {
				continue
			}
			n.potential = n.responsiveness * n.activate(nil, nil, n)
		}
		for _, n := range w.actionLayer {
			if n.responsiveness == 0 {
				continue
			}
			n.potential = n.responsiveness * n.activate(s, &s.wyrmas[iw], n)
		}
	}
}

func (s *Simulation) applySelection() int {
	// leave only those who ended up inside selection area
	survivors := len(s.wyrmas)
	for i, w := range s.wyrmas {
		if !s.selectionArea[w.x][w.y] {
			s.wyrmas[i].dead = true
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

func (s *Simulation) repopulate() {
	survivors := make([]*Wyrm, 0, len(s.wyrmas))
	for i := range s.wyrmas {
		if !s.wyrmas[i].dead {
			survivors = append(survivors, &s.wyrmas[i])
		}
	}
	genomes := s.breed(survivors, len(s.wyrmas))

	// clean the world
	for x := range s.world {
		for y := range s.world[x] {
			s.world[x][y] = nil
		}
	}

	// reuse previous generation by re-placing them randomly
	// and rewiring neurons using new genome
	for i, genome := range genomes {
		var x, y Dist
		for {
			x = Dist(rand.Intn(int(s.sizeX)))
			y = Dist(rand.Intn(int(s.sizeY)))
			if s.world[x][y] == nil {
				break
			}
		}
		sort.Slice(genome, func(i, j int) bool {
			return genome[i] < genome[j]
		})
		s.wyrmas[i].x = x
		s.wyrmas[i].y = y
		s.wyrmas[i].direction[0] = Dist(rand.Intn(3) - 1)
		s.wyrmas[i].direction[1] = Dist(rand.Intn(3) - 1)
		s.wyrmas[i].dead = false
		s.wyrmas[i].genome = genome
		s.wyrmas[i].wireNeurons()

		s.world[x][y] = &s.wyrmas[i]
	}
	s.createSelectionArea()
	s.tick = 0
}

func (s *Simulation) breed(wyrmas []*Wyrm, targetPopulation int) [][]Gene {
	genomes := make([][]Gene, 0, targetPopulation)
	pop := len(wyrmas)
	numChildren := targetPopulation / pop

	f := func(idx int) []Gene {
		genome := mixGenomes(wyrmas[idx].genome, wyrmas[(idx+1)%pop].genome)
		for i := range genome {
			if rand.Float64() <= s.mutationRate {
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
