package kubernetes

import (
	"fmt"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/platform"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func addObjectToGraph(graph platform.Graph, obj runtime.Object, includeOldReplicaSets bool) error {
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
				resource{
					kind: ResourceKindCoreV1Namespace,
					id:   fmt.Sprintf("kubernetes/%s", o.Namespace),
					name: o.Namespace,
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
				resource{
					kind: ResourceKindCoreV1Namespace,
					id:   fmt.Sprintf("kubernetes/%s", o.Namespace),
					name: o.Namespace,
				},
			)
		}
	case *appsv1.ReplicaSet:
		// Assume there's one replica
		replicas := 1
		if o.Spec.Replicas != nil {
			replicas = int(*o.Spec.Replicas)
		}

		if replicas == 0 && !includeOldReplicaSets {
			// Ignore the old replica set
			return nil
		}

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
				resource{
					kind: ResourceKindCoreV1Namespace,
					id:   fmt.Sprintf("kubernetes/%s", o.Namespace),
					name: o.Namespace,
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
				resource{
					kind: ResourceKindCoreV1Namespace,
					id:   fmt.Sprintf("kubernetes/%s", o.Namespace),
					name: o.Namespace,
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
				resource{
					kind: ResourceKindCoreV1Namespace,
					id:   fmt.Sprintf("kubernetes/%s", o.Namespace),
					name: o.Namespace,
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
				resource{
					kind: ResourceKindCoreV1Namespace,
					id:   fmt.Sprintf("kubernetes/%s", o.Namespace),
					name: o.Namespace,
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

			tree = append(tree, resource{
				kind: ResourceKindCoreV1Namespace,
				id:   fmt.Sprintf("kubernetes/%s", o.Namespace),
				name: o.Namespace,
			})

			graph.InsertTree(tree...)
		}
	default:
		// Panic as missing entries would be a programming issue, not runtime
		// bug
		panic("kubernetes: unsupported object kind <%s>")
	}

	return nil
}
