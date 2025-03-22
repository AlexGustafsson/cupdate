package docker

import (
	"fmt"

	"github.com/AlexGustafsson/cupdate/internal/platform"
)

// ResourceKind defines the types of resources exposed by the platform.
type ResourceKind string

const (
	ResourceKindHost           = "host"
	ResourceKindContainer      = "container"
	ResourceKindSwarmTask      = "swarm/task"
	ResourceKindSwarmService   = "swarm/service"
	ResourceKindSwarmNamespace = "swarm/namespace"
	ResourceKindComposeProject = "compose/project"
	ResourceKindComposeService = "compose/service"
)

// IsSupported returns whether or not the resource is supported.
func (r ResourceKind) IsSupported() bool {
	switch r {
	case ResourceKindHost, ResourceKindContainer, ResourceKindSwarmTask,
		ResourceKindSwarmService, ResourceKindSwarmNamespace,
		ResourceKindComposeProject, ResourceKindComposeService:
		return true
	default:
		return false
	}
}

// Resource is a Docker resource found on the platform.
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
	return "docker/" + string(r.kind)
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
	case ResourceKindHost:
		return "host"
	case ResourceKindContainer:
		return "container"
	case ResourceKindSwarmTask:
		return "task"
	case ResourceKindSwarmService:
		return "service"
	case ResourceKindSwarmNamespace:
		return "namespace"
	case ResourceKindComposeProject:
		return "project"
	case ResourceKindComposeService:
		return "service"
	default:
		// Panic as missing entries would be a programming issue, not runtime
		// bug
		panic(fmt.Sprintf("docker: tag for resource kind <%s> not implemented", kind))
	}
}
