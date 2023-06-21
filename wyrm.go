package main

import "math"

var (
	actions = []activationFn{aKill, aResp, aMoveF, aMoveEW, aMoveNS}
	sensors = []activationFn{sAge, sRand, sPop, sDistN, sDirN, sOsc, sLat, sLon}
)

type Wyrm struct {
	x, y        Dist
	age         int
	direction   Direction
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

func tanhActivation(_ *State, _ *Wyrm, n *Neuron) float64 {
	var sum float64
	for _, l := range n.inputs {
		sum += l.source.potential * l.weight
	}
	return math.Tanh(n.responsiveness * sum)
}
