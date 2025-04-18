package imageworkflow

import (
	"slices"
	"sync"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/oci"
)

// Result is a simple type useful to note values that may be empty both when
// successful and unsuccessful. The zero value of a result is an unsuccessful
// result.
type Result[T any] struct {
	OK    bool
	Value T
}

type Data struct {
	sync.Mutex
	ImageReference  oci.Reference
	Created         *time.Time
	Image           string
	LatestReference *oci.Reference
	LatestCreated   *time.Time
	Tags            []string
	Description     string
	FullDescription Result[*models.ImageDescription]
	ReleaseNotes    Result[*models.ImageReleaseNotes]
	Links           []models.ImageLink
	Vulnerabilities []models.ImageVulnerability
	Graph           models.Graph
	Scorecard       Result[*models.ImageScorecard]
	Provenance      Result[*models.ImageProvenance]
	SBOM            Result[*models.ImageSBOM]
	RegistryAuth    *httputil.AuthMux
}

func (d *Data) InsertTag(tag string) {
	d.Lock()
	defer d.Unlock()

	if !slices.Contains(d.Tags, tag) {
		d.Tags = append(d.Tags, tag)
	}
}

func (d *Data) InsertLinks(links []models.ImageLink) {
	d.Lock()
	defer d.Unlock()

	for _, link := range links {
		exists := false
		for _, other := range d.Links {
			if link.Type == other.Type && link.URL == other.URL {
				exists = true
				break
			}
		}

		if !exists {
			d.Links = append(d.Links, link)
		}
	}
}

func (d *Data) InsertLink(link models.ImageLink) {
	d.InsertLinks([]models.ImageLink{link})
}

func (d *Data) InsertVulnerabilities(vulnerabilities []models.ImageVulnerability) {
	d.Lock()
	defer d.Unlock()

	for _, vulnerability := range vulnerabilities {
		exists := false
		for _, other := range d.Vulnerabilities {
			if vulnerability.ID == other.ID {
				exists = true
				break
			}
		}

		if !exists {
			d.Vulnerabilities = append(d.Vulnerabilities, vulnerability)
		}
	}
}

func (d *Data) InsertVulnerability(vulnerability models.ImageVulnerability) {
	d.InsertVulnerabilities([]models.ImageVulnerability{vulnerability})
}
