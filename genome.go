package main

import "fmt"

// Gene only encodes wyrm's nn connections
// Connection is encoded as follows:
// 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
// │ │  7 bits   │ │ │   7 bits  │ │          16 bits            │
// │ └──src ID───┘ │ └──sink ID──┘ └─────────weight──────────────┘
// └─> src type    └─> sink type
// src or sink type of 1 means inner layer neuron
// 16 bit weight is normalized as float in range  (-4, 4]
type Gene uint32

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
