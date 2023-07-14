package main

import (
	"fmt"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"log"
	"math"
)

var (
	actions     = []activationFn{aKill, aResp, aMoveF, aMoveEW, aMoveNS}
	sensors     = []activationFn{sAge, sRand, sPop, sDistN, sDirN, sOsc, sLat, sLon}
	sensorNames = []string{"sAge", "sRand", "sPop", "sDistN", "sDirN", "sOsc", "sLat", "sLon"}
	actionNames = []string{"aKill", "aResp", "aMoveF", "aMoveEW", "aMoveNS"}
)

type Wyrm struct {
	x, y        Dist
	direction   Direction
	age         int
	dead        bool
	genome      []Gene
	sensorLayer []*Neuron
	innerLayer  []*Neuron
	actionLayer []*Neuron
}

type Link struct {
	weight float64
	source *Neuron
}

type Neuron struct {
	potential      float64
	responsiveness float64
	activate       activationFn
	inputs         []Link
}

func NewWyrm(x, y Dist, numInner int, genome []Gene) Wyrm {
	w := Wyrm{
		x: x, y: y,
		genome:      genome,
		sensorLayer: make([]*Neuron, len(sensors)),
		innerLayer:  make([]*Neuron, numInner),
		actionLayer: make([]*Neuron, len(actions)),
	}
	for i := range sensors {
		w.sensorLayer[i] = &Neuron{responsiveness: 1, activate: sensors[i]}
	}
	for i := range actions {
		w.actionLayer[i] = &Neuron{responsiveness: 1, activate: actions[i], inputs: make([]Link, 0, 1)}

	}
	for i := range w.innerLayer {
		w.innerLayer[i] = &Neuron{responsiveness: 1, activate: tanhActivation, inputs: make([]Link, 0, 1)}
	}

	var src, sink *Neuron
	for _, gene := range genome {
		if srcInner, id := gene.getSrc(); srcInner {
			src = w.innerLayer[id%byte(len(w.innerLayer))]
		} else {
			src = w.sensorLayer[id%byte(len(w.sensorLayer))]
		}
		if sinkInner, id := gene.getSink(); sinkInner {
			sink = w.innerLayer[id%byte(len(w.innerLayer))]
		} else {
			sink = w.actionLayer[id%byte(len(w.actionLayer))]
		}
		sink.inputs = append(sink.inputs, Link{weight: gene.getWeight(), source: src})
	}

	return w
}

func (w Wyrm) DumpGenome(filename string) {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		log.Fatal(err)
	}
	nodes := make(map[*Neuron]*cgraph.Node)

	for i, n := range w.sensorLayer {
		node, err := graph.CreateNode(sensorNames[i])
		if err != nil {
			log.Fatal(err)
		}
		nodes[n] = node
	}
	for i, n := range w.innerLayer {
		node, err := graph.CreateNode(fmt.Sprintf("inner-%d", i))
		if err != nil {
			log.Fatal(err)
		}
		nodes[n] = node
	}
	for i, n := range w.actionLayer {
		node, err := graph.CreateNode(actionNames[i])
		if err != nil {
			log.Fatal(err)
		}
		nodes[n] = node
	}
	for neuron := range nodes {
		for _, link := range neuron.inputs {
			name := fmt.Sprintf("%.2f", link.weight)
			e, _ := graph.CreateEdge(name, nodes[link.source], nodes[neuron])
			e.SetLabel(name)
		}
	}
	if err := g.RenderFilename(graph, graphviz.PNG, filename); err != nil {
		log.Fatal(err)
	}
}

func tanhActivation(_ *Simulation, _ *Wyrm, n *Neuron) float64 {
	var sum float64
	for _, l := range n.inputs {
		sum += l.source.potential * l.weight
	}
	return math.Tanh(n.responsiveness * sum)
}
