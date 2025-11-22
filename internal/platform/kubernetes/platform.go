package kubernetes

import (
	"context"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/platform"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var _ platform.ContinuousGrapher = (*Platform)(nil)

// Platform implements graphing for the Kubernetes platform.
type Platform struct {
	grapher *InformerGrapher
}

type Options struct {
	// IncludeOldReplicaSets will include all replica sets, no matter their age.
	// Defaults to false.
	IncludeOldReplicaSets bool
	// DebounceInterval is an interval controlling the minimum duration between
	// graphs.
	// Defaults to one minute.
	DebounceInterval time.Duration
}

// NewPlatform initializes a new [Platform].
//
//   - config hold information about how to connect to the Kubernetes APIs.
func NewPlatform(config *rest.Config, options *Options) (*Platform, error) {
	if options == nil {
		options = &Options{}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	debounceInterval := 1 * time.Minute
	if options.DebounceInterval > 0 {
		debounceInterval = options.DebounceInterval
	}

	grapher, err := NewInformerGrapher(clientset, options.IncludeOldReplicaSets, debounceInterval)
	if err != nil {
		return nil, err
	}

	p := &Platform{
		grapher: grapher,
	}

	grapher.Start()

	return p, nil
}

// Graph implements platform.ContinuousGrapher.
func (p *Platform) Graph(ctx context.Context) error {
	return p.grapher.Graph(ctx)
}

// Graphs implements platform.ContinuousGrapher.
func (p *Platform) Graphs() <-chan platform.Graph {
	return p.grapher.Graphs()
}

// Close closes the platform.
func (p *Platform) Close() error {
	p.grapher.Close()
	return nil
}
