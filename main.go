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

func main() {
	simulation := NewSimulation(128, 128, 5, 5,
		100, 30, 10, 1000, 0.01)

	window, err := sdl.CreateWindow("wyrmas", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(simulation.sizeX*cellSize), int32(simulation.sizeY*cellSize), sdl.WINDOW_SHOWN)
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

func visualSimulation(simulation Simulation, ticksPerGen int, renderer *sdl.Renderer) {
	var (
		running      = true
		pause        = false
		generation   = 0
		rates        = []int{0, 1, 10, 50, 100}
		rateIdx      = 1
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
			survivors := simulation.applySelection()
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
			simulation.repopulate()
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
