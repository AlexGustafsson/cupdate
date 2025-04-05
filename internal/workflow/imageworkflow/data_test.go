package imageworkflow

import (
	"testing"

	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestDataInsertTag(t *testing.T) {
	data := Data{
		Tags: make([]string, 0),
	}

	data.InsertTag("tag")
	assert.Equal(t, []string{"tag"}, data.Tags)

	// No duplicates
	data.InsertTag("tag")
	assert.Equal(t, []string{"tag"}, data.Tags)
}

func TestDataInsertLink(t *testing.T) {
	data := Data{
		Links: make([]models.ImageLink, 0),
	}

	data.InsertLink(models.ImageLink{
		Type: "generic",
		URL:  "https://example.com",
	})
	data.InsertLink(models.ImageLink{
		Type: "svc",
		URL:  "https://example.com/git",
	})
	assert.Equal(t, []models.ImageLink{
		{
			Type: "generic",
			URL:  "https://example.com",
		},
		{
			Type: "svc",
			URL:  "https://example.com/git",
		},
	}, data.Links)

	// No duplicates
	data.InsertLink(models.ImageLink{
		Type: "generic",
		URL:  "https://example.com",
	})
	assert.Equal(t, []models.ImageLink{
		{
			Type: "generic",
			URL:  "https://example.com",
		},
		{
			Type: "svc",
			URL:  "https://example.com/git",
		},
	}, data.Links)
}
