package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"math/rand"
	"time"
)

const (
	cellSize = 7
)

func visualSimulation(simulation Simulation, ticksPerGen int, renderer *sdl.Renderer) {
	var (
		running      = true
		pause        = false
		generation   = 0
		rates        = []int{0, 1, 10, 50, 100}
		rateIdx      = 4
		genStart     = time.Now()
		dumpSurvivor = false
	)

	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch ev := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if ev.Type != sdl.KEYUP {
					break
				}
				if ev.Keysym.Sym == sdl.K_SPACE {
					rateIdx = (rateIdx + 1) % len(rates)
				}
				if ev.Keysym.Sym == sdl.K_d {
					dumpSurvivor = true
				}
			}
		}

		if pause {
			sdl.Delay(50)
			continue
		}

		if rate := rates[rateIdx]; rate != 0 && simulation.tick%rate == 0 {
			renderSimulation(simulation, renderer)
		}

		simulation.simulationStep()
		if simulation.tick >= ticksPerGen {
			genTime := time.Now().Sub(genStart)
			generation++
			targetPop := len(simulation.wyrmas)
			survivors := simulation.selectionZone()
			fmt.Printf("generation %d, survived %.2f%%\n", generation, 100.0*float32(survivors)/float32(targetPop))
			fmt.Printf("generation simulation took %.3f sec, rate %.1f ticks/sec\n", genTime.Seconds(), float64(ticksPerGen)/genTime.Seconds())
			if dumpSurvivor {
				dumpSurvivor = false
				for {
					if w := simulation.wyrmas[rand.Intn(len(simulation.wyrmas))]; !w.dead {
						fmt.Printf("dumped survivor %v", w.genome)
						w.DumpGenomeGraph(fmt.Sprintf("survivor-gen%d.png", generation))
						break

					}
				}
			}
			simulation.tick = 0
			simulation.repopulate(targetPop)
			genStart = time.Now()
		}
	}
}

func renderSimulation(simulation Simulation, renderer *sdl.Renderer) {
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
	renderer.SetDrawColor(100, 100, 100, 255)
	for x := int32(0); x < int32(simulation.sizeX*cellSize); x += cellSize {
		renderer.DrawLine(x, 0, x, int32(simulation.sizeY*cellSize))
	}
	for y := int32(0); y < int32(simulation.sizeX*cellSize); y += cellSize {
		renderer.DrawLine(0, y, int32(simulation.sizeX*cellSize), y)
	}
	for cx := range simulation.world {
		for cy := range simulation.world[cx] {
			x := int32(cx * cellSize)
			y := int32(cy * cellSize)
			if simulation.selectionArea[cx][cy] {
				renderer.SetDrawColor(0, 50, 0, 255)
				renderer.FillRect(&sdl.Rect{X: x + 1, Y: y + 1, W: cellSize - 1, H: cellSize - 1})
			}
			if w := simulation.world[cx][cy]; w != nil {
				renderer.SetDrawColor(0, 200, 0, 255)
				renderer.FillRect(&sdl.Rect{X: x + 1, Y: y + 1, W: cellSize - 1, H: cellSize - 1})
			}
		}
	}
	renderer.Present()
}

func createSelectionArea(s *Simulation) {
	sx := int(s.sizeX)
	sy := int(s.sizeY)
	s.selectionArea = make([][]bool, sx)
	for x := range s.selectionArea {
		s.selectionArea[x] = make([]bool, sy)
	}

	for x := range s.selectionArea {
		for y := range s.selectionArea[x] {

			//if x >= sx/8 && x <= sx*7/8 && y >= sy/8 && y <= sy*7/8 {
			//	continue
			//}

			//if x%20 < 15 && y%20 < 15 {
			//	continue
			//}

			if !(x >= sx*3/8 && x < sx*5/8 && y >= sy*3/8 && y < sy*5/8) {
				continue
			}
			s.selectionArea[x][y] = true
		}
	}

	/*
			for x := sx * 7 / 8; x < sx; x++ {
				for y := range s.selectionArea[x] {
					s.selectionArea[x][y] = true
				}
		}
	*/
}

func main() {
	sizeX := 128
	sizeY := 128
	targetPopulation := 1000
	genomeLen := 30
	numInnerNeurons := 5
	mutationProbability := 0.09
	maxAge := 100
	maxDist := 30

	simulation := Simulation{
		sizeX: Dist(sizeX), sizeY: Dist(sizeY), oscPeriod: 5,
		mutationProbability: mutationProbability, numInnerNeurons: numInnerNeurons,
		maxAge: maxAge, maxDist: Dist(maxDist),
	}
	simulation.world = make([][]*Wyrm, sizeX)
	createSelectionArea(&simulation)
	for x := range simulation.world {
		simulation.world[x] = make([]*Wyrm, sizeY)
	}

	simulation.randomizePopulation(targetPopulation, genomeLen)
	/*
		for i, w := range simulation.wyrmas[:10] {
			w.DumpGenomeGraph(fmt.Sprintf("genomes/wyrm%d.png", i))
		}
	*/
	window, err := sdl.CreateWindow("wyrmas", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(sizeX*cellSize), int32(sizeY*cellSize), sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("creating window: %s", err)
	}
	defer window.Destroy()
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("creating renderer: %s", err)
	}
	visualSimulation(simulation, 300, renderer)
}
