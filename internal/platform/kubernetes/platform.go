package kubernetes

import (
	"context"

	"github.com/AlexGustafsson/cupdate/internal/platform"
	"golang.org/x/sync/errgroup"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/pager"
)

var _ platform.Grapher = (*Platform)(nil)
var _ platform.ContinousGrapher = (*Platform)(nil)

type Platform struct {
	clientset *kubernetes.Clientset

	includeOldReplicaSets bool
}

type Options struct {
	IncludeOldReplicaSets bool
}

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

func (p *Platform) Graph(ctx context.Context) (platform.Graph, error) {
	// TODO: Do we need to adhere to this interface if we only ever intend for
	// GraphContinously to be used? Could Graph use GraphContinously once?
	graph := platform.NewGraph()

	pageFuncs := []pager.ListPageFunc{
		func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
			return p.clientset.AppsV1().Deployments("").List(ctx, opts)
		},
		func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
			return p.clientset.AppsV1().DaemonSets("").List(ctx, opts)
		},
		func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
			return p.clientset.AppsV1().ReplicaSets("").List(ctx, opts)
		},
		func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
			return p.clientset.AppsV1().StatefulSets("").List(ctx, opts)
		},
		func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
			return p.clientset.BatchV1().CronJobs("").List(ctx, opts)
		},
		func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
			return p.clientset.BatchV1().Jobs("").List(ctx, opts)
		},
		func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
			return p.clientset.CoreV1().Pods("").List(ctx, opts)
		},
	}

	var wg errgroup.Group
	for _, pageFunc := range pageFuncs {
		pager := pager.New(pageFunc)
		wg.Go(func() error {
			return pager.EachListItem(ctx, metav1.ListOptions{}, func(obj runtime.Object) error {
				return addObjectToGraph(graph, obj, p.includeOldReplicaSets)
			})
		})
	}

	if err := wg.Wait(); err != nil {
		return nil, err
	}

	return graph, nil
}

func (p *Platform) GraphContinously(ctx context.Context) (<-chan platform.Graph, error) {
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
