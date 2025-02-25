package graph

import (
	"strings"
	"sync"
)

type Node interface {
	// ID uniquely identifies this node. Any node with the same ID is expected to
	// contain the exact same data.
	ID() string
}

type stringer interface {
	String() string
}

// Graph is a directed, cyclic, unweighted graph.
type Graph[T Node] struct {
	// edges holds a set of adjacent node ids, mapped by the node's id. The bool
	// describes whether or not a is the parent of b.
	//   - true: a->b
	//   - false: a<-b
	edges map[string]map[string]bool
	// nodes holds a set of nodes, mapped by their ids.
	nodes map[string]T
	mutex sync.Mutex
}

func New[T Node]() *Graph[T] {
	return &Graph[T]{
		edges: make(map[string]map[string]bool),
		nodes: make(map[string]T),
	}
}

// InsertTree inserts nodes of a tree, ordered root first, leaf last.
func (g *Graph[T]) InsertTree(nodes ...T) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	for i := 0; i < len(nodes); i++ {
		g.insertNode(nodes[i])
		if i > 0 {
			g.insertEdge(nodes[i-1].ID(), nodes[i].ID(), true)
			g.insertEdge(nodes[i].ID(), nodes[i-1].ID(), false)
		}
	}
}

// InsertGraph merges other into g.
func (g *Graph[T]) InsertGraph(other *Graph[T]) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	for from, edges := range other.edges {
		for to, direction := range edges {
			g.insertEdge(from, to, direction)
		}
	}

	for _, node := range other.nodes {
		g.insertNode(node)
	}
}

// insertNode inserts the node into the graph.
func (g *Graph[T]) insertNode(n T) {
	g.nodes[n.ID()] = n
	if _, ok := g.edges[n.ID()]; !ok {
		g.edges[n.ID()] = make(map[string]bool)
	}
}

// insertEdge inserts an edge from a to b with the specified direction.
func (g *Graph[T]) insertEdge(a string, b string, direction bool) {
	if _, ok := g.edges[a]; !ok {
		g.edges[a] = make(map[string]bool)
	}
	g.edges[a][b] = direction
}

func (g *Graph[T]) Roots() []T {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	return g.roots()
}

func (g *Graph[T]) roots() []T {
	roots := make([]T, 0)
	for nodeID, node := range g.nodes {
		parents := 0
		for _, isParent := range g.edges[nodeID] {
			if !isParent {
				parents++
			}
		}
		if parents == 0 {
			roots = append(roots, node)
		}
	}

	return roots
}

// String returns a string describing the graph.
func (g *Graph[T]) String() string {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	var result strings.Builder

	roots := g.roots()
	for i := 0; i < len(roots); i++ {
		result.WriteString(g.describeFromRoot(roots[i].ID()))
		if i < len(roots)-1 {
			result.WriteByte('\n')
		}
	}

	return result.String()
}

func (g *Graph[T]) describeFromRoot(rootID string) string {
	var result strings.Builder

	paths := g.traverse(rootID)
	for i := 0; i < len(paths); i++ {
		labels := make([]string, 0)
		for _, nodeID := range paths[i] {
			node := g.nodes[nodeID]
			if named, ok := any(node).(stringer); ok {
				labels = append(labels, named.String())
			} else {
				labels = append(labels, node.ID())
			}
		}

		result.WriteString(strings.Join(labels, "->"))
		if i < len(paths)-1 {
			result.WriteByte('\n')
		}
	}

	return result.String()
}

func (g *Graph[T]) children(nodeID string) []string {
	childrenIDs := make([]string, 0)
	for adjacentID, isParent := range g.edges[nodeID] {
		if isParent {
			childrenIDs = append(childrenIDs, adjacentID)
		}
	}
	return childrenIDs
}

// Subgraph returns a graph including all nodes and edges reachable from the
// specified root.
func (g *Graph[T]) Subgraph(rootID string) *Graph[T] {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	subgraph := New[T]()

	visited := make(map[string]struct{})
	queue := []string{rootID}
	for len(queue) > 0 {
		root := queue[0]
		queue = queue[1:]
		subgraph.insertNode(g.nodes[root])

		children := g.children(root)
		for _, child := range children {
			if _, ok := visited[child]; !ok {
				queue = append(queue, child)
				subgraph.insertEdge(root, child, true)
				subgraph.insertEdge(child, root, false)
			}
		}
	}

	return subgraph
}

// Edges returns all of the graph's edges.
func (g *Graph[T]) Edges() map[string]map[string]bool {
	return g.edges
}

// Nodes returns all of the graph's nodes.
func (g *Graph[T]) Nodes() []T {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	nodes := make([]T, 0)
	for _, node := range g.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

func (g *Graph[T]) traverse(rootID string) [][]string {
	children := g.children(rootID)
	if len(children) == 0 {
		return [][]string{{rootID}}
	}

	paths := make([][]string, 0)
	for _, child := range children {
		for _, path := range g.traverse(child) {
			path = append([]string{rootID}, path...)
			paths = append(paths, path)
		}
	}

	return paths
}

// DeleteFunc deletes all nodes for which del returns true.
// If removing a node causes the graph to change, the deletion is recursive,
// removing all nodes that are no longer referenced.
func (g *Graph[T]) DeleteFunc(del func(T) bool) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// Delete all nodes (and their edges) matching by the func
	for id, node := range g.nodes {
		if del(node) {
			g.deleteNodeRecursive(id)
		}
	}
}

func (g *Graph[T]) deleteNodeRecursive(id string) {
	// Delete the node itself
	delete(g.nodes, id)

	for adjacent, removedParentFromChild := range g.edges[id] {
		// Delete all relationships the node had
		delete(g.edges[adjacent], id)

		// Recursively delete adjacent nodes if their relationships changed
		if removedParentFromChild {
			hasAdditionalParent := false
			for _, isParent := range g.edges[adjacent] {
				if !isParent {
					hasAdditionalParent = true
					break
				}
			}
			if !hasAdditionalParent {
				g.deleteNodeRecursive(adjacent)
			}
		} else {
			hasAdditionalChild := false
			for _, isParent := range g.edges[adjacent] {
				if isParent {
					hasAdditionalChild = true
					break
				}
			}
			if !hasAdditionalChild {
				g.deleteNodeRecursive(adjacent)
			}
		}
	}

	// Delete the map of relationships
	delete(g.edges, id)
}
