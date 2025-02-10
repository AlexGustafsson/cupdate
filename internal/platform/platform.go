package platform

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/graph"
	"github.com/AlexGustafsson/cupdate/internal/oci"
)

type Node interface {
	ID() string
	Type() string
	Labels() Labels
}

type Labels map[string]string

func (l Labels) Ignore() bool {
	if v, ok := l["config.cupdate/ignore"]; ok {
		return v == "true"
	}

	if v, ok := l["cupdate.config.ignore"]; ok {
		return v == "true"
	}

	return false
}

type Graph = *graph.Graph[Node]

var _ Node = (*ImageNode)(nil)

type ImageNode struct {
	Reference oci.Reference
}

func (n ImageNode) ID() string {
	return fmt.Sprintf("oci/image/%s", n.Reference)
}

func (n ImageNode) Type() string {
	return "image"
}

func (n ImageNode) Labels() Labels {
	return nil
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
	ch := make(chan Graph, 1)

	slog.DebugContext(ctx, "Polling graph")
	graph, err := g.Grapher.Graph(ctx)
	if err != nil {
		return nil, err
	}
	ch <- graph

	go func() {
		defer close(ch)

		ticker := time.NewTicker(g.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				slog.DebugContext(ctx, "Polling graph")
				graph, err := g.Grapher.Graph(ctx)
				if err != nil {
					slog.ErrorContext(ctx, "Failed to poll graph", slog.Any("error", err))
					continue
				}

				ch <- graph
			}
		}
	}()

	return ch, nil
}

// CompoundGrapher creates a graph from one or more [Grapher] simultaneously.
type CompoundGrapher struct {
	Graphers []Grapher
}

func (g *CompoundGrapher) Graph(ctx context.Context) (Graph, error) {
	graphs := make([]Graph, len(g.Graphers))
	errs := make([]error, len(g.Graphers))

	// Don't use ErrGroup in order to retain all errors to help users debug any
	// issues
	var wg sync.WaitGroup
	for i, grapher := range g.Graphers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			graphs[i], errs[i] = grapher.Graph(ctx)
		}()
	}
	wg.Wait()

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	compoundGraph := NewGraph()
	for _, graph := range graphs {
		compoundGraph.InsertGraph(graph)
	}

	return compoundGraph, nil
}

func NewGraph() Graph {
	return graph.New[Node]()
}
