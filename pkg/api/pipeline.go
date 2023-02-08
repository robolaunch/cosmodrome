package api

import "github.com/dominikbraun/graph"

type PipelineInterface interface {
	Build()
	Push()
	View()
}

type PipelineAbstract struct {
	PipelineInterface
	Name string
	DAG  graph.Graph[string, BuildComponent]
}

type Pipeline struct {
	PipelineAbstract
}
