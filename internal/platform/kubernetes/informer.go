package kubernetes

import (
	"context"
	"log/slog"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/platform"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type InformerGrapher struct {
	clientset       *kubernetes.Clientset
	informerFactory informers.SharedInformerFactory

	events chan struct{}
	ch     chan platform.Graph

	close chan struct{}

	includeOldReplicaSets bool
}

func NewInformerGrapher(clientset *kubernetes.Clientset, includeOldReplicaSets bool) (*InformerGrapher, error) {
	grapher := &InformerGrapher{
		clientset: clientset,
		// TODO: Make resync configurable
		informerFactory: informers.NewSharedInformerFactory(clientset, 30*time.Minute),

		events: make(chan struct{}),
		ch:     make(chan platform.Graph),

		includeOldReplicaSets: includeOldReplicaSets,
	}

	informerFuncs := []cache.SharedIndexInformer{
		grapher.informerFactory.Apps().V1().Deployments().Informer(),
		grapher.informerFactory.Apps().V1().DaemonSets().Informer(),
		grapher.informerFactory.Apps().V1().ReplicaSets().Informer(),
		grapher.informerFactory.Apps().V1().StatefulSets().Informer(),
		grapher.informerFactory.Batch().V1().CronJobs().Informer(),
		grapher.informerFactory.Batch().V1().Jobs().Informer(),
		grapher.informerFactory.Core().V1().Pods().Informer(),
	}

	for _, informerFunc := range informerFuncs {
		_, err := informerFunc.AddEventHandler(grapher)
		if err != nil {
			return nil, err
		}
	}

	go func() {
		defer close(grapher.ch)

		for range grapher.events {
			// TODO: Make the timeout configurable? 30s should be plenty, but who
			// knows
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
			graph, err := grapher.Graph(ctx)
			cancel()
			if err != nil {
				slog.ErrorContext(ctx, "Failed to graph informer", slog.Any("error", err))
				continue
			}

			grapher.ch <- graph
		}
	}()

	return grapher, nil
}

func (g *InformerGrapher) Start() {

	g.close = make(chan struct{})
	// Trigger once after sync
	go func() {
		g.informerFactory.WaitForCacheSync(g.close)
		g.events <- struct{}{}
	}()
	g.informerFactory.Start(g.close)
}

func (g *InformerGrapher) Stop() {
	close(g.close)
	g.informerFactory.Shutdown()
	close(g.events)
}

func (g *InformerGrapher) Graph(ctx context.Context) (platform.Graph, error) {
	graph := platform.NewGraph()

	deployments, err := g.informerFactory.Apps().V1().Deployments().Lister().List(labels.Everything())
	if err != nil {
		return nil, err
	}
	for _, object := range deployments {
		addObjectToGraph(graph, object, g.includeOldReplicaSets)
	}

	daemonSets, err := g.informerFactory.Apps().V1().DaemonSets().Lister().List(labels.Everything())
	if err != nil {
		return nil, err
	}
	for _, object := range daemonSets {
		addObjectToGraph(graph, object, g.includeOldReplicaSets)
	}

	replicaSets, err := g.informerFactory.Apps().V1().ReplicaSets().Lister().List(labels.Everything())
	if err != nil {
		return nil, err
	}
	for _, object := range replicaSets {
		addObjectToGraph(graph, object, g.includeOldReplicaSets)
	}

	statefulSets, err := g.informerFactory.Apps().V1().StatefulSets().Lister().List(labels.Everything())
	if err != nil {
		return nil, err
	}
	for _, object := range statefulSets {
		addObjectToGraph(graph, object, g.includeOldReplicaSets)
	}

	cronJobs, err := g.informerFactory.Batch().V1().CronJobs().Lister().List(labels.Everything())
	if err != nil {
		return nil, err
	}
	for _, object := range cronJobs {
		addObjectToGraph(graph, object, g.includeOldReplicaSets)
	}

	jobs, err := g.informerFactory.Batch().V1().Jobs().Lister().List(labels.Everything())
	if err != nil {
		return nil, err
	}
	for _, object := range jobs {
		addObjectToGraph(graph, object, g.includeOldReplicaSets)
	}

	pods, err := g.informerFactory.Core().V1().Pods().Lister().List(labels.Everything())
	if err != nil {
		return nil, err
	}
	for _, object := range pods {
		addObjectToGraph(graph, object, g.includeOldReplicaSets)
	}

	return graph, nil
}

func (g *InformerGrapher) Graphs() <-chan platform.Graph {
	return g.ch
}

func (g *InformerGrapher) onEvent(isInitialList bool) {
	if isInitialList {
		return
	}

	g.events <- struct{}{}
}

func (g *InformerGrapher) OnAdd(object any, isInitialList bool) {
	g.onEvent(isInitialList)
}

func (g *InformerGrapher) OnUpdate(oldObject any, newObject any) {
	g.onEvent(false)
}

func (g *InformerGrapher) OnDelete(object any) {
	g.onEvent(false)
}
