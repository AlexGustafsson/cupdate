package api

import (
	"time"

	"github.com/AlexGustafsson/cupdate/internal/models"
)

type Resource struct {
	Embedded string `json:"_embedded,omitempty"`
	Links    Links  `json:"_links,omitempty"`
}

type Link struct {
	Href      string `json:"href"`
	Name      string `json:"name,omitempty"`
	Templated bool   `json:"templated,omitempty"`
	Title     string `json:"title,omitempty"`
	Type      string `json:"type,omitempty"`
}

type Links map[string][]Link

type ImageResource struct {
	Resource
	Name         string    `json:"name"`
	Tag          string    `json:"tag,omitempty"`
	Digest       string    `json:"digest,omitempty"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
}

type AnnotationsResource struct {
	Resource
	Annotations       map[string]string `json:"annotations"`
	LatestAnnotations map[string]string `json:"latestAnnotations"`
}

type VulnerabilitiesResource struct {
	Resource
	Total int `json:"total"`
}

type Collection struct {
	Resource
	Page  int
	Size  int
	Count int
	Total int
}

type TagsResource struct {
	Resource
	Tags []string
}

type SBOMAttestationResource struct {
	Resource
	SBOM models.ImageSBOM
}
