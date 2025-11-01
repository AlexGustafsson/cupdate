package schema

type Blob[T any] struct {
	Type string `sql:"type"`
	Data T      `sql:"data"`
}
