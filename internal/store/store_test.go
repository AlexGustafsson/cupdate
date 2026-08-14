package store

import (
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

	err := Initialize(t.Context(), uri)
	require.NoError(t, err)

	store, err := New(t.Context(), uri, readOnly)
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

	_, err := store.InsertRawImage(t.Context(), &expected)
	require.NoError(t, err)

	actual, err := store.ListRawImages(t.Context(), nil)
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

	_, err := store.InsertRawImage(t.Context(), &models.RawImage{
		Reference: expected.Reference,
	})
	require.NoError(t, err)

	err = store.InsertImage(t.Context(), expected, false)
	require.NoError(t, err)

	actual, err := store.GetImage(t.Context(), "mongo:4")
	require.NoError(t, err)
	assert.EqualValues(t, expected, actual)

	// Make sure triggers don't complain when upserting
	err = store.InsertImage(t.Context(), expected, false)
	require.NoError(t, err)

	changes, err := store.GetChanges(t.Context(), nil)
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

	url, err := store.GetImageLogo(t.Context(), expected.Reference)
	require.NoError(t, err)
	assert.Equal(t, expected.Image, url)
}

func TestStoreTags(t *testing.T) {
	store := newStore(t, false)
	defer store.Close()

	_, err := store.InsertRawImage(t.Context(), &models.RawImage{
		Reference: "mongo:4",
	})
	require.NoError(t, err)

	err = store.InsertImage(t.Context(), &models.Image{
		Reference: "mongo:4",
		Tags:      []string{"docker"},
	}, false)
	require.NoError(t, err)

	actual, err := store.GetTags(t.Context())
	require.NoError(t, err)
	assert.Equal(t, []string{"docker"}, actual)
}

func TestStoreImageDescription(t *testing.T) {
	store := newStore(t, false)
	defer store.Close()

	expected := models.ImageDescription{
		Markdown: "# Release",
	}

	_, err := store.InsertRawImage(t.Context(), &models.RawImage{
		Reference: "mongo:4",
	})
	require.NoError(t, err)

	err = store.InsertImage(t.Context(), &models.Image{
		Reference: "mongo:4",
	}, false)
	require.NoError(t, err)

	err = store.InsertImageDescription(t.Context(), "mongo:4", &expected)
	require.NoError(t, err)

	actual, err := store.GetImageDescription(t.Context(), "mongo:4")
	require.NoError(t, err)
	assert.Equal(t, &expected, actual)

	changes, err := store.GetChanges(t.Context(), nil)
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

	_, err := store.InsertRawImage(t.Context(), &models.RawImage{
		Reference: "mongo:4",
	})
	require.NoError(t, err)

	err = store.InsertImage(t.Context(), &models.Image{
		Reference: "mongo:4",
	}, false)
	require.NoError(t, err)

	err = store.InsertImageReleaseNotes(t.Context(), "mongo:4", &expected)
	require.NoError(t, err)

	actual, err := store.GetImageReleaseNotes(t.Context(), "mongo:4")
	require.NoError(t, err)
	assert.Equal(t, &expected, actual)

	changes, err := store.GetChanges(t.Context(), nil)
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

	_, err := store.InsertRawImage(t.Context(), &models.RawImage{
		Reference: "mongo:4",
	})
	require.NoError(t, err)

	err = store.InsertImage(t.Context(), &models.Image{
		Reference: "mongo:4",
	}, false)
	require.NoError(t, err)

	err = store.InsertImageGraph(t.Context(), "mongo:4", &expected)
	require.NoError(t, err)

	actual, err := store.GetImageGraph(t.Context(), "mongo:4")
	require.NoError(t, err)
	assert.Equal(t, &expected, actual)

	changes, err := store.GetChanges(t.Context(), nil)
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
		_, err := store.InsertRawImage(t.Context(), &models.RawImage{
			Reference: image.Reference,
		})
		require.NoError(t, err)

		err = store.InsertImage(t.Context(), &image, false)
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
	actualPage, err := store.ListImages(t.Context(), &ListImageOptions{Page: 0, Limit: 1})
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
	actualPage, err = store.ListImages(t.Context(), &ListImageOptions{Page: 1, Limit: 1})
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
		_, err := store.InsertRawImage(t.Context(), &models.RawImage{
			Reference: image.Reference,
		})
		require.NoError(t, err)

		err = store.InsertImage(t.Context(), &image, false)
		require.NoError(t, err)
	}

	page, err := store.ListImages(t.Context(), &ListImageOptions{
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
		_, err := store.InsertRawImage(t.Context(), &models.RawImage{
			Reference: image.Reference,
		})
		require.NoError(t, err)

		err = store.InsertImage(t.Context(), image, false)
		require.NoError(t, err)
	}

	removed, err := store.DeleteNonPresent(t.Context(), []string{"mongo:4"})
	require.NoError(t, err)
	assert.Equal(t, int64(3), removed)

	actual, err := store.ListImages(t.Context(), nil)
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

	_, err := store.InsertRawImage(t.Context(), &models.RawImage{
		Reference: image.Reference,
	})
	require.NoError(t, err)

	err = store.InsertImage(t.Context(), image, false)
	require.NoError(t, err)

	image.LatestReference = "mongo:5"
	err = store.InsertImage(t.Context(), image, false)
	require.NoError(t, err)

	changes, err := store.GetChanges(t.Context(), nil)
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

	_, err := store.InsertRawImage(t.Context(), rawImage)
	require.NoError(t, err)

	image := &models.Image{
		Reference:       "mongo:4",
		Tags:            []string{},
		Links:           []models.ImageLink{},
		Vulnerabilities: 0,
		LastModified:    time.Date(2024, 10, 05, 18, 39, 0, 0, time.Local).UTC(),
	}

	err = store.InsertImage(t.Context(), image, false)
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

	err = store.InsertWorkflowRun(t.Context(), "mongo:4", expected)
	require.NoError(t, err)

	actual, err := store.GetLatestWorkflowRun(t.Context(), "mongo:4")
	require.NoError(t, err)
	assert.EqualValues(t, &expected, actual)

	// Insert a later job, expect it to be the latest
	expected.Started = time.Date(2025, 02, 01, 17, 40, 0, 0, time.Local).UTC()

	err = store.InsertWorkflowRun(t.Context(), "mongo:4", expected)
	require.NoError(t, err)

	actual, err = store.GetLatestWorkflowRun(t.Context(), "mongo:4")
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

	_, err := store.InsertRawImage(t.Context(), &models.RawImage{
		Reference: image.Reference,
	})
	require.NoError(t, err)

	err = store.InsertImage(t.Context(), image, false)
	require.NoError(t, err)

	err = store.InsertImageDescription(t.Context(), image.Reference, &models.ImageDescription{
		Markdown: "# Image",
	})
	require.NoError(t, err)

	err = store.InsertImageReleaseNotes(t.Context(), image.Reference, &models.ImageReleaseNotes{
		Markdown: "# Release",
	})
	require.NoError(t, err)

	err = store.InsertImageGraph(t.Context(), image.Reference, &models.Graph{
		Edges: make(map[string]map[string]bool),
		Nodes: make(map[string]models.GraphNode),
	})
	require.NoError(t, err)

	err = store.InsertWorkflowRun(t.Context(), image.Reference, models.WorkflowRun{
		TraceID: "1234",
	})
	require.NoError(t, err)

	// Remove the raw image and expect all data to be removed with it
	removed, err := store.DeleteNonPresent(t.Context(), []string{})
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

func TestStoreGetUpdates(t *testing.T) {
	store := newStore(t, false)
	defer store.Close()

	// Insert two similar base image (test INSERT)
	update1ReleasedAt := time.Date(2025, 10, 27, 16, 31, 0, 0, time.Local)
	update2ReleasedAt := time.Date(2025, 10, 28, 16, 31, 0, 0, time.Local)

	base1 := &models.Image{
		Reference: "mongo:4.0.0",
		Annotations: oci.Annotations{
			"version": "4.0.0",
		},
		LatestReference: "mongo:5.0.0",
		LatestAnnotations: oci.Annotations{
			"version": "5.0.0",
		},
		LatestCreated:       &update1ReleasedAt,
		VersionDiffSortable: 1,
	}

	_, err := store.InsertRawImage(t.Context(), &models.RawImage{
		Reference: base1.Reference,
	})
	require.NoError(t, err)

	err = store.InsertImage(t.Context(), base1, false)
	require.NoError(t, err)

	base2 := &models.Image{
		Reference: "mongo:4.0.1",
		Annotations: oci.Annotations{
			"version": "4.0.1",
		},
		LatestReference: "mongo:5.0.0",
		LatestAnnotations: oci.Annotations{
			"version": "5.0.0",
		},
		LatestCreated:       &update1ReleasedAt,
		VersionDiffSortable: 1,
	}

	_, err = store.InsertRawImage(t.Context(), &models.RawImage{
		Reference: base2.Reference,
	})
	require.NoError(t, err)

	err = store.InsertImage(t.Context(), base2, false)
	require.NoError(t, err)

	// Insert image updates (test UPDATE)
	updated1 := &models.Image{
		Reference: "mongo:4.0.0",
		Annotations: oci.Annotations{
			"version": "4.0.0",
		},
		LatestReference: "mongo:5.0.1",
		LatestAnnotations: oci.Annotations{
			"version": "5.0.1",
		},
		LatestCreated:       &update2ReleasedAt,
		VersionDiffSortable: 1,
	}

	_, err = store.InsertRawImage(t.Context(), &models.RawImage{
		Reference: updated1.Reference,
	})
	require.NoError(t, err)

	err = store.InsertImage(t.Context(), updated1, false)
	require.NoError(t, err)

	updated2 := &models.Image{
		Reference: "mongo:4.0.1",
		Annotations: oci.Annotations{
			"version": "4.0.1",
		},
		LatestReference: "mongo:5.0.1",
		LatestAnnotations: oci.Annotations{
			"version": "5.0.1",
		},
		LatestCreated:       &update2ReleasedAt,
		VersionDiffSortable: 1,
	}

	_, err = store.InsertRawImage(t.Context(), &models.RawImage{
		Reference: updated2.Reference,
	})
	require.NoError(t, err)

	err = store.InsertImage(t.Context(), updated2, false)
	require.NoError(t, err)

	// Expect there to be updates
	updates, err := store.GetUpdates(t.Context(), nil)
	require.NoError(t, err)

	assert.EqualValues(t, []models.ImageUpdate{
		{
			NewReference: "mongo:5.0.1",
			NewAnnotations: oci.Annotations{
				"version": "5.0.1",
			},
			OldReference: "mongo:4.0.1",
			OldAnnotations: oci.Annotations{
				"version": "4.0.1",
			},
			Identified: updates[0].Identified, // Difficult to assert
			Released:   &update2ReleasedAt,
		},
		{
			NewReference: "mongo:5.0.1",
			NewAnnotations: oci.Annotations{
				"version": "5.0.1",
			},
			OldReference: "mongo:4.0.0",
			OldAnnotations: oci.Annotations{
				"version": "4.0.0",
			},
			Identified: updates[1].Identified, // Difficult to assert
			Released:   &update2ReleasedAt,
		},
		{
			NewReference: "mongo:5.0.0",
			NewAnnotations: oci.Annotations{
				"version": "5.0.0",
			},
			OldReference: "mongo:4.0.1",
			OldAnnotations: oci.Annotations{
				"version": "4.0.1",
			},
			Identified: updates[2].Identified, // Difficult to assert
			Released:   &update1ReleasedAt,
		},
		{
			NewReference: "mongo:5.0.0",
			NewAnnotations: oci.Annotations{
				"version": "5.0.0",
			},
			OldReference: "mongo:4.0.0",
			OldAnnotations: oci.Annotations{
				"version": "4.0.0",
			},
			Identified: updates[3].Identified, // Difficult to assert
			Released:   &update1ReleasedAt,
		},
	}, updates)
}

func TestStoreInsertUnprocessedImage(t *testing.T) {
	store := newStore(t, false)
	defer store.Close()

	_, err := store.InsertRawImage(t.Context(), &models.RawImage{
		Reference: "mongo:4",
	})
	require.NoError(t, err)

	// Insert a processed image
	err = store.InsertImage(t.Context(), &models.Image{
		Reference: "mongo:4",
		Tags:      []string{"some tag"},
	}, false)
	require.NoError(t, err)

	// Try to insert an unprocessed image, expecting it to succeed but change
	// nothing
	err = store.InsertImage(t.Context(), &models.Image{
		Reference: "mongo:4",
		Tags:      []string{"unprocessed"},
	}, true)
	require.NoError(t, err)

	actual, err := store.GetImage(t.Context(), "mongo:4")
	require.NoError(t, err)
	assert.EqualValues(t, &models.Image{
		Reference:    "mongo:4",
		Tags:         []string{"some tag"},
		Links:        []models.ImageLink{},
		LastModified: actual.LastModified, // Difficult to assert
	}, actual)

	// Try to insert another unprocessed image, expecting it to succeed as there's
	// no existing image
	_, err = store.InsertRawImage(t.Context(), &models.RawImage{
		Reference: "mongo:5",
	})
	require.NoError(t, err)

	err = store.InsertImage(t.Context(), &models.Image{
		Reference: "mongo:5",
		Tags:      []string{"unprocessed"},
	}, true)
	require.NoError(t, err)

	actual, err = store.GetImage(t.Context(), "mongo:5")
	require.NoError(t, err)
	assert.EqualValues(t, &models.Image{
		Reference:    "mongo:5",
		Tags:         []string{"unprocessed"},
		Links:        []models.ImageLink{},
		LastModified: actual.LastModified, // Difficult to assert
	}, actual)
}
