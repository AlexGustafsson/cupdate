package imageworkflow

import (
	"slices"
	"sync"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/oci"
)

type Data struct {
	sync.Mutex
	ImageReference  oci.Reference
	Created         *time.Time
	Image           string
	LatestReference *oci.Reference
	LatestCreated   *time.Time
	Tags            []string
	Description     string
	FullDescription *models.ImageDescription
	ReleaseNotes    *models.ImageReleaseNotes
	Links           []models.ImageLink
	Vulnerabilities []models.ImageVulnerability
	Graph           models.Graph
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

	d.Vulnerabilities = append(d.Vulnerabilities, vulnerabilities...)
}

func (d *Data) InsertVulnerability(vulnerability models.ImageVulnerability) {
	d.InsertVulnerabilities([]models.ImageVulnerability{vulnerability})
}
