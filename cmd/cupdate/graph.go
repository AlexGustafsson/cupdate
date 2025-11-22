package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/AlexGustafsson/cupdate/internal/events"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/platform"
	"github.com/AlexGustafsson/cupdate/internal/platform/docker"
	"github.com/AlexGustafsson/cupdate/internal/platform/kubernetes"
	"github.com/AlexGustafsson/cupdate/internal/store"
	"github.com/AlexGustafsson/cupdate/internal/worker"
)

func HandleGraphs(ctx context.Context, targetPlatform platform.ContinuousGrapher, platformEvents *events.Hub[models.PlatformEvent], writeStore *store.Store, processQueue *worker.Queue[oci.Reference]) {
	defer slog.Debug("Closed graphing")
	for {
		select {
		case <-ctx.Done():
			return
		case graph, ok := <-targetPlatform.Graphs():
			if !ok {
				return
			}

			slog.Debug("Got updated platform graph")

			// Delete ignored images / trees
			graph.DeleteFunc(func(n platform.Node) bool {
				return n.Labels().Ignore()
			})

			roots := graph.Roots()

			totalInserted := 0
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
				slog.Debug("Inserting raw image", slog.String("reference", rawImage.Reference))
				inserted, err := writeStore.InsertRawImage(ctx, rawImage)
				if err != nil {
					slog.Error("Failed to insert raw image", slog.Any("error", err))
					continue
				}

				// Try to schedule the image for processing
				if inserted {
					slog.Debug("Raw image inserted for first time - scheduling for processing")
					totalInserted++
					processQueue.PushBack(imageNode.Reference)
				}
			}

			allReferences := make([]string, 0)
			for _, root := range roots {
				imageNode := root.(platform.ImageNode)
				allReferences = append(allReferences, imageNode.Reference.String())
			}

			slog.Debug("Cleaning up removed images")
			totalRemoved, err := writeStore.DeleteNonPresent(ctx, allReferences)
			if err == nil {
				slog.Debug("Cleaned up removed images successfully", slog.Int64("removed", totalRemoved))
			} else {
				slog.Error("Failed to clean up removed images", slog.Any("error", err))
			}

			if totalInserted > 0 || totalRemoved > 0 {
				platformEvents.Broadcast(ctx, models.PlatformEvent{Type: models.EventTypeGraphUpdated})
			}
		}
	}
}
