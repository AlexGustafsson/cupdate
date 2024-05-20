package source

import "context"

type Source interface {
	EachListItem(context.Context, func(Entry) error) error
}

type Entry struct {
	Image   string
	Version string
	Origin  Origin
}

type Origin interface {
	Kind() string
}
