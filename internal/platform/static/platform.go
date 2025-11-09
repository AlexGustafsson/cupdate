package static

import (
	"bufio"
	"context"
	"os"

	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/platform"
)

var _ platform.Grapher = (*Platform)(nil)

type Platform struct {
	FilePath string
}

// Graph implements platform.Grapher.
func (p *Platform) Graph(context.Context) (platform.Graph, error) {
	file, err := os.Open(p.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	graph := platform.NewGraph()

	reader := bufio.NewScanner(file)
	for reader.Scan() {
		reference, err := oci.ParseReference(reader.Text())
		if err != nil {
			// TODO: Provide a helpful error with line number?
			return nil, err
		}

		graph.InsertTree(platform.ImageNode{
			Reference: reference,
		})
	}

	return graph, nil
}
