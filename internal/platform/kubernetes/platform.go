package kubernetes

import (
	"context"
	"fmt"

	"github.com/AlexGustafsson/cupdate/internal/platform"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var _ platform.Grapher = (*Platform)(nil)
var _ platform.ContinuousGrapher = (*Platform)(nil)

// Platform implements graphing for the Kubernetes platform.
type Platform struct {
	clientset *kubernetes.Clientset

	includeOldReplicaSets bool
}

type Options struct {
	// IncludeOldReplicaSets will include all replica sets, no matter their age.
	IncludeOldReplicaSets bool
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

	return &Platform{
		clientset: clientset,

		includeOldReplicaSets: options.IncludeOldReplicaSets,
	}, nil
}

// Graph implements platform.Platform.
func (p *Platform) Graph(ctx context.Context) (platform.Graph, error) {
	// TODO: Do we need to adhere to this interface if we only ever intend for
	// GraphContinuously to be used? Could Graph use GraphContinuously once?
	return nil, fmt.Errorf("not implemented")
}

func (p *Platform) GraphContinuously(ctx context.Context) (<-chan platform.Graph, error) {
	grapher, err := NewInformerGrapher(p.clientset, p.includeOldReplicaSets)
	if err != nil {
		return nil, err
	}

	grapher.Start()

	go func() {
		<-ctx.Done()
		grapher.Stop()
	}()

	return grapher.Graphs(), nil
}
