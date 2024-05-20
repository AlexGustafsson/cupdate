package k8s

import (
	"context"
	"strings"

	"github.com/AlexGustafsson/k8s-image-feed/internal/source"
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

var _ source.Source = (*Source)(nil)

type Source struct {
	clientset *kubernetes.Clientset
}

func New(config *rest.Config) (*Source, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Source{
		clientset: clientset,
	}, nil
}

func (s *Source) EachListItem(ctx context.Context, fn func(source.Entry) error) error {
	pageFuncs := []pager.ListPageFunc{
		func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
			return s.clientset.AppsV1().Deployments("").List(ctx, opts)
		},
		func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
			return s.clientset.AppsV1().DaemonSets("").List(ctx, opts)
		},
		func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
			return s.clientset.AppsV1().ReplicaSets("").List(ctx, opts)
		},
		func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
			return s.clientset.AppsV1().StatefulSets("").List(ctx, opts)
		},
		func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
			return s.clientset.BatchV1().CronJobs("").List(ctx, opts)
		},
		func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
			return s.clientset.BatchV1().Jobs("").List(ctx, opts)
		},
		func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
			return s.clientset.CoreV1().Pods("").List(ctx, opts)
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
						image, version, _ := strings.Cut(container.Image, ":")
						fn(source.Entry{
							Image:   image,
							Version: version,
							Origin: &Origin{
								ResourceKind:  "apps/v1/Deployment",
								Namespace:     o.Namespace,
								Name:          o.Name,
								Created:       o.CreationTimestamp.UTC(),
								ContainerName: container.Name,
							},
						})
					}
				case *appsv1.DaemonSet:
					for _, container := range o.Spec.Template.Spec.Containers {
						image, version, _ := strings.Cut(container.Image, ":")
						fn(source.Entry{
							Image:   image,
							Version: version,
							Origin: &Origin{
								ResourceKind:  "apps/v1/DaemonSet",
								Namespace:     o.Namespace,
								Name:          o.Name,
								Created:       o.CreationTimestamp.UTC(),
								ContainerName: container.Name,
							},
						})
					}
				case *appsv1.ReplicaSet:
					for _, container := range o.Spec.Template.Spec.Containers {
						image, version, _ := strings.Cut(container.Image, ":")
						fn(source.Entry{
							Image:   image,
							Version: version,
							Origin: &Origin{
								ResourceKind:  "apps/v1/ReplicaSet",
								Namespace:     o.Namespace,
								Name:          o.Name,
								Created:       o.CreationTimestamp.UTC(),
								ContainerName: container.Name,
							},
						})
					}
				case *appsv1.StatefulSet:
					for _, container := range o.Spec.Template.Spec.Containers {
						image, version, _ := strings.Cut(container.Image, ":")
						fn(source.Entry{
							Image:   image,
							Version: version,
							Origin: &Origin{
								ResourceKind:  "apps/v1/StatefulSet",
								Namespace:     o.Namespace,
								Name:          o.Name,
								Created:       o.CreationTimestamp.UTC(),
								ContainerName: container.Name,
							},
						})
					}
				case *batchv1.CronJob:
					for _, container := range o.Spec.JobTemplate.Spec.Template.Spec.Containers {
						image, version, _ := strings.Cut(container.Image, ":")
						fn(source.Entry{
							Image:   image,
							Version: version,
							Origin: &Origin{
								ResourceKind:  "batch/v1/CronJob",
								Namespace:     o.Namespace,
								Name:          o.Name,
								Created:       o.CreationTimestamp.UTC(),
								ContainerName: container.Name,
							},
						})
					}
				case *batchv1.Job:
					for _, container := range o.Spec.Template.Spec.Containers {
						image, version, _ := strings.Cut(container.Image, ":")
						fn(source.Entry{
							Image:   image,
							Version: version,
							Origin: &Origin{
								ResourceKind:  "batch/v1/Job",
								Namespace:     o.Namespace,
								Name:          o.Name,
								Created:       o.CreationTimestamp.UTC(),
								ContainerName: container.Name,
							},
						})
					}
				case *corev1.Pod:
					for _, container := range o.Spec.Containers {
						image, version, _ := strings.Cut(container.Image, ":")
						fn(source.Entry{
							Image:   image,
							Version: version,
							Origin: &Origin{
								ResourceKind:  "core/v1/Pod",
								Namespace:     o.Namespace,
								Name:          o.Name,
								Created:       o.CreationTimestamp.UTC(),
								ContainerName: container.Name,
							},
						})
					}
				default:
					panic("unsupported object kind")
				}

				return nil
			})
		})
	}

	return wg.Wait()
}
