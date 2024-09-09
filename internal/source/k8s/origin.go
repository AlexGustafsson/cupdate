package k8s

import (
	"time"

	"github.com/AlexGustafsson/cupdate/internal/source"
)

var _ source.Origin = (*Origin)(nil)

type Origin struct {
	ResourceKind  string
	Namespace     string
	Name          string
	Created       time.Time
	ContainerName string
	Owners        []Origin
	Parents       []Parent
}

func (o *Origin) Kind() string {
	return "k8s"
}

type Parent struct {
	ResourceKind string
	Name         string
}
