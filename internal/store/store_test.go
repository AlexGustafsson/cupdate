package store

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newStore(t *testing.T, readOnly bool) *Store {
	uri := "file://" + t.TempDir() + "/sqlite.db"

	err := Initialize(context.TODO(), uri)
	require.NoError(t, err)

	store, err := New(uri, readOnly)
	require.NoError(t, err)

	return store
}

func TestStoreInsertRawImage(t *testing.T) {
	store := newStore(t, false)
	defer store.Close()

	expected := models.RawImage{
		Reference: "mongo:4",
		Tags:      []string{"docker"},
		Graph: models.Graph{
			Edges: map[string]map[string]bool{},
			Nodes: map[string]models.GraphNode{},
		},
		LastProcessed: time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
	}

	_, err := store.InsertRawImage(context.TODO(), &expected)
	require.NoError(t, err)

	actual, err := store.ListRawImages(context.TODO(), nil)
	require.NoError(t, err)
	assert.EqualValues(t, expected, actual[0])
}

func TestStoreInsertImage(t *testing.T) {
	store := newStore(t, false)
	defer store.Close()

	expected := &models.Image{
		Reference: "mongo:4",
		Annotations: oci.Annotations{
			"version": "4.0.0",
		},
		LatestReference: "mongo:4",
		LatestAnnotations: oci.Annotations{
			"version": "4.0.0",
		},
		Description: "Mongo is a database",
		Tags:        []string{"docker"},
		Links: []models.ImageLink{
			{
				Type: "docker",
				URL:  "https://docker.com/_/mongo",
			},
		},
		Vulnerabilities: 0,
		LastModified:    time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
		Image:           "https://example.com/logo.png",
	}

	_, err := store.InsertRawImage(context.TODO(), &models.RawImage{
		Reference: expected.Reference,
	})
	require.NoError(t, err)

	err = store.InsertImage(context.TODO(), expected)
	require.NoError(t, err)

	actual, err := store.GetImage(context.TODO(), "mongo:4")
	require.NoError(t, err)
	assert.EqualValues(t, expected, actual)

	// Make sure triggers don't complain when upserting
	err = store.InsertImage(context.TODO(), expected)
	require.NoError(t, err)

	changes, err := store.GetChanges(context.TODO(), nil)
	require.NoError(t, err)
	assert.EqualValues(t, []Change{
		{
			Reference:    "mongo:4",
			Time:         changes[0].Time,
			Type:         "insert",
			ChangedBasic: true,
		},
		{
			Reference:    "mongo:4",
			Time:         changes[1].Time,
			Type:         "insert",
			ChangedLinks: true,
		},
	}, changes)

	url, err := store.GetImageLogo(context.TODO(), expected.Reference)
	require.NoError(t, err)
	assert.Equal(t, expected.Image, url)
}

func TestStoreTags(t *testing.T) {
	store := newStore(t, false)
	defer store.Close()

	_, err := store.InsertRawImage(context.TODO(), &models.RawImage{
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
	store := newStore(t, false)
	defer store.Close()

	expected := models.ImageDescription{
		Markdown: "# Release",
	}

	_, err := store.InsertRawImage(context.TODO(), &models.RawImage{
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

	changes, err := store.GetChanges(context.TODO(), nil)
	require.NoError(t, err)
	assert.EqualValues(t, []Change{
		{
			Reference:    "mongo:4",
			Time:         changes[0].Time,
			Type:         "insert",
			ChangedBasic: true,
		},
		{
			Reference:    "mongo:4",
			Time:         changes[1].Time,
			Type:         "insert",
			ChangedLinks: true,
		},
		{
			Reference:          "mongo:4",
			Time:               changes[2].Time,
			Type:               "insert",
			ChangedDescription: true,
		},
	}, changes)
}

func TestStoreImageReleaseNotes(t *testing.T) {
	store := newStore(t, false)
	defer store.Close()

	expected := models.ImageReleaseNotes{
		Title:    "Release",
		Markdown: "# Release",
		Released: time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
	}

	_, err := store.InsertRawImage(context.TODO(), &models.RawImage{
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

	changes, err := store.GetChanges(context.TODO(), nil)
	require.NoError(t, err)
	assert.EqualValues(t, []Change{
		{
			Reference:    "mongo:4",
			Time:         changes[0].Time,
			Type:         "insert",
			ChangedBasic: true,
		},
		{
			Reference:    "mongo:4",
			Time:         changes[1].Time,
			Type:         "insert",
			ChangedLinks: true,
		},
		{
			Reference:           "mongo:4",
			Time:                changes[2].Time,
			Type:                "insert",
			ChangedReleaseNotes: true,
		},
	}, changes)
}

func TestStoreImageGraph(t *testing.T) {
	store := newStore(t, false)
	defer store.Close()

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

	_, err := store.InsertRawImage(context.TODO(), &models.RawImage{
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

	changes, err := store.GetChanges(context.TODO(), nil)
	require.NoError(t, err)
	assert.EqualValues(t, []Change{
		{
			Reference:    "mongo:4",
			Time:         changes[0].Time,
			Type:         "insert",
			ChangedBasic: true,
		},
		{
			Reference:    "mongo:4",
			Time:         changes[1].Time,
			Type:         "insert",
			ChangedLinks: true,
		},
		{
			Reference:    "mongo:4",
			Time:         changes[2].Time,
			Type:         "insert",
			ChangedGraph: true,
		},
	}, changes)
}

func TestListImages(t *testing.T) {
	store := newStore(t, false)
	defer store.Close()

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
			Vulnerabilities: 0,
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
			Vulnerabilities: 0,
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
			Page:     1,
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
			Page:     2,
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
	store := newStore(t, false)
	defer store.Close()

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
			Vulnerabilities: 0,
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
				Vulnerabilities: 0,
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
			Page:     1,
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
	store := newStore(t, false)
	defer store.Close()

	images := []*models.Image{
		{
			Reference:       "mongo:1",
			LatestReference: "mongo:1",
			Tags:            []string{},
			Links:           []models.ImageLink{},
			Vulnerabilities: 0,
			LastModified:    time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
		},
		{
			Reference:       "mongo:2",
			LatestReference: "mongo:2",
			Tags:            []string{},
			Links:           []models.ImageLink{},
			Vulnerabilities: 0,
			LastModified:    time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
		},
		{
			Reference:       "mongo:3",
			LatestReference: "mongo:3",
			Tags:            []string{},
			Links:           []models.ImageLink{},
			Vulnerabilities: 0,
			LastModified:    time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
		},
		{
			Reference:       "mongo:4",
			LatestReference: "mongo:4",
			Tags:            []string{},
			Links:           []models.ImageLink{},
			Vulnerabilities: 0,
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
				Vulnerabilities: 0,
				LastModified:    time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
			},
		},
		Summary: models.ImagePageSummary{
			Images: 1,
		},
		Pagination: models.PaginationMetadata{
			Total: 1,
			Page:  1,
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

func TestStoreUpdateImageReference(t *testing.T) {
	store := newStore(t, false)
	defer store.Close()

	image := &models.Image{
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
		Vulnerabilities: 0,
		LastModified:    time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local),
		Image:           "https://example.com/logo.png",
	}

	_, err := store.InsertRawImage(context.TODO(), &models.RawImage{
		Reference: image.Reference,
	})
	require.NoError(t, err)

	err = store.InsertImage(context.TODO(), image)
	require.NoError(t, err)

	image.LatestReference = "mongo:5"
	err = store.InsertImage(context.TODO(), image)
	require.NoError(t, err)

	changes, err := store.GetChanges(context.TODO(), nil)
	require.NoError(t, err)
	assert.EqualValues(t, []Change{
		{
			Reference:    "mongo:4",
			Time:         changes[0].Time,
			Type:         "insert",
			ChangedBasic: true,
		},
		{
			Reference:    "mongo:4",
			Time:         changes[1].Time,
			Type:         "insert",
			ChangedLinks: true,
		},
		{
			Reference:    "mongo:4",
			Time:         changes[2].Time,
			Type:         "update",
			ChangedBasic: true,
		},
	}, changes)
}

// TODO: Times are set to a fixed zone as there are issues when comparing the
// times cross-platform. Using the local time will not work if tested with UTC.
// Use time.Local).UTC() as a workaround.
// SEE: https://github.com/stretchr/testify/issues/843#issuecomment-1952362012
func TestInsertWorkflowRun(t *testing.T) {
	store := newStore(t, false)
	defer store.Close()

	rawImage := &models.RawImage{
		Reference: "mongo:4",
	}

	_, err := store.InsertRawImage(context.TODO(), rawImage)
	require.NoError(t, err)

	image := &models.Image{
		Reference:       "mongo:4",
		Tags:            []string{},
		Links:           []models.ImageLink{},
		Vulnerabilities: 0,
		LastModified:    time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local).UTC(),
	}

	err = store.InsertImage(context.TODO(), image)
	require.NoError(t, err)

	expected := models.WorkflowRun{
		TraceID:         "trace-123",
		Started:         time.Date(2025, 02, 01, 17, 35, 0, 0, time.Local).UTC(),
		DurationSeconds: 25.0,
		Result:          models.WorkflowRunResultSucceeded,
		Jobs: []models.JobRun{
			{
				Result: models.JobRunResultSucceeded,
				Steps: []models.StepRun{
					{
						Result:          models.StepRunResultSucceeded,
						StepName:        "test step",
						Started:         time.Date(2025, 02, 01, 17, 35, 0, 0, time.Local).UTC(),
						DurationSeconds: 25.0,
					},
				},
				DependsOn:       []string{},
				JobID:           "test-job",
				JobName:         "test job",
				Started:         time.Date(2025, 02, 01, 17, 35, 0, 0, time.Local).UTC(),
				DurationSeconds: 25.0,
			},
		},
	}

	err = store.InsertWorkflowRun(context.TODO(), "mongo:4", expected)
	require.NoError(t, err)

	actual, err := store.GetLatestWorkflowRun(context.TODO(), "mongo:4")
	require.NoError(t, err)
	assert.EqualValues(t, &expected, actual)

	// Insert a later job, expect it to be the latest
	expected.Started = time.Date(2025, 02, 01, 17, 40, 0, 0, time.Local).UTC()

	err = store.InsertWorkflowRun(context.TODO(), "mongo:4", expected)
	require.NoError(t, err)

	actual, err = store.GetLatestWorkflowRun(context.TODO(), "mongo:4")
	require.NoError(t, err)
	assert.EqualValues(t, &expected, actual)
}

func TestCascadeDelete(t *testing.T) {
	store := newStore(t, false)
	defer store.Close()

	image := &models.Image{
		Reference: "mongo:4",
		Tags:      []string{"docker"},
		Links: []models.ImageLink{
			{
				Type: "docker",
				URL:  "https://docker.com/_/mongo",
			},
		},
		Vulnerabilities: 0,
	}

	_, err := store.InsertRawImage(context.TODO(), &models.RawImage{
		Reference: image.Reference,
	})
	require.NoError(t, err)

	err = store.InsertImage(context.TODO(), image)
	require.NoError(t, err)

	err = store.InsertImageDescription(context.TODO(), image.Reference, &models.ImageDescription{
		Markdown: "# Image",
	})
	require.NoError(t, err)

	err = store.InsertImageReleaseNotes(context.TODO(), image.Reference, &models.ImageReleaseNotes{
		Markdown: "# Release",
	})
	require.NoError(t, err)

	err = store.InsertImageGraph(context.TODO(), image.Reference, &models.Graph{
		Edges: make(map[string]map[string]bool),
		Nodes: make(map[string]models.GraphNode),
	})
	require.NoError(t, err)

	err = store.InsertWorkflowRun(context.TODO(), image.Reference, models.WorkflowRun{
		TraceID: "1234",
	})
	require.NoError(t, err)

	// Remove the raw image and expect all data to be removed with it
	removed, err := store.DeleteNonPresent(context.TODO(), []string{})
	require.NoError(t, err)
	assert.Equal(t, int64(1), removed)

	res, err := store.db.Query("SELECT name FROM sqlite_master WHERE type='table';")
	require.NoError(t, err)

	ftsTables := []string{"images_fts"}
	ignoredTables := []string{"revision"}

	for res.Next() {
		var tableName string
		require.NoError(t, res.Scan(&tableName))

		if slices.Contains(ignoredTables, tableName) {
			continue
		}

		// Ignore tables created and used by FTS
		isFTS := false
		for _, ftsTable := range ftsTables {
			if strings.HasPrefix(tableName, ftsTable+"_") {
				isFTS = true
				break
			}
		}
		if isFTS {
			continue
		}

		res := store.db.QueryRow(fmt.Sprintf("SELECT COUNT(1) FROM %s;", tableName))

		var count int
		require.NoError(t, res.Scan(&count))

		assert.Equal(t, 0, count, "Table %s should be empty", tableName)
	}
	require.NoError(t, res.Err())
	require.NoError(t, res.Close())
}
