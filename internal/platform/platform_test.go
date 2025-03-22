package platform

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var _ Grapher = (*MockGrapher)(nil)

type MockGrapher struct {
	mock.Mock
}

// Graph implements Grapher.
func (m *MockGrapher) Graph(ctx context.Context) (Graph, error) {
	args := m.Called(ctx)
	return args.Get(0).(Graph), args.Error(1)
}

var _ Node = (*TestNode)(nil)

type TestNode struct {
	id       string
	nodeType string
}

// ID implements Node.
func (m *TestNode) ID() string {
	return m.id
}

// Type implements Node.
func (m *TestNode) Type() string {
	return m.nodeType
}

// Labels implements Node.
func (m *TestNode) Labels() Labels {
	return nil
}

// InternalLabels implements Node.
func (m *TestNode) InternalLabels() InternalLabels {
	return nil
}

func TestCompoundGrapherHappyPath(t *testing.T) {
	graph1 := NewGraph()
	graph1.InsertTree(
		&TestNode{id: "graph1/node1", nodeType: "testnode"},
		&TestNode{id: "graph1/node2", nodeType: "testnode"},
	)

	grapher1 := &MockGrapher{}
	grapher1.On("Graph", mock.Anything).Return(graph1, nil)

	graph2 := NewGraph()
	graph2.InsertTree(
		&TestNode{id: "graph2/node1", nodeType: "testnode"},
		&TestNode{id: "graph2/node2", nodeType: "testnode"},
	)

	grapher2 := &MockGrapher{}
	grapher2.On("Graph", mock.Anything).Return(graph2, nil)

	compoundGrapher := CompoundGrapher{
		Graphers: []Grapher{grapher1, grapher2},
	}

	graph, err := compoundGrapher.Graph(context.TODO())
	require.NoError(t, err)
	grapher1.AssertExpectations(t)
	grapher2.AssertExpectations(t)

	expectedString := `graph1/node1->graph1/node2
graph2/node1->graph2/node2`

	actualString := graph.String()

	// Ignore order when matching
	expected := strings.Split(expectedString, "\n")
	actual := strings.Split(actualString, "\n")

	assert.ElementsMatch(t, expected, actual)
}

func TestCompoundGrapherError(t *testing.T) {
	testCases := []struct {
		Name        string
		Err1        error
		Err2        error
		ExpectedErr string
	}{
		{
			Name:        "grapher 1 fails",
			Err1:        errors.New("failed to graph"),
			Err2:        nil,
			ExpectedErr: "failed to graph",
		},
		{
			Name:        "grapher 2 fails",
			Err1:        nil,
			Err2:        errors.New("failed to graph"),
			ExpectedErr: "failed to graph",
		},
		{
			Name:        "grapher 1 and 2 fail",
			Err1:        errors.New("failed to graph"),
			Err2:        errors.New("failed to graph"),
			ExpectedErr: "failed to graph\nfailed to graph",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			grapher1 := &MockGrapher{}
			grapher1.On("Graph", mock.Anything).Return(Graph(nil), testCase.Err1)

			grapher2 := &MockGrapher{}
			grapher2.On("Graph", mock.Anything).Return(Graph(nil), testCase.Err2)

			compoundGrapher := CompoundGrapher{
				Graphers: []Grapher{grapher1, grapher2},
			}

			_, err := compoundGrapher.Graph(context.TODO())
			require.EqualError(t, err, testCase.ExpectedErr)
			grapher1.AssertExpectations(t)
			grapher2.AssertExpectations(t)
		})
	}
}
