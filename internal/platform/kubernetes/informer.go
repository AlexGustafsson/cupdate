package kubernetes

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/platform"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
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

	mutex sync.Mutex
}

func NewInformerGrapher(clientset *kubernetes.Clientset, includeOldReplicaSets bool) (*InformerGrapher, error) {
	grapher := &InformerGrapher{
		clientset: clientset,
		// TODO: Make resync configurable
		informerFactory:       informers.NewSharedInformerFactory(clientset, 30*time.Minute),
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

	return grapher, nil
}

func (g *InformerGrapher) Start() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.ch = make(chan platform.Graph)
	g.events = make(chan struct{})
	g.close = make(chan struct{})

	// Handle events and produce graphs
	go func() {
		defer close(g.ch)

		for range g.events {
			slog.Debug("Got informer event from Kubernetes, graphing")
			// TODO: Make the timeout configurable? 30s should be plenty, but who
			// knows
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			graph, err := g.Graph(ctx)
			cancel()
			if err != nil {
				slog.ErrorContext(ctx, "Failed to graph informer", slog.Any("error", err))
				continue
			}

			g.ch <- graph
		}
	}()

	// Trigger once after sync
	go func() {
		g.informerFactory.WaitForCacheSync(g.close)
		g.events <- struct{}{}
	}()

	g.informerFactory.Start(g.close)
}

func (g *InformerGrapher) Stop() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.close != nil {
		close(g.close)
		g.close = nil
	}

	g.informerFactory.Shutdown()

	if g.events != nil {
		close(g.events)
		g.events = nil
	}
}

func (g *InformerGrapher) Graph(ctx context.Context) (platform.Graph, error) {
	resources := make(map[types.UID]v1.Object)

	deployments, err := g.informerFactory.Apps().V1().Deployments().Lister().List(labels.Everything())
	if err != nil {
		return nil, err
	}
	for _, object := range deployments {
		resources[object.UID] = object
	}

	daemonSets, err := g.informerFactory.Apps().V1().DaemonSets().Lister().List(labels.Everything())
	if err != nil {
		return nil, err
	}
	for _, object := range daemonSets {
		resources[object.UID] = object
	}

	replicaSets, err := g.informerFactory.Apps().V1().ReplicaSets().Lister().List(labels.Everything())
	if err != nil {
		return nil, err
	}
	for _, object := range replicaSets {
		resources[object.UID] = object
	}

	statefulSets, err := g.informerFactory.Apps().V1().StatefulSets().Lister().List(labels.Everything())
	if err != nil {
		return nil, err
	}
	for _, object := range statefulSets {
		resources[object.UID] = object
	}

	cronJobs, err := g.informerFactory.Batch().V1().CronJobs().Lister().List(labels.Everything())
	if err != nil {
		return nil, err
	}
	for _, object := range cronJobs {
		resources[object.UID] = object
	}

	jobs, err := g.informerFactory.Batch().V1().Jobs().Lister().List(labels.Everything())
	if err != nil {
		return nil, err
	}
	for _, object := range jobs {
		resources[object.UID] = object
	}

	pods, err := g.informerFactory.Core().V1().Pods().Lister().List(labels.Everything())
	if err != nil {
		return nil, err
	}
	for _, object := range pods {
		resources[object.UID] = object
	}

	graph := platform.NewGraph()
	for _, pod := range pods {
		didAddImage := false
		for _, containerSpec := range pod.Spec.Containers {

			// Resolve the container's image reference
			specImage := containerSpec.Image
			var statusImage, statusImageID string
			// Note that container statuses are not well-defined. See SDK docs.
			// Note that container statuses are not always available, for example,
			// before a cron job has created it.
			// Use the first match
			for _, status := range pod.Status.ContainerStatuses {
				if status.Name == containerSpec.Name {
					statusImage = status.Image
					statusImageID = status.ImageID
					break
				}
			}

			ref, err := getImageReference(specImage, statusImage, statusImageID)
			if err != nil {
				slog.ErrorContext(ctx, "Failed to identify a valid image reference for container", slog.String("pod", pod.Name), slog.String("container", containerSpec.Name))
				continue
			}

			// Just when Kubernetes has started a pod, the runtime won't have resolved
			// the reference, meaning the digest is empty. Immediately after, the
			// digest is resolved and the change is processed by the informer again.
			// Therefore ignore references without digests and assume that we'll get
			// them soon, with a digest
			if !ref.HasDigest {
				slog.DebugContext(ctx, "Ignoring reference without digest", slog.String("reference", ref.Reference()), slog.String("pod", pod.Name), slog.String("container", containerSpec.Name))
				continue
			}

			graph.InsertTree(
				platform.ImageNode{
					Reference: ref,
				},
				resource{
					id:   fmt.Sprintf("kubernetes/%s/container/%s", pod.UID, containerSpec.Name),
					kind: ResourceKindCoreV1Container,
					name: containerSpec.Name,
				},
				// This node is technically already added by addObjectToGraph later on,
				// but we still need to reference the resource to connect the relation
				resource{
					kind: ResourceKindCoreV1Pod,
					id:   fmt.Sprintf("kubernetes/%s", pod.UID),
					name: pod.Name,
				},
			)
			didAddImage = true
		}

		// If we found and added a valid image, resolve and add the rest of the
		// pod's hierarchy
		if didAddImage {
			addObjectToGraph(graph, resources, pod)
		}
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
