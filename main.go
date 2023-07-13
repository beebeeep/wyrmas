package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"log"
)

const (
	cellSize = 7
)

func visualSimulation(simulation Simulation, renderer *sdl.Renderer) {
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}
		simulation.simulationStep()
		fmt.Printf("tick %d\n", simulation.tick)
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()
		for cx := range simulation.world {
			for cy := range simulation.world[cx] {
				x := int32(cx * cellSize)
				y := int32(cy * cellSize)
				renderer.SetDrawColor(100, 100, 100, 255)
				renderer.DrawRect(&sdl.Rect{X: x, Y: y, W: cellSize, H: cellSize})
				if simulation.world[cx][cy] != nil {
					renderer.SetDrawColor(0, 200, 0, 255)
					renderer.FillRect(&sdl.Rect{X: x + 1, Y: y + 1, W: cellSize - 1, H: cellSize - 1})
				}
			}
		}
		renderer.Present()

	}
}

func main() {
	sizeX := 128
	sizeY := 128
	simulation := Simulation{
		sizeX: Dist(sizeX), sizeY: Dist(sizeY), oscPeriod: 5,
		mutationProbability: 0.05, numInnerNeurons: 1,
		maxAge: 100,
	}
	simulation.world = make([][]*Wyrm, sizeX)
	for x := range simulation.world {
		simulation.world[x] = make([]*Wyrm, sizeY)
	}
	simulation.randomizePopulation(1000, 4)
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
	visualSimulation(simulation, renderer)
}
