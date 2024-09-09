package k8s

import (
	"time"

	"github.com/AlexGustafsson/cupdate/internal/source"
)

var _ source.Origin = (*Origin)(nil)

type ResourceKind string

const (
	ResourceKindAppsV1Deployment  = "apps/v1/Deployment"
	ResourceKindAppsV1DaemonSet   = "apps/v1/DaemonSet"
	ResourceKindAppsV1ReplicaSet  = "apps/v1/ReplicaSet"
	ResourceKindAppsV1StatefulSet = "apps/v1/StatefulSet"
	ResourceKindBatchV1CronJob    = "batch/v1/CronJob"
	ResourceKindBatchV1Job        = "batch/v1/Job"
	ResourceKindCoreV1Pod         = "core/v1/Pod"
)

type Parent struct {
	ResourceKind ResourceKind
	Namespace    string
	Name         string
	Parent       *Parent
}

type Pod struct {
	// Name might be empty is the resource it was found in did not specify a name.
	// Can happen if it was discovered in a deployment and the author did not
	// specify an explicit name, giving it a name at runtime.
	// Always defined if IsTemplate is false.
	Name      string
	Namespace string
	Created   time.Time
	// IsTemplate is true if the pod was found in a template, rather than an
	// actual running pod.
	IsTemplate bool
	Parent     *Parent
}

type Container struct {
	Name      string
	Namespace string
	Pod       *Pod
}

type Origin struct {
	Container *Container
}

func (o *Origin) Kind() string {
	return "k8s"
}
