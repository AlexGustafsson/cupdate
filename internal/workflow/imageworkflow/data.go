package imageworkflow

import (
	"sync"

	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
)

type Data struct {
	sync.Mutex
	ImageReference oci.Reference
	Image          string
	LatestVersion  *oci.Reference
	Tags           []string
	Description    *models.ImageDescription
	ReleaseNotes   *models.ImageReleaseNotes
	Links          []models.ImageLink
}
