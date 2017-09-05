package gogobar

type Graph struct {
	min    float64
	max    float64
	values []float64
}

func NewGraph(width uint64) *Graph {
	graph := new(Graph)
	graph.min = 0.0
	graph.max = 1.0
	graph.values = make([]float64, width)
	return graph
}

func (graph *Graph) AddValue(val float64) {
	graph.values = append(graph.values[1:], []float64{val}...)
}

func (graph *Graph) Call() {
	diff := graph.max - graph.min
	for _, element := range graph.values {
		ratio := (element - graph.min) / diff
		if ratio <= 0.125 {
			buffer.WriteRune('▁')
			continue
		}
		if ratio <= 0.25 {
			buffer.WriteRune('▂')
			continue
		}
		if ratio <= 0.375 {
			buffer.WriteRune('▃')
			continue
		}
		if ratio <= 0.5 {
			buffer.WriteRune('▄')
			continue
		}
		if ratio <= 0.625 {
			buffer.WriteRune('▅')
			continue
		}
		if ratio <= 0.75 {
			buffer.WriteRune('▆')
			continue
		}
		if ratio <= 0.875 {
			buffer.WriteRune('▇')
			continue
		}
		buffer.WriteRune('█')
	}
}
