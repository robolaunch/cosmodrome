package api

import "github.com/dominikbraun/graph"

type PipelineInterface interface {
	Build()
	Push()
	View()
}

type PipelineAbstract struct {
	PipelineInterface
	Registry       string
	Name           string
	PushComponents bool
	DAG            graph.Graph[string, BuildComponent]
}

type Pipeline struct {
	PipelineAbstract
}
