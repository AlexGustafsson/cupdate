package store

import (
	"context"
	"testing"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreInsertImage(t *testing.T) {
	store, err := New("file://"+t.TempDir()+"/sqlite.db", false)
	require.NoError(t, err)

	expected := &models.Image{
		Reference:       "mongo:4",
		LatestReference: "mongo:4",
		Description:     "Mongo is a database",
		Tags:            []string{"docker"},
		Links: []models.ImageLink{
			{
				Type: "docker",
				URL:  "https://docker.com/_/mongo",
			},
		},
		LastModified: time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
		Image:        "https://example.com/logo.png",
	}

	err = store.InsertTag(context.TODO(), &models.Tag{
		Name:        "docker",
		Description: "Docker",
		Color:       "#0000ff",
	})
	require.NoError(t, err)

	err = store.InsertImage(context.TODO(), expected)
	require.NoError(t, err)

	actual, err := store.GetImage(context.TODO(), "mongo:4")
	require.NoError(t, err)
	assert.EqualValues(t, expected, actual)
}

func TestStoreTags(t *testing.T) {
	store, err := New("file://"+t.TempDir()+"/sqlite.db", false)
	require.NoError(t, err)

	expected := models.Tag{
		Name:        "docker",
		Description: "Docker",
		Color:       "#0000ff",
	}

	err = store.InsertTag(context.TODO(), &expected)
	require.NoError(t, err)

	actual, err := store.GetTags(context.TODO())
	require.NoError(t, err)
	assert.Equal(t, []models.Tag{expected}, actual)
}

func TestStoreImageDescription(t *testing.T) {
	store, err := New("file://"+t.TempDir()+"/sqlite.db", false)
	require.NoError(t, err)

	expected := models.ImageDescription{
		Markdown: "# Release",
	}

	err = store.InsertImage(context.TODO(), &models.Image{
		Reference: "mongo:4",
	})
	require.NoError(t, err)

	err = store.InsertImageDescription(context.TODO(), "mongo:4", &expected)
	require.NoError(t, err)

	actual, err := store.GetImageDescription(context.TODO(), "mongo:4")
	require.NoError(t, err)
	assert.Equal(t, &expected, actual)
}

func TestStoreImageReleaseNotes(t *testing.T) {
	store, err := New("file://"+t.TempDir()+"/sqlite.db", false)
	require.NoError(t, err)

	expected := models.ImageReleaseNotes{
		Title:    "Release",
		Markdown: "# Release",
		Released: time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
	}

	err = store.InsertImage(context.TODO(), &models.Image{
		Reference: "mongo:4",
	})
	require.NoError(t, err)

	err = store.InsertImageReleaseNotes(context.TODO(), "mongo:4", &expected)
	require.NoError(t, err)

	actual, err := store.GetImageReleaseNotes(context.TODO(), "mongo:4")
	require.NoError(t, err)
	assert.Equal(t, &expected, actual)
}

func TestStoreImageGraph(t *testing.T) {
	store, err := New("file://"+t.TempDir()+"/sqlite.db", false)
	require.NoError(t, err)

	expected := models.Graph{
		Edges: map[string]map[string]bool{
			"mongo:4": {
				"pod": true,
			},
		},
		Nodes: map[string]models.GraphNode{
			"mongo:4": {
				Domain: "oci",
				Type:   "image",
				Name:   "mongo:4",
			},
			"mongo": {
				Domain: "kubernetes",
				Type:   "pod",
				Name:   "mongo",
			},
		},
	}

	err = store.InsertImage(context.TODO(), &models.Image{
		Reference: "mongo:4",
	})
	require.NoError(t, err)

	err = store.InsertImageGraph(context.TODO(), "mongo:4", &expected)
	require.NoError(t, err)

	actual, err := store.GetImageGraph(context.TODO(), "mongo:4")
	require.NoError(t, err)
	assert.Equal(t, &expected, actual)
}

func TestListImages(t *testing.T) {
	store, err := New("file://"+t.TempDir()+"/sqlite.db", false)
	require.NoError(t, err)

	err = store.InsertTag(context.TODO(), &models.Tag{
		Name:        "docker",
		Description: "Docker",
		Color:       "#0000ff",
	})
	require.NoError(t, err)

	expectedImages := []models.Image{
		{
			Reference:       "mongo:3",
			LatestReference: "mongo:4",
			Description:     "Mongo is a database",
			Tags:            []string{"docker"},
			Links: []models.ImageLink{
				{
					Type: "docker",
					URL:  "https://docker.com/_/mongo",
				},
			},
			LastModified: time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
			Image:        "https://example.com/logo.png",
		},
		{
			Reference:       "mongo:4",
			LatestReference: "mongo:4",
			Description:     "Mongo is a database",
			Tags:            []string{"docker"},
			Links: []models.ImageLink{
				{
					Type: "docker",
					URL:  "https://docker.com/_/mongo",
				},
			},
			LastModified: time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
			Image:        "https://example.com/logo.png",
		},
	}

	for _, image := range expectedImages {
		err = store.InsertImage(context.TODO(), &image)
		require.NoError(t, err)
	}

	expectedPage := &models.ImagePage{
		Images: []models.Image{
			expectedImages[0],
		},
		Summary: models.ImagePageSummary{
			Images:   2,
			Outdated: 1,
		},
		Pagination: models.PaginationMetadata{
			Total:    2,
			Page:     0,
			Size:     1,
			Next:     "",
			Previous: "",
		},
	}
	actualPage, err := store.ListImages(context.TODO(), &ListImageOptions{Page: 0, Limit: 1})
	require.NoError(t, err)
	assert.Equal(t, expectedPage, actualPage)

	expectedPage = &models.ImagePage{
		Images: []models.Image{
			expectedImages[1],
		},
		Summary: models.ImagePageSummary{
			Images:   2,
			Outdated: 1,
		},
		Pagination: models.PaginationMetadata{
			Total:    2,
			Page:     1,
			Size:     1,
			Next:     "",
			Previous: "",
		},
	}
	actualPage, err = store.ListImages(context.TODO(), &ListImageOptions{Page: 0, Limit: 1})
	require.NoError(t, err)
	assert.Equal(t, expectedPage, actualPage)

}

func TestStoreDeleteNonPresent(t *testing.T) {
	store, err := New("file://"+t.TempDir()+"/sqlite.db", false)
	require.NoError(t, err)

	images := []*models.Image{
		{
			Reference:       "mongo:1",
			LatestReference: "mongo:1",
			Tags:            []string{},
			Links:           []models.ImageLink{},
			LastModified:    time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
		},
		{
			Reference:       "mongo:2",
			LatestReference: "mongo:2",
			Tags:            []string{},
			Links:           []models.ImageLink{},
			LastModified:    time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
		},
		{
			Reference:       "mongo:3",
			LatestReference: "mongo:3",
			Tags:            []string{},
			Links:           []models.ImageLink{},
			LastModified:    time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
		},
		{
			Reference:       "mongo:4",
			LatestReference: "mongo:4",
			Tags:            []string{},
			Links:           []models.ImageLink{},
			LastModified:    time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
		},
	}

	expected := &models.ImagePage{
		Images: []models.Image{
			{
				Reference:       "mongo:4",
				LatestReference: "mongo:4",
				Tags:            []string{},
				Links:           []models.ImageLink{},
				LastModified:    time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
			},
		},
		Summary: models.ImagePageSummary{
			Images: 1,
		},
		Pagination: models.PaginationMetadata{
			Total: 1,
			Page:  0,
			Size:  30,
		},
	}

	for _, image := range images {
		err := store.InsertImage(context.TODO(), image)
		require.NoError(t, err)
	}

	removed, err := store.DeleteNonPresent(context.TODO(), []string{"mongo:4"})
	require.NoError(t, err)
	assert.Equal(t, int64(1), removed)

	actual, err := store.ListImages(context.TODO(), nil)
	require.NoError(t, err)
	assert.EqualValues(t, expected, actual)
}
