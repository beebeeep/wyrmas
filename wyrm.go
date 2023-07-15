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
	sensors     = []activationFn{sAge, sRand, sPop, sDistN, sDirN, sOsc, sLat, sLon, sGoodD, sGoodCA, sGoodCF}
	sensorNames = []string{"sAge", "sRand", "sPop", "sDistN", "sDirN", "sOsc", "sLat", "sLon", "sGoodCD", "sGoodCA", "sGoodCF"}
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
	w.wireNeurons()

	return w
}

func (w *Wyrm) wireNeurons() {
	var src, sink *Neuron
	// reset neuron links
	for _, n := range w.innerLayer {
		n.inputs = n.inputs[:0]
	}
	for _, n := range w.actionLayer {
		n.inputs = n.inputs[:0]
	}

	// wire links according genome
	for _, gene := range w.genome {
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
}

func (w Wyrm) DumpGenomeGraph(filename string) {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		log.Fatal(err)
	}
	nodes := make(map[*Neuron]*cgraph.Node)

	for i, n := range w.sensorLayer {
		if w.countSinks(n) == 0 {
			continue
		}
		node, err := graph.CreateNode(sensorNames[i])
		if err != nil {
			log.Fatal(err)
		}
		nodes[n] = node
	}
	for i, n := range w.innerLayer {
		if w.countSinks(n) == 0 {
			continue
		}
		node, err := graph.CreateNode(fmt.Sprintf("inner-%d", i))
		if err != nil {
			log.Fatal(err)
		}
		nodes[n] = node
	}
	for i, n := range w.actionLayer {
		if len(n.inputs) == 0 {
			continue
		}
		node, err := graph.CreateNode(actionNames[i])
		if err != nil {
			log.Fatal(err)
		}
		nodes[n] = node
	}
	for neuron := range nodes {
		for _, link := range neuron.inputs {
			name := fmt.Sprintf("%p-%p", nodes[link.source], nodes[neuron])
			e, _ := graph.CreateEdge(name, nodes[link.source], nodes[neuron])
			e.SetLabel(fmt.Sprintf("%.2f", link.weight))
			e.SetLabelFontSize(4)
		}
	}
	if err := g.RenderFilename(graph, graphviz.PNG, filename); err != nil {
		log.Fatal(err)
	}
}

func (w Wyrm) countSinks(t *Neuron) int {
	count := 0
	for _, n := range w.innerLayer {
		for _, in := range n.inputs {
			if in.source == t {
				count++
			}
		}
	}
	for _, n := range w.actionLayer {
		for _, in := range n.inputs {
			if in.source == t {
				count++
			}
		}
	}
	return count
}

func tanhActivation(_ *Simulation, _ *Wyrm, n *Neuron) float64 {
	var sum float64
	for _, l := range n.inputs {
		sum += l.source.potential * l.weight
	}
	return math.Tanh(n.responsiveness * sum)
}
