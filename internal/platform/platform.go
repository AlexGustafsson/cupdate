package platform

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"strings"
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

// Labels holds labels / annotations found by platform implementations, which
// map to things like Docker labels or Kubernetes resource annotations.
type Labels map[string]string

// Ignore returns true if the Cupdate ignore label is set to true.
func (l Labels) Ignore() bool {
	if l == nil {
		return false
	}

	if v, ok := l["config.cupdate/ignore"]; ok {
		return v == "true"
	}

	if v, ok := l["cupdate.config.ignore"]; ok {
		return v == "true"
	}

	return false
}

// Ignore returns true if the Cupdate pin label is set to true.
func (l Labels) Pin() bool {
	if l == nil {
		return false
	}

	if v, ok := l["config.cupdate/pin"]; ok {
		return v == "true"
	}

	if v, ok := l["cupdate.config.pin"]; ok {
		return v == "true"
	}

	return false
}

// Ignore returns true if the Cupdate stay-on-current-major label is set to
// true.
func (l Labels) StayOnCurrentMajor() bool {
	if l == nil {
		return false
	}

	if v, ok := l["config.cupdate/stay-on-current-major"]; ok {
		return v == "true"
	}

	if v, ok := l["cupdate.config.stay-on-current-major"]; ok {
		return v == "true"
	}

	return false
}

// RemoveUnsupported removes unsupported labels.
func (l Labels) RemoveUnsupported() Labels {
	clone := maps.Clone(l)
	for k := range l {
		if !strings.HasPrefix(k, "config.cupdate/") && !strings.HasPrefix(k, "cupdate.config.") {
			delete(clone, k)
		}
	}
	return clone
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
	// Graph returns a graph of all images found on the platform.
	// The graph's roots are [ImageNode]s.
	Graph(context.Context) (Graph, error)
}

// Grapher provides asynchronous updates of graphs of resources in a platform.
type ContinuousGrapher interface {
	// GraphContinuously returns a channel which will receive a graph of all
	// images found on the platform whenever the graph changes.
	// The graph's roots are [ImageNode]s.
	GraphContinuously(context.Context) (<-chan Graph, error)
}

var _ (ContinuousGrapher) = (*PollGrapher)(nil)

// PollGrapher is a [ContinuousGrapher] implementation which polls an underlying
// [Grapher] implementation.
type PollGrapher struct {
	Grapher  Grapher
	Interval time.Duration
}

// GraphContinuously implements ContinuousGrapher.
func (g *PollGrapher) GraphContinuously(ctx context.Context) (<-chan Graph, error) {
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

// Graph implements Grapher.
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

// NewGraph returns a new [Graph].
func NewGraph() Graph {
	return graph.New[Node]()
}
