package docker

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log/slog"
	"maps"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/graph"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/platform"
)

var _ platform.Grapher = (*Platform)(nil)

// Platform implements graphing for the Docker platform.
type Platform struct {
	client   *http.Client
	basePath string

	includeAllContainers bool
	version              Version
	dockerURI            string
}

type Options struct {
	// IncludeAllContainers will graph all containers, no matter their state.
	// Defaults to false - only include running containers.
	IncludeAllContainers bool
	// TLSClientConfig sets the TLSClientConfig of the underlying client.
	// Only used if scheme is tcp or https.
	TLSClientConfig *tls.Config
}

// NewPlatform initializes a new [Platform].
//
//   - dockerURI is the URI to the docker socket. Such as unix://docker.sock,
//     tcp://127.0.0.1:8080, http://127.0.0.1:8080 or https://127.0.0.1:8080.
func NewPlatform(ctx context.Context, dockerURI string, options *Options) (*Platform, error) {
	if options == nil {
		options = &Options{}
	}

	scheme, _, ok := strings.Cut(dockerURI, "://")
	if !ok {
		return nil, fmt.Errorf("invalid docker URI")
	}

	basePath := ""
	transport := httputil.NewTransport()
	switch scheme {
	case "unix":
		host := strings.TrimPrefix(dockerURI, "unix://")
		if _, err := os.Stat(host); err != nil {
			return nil, err
		}

		basePath = "http://_"
		transport.DialContext = func(ctx context.Context, _, _ string) (net.Conn, error) {
			return (&net.Dialer{
				Timeout: 5 * time.Second,
			}).DialContext(ctx, "unix", host)
		}
	case "tcp", "http", "https":
		url, err := url.Parse(dockerURI)
		if err != nil {
			return nil, fmt.Errorf("invalid docker URI: %w", err)
		}

		if url.Scheme == "tcp" {
			if options.TLSClientConfig == nil {
				url.Scheme = "http"
			} else {
				url.Scheme = "https"
			}
		}

		basePath = url.String()

		if url.Scheme == "tcp" || url.Scheme == "https" {
			transport.TLSClientConfig = options.TLSClientConfig
		}
	default:
		return nil, fmt.Errorf("unsupported docker URI: %s", dockerURI)
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	p := &Platform{
		client:   client,
		basePath: basePath,

		includeAllContainers: options.IncludeAllContainers,
	}

	// Make sure that we can connect to the host.
	// For now, we probably support most API versions - no need to limit the use
	// or pin to specific API versions using docker's versioned path prefix
	version, err := p.getVersion(ctx)
	if err != nil {
		return nil, err
	}

	p.version = *version
	p.dockerURI = dockerURI

	return p, nil
}

// Version returns the version of the platform.
func (p *Platform) Version() Version {
	return p.version
}

// getVersion returns the api version and minimum supported api version of the
// Docker runtime.
func (p *Platform) getVersion(ctx context.Context) (*Version, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.basePath+"/version", nil)
	if err != nil {
		return nil, err
	}

	res, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err := httputil.AssertStatusCode(res, http.StatusOK); err != nil {
		return nil, err
	}

	var version *Version
	if err := json.NewDecoder(res.Body).Decode(&version); err != nil {
		return nil, err
	}

	return version, nil
}

type GetContainersOptions struct {
	// All maps to the all query parameter of the containers API, returning all
	// containers no matter their state.
	All bool
	// Filters maps to the filters query parameter of the containers API,
	// filtering containers to include.
	Filters map[string][]string
}

// GetContainers retrieves container information from the Docker runtime.
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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.basePath+"/containers/json?"+query.Encode(), nil)
	if err != nil {
		return nil, err
	}

	res, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err := httputil.AssertStatusCode(res, http.StatusOK); err != nil {
		return nil, err
	}

	var result []Container
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetImage retrieves container image information from the Docker runtime by
// name or id. Returns an error if the image does not exist.
func (p *Platform) GetImage(ctx context.Context, nameOrID string) (*Image, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.basePath+"/images/"+url.PathEscape(nameOrID)+"/json", nil)
	if err != nil {
		return nil, err
	}

	res, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err := httputil.AssertStatusCode(res, http.StatusOK); err != nil {
		return nil, err
	}

	var result Image
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Graph implements platform.Platform.
//
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
		// When running Podman pods, containers can be reported without an image.
		// Ignore such containers as the Podman image API will throw internal server
		// errors when requests are made using these empty ids.
		if container.ImageID == "" || container.ImageID == "sha256:" {
			slog.Warn("Ignoring container with invalid image id", slog.String("containerId", container.ID))
			continue
		}

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
		var repoDigests []string
		if image, ok := images[container.ImageID]; ok {
			repoDigests = image.RepoDigests
		}

		// Identify the image reference used, containing a digest, if available.
		// At times, Docker compose will reference images in arguably weird ways,
		// potentially just by their digest. This is generally not a valid OCI
		// reference and we can't infer / build the reference ourselves. Here, as
		// the fault is obvious and cannot collide with the arguably correct
		// reference "sha256:1234...", let's make sure no "sha256" tag makes it
		// through
		ref, err := getImageReference(container.Image, repoDigests)
		if err != nil || ref.Domain == "docker.io" && ref.Path == "sha256" {
			slog.ErrorContext(ctx, "Failed to identify a valid image reference for container", slog.String("container", container.ID))
			continue
		}

		tree := []platform.Node{
			platform.ImageNode{
				Reference: ref,
			},
			resource{
				kind:   ResourceKindContainer,
				id:     fmt.Sprintf("docker/containers/%s", container.ID),
				name:   container.Name(),
				labels: maps.Clone(container.Labels),
			},
		}

		// Add graph nodes for Docker Swarm and Compose, if available
		if container.Labels != nil {
			if taskID, ok := container.Labels["com.docker.swarm.task.id"]; ok {
				taskName, ok := container.Labels["com.docker.swarm.task.name"]
				if !ok {
					taskName = taskID
				}

				tree = append(tree, resource{
					kind: ResourceKindSwarmTask,
					id:   fmt.Sprintf("docker/swarm/task/%s", taskID),
					name: taskName,
				})
			}

			if serviceID, ok := container.Labels["com.docker.swarm.service.id"]; ok {
				serviceName, ok := container.Labels["com.docker.swarm.service.name"]
				if !ok {
					serviceName = serviceID
				}

				tree = append(tree, resource{
					kind: ResourceKindSwarmService,
					id:   fmt.Sprintf("docker/swarm/service/%s", serviceID),
					name: serviceName,
				})
			} else if service, ok := container.Labels["com.docker.compose.service"]; ok {
				tree = append(tree, resource{
					kind: ResourceKindComposeService,
					id:   fmt.Sprintf("docker/compose/service/%s", service),
					name: service,
				})
			}

			if namespace, ok := container.Labels["com.docker.stack.namespace"]; ok {
				tree = append(tree, resource{
					kind: ResourceKindSwarmNamespace,
					id:   fmt.Sprintf("docker/swarm/namespace/%s", namespace),
					name: namespace,
				})
			} else if project, ok := container.Labels["com.docker.compose.project"]; ok {
				tree = append(tree, resource{
					kind: ResourceKindComposeProject,
					id:   fmt.Sprintf("docker/compose/project/%s", project),
					name: project,
				})
			}
		}

		// Add a graph node for the host
		tree = append(tree, resource{
			kind:   ResourceKindHost,
			id:     fmt.Sprintf("docker/host/%s", p.dockerURI),
			name:   p.dockerURI, // TODO: Use hostname?
			labels: nil,
			internalLabels: platform.InternalLabels{
				platform.InternalLabelHostArchitecture: p.version.Architecture,
				platform.InternalLabelOperatingSystem:  p.version.OS,
			},
		})

		graph.InsertTree(tree...)
	}

	return graph, nil
}

// Container is a container as defined by the Docker runtime API.
type Container struct {
	ID      string `json:"Id"`
	Names   []string
	Image   string
	ImageID string
	Labels  map[string]string

	// ... other ignored fields
}

// Name returns the name of the container, or its ID if no name is found.
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

// Image is an image as defined by the Docker runtime API.
type Image struct {
	ID          string `json:"Id"`
	RepoTags    []string
	RepoDigests []string

	// ... other ignored fields
}
