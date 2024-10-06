package imageworkflow

import (
	"slices"
	"sync"

	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
)

type Data struct {
	sync.Mutex
	ImageReference  oci.Reference
	Image           string
	LatestReference oci.Reference
	Tags            []string
	Description     *models.ImageDescription
	ReleaseNotes    *models.ImageReleaseNotes
	Links           []models.ImageLink
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
