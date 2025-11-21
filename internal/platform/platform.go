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

// Node is a platform resource represented as a graph node.
type Node interface {
	// ID implements [graph.Node].
	ID() string
	// Type returns a type uniquely describing the node.
	Type() string
	// Labels returns labels set on the resource represented by the node.
	Labels() Labels
	// InternalLabels returns labels set by Cupdate.
	InternalLabels() InternalLabels
}

const (
	InternalLabelHostArchitecture string = "host-architecture"
	InternalLabelOperatingSystem  string = "host-operating-system"
)

// InternalLabels holds labels maintained by Cupdate for use by Cupdate.
type InternalLabels map[string]string

func (l InternalLabels) InternalCupdateArchitecture() string {
	if l == nil {
		return ""
	}

	return l[InternalLabelHostArchitecture]
}

func (l InternalLabels) InternalCupdateOperatingSystem() string {
	if l == nil {
		return ""
	}

	return l[InternalLabelOperatingSystem]
}

// Graph is a graph implementation holding [Nodes].
type Graph = *graph.Graph[Node]

var _ Node = (*ImageNode)(nil)

// ImageNode represents a resource common to all platforms - the OCI image.
type ImageNode struct {
	Reference oci.Reference
}

// ID implements Node.
func (n ImageNode) ID() string {
	return fmt.Sprintf("oci/image/%s", n.Reference)
}

// Type implements Node.
func (n ImageNode) Type() string {
	return "image"
}

// Labels implements Node.
func (n ImageNode) Labels() Labels {
	return nil
}

// InternalLabels implements Node.
func (n ImageNode) InternalLabels() InternalLabels {
	return nil
}

// String implements Node.
func (n ImageNode) String() string {
	return fmt.Sprintf("%s<%s>", n.Type(), n.Reference.String())
}

// Grapher provides graphs of resources in a platform.
type Grapher interface {
	// Graph graphs all images found on the platform.
	// The graph's roots are [ImageNode]s.
	Graph(context.Context) (Graph, error)
}

// Grapher provides asynchronous updates of graphs of resources in a platform.
type ContinuousGrapher interface {
	// Graph graphs all images found on the platform.
	// The graph is published on the channel returned by Graphs.
	// The graph's roots are [ImageNode]s.
	Graph(context.Context) error
	// Graphs returns a channel which will receive a graph of all images found on
	// the platform whenever the graph changes.
	// The graph's roots are [ImageNode]s.
	Graphs() <-chan Graph
	// Close closes the grapher.
	Close() error
}

var _ (ContinuousGrapher) = (*PollGrapher)(nil)

// PollGrapher is a [ContinuousGrapher] implementation which polls an underlying
// [Grapher] implementation.
type PollGrapher struct {
	grapher Grapher
	graphs  chan Graph
	ticker  *time.Ticker

	close chan struct{}
	done  chan struct{}
}

func NewPollGrapher(grapher Grapher, interval time.Duration) *PollGrapher {
	g := &PollGrapher{
		grapher: grapher,
		graphs:  make(chan Graph),
		ticker:  time.NewTicker(interval),

		close: make(chan struct{}),
		done:  make(chan struct{}),
	}

	go func() {
		defer close(g.done)
		for {
			select {
			case <-g.close:
				// time.Ticker's channel is not closed on stop
				return
			case <-g.ticker.C:
				// TODO: Make the timeout configurable? 30s should be plenty, but who
				// knows
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				slog.DebugContext(ctx, "Polling graph")
				err := g.Graph(ctx)
				cancel()
				if err != nil {
					slog.ErrorContext(ctx, "Failed to poll graph", slog.Any("error", err))
					continue
				}
			}
		}
	}()

	return g
}

// Graph implements ContinuousGrapher.
func (g *PollGrapher) Graph(ctx context.Context) error {
	graph, err := g.grapher.Graph(ctx)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case g.graphs <- graph:
		return nil
	}
}

// Graphs implements ContinuousGrapher.
func (g *PollGrapher) Graphs() <-chan Graph {
	return g.graphs
}

func (g *PollGrapher) Close() error {
	g.ticker.Stop()
	close(g.close)
	<-g.done
	return nil
}

// CompoundGrapher creates a graph from one or more [Grapher] simultaneously.
type CompoundGrapher struct {
	Graphers []Grapher
}

// Graph implements Grapher.
func (g *CompoundGrapher) Graph(ctx context.Context) (Graph, error) {
	graphs := make([]Graph, len(g.Graphers))
	errs := make([]error, len(g.Graphers))

	// Don't use ErrGroup in order to retain all errors to help users debug any
	// issues
	var wg sync.WaitGroup
	for i, grapher := range g.Graphers {
		wg.Go(func() {
			graphs[i], errs[i] = grapher.Graph(ctx)
		})
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

// NewGraph returns a new [Graph].
func NewGraph() Graph {
	return graph.New[Node]()
}
