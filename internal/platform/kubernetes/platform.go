package kubernetes

import (
	"context"
	"fmt"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/platform"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"golang.org/x/sync/errgroup"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/pager"
)

var _ platform.Platform = (*Platform)(nil)

type Platform struct {
	clientset *kubernetes.Clientset
}

func NewPlatform(config *rest.Config) (*Platform, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Platform{
		clientset: clientset,
	}, nil
}

func (p *Platform) Graph(ctx context.Context) (platform.Graph, error) {
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
				switch o := obj.(type) {
				case *appsv1.Deployment:
					for _, container := range o.Spec.Template.Spec.Containers {
						reference, err := oci.ParseReference(container.Image)
						if err != nil {
							return err
						}

						graph.InsertTree(
							platform.ImageNode{
								Reference: reference,
							},
							resource{
								kind: ResourceKindCoreV1Container,
								id:   fmt.Sprintf("kubernetes/%s/container/%s", o.UID, container.Name),
								name: container.Name,
							},
							resource{
								kind: ResourceKindCoreV1Pod,
								id:   fmt.Sprintf("kubernetes/%s/pod", o.UID),
								name: o.Spec.Template.Name,
							},
							resource{
								kind: ResourceKindAppsV1Deployment,
								id:   fmt.Sprintf("kubernetes/%s", o.UID),
								name: o.Name,
							},
						)
					}
				case *appsv1.DaemonSet:
					for _, container := range o.Spec.Template.Spec.Containers {
						reference, err := oci.ParseReference(container.Image)
						if err != nil {
							return err
						}

						graph.InsertTree(
							platform.ImageNode{
								Reference: reference,
							},
							resource{
								kind: ResourceKindCoreV1Container,
								id:   fmt.Sprintf("kubernetes/%s/container/%s", o.UID, container.Name),
								name: container.Name,
							},
							resource{
								kind: ResourceKindCoreV1Pod,
								id:   fmt.Sprintf("kubernetes/%s/pod", o.UID),
								name: o.Spec.Template.Name,
							},
							resource{
								kind: ResourceKindAppsV1DaemonSet,
								id:   fmt.Sprintf("kubernetes/%s", o.UID),
								name: o.Name,
							},
						)
					}
				case *appsv1.ReplicaSet:
					for _, container := range o.Spec.Template.Spec.Containers {
						reference, err := oci.ParseReference(container.Image)
						if err != nil {
							return err
						}

						graph.InsertTree(
							platform.ImageNode{
								Reference: reference,
							},
							resource{
								kind: ResourceKindCoreV1Container,
								id:   fmt.Sprintf("kubernetes/%s/container/%s", o.UID, container.Name),
								name: container.Name,
							},
							resource{
								kind: ResourceKindCoreV1Pod,
								id:   fmt.Sprintf("kubernetes/%s/pod", o.UID),
								name: o.Spec.Template.Name,
							},
							resource{
								kind: ResourceKindAppsV1ReplicaSet,
								id:   fmt.Sprintf("kubernetes/%s", o.UID),
								name: o.Name,
							},
						)
					}
				case *appsv1.StatefulSet:
					for _, container := range o.Spec.Template.Spec.Containers {
						reference, err := oci.ParseReference(container.Image)
						if err != nil {
							return err
						}

						graph.InsertTree(
							platform.ImageNode{
								Reference: reference,
							},
							resource{
								kind: ResourceKindCoreV1Container,
								id:   fmt.Sprintf("kubernetes/%s/container/%s", o.UID, container.Name),
								name: container.Name,
							},
							resource{
								kind: ResourceKindCoreV1Pod,
								id:   fmt.Sprintf("kubernetes/%s/pod", o.UID),
								name: o.Spec.Template.Name,
							},
							resource{
								kind: ResourceKindAppsV1StatefulSet,
								id:   fmt.Sprintf("kubernetes/%s", o.UID),
								name: o.Name,
							},
						)
					}
				case *batchv1.CronJob:
					for _, container := range o.Spec.JobTemplate.Spec.Template.Spec.Containers {
						reference, err := oci.ParseReference(container.Image)
						if err != nil {
							return err
						}

						graph.InsertTree(
							platform.ImageNode{
								Reference: reference,
							},
							resource{
								kind: ResourceKindCoreV1Container,
								id:   fmt.Sprintf("kubernetes/%s/container/%s", o.UID, container.Name),
								name: container.Name,
							},
							resource{
								kind: ResourceKindCoreV1Pod,
								id:   fmt.Sprintf("kubernetes/%s/pod", o.UID),
								name: o.Spec.JobTemplate.Name,
							},
							resource{
								kind: ResourceKindBatchV1CronJob,
								id:   fmt.Sprintf("kubernetes/%s", o.UID),
								name: o.Name,
							},
						)
					}
				case *batchv1.Job:
					for _, container := range o.Spec.Template.Spec.Containers {
						reference, err := oci.ParseReference(container.Image)
						if err != nil {
							return err
						}

						graph.InsertTree(
							platform.ImageNode{
								Reference: reference,
							},
							resource{
								kind: ResourceKindCoreV1Container,
								id:   fmt.Sprintf("kubernetes/%s/container/%s", o.UID, container.Name),
								name: container.Name,
							},
							resource{
								kind: ResourceKindCoreV1Pod,
								id:   fmt.Sprintf("kubernetes/%s/pod", o.UID),
								name: o.Spec.Template.Name,
							},
							resource{
								kind: ResourceKindBatchV1Job,
								id:   fmt.Sprintf("kubernetes/%s", o.UID),
								name: o.Name,
							},
						)
					}
				case *corev1.Pod:
					for _, container := range o.Spec.Containers {
						// For now, let's assume a pod only has one owning reference
						var parent Resource
						if len(o.OwnerReferences) > 0 {
							parent = resource{
								id:   fmt.Sprintf("kubernetes/%s", o.OwnerReferences[0].UID),
								kind: ResourceKind(strings.ToLower(o.OwnerReferences[0].APIVersion + "/" + o.OwnerReferences[0].Kind)),
								name: o.OwnerReferences[0].Name,
							}
						}

						reference, err := oci.ParseReference(container.Image)
						if err != nil {
							return err
						}

						tree := []platform.Node{
							platform.ImageNode{
								Reference: reference,
							},
							resource{
								kind: ResourceKindCoreV1Container,
								id:   fmt.Sprintf("kubernetes/%s/container/%s", o.UID, container.Name),
								name: container.Name,
							},
							resource{
								kind: ResourceKindCoreV1Pod,
								id:   fmt.Sprintf("kubernetes/%s/pod", o.UID),
								name: o.Name,
							},
						}

						if parent != nil {
							tree = append(tree, parent)
						}

						graph.InsertTree(tree...)
					}
				default:
					// Panic as missing entries would be a programming issue, not runtime
					// bug
					panic("unsupported object kind")
				}

				return nil
			})
		})
	}

	if err := wg.Wait(); err != nil {
		return nil, err
	}

	return graph, nil
}
