package jobs

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/AlexGustafsson/cupdate/internal/cache"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/pipeline"
	"github.com/AlexGustafsson/cupdate/internal/platform"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"github.com/stretchr/testify/require"
)

func TestDefaultPipeline(t *testing.T) {
	cache, err := cache.NewDiskCache("cache")
	require.NoError(t, err)

	ref, err := oci.ParseReference("renovate/renovate:38.70.2")
	require.NoError(t, err)

	image := ""
	latestVersion, err := oci.ParseReference("renovate/renovate:38.70.2")
	require.NoError(t, err)
	tags := make([]string, 0)
	description := ""
	releaseNotes := ""
	graph := platform.NewGraph() // TODO
	links := make([]models.ImageLink, 0)

	data := ImageData{
		ImageReference: ref,
		Image:          &image,
		LatestVersion:  &latestVersion,
		Tags:           &tags,
		Description:    &description,
		ReleaseNotes:   &releaseNotes,
		Graph:          &graph,
		Links:          &links,
	}

	pipeline := pipeline.New(cache, DefaultJobs())
	pipeline.Run(context.TODO(), data)

	// NOTE: Turns out the JSON cache doesn't like encoding reference.Reference
	fmt.Println("image", image)
	fmt.Println("latestVersion", latestVersion)
	fmt.Println("tags", tags)
	fmt.Println("description", strings.ReplaceAll(description[:min(100, len(description))], "\n", "\t"))
	fmt.Println("releaseNotes", strings.ReplaceAll(releaseNotes[:min(100, len(releaseNotes))], "\n", "\t"))
	// fmt.Println("graph", graph) // TODO
	fmt.Printf("links %+v\n", links)
}
