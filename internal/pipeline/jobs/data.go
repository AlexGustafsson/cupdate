package jobs

import (
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/platform"
	"github.com/distribution/reference"
)

type ImageData struct {
	ImageReference reference.Named
	Image          *string
	LatestVersion  *reference.Named
	Tags           *[]string
	Description    *string
	ReleaseNotes   *string
	Graph          *platform.Graph
	Links          *[]models.ImageLink
}