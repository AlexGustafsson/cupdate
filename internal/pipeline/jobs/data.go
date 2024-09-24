package jobs

import (
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/platform"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
)

type ImageData struct {
	ImageReference oci.Reference
	Image          *string
	LatestVersion  *oci.Reference
	Tags           *[]string
	Description    *string
	ReleaseNotes   *string
	Graph          *platform.Graph
	Links          *[]models.ImageLink
}
