package platform

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/graph"
	"github.com/AlexGustafsson/cupdate/internal/oci"
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

type Grapher interface {
	// Graph returns a graph of all images found on the platform.
	// The graph's roots are [ImageNode]s.
	Graph(context.Context) (Graph, error)
}

type ContinousGrapher interface {
	// GraphContinously returns a channel which will receive a graph of all
	// images found on the platform whenever the graph changes.
	// The graph's roots are [ImageNode]s.
	GraphContinously(context.Context) (<-chan Graph, error)
}

var _ (ContinousGrapher) = (*PollGrapher)(nil)

type PollGrapher struct {
	Grapher  Grapher
	Interval time.Duration
}

func (g *PollGrapher) GraphContinously(ctx context.Context) (<-chan Graph, error) {

	ch := make(chan Graph)

	go func() {
		defer close(ch)

		ticker := time.NewTicker(g.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				slog.Debug("Polling graph")
				graph, err := g.Grapher.Graph(ctx)
				if err != nil {
					slog.Error("Failed to poll graph", slog.Any("error", err))
					continue
				}

				ch <- graph
			}
		}
	}()

	return ch, nil
}

func NewGraph() Graph {
	return graph.New[Node]()
}
