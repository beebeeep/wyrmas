package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"math/rand"
	"os/exec"
	"time"
)

const (
	cellSize = 7
)

func main() {
	simulation := NewSimulation(128, 128, 5, 10,
		100, 30, 30, 1000, 0.05)

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
		running    = true
		generation = 0
		rates      = []int{0, 1, 10, 50, 100}
		rateIdx    = 1
		genStart   = time.Now()
		dump       = false
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
					dump = true
				}
			}
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
			if dump {
				dump = false
				dumpSurvivor(simulation, generation)
			}
			simulation.tick = 0
			simulation.repopulate()
			genStart = time.Now()
		}
	}
}

func dumpSurvivor(simulation Simulation, generation int) {
	for {
		if w := simulation.wyrmas[rand.Intn(len(simulation.wyrmas))]; !w.dead {
			fmt.Printf("dumped survivor %v", w.genome)
			fname := fmt.Sprintf("survivor-gen%d.png", generation)
			w.DumpGenomeGraph(fname)
			cmd := exec.Command("open", fname)
			go cmd.Run()
			return
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
