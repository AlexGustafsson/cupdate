package docker

import (
	"fmt"

	"github.com/AlexGustafsson/cupdate/internal/platform"
)

type ResourceKind string

const (
	ResourceKindContainer      = "container"
	ResourceKindSwarmTask      = "swarm/task"
	ResourceKindSwarmService   = "swarm/service"
	ResourceKindSwarmNamespace = "swarm/namespace"
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
	return "docker/" + string(r.kind)
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

func TagName(kind ResourceKind) string {
	switch kind {
	case ResourceKindContainer:
		return "container"
	case ResourceKindSwarmTask:
		return "task"
	case ResourceKindSwarmService:
		return "service"
	case ResourceKindSwarmNamespace:
		return "namespace"
	default:
		// Panic as missing entries would be a programming issue, not runtime
		// bug
		panic(fmt.Sprintf("docker: tag for resource kind <%s> not implemented", kind))
	}
}
