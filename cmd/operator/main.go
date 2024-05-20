package main

import (
	"context"
	"fmt"

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

func main() {
	// config, err := rest.InClusterConfig()
	// if err != nil {
	// 	panic(err.Error())
	// }
	config := &rest.Config{
		Host: "http://localhost:8001",
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pageFuncs := []pager.ListPageFunc{
		func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
			return clientset.AppsV1().Deployments("").List(ctx, opts)
		},
		func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
			return clientset.AppsV1().DaemonSets("").List(ctx, opts)
		},
		func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
			return clientset.AppsV1().ReplicaSets("").List(ctx, opts)
		},
		func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
			return clientset.AppsV1().StatefulSets("").List(ctx, opts)
		},
		func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
			return clientset.BatchV1().CronJobs("").List(ctx, opts)
		},
		func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
			return clientset.BatchV1().Jobs("").List(ctx, opts)
		},
		func(ctx context.Context, opts metav1.ListOptions) (runtime.Object, error) {
			return clientset.CoreV1().Pods("").List(ctx, opts)
		},
	}

	var wg errgroup.Group
	for _, pageFunc := range pageFuncs {
		pager := pager.New(pageFunc)
		wg.Go(func() error {
			return pager.EachListItem(context.TODO(), metav1.ListOptions{}, printImages)
		})
	}

	if err := wg.Wait(); err != nil {
		panic(err.Error())
	}
}

func printImages(obj runtime.Object) error {
	switch o := obj.(type) {
	case *appsv1.Deployment:
		for _, container := range o.Spec.Template.Spec.Containers {
			fmt.Println(o.Name, container.Name, container.Image)
		}
	case *appsv1.DaemonSet:
		for _, container := range o.Spec.Template.Spec.Containers {
			fmt.Println(o.Name, container.Name, container.Image)
		}
	case *appsv1.ReplicaSet:
		for _, container := range o.Spec.Template.Spec.Containers {
			fmt.Println(o.Name, container.Name, container.Image)
		}
	case *appsv1.StatefulSet:
		for _, container := range o.Spec.Template.Spec.Containers {
			fmt.Println(o.Name, container.Name, container.Image)
		}
	case *batchv1.CronJob:
		for _, container := range o.Spec.JobTemplate.Spec.Template.Spec.Containers {
			fmt.Println(o.Name, container.Name, container.Image)
		}
	case *batchv1.Job:
		for _, container := range o.Spec.Template.Spec.Containers {
			fmt.Println(o.Name, container.Name, container.Image)
		}
	case *corev1.Pod:
		for _, container := range o.Spec.Containers {
			fmt.Println(o.Name, container.Name, container.Image)
		}
	default:
		panic("unsupported object kind")
	}

	return nil
}
