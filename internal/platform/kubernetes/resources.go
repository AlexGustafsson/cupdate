package kubernetes

import (
	"fmt"

	"github.com/AlexGustafsson/cupdate/internal/platform"
)

type ResourceKind string

const (
	ResourceKindAppsV1Deployment  = "apps/v1/deployment"
	ResourceKindAppsV1DaemonSet   = "apps/v1/daemonset"
	ResourceKindAppsV1ReplicaSet  = "apps/v1/replicaset"
	ResourceKindAppsV1StatefulSet = "apps/v1/statefulset"
	ResourceKindBatchV1CronJob    = "batch/v1/cronjob"
	ResourceKindBatchV1Job        = "batch/v1/job"
	ResourceKindCoreV1Pod         = "core/v1/pod"
	ResourceKindCoreV1Container   = "core/v1/container"
)

type Resource interface {
	platform.Node
	Kind() ResourceKind
	Name() string
	String() string
}

type resource struct {
	id   string
	kind ResourceKind
	name string
}

func (r resource) ID() string {
	return r.id
}

func (r resource) Type() string {
	return "kubernetes/" + string(r.kind)
}

func (r resource) Kind() ResourceKind {
	return r.kind
}

func (r resource) Name() string {
	return r.name
}

func (r resource) String() string {
	return fmt.Sprintf("%s<%s>", r.kind, r.name)
}
