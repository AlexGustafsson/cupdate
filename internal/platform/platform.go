package platform

import (
	"context"
	"fmt"

	"github.com/AlexGustafsson/cupdate/internal/graph"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
)

type Node interface {
	ID() string
	Type() string
}

type Graph = *graph.Graph[Node]

type ImageNode struct {
	Reference oci.Reference
}

func (n ImageNode) ID() string {
	return fmt.Sprintf("oci/image/%s", n.Reference)
}

func (n ImageNode) Type() string {
	return "image"
}

func (n ImageNode) String() string {
	return fmt.Sprintf("%s<%s>", n.Type(), n.Reference.String())
}

type Platform interface {
	// Graph returns a graph of all images found on the platform.
	// The graph's roots are [ImageNode]s.
	Graph(context.Context) (Graph, error)
}

func NewGraph() Graph {
	return graph.New[Node]()
}
