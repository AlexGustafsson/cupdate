package kubernetes

import (
	"fmt"
	"maps"

	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/platform"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func getImageReference(specImage string, statusImage string, statusImageID string) (oci.Reference, error) {
	// Parse the image as defined by the spec
	specRef, err := oci.ParseReference(specImage)
	if err != nil {
		return oci.Reference{}, err
	}

	// Try to parse the image as defined by the runtime, which can contain the
	// resolved ids. Note that these are not well-defined.
	// See: https://github.com/kubernetes/kubernetes/issues/115199
	//
	// Example from a container's spec:
	// "image": "yooooomi/your_spotify_client:1.12.0"
	//
	// Examples from a container's status in microk8s/containerd:
	// "image": "docker.io/yooooomi/your_spotify_client:1.12.0"
	// "imageID": "docker.io/yooooomi/your_spotify_client@sha256:9c20e607b82fc66a0b81a45c04d6ccd8fd056add3a3adacc0bb7a6b99460fdb0"
	// "containerID": "containerd://dadf5d1aca357f514cd558b11140786e46e729f25ddc7f847382fdff127a44b8"
	//
	// Note that the digest is not always to the manifest list, but might as well
	// be to the actual manifest in use. It is not well-defined
	statusRef, statusErr := oci.ParseReference(statusImage)
	statusRuntimeRef, statusRuntimeErr := oci.ParseReference(statusImageID)

	// In the best of worlds, the status' image id contains a tag and a digest, if
	// so, try to use it
	if statusRuntimeErr == nil && statusRuntimeRef.HasDigest {
		// Fall back to the tag specified by the spec
		if !statusRuntimeRef.HasTag && specRef.HasTag {
			statusRuntimeRef.HasTag = true
			statusRuntimeRef.Tag = specRef.Tag
		}
		return statusRuntimeRef, nil
	}

	// If the image id doesn't have a digest, it's not much more useful than just
	// the plain image reference. Try to use the resolved one anyway, but ensure
	// that it contains at least as much information as the one from the spec
	if statusErr == nil && statusRef.HasTag == specRef.HasTag && statusRef.HasDigest == specRef.HasDigest {
		return statusRef, nil
	}

	// If we can't find a resolved reference, use the one from the spec
	return specRef, nil
}

func addObjectToGraph(graph platform.Graph, nodeResource resource, resources map[types.UID]v1.Object, object v1.Object) {
	objectResource := mapObjectToResource(object)

	ownerReferences := object.GetOwnerReferences()

	// The object is a leaf node, no need to continue to traverse the hierarchy
	if len(ownerReferences) == 0 {
		namespace := object.GetNamespace()

		// The resource has no owning entity, just add the resource to the graph
		if namespace == "" {
			graph.InsertTree(objectResource, nodeResource)
			return
		}

		// The resource has a namespace, add a meta resource for it
		graph.InsertTree(
			objectResource,
			resource{
				kind: ResourceKindCoreV1Namespace,
				id:   fmt.Sprintf("kubernetes/%s", namespace),
				name: namespace,
			},
			nodeResource,
		)

		return
	}

	for _, ownerReference := range ownerReferences {
		ownerObject, ok := resources[ownerReference.UID]

		// We haven't seen this owner reference when watching Kubernetes'
		// resources - it's likely that it's a CRD we don't support
		if !ok {
			graph.InsertTree(
				objectResource,
				resource{
					kind: ResourceKindUnknown,
					id:   fmt.Sprintf("kubernetes/%s", ownerReference.UID),
					name: ownerReference.Name,
				},
			)
			continue
		}

		// Insert the relationship between the resource and its owner
		ownerResource := mapObjectToResource(ownerObject)
		graph.InsertTree(
			objectResource,
			ownerResource,
		)

		// Assuming there are no cycles, continue to traverse the hierarchy
		addObjectToGraph(graph, nodeResource, resources, ownerObject)
	}
}

func mapObjectToResource(object v1.Object) resource {
	switch o := object.(type) {
	case *appsv1.Deployment:
		return resource{
			kind:   ResourceKindAppsV1Deployment,
			id:     fmt.Sprintf("kubernetes/%s", o.UID),
			name:   o.Name,
			labels: maps.Clone(o.Labels),
		}
	case *appsv1.DaemonSet:
		return resource{
			kind:   ResourceKindAppsV1DaemonSet,
			id:     fmt.Sprintf("kubernetes/%s", o.UID),
			name:   o.Name,
			labels: maps.Clone(o.Labels),
		}
	case *appsv1.ReplicaSet:
		return resource{
			kind:   ResourceKindAppsV1ReplicaSet,
			id:     fmt.Sprintf("kubernetes/%s", o.UID),
			name:   o.Name,
			labels: maps.Clone(o.Labels),
		}
	case *appsv1.StatefulSet:
		return resource{
			kind:   ResourceKindAppsV1StatefulSet,
			id:     fmt.Sprintf("kubernetes/%s", o.UID),
			name:   o.Name,
			labels: maps.Clone(o.Labels),
		}
	case *batchv1.CronJob:
		return resource{
			kind:   ResourceKindBatchV1CronJob,
			id:     fmt.Sprintf("kubernetes/%s", o.UID),
			name:   o.Name,
			labels: maps.Clone(o.Labels),
		}
	case *batchv1.Job:
		return resource{
			kind:   ResourceKindBatchV1Job,
			id:     fmt.Sprintf("kubernetes/%s", o.UID),
			name:   o.Name,
			labels: maps.Clone(o.Labels),
		}
	case *corev1.Pod:
		return resource{
			kind:   ResourceKindCoreV1Pod,
			id:     fmt.Sprintf("kubernetes/%s", o.UID),
			name:   o.Name,
			labels: maps.Clone(o.Labels),
		}
	default:
		// The object is not something common that we support, return a catch all
		// resource
		return resource{
			kind: ResourceKindUnknown,
			id:   fmt.Sprintf("kubernetes/%s", object.GetUID()),
			name: object.GetName(),
		}
	}
}
