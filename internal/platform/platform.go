package platform

import (
	"context"

	"github.com/distribution/reference"
)

type Platform interface {
	// Images returns all unique images in use or referenced within the platform
	// as well as a Graph describing in what ways the images are used.
	Images(context.Context) ([]reference.Named, Graph, error)
}

type Origin interface {
	Kind() string
}

type Graph map[string][]Origin

func (g Graph) AddOrigin(reference reference.Reference, origin Origin) {
	key := reference.String()

	origins := g[key]
	if origins == nil {
		origins = []Origin{origin}
	} else {
		origins = append(origins, origin)
	}

	g[key] = origins
}

func (g Graph) Origins(reference reference.Reference) []Origin {
	key := reference.String()

	return g[key]
}
