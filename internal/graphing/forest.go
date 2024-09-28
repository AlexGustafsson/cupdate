package graphing

import (
	"strings"
)

type Node interface {
	// ID uniquely identifies this node. Any node with the same ID is expected to
	// contain the exact same data.
	ID() string
}

// Forest is a directed, unweighted graph describing how an image is used.
type Forest[T Node] struct {
	// edges holds a set of adjacent node ids, mapped by the node's id. The bool
	// describes whether or not a is the parent of b.
	//   - true: a->b
	//   - false: a<-b
	edges map[string]map[string]bool
	// nodes holds a set of nodes, mapped by their ids.
	nodes map[string]T
}

func NewForest[T Node]() *Forest[T] {
	return &Forest[T]{
		edges: make(map[string]map[string]bool),
		nodes: make(map[string]T),
	}
}

// InsertTree inserts nodes of a tree, ordered root first, leaf last.
func (f *Forest[T]) InsertTree(nodes ...T) {
	for i := 0; i < len(nodes); i++ {
		f.insertNode(nodes[i])
		if i > 0 {
			f.insertEdge(nodes[i-1].ID(), nodes[i].ID(), true)
			f.insertEdge(nodes[i].ID(), nodes[i-1].ID(), false)
		}
	}
}

// InsertForest merges other into f.
func (f *Forest[T]) InsertForest(other *Forest[T]) {
	for from, edges := range other.edges {
		for to, direction := range edges {
			f.insertEdge(from, to, direction)
		}
	}

	for _, node := range other.nodes {
		f.insertNode(node)
	}
}

// insertNode inserts the node into the forest.
func (f *Forest[T]) insertNode(n T) {
	f.nodes[n.ID()] = n
}

// insertEdge inserts an edge from a to b with the specified direction.
func (f *Forest[T]) insertEdge(a string, b string, direction bool) {
	if _, ok := f.edges[a]; !ok {
		f.edges[a] = make(map[string]bool)
	}
	f.edges[a][b] = direction
}

func (f *Forest[T]) Roots() []T {
	roots := make([]T, 0)
	for nodeID, node := range f.nodes {
		parents := 0
		for _, isParent := range f.edges[nodeID] {
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

func (f *Forest[T]) String() string {
	var result strings.Builder

	roots := f.Roots()
	for i := 0; i < len(roots); i++ {
		result.WriteString(f.describeFromRoot(roots[i].ID()))
		if i < len(roots)-1 {
			result.WriteByte('\n')
		}
	}

	return result.String()
}

func (f *Forest[T]) describeFromRoot(rootID string) string {
	var result strings.Builder

	paths := f.traverse(rootID)
	for i := 0; i < len(paths); i++ {
		result.WriteString(strings.Join(paths[i], "->"))
		if i < len(paths)-1 {
			result.WriteByte('\n')
		}
	}

	return result.String()
}

func (f *Forest[T]) childrenIDs(nodeID string) []string {
	childrenIDs := make([]string, 0)
	for adjacentID, isParent := range f.edges[nodeID] {
		if isParent {
			childrenIDs = append(childrenIDs, adjacentID)
		}
	}
	return childrenIDs
}

func (f *Forest[T]) traverse(rootID string) [][]string {
	children := f.childrenIDs(rootID)
	if len(children) == 0 {
		return [][]string{{rootID}}
	}

	paths := make([][]string, 0)
	for _, child := range children {
		for _, path := range f.traverse(child) {
			path = append([]string{rootID}, path...)
			paths = append(paths, path)
		}
	}

	return paths
}
