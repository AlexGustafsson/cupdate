package worker

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/platform"
	"github.com/AlexGustafsson/cupdate/internal/platform/docker"
	"github.com/AlexGustafsson/cupdate/internal/platform/kubernetes"
	"github.com/AlexGustafsson/cupdate/internal/store"
)

type GraphWorker struct {
	grapher platform.ContinuousGrapher
	store   *store.Store
	queue   *Queue[oci.Reference]
}

func NewGraphWorker(grapher platform.ContinuousGrapher, store *store.Store, queue *Queue[oci.Reference]) *GraphWorker {
	return &GraphWorker{
		grapher: grapher,
		store:   store,
		queue:   queue,
	}
}

func (w *GraphWorker) GraphContinuously(ctx context.Context) error {
	graphs, err := w.grapher.GraphContinuously(ctx)
	if err != nil {
		slog.Error("Failed to start graphing platform", slog.Any("error", err))
		return err
	}

	for graph := range graphs {
		slog.Debug("Got updated platform graph")

		// Delete ignored images / trees
		graph.DeleteFunc(func(n platform.Node) bool {
			return n.Labels().Ignore()
		})

		roots := graph.Roots()

		for _, root := range roots {
			imageNode := root.(platform.ImageNode)

			subgraph := graph.Subgraph(root.ID())

			edges := subgraph.Edges()
			nodes := subgraph.Nodes()

			var namespaceNode *platform.Node

			mappedNodes := make(map[string]models.GraphNode)
			for _, node := range nodes {
				switch n := node.(type) {
				case kubernetes.Resource:
					mappedNodes[node.ID()] = models.GraphNode{
						Domain:         "kubernetes",
						Type:           string(n.Kind()),
						Name:           n.Name(),
						Labels:         n.Labels().RemoveUnsupported(),
						InternalLabels: n.InternalLabels(),
					}
					if node.Type() == "kubernetes/"+kubernetes.ResourceKindCoreV1Namespace {
						namespaceNode = &node
					}
				case docker.Resource:
					mappedNodes[node.ID()] = models.GraphNode{
						Domain:         "docker",
						Type:           string(n.Kind()),
						Name:           n.Name(),
						Labels:         n.Labels().RemoveUnsupported(),
						InternalLabels: n.InternalLabels(),
					}
					if node.Type() == "docker/"+docker.ResourceKindSwarmNamespace || node.Type() == "docker/"+docker.ResourceKindComposeProject {
						namespaceNode = &node
					}
				case platform.ImageNode:
					// This node is added later on
				default:
					panic(fmt.Sprintf("mapping unimplemented node type: %s", node.Type()))
				}
			}

			// Resolve labels for the image node. The nearest label takes precedence
			resolvedLabels := make(map[string]string)
			queue := []string{root.ID()}
			for len(queue) > 0 {
				id := queue[0]
				queue = queue[1:]

				for k, v := range mappedNodes[id].Labels {
					_, ok := resolvedLabels[k]
					if !ok {
						resolvedLabels[k] = v
					}
				}

				for adjacent, isChild := range edges[id] {
					if isChild {
						queue = append(queue, adjacent)
					}
				}
			}
			mappedNodes[root.ID()] = models.GraphNode{
				Domain:         "oci",
				Type:           "image",
				Name:           imageNode.Reference.String(),
				Labels:         resolvedLabels,
				InternalLabels: nil,
			}

			tags := []string{}

			// Set tags for resources
			if namespaceNode != nil {
				children := edges[(*namespaceNode).ID()]
				for childID, isParent := range children {
					if isParent {
						continue
					}

					var childNode *platform.Node
					for _, node := range nodes {
						var n = node
						if node.ID() == childID {
							childNode = &n
							break
						}
					}

					if childNode != nil {
						switch resource := (*childNode).(type) {
						case kubernetes.Resource:
							kind := resource.Kind()
							if kind.IsSupported() {
								tags = append(tags, kubernetes.TagName(resource.Kind()))
							}
						case docker.Resource:
							tags = append(tags, docker.TagName(resource.Kind()))
						}
					}
				}
			}

			mappedGraph := models.Graph{
				Edges: edges,
				Nodes: mappedNodes,
			}

			rawImage := &models.RawImage{
				Reference: imageNode.Reference.String(),
				Tags:      tags,
				Graph:     mappedGraph,
			}

			// TODO: Do this inside of the worker as well?
			slog.DebugContext(ctx, "Inserting raw image", slog.String("reference", rawImage.Reference))
			inserted, err := w.store.InsertRawImage(context.TODO(), rawImage)
			if err != nil {
				slog.ErrorContext(ctx, "Failed to insert raw image", slog.Any("error", err))
				continue
			}

			// Try to schedule the image for processing
			if inserted {
				slog.DebugContext(ctx, "Raw image inserted for first time - scheduling for processing")
				w.queue.Push(imageNode.Reference)
			}
		}

		allReferences := make([]string, 0)
		for _, root := range roots {
			imageNode := root.(platform.ImageNode)
			allReferences = append(allReferences, imageNode.Reference.String())
		}

		slog.DebugContext(ctx, "Cleaning up removed images")
		removed, err := w.store.DeleteNonPresent(context.TODO(), allReferences)
		if err == nil {
			slog.DebugContext(ctx, "Cleaned up removed images successfully", slog.Int64("removed", removed))
		} else {
			slog.ErrorContext(ctx, "Failed to clean up removed images", slog.Any("error", err))
		}
	}

	return nil
}
