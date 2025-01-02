package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/graph"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/platform"
)

var _ platform.Grapher = (*Platform)(nil)

type Platform struct {
	client *http.Client

	includeAllContainers bool
}

type Options struct {
	IncludeAllContainers bool
}

func NewPlatform(ctx context.Context, host string, options *Options) (*Platform, error) {
	if options == nil {
		options = &Options{}
	}

	if !strings.HasPrefix(host, "unix://") {
		return nil, fmt.Errorf("unexpected docker host - expected a unix socket")
	}
	path := strings.TrimPrefix(host, "unix://")

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", path)
			},
		},
	}

	p := &Platform{
		client: client,

		includeAllContainers: options.IncludeAllContainers,
	}

	// Make sure that we can connect to the socket.
	// For now, we probably support most API versions - no need to limit the use
	// or pin to specific API versions using docker's versioned path prefix
	_, _, err := p.GetVersion(ctx)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Platform) GetVersion(ctx context.Context) (string, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://unix/version", nil)
	if err != nil {
		return "", "", err
	}

	res, err := p.client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var body struct {
		APIVersion    string `json:"ApiVersion"`
		MinAPIVersion string `json:"MinAPIVersion"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return "", "", err
	}

	return body.APIVersion, body.MinAPIVersion, nil
}

type GetContainersOptions struct {
	All     bool
	Filters map[string][]string
}

func (p *Platform) GetContainers(ctx context.Context, options *GetContainersOptions) ([]Container, error) {
	query := make(url.Values)
	if options != nil && options.All {
		query.Set("all", "true")
	}
	if options != nil && options.Filters != nil {
		filters, err := json.Marshal(options.Filters)
		if err != nil {
			return nil, err
		}

		query.Set("filters", string(filters))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://unix/containers/json?"+query.Encode(), nil)
	if err != nil {
		return nil, err
	}

	res, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var result []Container
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *Platform) GetImage(ctx context.Context, nameOrID string) (*Image, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://unix/images/"+url.PathEscape(nameOrID)+"/json", nil)
	if err != nil {
		return nil, err
	}

	res, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var result Image
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Graph implements platform.Platform.
// SEE: https://docs.docker.com/reference/api/engine/version/v1.47/
func (p *Platform) Graph(ctx context.Context) (*graph.Graph[platform.Node], error) {
	options := &GetContainersOptions{
		All: p.includeAllContainers,
	}
	containers, err := p.GetContainers(ctx, options)
	if err != nil {
		return nil, err
	}

	images := make(map[string]*Image)
	for _, container := range containers {
		_, ok := images[container.ImageID]
		if !ok {
			image, err := p.GetImage(ctx, container.ImageID)
			if err != nil {
				return nil, err
			}

			images[image.ID] = image
		}
	}

	graph := platform.NewGraph()

	for _, container := range containers {
		image := images[container.ImageID]

		reference, ok := image.Reference()
		if !ok {
			slog.WarnContext(ctx, "Failed to identify a reference for image", slog.String("id", image.ID))
			continue
		}

		graph.InsertTree(
			platform.ImageNode{
				Reference: reference,
			},
			resource{
				kind: ResourceKindContainer,
				id:   fmt.Sprintf("docker/containers/%s", container.ID),
				name: container.Name(),
			},
		)
	}

	return graph, nil
}

type Container struct {
	ID      string `json:"Id"`
	Names   []string
	Image   string
	ImageID string

	// ... other ignored fields
}

func (c Container) Name() string {
	for _, name := range c.Names {
		// For whatever reason, names are prefixed with "/"
		name = strings.TrimPrefix(name, "/")
		if name != "" {
			return name
		}
	}

	return c.ID
}

type Image struct {
	ID          string `json:"Id"`
	RepoTags    []string
	RepoDigests []string

	// ... other ignored fields
}

func (i Image) Reference() (oci.Reference, bool) {
	for _, tagged := range i.RepoTags {
		ref, err := oci.ParseReference(tagged)
		if err == nil {
			return ref, true
		}
	}
	for _, digested := range i.RepoDigests {
		ref, err := oci.ParseReference(digested)
		if err == nil {
			return ref, true
		}
	}

	return oci.Reference{}, false
}
