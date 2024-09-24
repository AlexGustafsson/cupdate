package kubernetes

import (
	"context"
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

func (p *Platform) Images(ctx context.Context) ([]oci.Reference, platform.Graph, error) {
	references := make([]oci.Reference, 0)
	graph := make(platform.Graph)
	return references, graph, p.EachListItem(ctx, func(reference oci.Reference, origin platform.Origin) {
		references = append(references, reference)
		graph.AddOrigin(reference, origin)
	})
}

func (p *Platform) EachListItem(ctx context.Context, fn func(oci.Reference, platform.Origin)) error {
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

	// TODO: Build ImageID for resources that don't have them.
	// TODO: Connect pods with other resources (jobs etc.) to sort out duplicates.
	// We still want to cover deployments, jobs etc. that haven't run but still
	// refer to an image, though.

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

						fn(reference, &Origin{
							Container: Container{
								Name:      container.Name,
								Namespace: o.Namespace,
								Pod: &Pod{
									Name:       o.Spec.Template.Name,
									Namespace:  o.Namespace,
									IsTemplate: true,
									Parent: &Parent{
										ResourceKind: ResourceKindAppsV1Deployment,
										Namespace:    o.Namespace,
										Name:         o.Name,
									},
								},
							},
						})
					}
				case *appsv1.DaemonSet:
					for _, container := range o.Spec.Template.Spec.Containers {
						reference, err := oci.ParseReference(container.Image)
						if err != nil {
							return err
						}

						fn(reference, &Origin{
							Container: Container{
								Name:      container.Name,
								Namespace: o.Namespace,
								Pod: &Pod{
									Name:       o.Spec.Template.Name,
									Namespace:  o.Namespace,
									IsTemplate: true,
									Parent: &Parent{
										ResourceKind: ResourceKindAppsV1DaemonSet,
										Namespace:    o.Namespace,
										Name:         o.Name,
									},
								},
							},
						})
					}
				case *appsv1.ReplicaSet:
					for _, container := range o.Spec.Template.Spec.Containers {
						reference, err := oci.ParseReference(container.Image)
						if err != nil {
							return err
						}

						fn(reference, &Origin{
							Container: Container{
								Name:      container.Name,
								Namespace: o.Namespace,
								Pod: &Pod{
									Name:       o.Spec.Template.Name,
									Namespace:  o.Namespace,
									IsTemplate: true,
									Parent: &Parent{
										ResourceKind: ResourceKindAppsV1ReplicaSet,
										Namespace:    o.Namespace,
										Name:         o.Name,
									},
								},
							},
						})
					}
				case *appsv1.StatefulSet:
					for _, container := range o.Spec.Template.Spec.Containers {
						reference, err := oci.ParseReference(container.Image)
						if err != nil {
							return err
						}

						fn(reference, &Origin{
							Container: Container{
								Name:      container.Name,
								Namespace: o.Namespace,
								Pod: &Pod{
									Name:       o.Spec.Template.Name,
									Namespace:  o.Namespace,
									IsTemplate: true,
									Parent: &Parent{
										ResourceKind: ResourceKindAppsV1StatefulSet,
										Namespace:    o.Namespace,
										Name:         o.Name,
									},
								},
							},
						})
					}
				case *batchv1.CronJob:
					for _, container := range o.Spec.JobTemplate.Spec.Template.Spec.Containers {
						reference, err := oci.ParseReference(container.Image)
						if err != nil {
							return err
						}

						fn(reference, &Origin{
							Container: Container{
								Name:      container.Name,
								Namespace: o.Namespace,
								Pod: &Pod{
									Name:       o.Spec.JobTemplate.Spec.Template.Name,
									Namespace:  o.Namespace,
									IsTemplate: true,
									Parent: &Parent{
										ResourceKind: ResourceKindBatchV1CronJob,
										Namespace:    o.Namespace,
										Name:         o.Name,
									},
								},
							},
						})
					}
				case *batchv1.Job:
					for _, container := range o.Spec.Template.Spec.Containers {
						reference, err := oci.ParseReference(container.Image)
						if err != nil {
							return err
						}

						fn(reference, &Origin{
							Container: Container{
								Name:      container.Name,
								Namespace: o.Namespace,
								Pod: &Pod{
									Name:       o.Spec.Template.Name,
									Namespace:  o.Namespace,
									IsTemplate: true,
									Parent: &Parent{
										ResourceKind: ResourceKindBatchV1Job,
										Namespace:    o.Namespace,
										Name:         o.Name,
									},
								},
							},
						})
					}
				case *corev1.Pod:
					for _, container := range o.Spec.Containers {
						// For now, let's assume a pod only has one owning reference
						var parent *Parent
						if len(o.OwnerReferences) > 0 {
							parent = &Parent{
								ResourceKind: ResourceKind(strings.ToLower(o.OwnerReferences[0].APIVersion + "/" + o.OwnerReferences[0].Kind)),
								Namespace:    o.Namespace,
								Name:         o.Name,
							}
						}

						reference, err := oci.ParseReference(container.Image)
						if err != nil {
							return err
						}

						fn(reference, &Origin{
							Container: Container{
								Name:      container.Name,
								Namespace: o.Namespace,
								Pod: &Pod{
									Name:       o.Name,
									Namespace:  o.Namespace,
									IsTemplate: false,
									Created:    o.CreationTimestamp.UTC(),
									Parent:     parent,
								},
							},
						})
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

	return wg.Wait()
}
