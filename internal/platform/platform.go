package platform

import (
	"context"

	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
)

type Platform interface {
	// Images returns all unique images in use or referenced within the platform
	// as well as a Graph describing in what ways the images are used.
	Images(context.Context) ([]oci.Reference, Graph, error)
}

type Origin interface {
	Kind() string
}

type Graph map[string][]Origin

func (g Graph) AddOrigin(reference oci.Reference, origin Origin) {
	key := reference.String()

	origins := g[key]
	if origins == nil {
		origins = []Origin{origin}
	} else {
		origins = append(origins, origin)
	}

	g[key] = origins
}

func (g Graph) Origins(reference oci.Reference) []Origin {
	key := reference.String()

	return g[key]
}
