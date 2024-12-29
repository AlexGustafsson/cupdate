package store

import (
	"context"
	"testing"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreInsertRawImage(t *testing.T) {
	store, err := New("file://"+t.TempDir()+"/sqlite.db", false)
	require.NoError(t, err)

	expected := models.RawImage{
		Reference: "mongo:4",
		Tags:      []string{"docker"},
		Graph: models.Graph{
			Edges: map[string]map[string]bool{},
			Nodes: map[string]models.GraphNode{},
		},
		LastProcessed: time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
	}

	ctx, cancel := context.WithCancel(context.TODO())
	go func() {
		defer cancel()
		ch := store.Subscribe(ctx)

		assert.Equal(t, Event{Reference: "mongo:4", Type: EventTypeUpdated}, <-ch)
	}()

	_, err = store.InsertRawImage(context.TODO(), &expected)
	require.NoError(t, err)

	actual, err := store.ListRawImages(context.TODO(), nil)
	require.NoError(t, err)
	assert.EqualValues(t, expected, actual[0])

	<-ctx.Done()
}

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
		Vulnerabilities: []models.ImageVulnerability{
			{
				ID:          1234, // Should not be respected
				Severity:    "low",
				Authority:   "test",
				Description: "Some CVE",
				Links:       []string{"https://example.com"},
			},
		},
		LastModified: time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
		Image:        "https://example.com/logo.png",
	}

	_, err = store.InsertRawImage(context.TODO(), &models.RawImage{
		Reference: expected.Reference,
	})
	require.NoError(t, err)

	err = store.InsertImage(context.TODO(), expected)
	require.NoError(t, err)

	// ID should not be respected
	expected.Vulnerabilities[0].ID = 1

	actual, err := store.GetImage(context.TODO(), "mongo:4")
	require.NoError(t, err)
	assert.EqualValues(t, expected, actual)

	// Make sure triggers don't complain when upserting
	err = store.InsertImage(context.TODO(), expected)
	require.NoError(t, err)
}

func TestStoreTags(t *testing.T) {
	store, err := New("file://"+t.TempDir()+"/sqlite.db", false)
	require.NoError(t, err)

	_, err = store.InsertRawImage(context.TODO(), &models.RawImage{
		Reference: "mongo:4",
	})
	require.NoError(t, err)

	err = store.InsertImage(context.TODO(), &models.Image{
		Reference: "mongo:4",
		Tags:      []string{"docker"},
	})
	require.NoError(t, err)

	actual, err := store.GetTags(context.TODO())
	require.NoError(t, err)
	assert.Equal(t, []string{"docker"}, actual)
}

func TestStoreImageDescription(t *testing.T) {
	store, err := New("file://"+t.TempDir()+"/sqlite.db", false)
	require.NoError(t, err)

	expected := models.ImageDescription{
		Markdown: "# Release",
	}

	_, err = store.InsertRawImage(context.TODO(), &models.RawImage{
		Reference: "mongo:4",
	})
	require.NoError(t, err)

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

	_, err = store.InsertRawImage(context.TODO(), &models.RawImage{
		Reference: "mongo:4",
	})
	require.NoError(t, err)

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

	_, err = store.InsertRawImage(context.TODO(), &models.RawImage{
		Reference: "mongo:4",
	})
	require.NoError(t, err)

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
			Vulnerabilities: []models.ImageVulnerability{},
			LastModified:    time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
			Image:           "https://example.com/logo.png",
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
			Vulnerabilities: []models.ImageVulnerability{},
			LastModified:    time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
			Image:           "https://example.com/logo.png",
		},
	}

	for _, image := range expectedImages {
		_, err := store.InsertRawImage(context.TODO(), &models.RawImage{
			Reference: image.Reference,
		})
		require.NoError(t, err)

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
	actualPage, err = store.ListImages(context.TODO(), &ListImageOptions{Page: 1, Limit: 1})
	require.NoError(t, err)
	assert.Equal(t, expectedPage, actualPage)
}

func TestListImagesQuery(t *testing.T) {
	store, err := New("file://"+t.TempDir()+"/sqlite.db", false)
	require.NoError(t, err)

	images := []models.Image{
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
			Vulnerabilities: []models.ImageVulnerability{},
			LastModified:    time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
			Image:           "https://example.com/logo.png",
		},
	}

	expectedPage := &models.ImagePage{
		Images: []models.Image{
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
				Vulnerabilities: []models.ImageVulnerability{},
				LastModified:    time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
				Image:           "https://example.com/logo.png",
			},
		},
		Summary: models.ImagePageSummary{
			Images:     1,
			Outdated:   1,
			Vulnerable: 0,
			Processing: 0,
		},
		Pagination: models.PaginationMetadata{
			Total:    1,
			Page:     0,
			Size:     30,
			Next:     "",
			Previous: "",
		},
	}

	for _, image := range images {
		_, err := store.InsertRawImage(context.TODO(), &models.RawImage{
			Reference: image.Reference,
		})
		require.NoError(t, err)

		err = store.InsertImage(context.TODO(), &image)
		require.NoError(t, err)
	}

	page, err := store.ListImages(context.TODO(), &ListImageOptions{
		Query: "database",
	})
	require.NoError(t, err)

	assert.Equal(t, expectedPage, page)
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
			Vulnerabilities: []models.ImageVulnerability{},
			LastModified:    time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
		},
		{
			Reference:       "mongo:2",
			LatestReference: "mongo:2",
			Tags:            []string{},
			Links:           []models.ImageLink{},
			Vulnerabilities: []models.ImageVulnerability{},
			LastModified:    time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
		},
		{
			Reference:       "mongo:3",
			LatestReference: "mongo:3",
			Tags:            []string{},
			Links:           []models.ImageLink{},
			Vulnerabilities: []models.ImageVulnerability{},
			LastModified:    time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
		},
		{
			Reference:       "mongo:4",
			LatestReference: "mongo:4",
			Tags:            []string{},
			Links:           []models.ImageLink{},
			Vulnerabilities: []models.ImageVulnerability{},
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
				Vulnerabilities: []models.ImageVulnerability{},
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
		_, err := store.InsertRawImage(context.TODO(), &models.RawImage{
			Reference: image.Reference,
		})
		require.NoError(t, err)

		err = store.InsertImage(context.TODO(), image)
		require.NoError(t, err)
	}

	removed, err := store.DeleteNonPresent(context.TODO(), []string{"mongo:4"})
	require.NoError(t, err)
	assert.Equal(t, int64(3), removed)

	actual, err := store.ListImages(context.TODO(), nil)
	require.NoError(t, err)
	assert.EqualValues(t, expected, actual)
}
