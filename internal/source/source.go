package source

import "context"

type Source interface {
	Entries(context.Context) ([]Entry, error)
}

type Entry struct {
	Image   string
	Version string
	ImageID string
	Origin  Origin
}

type Origin interface {
	Kind() string
}
