package main

import (
	"math"
	"math/rand"
	"testing"
)

const _blen = 1000000

func Benchmark_float32(b *testing.B) {
	b.Skip()
	b.StopTimer()
	a := make([]float32, _blen)
	for i := range a {
		a[i] = rand.Float32()
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		math.Tanh(float64(a[i%_blen]))
		math.Sin(float64(a[i%_blen]))
		math.Cos(float64(a[i%_blen]))
		math.Sinh(float64(a[i%_blen]))
	}
}

func Benchmark_float64(b *testing.B) {
	b.Skip()
	b.StopTimer()
	a := make([]float64, _blen)
	for i := range a {
		a[i] = rand.Float64()
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		math.Tanh(a[i%_blen])
		math.Sin(a[i%_blen])
		math.Cos(a[i%_blen])
		math.Sinh(a[i%_blen])
	}
}

func Benchmark_Simulation(b *testing.B) {
	simulation := NewSimulation(128, 128, 5, 5,
		100, 30, 10, 1000, 0.05)
	for n := 0; n < b.N; n++ {
		simulation.simulationStep()
		if simulation.tick >= 300 {
			simulation.tick = 0
			simulation.applySelection()
			simulation.repopulate()
		}
	}
}
