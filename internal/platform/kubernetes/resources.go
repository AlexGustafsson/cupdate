package kubernetes

import (
	"fmt"

	"github.com/AlexGustafsson/cupdate/internal/platform"
)

// ResourceKind defines the types of resources exposed by the platform.
type ResourceKind string

const (
	ResourceKindAppsV1Deployment  = "apps/v1/deployment"
	ResourceKindAppsV1DaemonSet   = "apps/v1/daemonset"
	ResourceKindAppsV1ReplicaSet  = "apps/v1/replicaset"
	ResourceKindAppsV1StatefulSet = "apps/v1/statefulset"
	ResourceKindBatchV1CronJob    = "batch/v1/cronjob"
	ResourceKindBatchV1Job        = "batch/v1/job"
	ResourceKindCoreV1Node        = "core/v1/node"
	ResourceKindCoreV1Namespace   = "core/v1/namespace"
	ResourceKindCoreV1Pod         = "core/v1/pod"
	ResourceKindCoreV1Container   = "core/v1/container"
	ResourceKindUnknown           = "unknown"
)

// IsSupported returns whether or not the resource is supported.
// Filters out custom resource definitions.
func (r ResourceKind) IsSupported() bool {
	switch r {
	case ResourceKindAppsV1Deployment, ResourceKindAppsV1DaemonSet,
		ResourceKindAppsV1ReplicaSet, ResourceKindAppsV1StatefulSet,
		ResourceKindBatchV1CronJob, ResourceKindBatchV1Job,
		ResourceKindCoreV1Node, ResourceKindCoreV1Namespace,
		ResourceKindCoreV1Pod, ResourceKindCoreV1Container:
		return true
	default:
		return false
	}
}

// Resource is a Kubernetes resource found on the platform.
type Resource interface {
	platform.Node
	// Kind returns the type of resource.
	Kind() ResourceKind
	// Name returns the name of the resource.
	Name() string
	// String returns a textual representation of the resource.
	String() string
}

var _ Resource = (*resource)(nil)

type resource struct {
	id             string
	kind           ResourceKind
	name           string
	labels         platform.Labels
	internalLabels platform.InternalLabels
}

// ID implements platform.Node.
func (r resource) ID() string {
	return r.id
}

// Type implements platform.Node.
func (r resource) Type() string {
	return "kubernetes/" + string(r.kind)
}

// Kind implements Resource.
func (r resource) Kind() ResourceKind {
	return r.kind
}

// Name implements Resource.
func (r resource) Name() string {
	return r.name
}

// Labels implements platform.Node.
func (r resource) Labels() platform.Labels {
	return r.labels
}

// InternalLabels implements platform.Node.
func (r resource) InternalLabels() platform.InternalLabels {
	return r.internalLabels
}

// String implements Resource.
func (r resource) String() string {
	return fmt.Sprintf("%s<%s>", r.kind, r.name)
}

// TagName returns the human-readable name of a tag representing the resource.
// Panics if [ResourceKind.IsSupported] returns false.
func TagName(kind ResourceKind) string {
	switch kind {
	case ResourceKindAppsV1Deployment:
		return "deployment"
	case ResourceKindAppsV1DaemonSet:
		return "daemon set"
	case ResourceKindAppsV1ReplicaSet:
		return "replica set"
	case ResourceKindAppsV1StatefulSet:
		return "stateful set"
	case ResourceKindBatchV1CronJob:
		return "cron job"
	case ResourceKindBatchV1Job:
		return "job"
	case ResourceKindCoreV1Node:
		return "node"
	case ResourceKindCoreV1Namespace:
		return "namespace"
	case ResourceKindCoreV1Pod:
		return "pod"
	case ResourceKindCoreV1Container:
		return "container"
	default:
		// Panic as missing entries would be a programming issue, not runtime
		// bug
		panic(fmt.Sprintf("kubernetes: tag for resource kind <%s> not implemented", kind))
	}
}
