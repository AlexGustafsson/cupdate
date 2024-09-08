package ghcr

import (
	"context"
	"fmt"
	"net/http"

	"github.com/AlexGustafsson/k8s-image-feed/internal/registry"
)

var _ registry.Registry = (*Registry)(nil)

type Registry struct {
	Client http.Client
}

// Get implements registry.Registry.
func (r *Registry) Get(ctx context.Context, name string, version string) (*registry.Image, error) {
	return nil, fmt.Errorf("not implemented")
}

// GetLatestVersion implements registry.Registry.
func (r *Registry) GetLatestVersion(ctx context.Context, name string, currentTag string) (*registry.Image, error) {
	// TODO: Use filter?
	// ...?filters%5Bversion_type%5D=tagged
	// res, err := r.Client.Get(fmt.Sprintf("https://github.com/%s/pkgs/container/trivy/versions", name))
	// if err != nil {
	// 	return nil, err
	// }

	return nil, fmt.Errorf("not implemented")
}
