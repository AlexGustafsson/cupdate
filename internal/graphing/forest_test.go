package graphing

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testNode struct {
	id string
}

func (n testNode) ID() string {
	return n.id
}

func TestForestLeaves(t *testing.T) {
	f := NewForest[testNode]()

	f.InsertTree(
		testNode{id: "image a"},
		testNode{id: "container a"},
		testNode{id: "pod a"},
		testNode{id: "deployment a"},
		testNode{id: "namespace a"},
	)

	f.InsertTree(
		testNode{id: "image a"},
		testNode{id: "container b"},
		testNode{id: "pod a"},
		testNode{id: "deployment a"},
		testNode{id: "namespace a"},
	)

	f.InsertTree(
		testNode{id: "image a"},
		testNode{id: "container c"},
		testNode{id: "pod b"},
		testNode{id: "deployment b"},
		testNode{id: "namespace b"},
	)

	f.InsertTree(
		testNode{id: "image b"},
		testNode{id: "container d"},
		testNode{id: "pod c"},
		testNode{id: "deployment c"},
		testNode{id: "namespace c"},
	)

	f.InsertTree(
		testNode{id: "image c"},
	)

	expectedRoots := []Node{
		testNode{id: "image a"},
		testNode{id: "image b"},
		testNode{id: "image c"},
	}

	assert.ElementsMatch(t, expectedRoots, f.Roots())
}

func TestForestString(t *testing.T) {
	f := NewForest[testNode]()

	f.InsertTree(
		testNode{id: "image a"},
		testNode{id: "container a"},
		testNode{id: "pod a"},
		testNode{id: "deployment a"},
		testNode{id: "namespace a"},
	)

	f.InsertTree(
		testNode{id: "image a"},
		testNode{id: "container b"},
		testNode{id: "pod a"},
		testNode{id: "deployment a"},
		testNode{id: "namespace a"},
	)

	f.InsertTree(
		testNode{id: "image a"},
		testNode{id: "container c"},
		testNode{id: "pod b"},
		testNode{id: "deployment b"},
		testNode{id: "namespace b"},
	)

	f.InsertTree(
		testNode{id: "image b"},
		testNode{id: "container d"},
		testNode{id: "pod c"},
		testNode{id: "deployment c"},
		testNode{id: "namespace c"},
	)

	f.InsertTree(
		testNode{id: "image c"},
	)

	expectedString := `image a->container a->pod a->deployment a->namespace a
image a->container b->pod a->deployment a->namespace a
image a->container c->pod b->deployment b->namespace b
image b->container d->pod c->deployment c->namespace c
image c`

	actualString := f.String()

	// Ignore order when matching
	expected := strings.Split(expectedString, "\n")
	actual := strings.Split(actualString, "\n")

	assert.ElementsMatch(t, expected, actual)
}
