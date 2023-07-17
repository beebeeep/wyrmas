package main

import (
	"fmt"
	"math"
	"math/bits"
	"math/rand"
)

// Gene only encodes wyrm's nn connections
// Connection is encoded as follows:
// 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
// │ │  7 bits   │ │ │   7 bits  │ │          16 bits            │
// │ └──src ID───┘ │ └──sink ID──┘ └─────────weight──────────────┘
// └─> src type    └─> sink type
// src or sink type of 1 means inner layer neuron
// 16 bit weight is normalized as float in range  (-4, 4]
// note: endiannes does not matter here
type Gene uint32

var _neutralGenome []Gene

func (g Gene) getSrc() (isInner bool, id byte) {
	return g>>31 == 1, byte(g >> 24 & 127)
}
func (g Gene) getSink() (isInner bool, id byte) {
	return g>>23&1 == 1, byte(g >> 16 & 127)
}

func (g Gene) getWeight() float64 {
	return float64(int16(g&65535)-32767) / 8192.0
}

func (g Gene) String() string {
	return fmt.Sprintf("%x", uint32(g))
}

func mixGenomes(a, b []Gene) []Gene {
	v := [2][]Gene{a, b}
	r := make([]Gene, len(a))
	for i, idx := range rand.Perm(len(a)) {
		r[i] = v[i%2][idx]
	}
	return r
}

func (g *Gene) mutate() {
	// flip from 1 to 3 random bits
	for i := 0; i <= rand.Intn(3); i++ {
		*g = *g ^ (1 << rand.Intn(32))
	}
}

// genomeDiff returns difference between two genomes as number from (0, 1]
// calculated normalized sum of hamming distance between of each genes
func genomeDiff(a, b []Gene) float64 {
	diff := 0
	for i := range a {
		diff += bits.OnesCount(uint(a[i]) ^ uint(b[i]))
	}
	return float64(diff) / float64(32*len(a))
}

func genomeHash(g []Gene) float64 {
	// TODO choose random gene in population as a basis?
	if len(g) != len(_neutralGenome) {
		_neutralGenome = _neutralGenome[:0]
		for range g {
			_neutralGenome = append(_neutralGenome, 0xaaaaaaaa)
			//_neutralGenome = append(_neutralGenome, 0xaaaaaaaa)
		}
	}
	return genomeDiff(g, _neutralGenome)
}

func genomeColor(genome []Gene) (uint8, uint8, uint8) {
	// interpret genome hash as a hue [0, 1] and return rgb color assuming s=v=1
	h := genomeHash(genome)
	kr := math.Mod(5+h*6, 6)
	kg := math.Mod(3+h*6, 6)
	kb := math.Mod(1+h*6, 6)

	r := 1 - math.Max(min3(kr, 4-kr, 1), 0)
	g := 1 - math.Max(min3(kg, 4-kg, 1), 0)
	b := 1 - math.Max(min3(kb, 4-kb, 1), 0)

	return uint8(255 * r), uint8(255 * g), uint8(255 * b)
}

func min3(a, b, c float64) float64 {
	return math.Min(math.Min(a, b), c)
}
