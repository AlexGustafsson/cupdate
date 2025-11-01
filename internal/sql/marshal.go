package sql

type Scanner interface {
	Scan(dest ...any) error
}

type Marshaler interface {
	MarshalSQL() ([]any, error)
}

type Unmarshaler interface {
	UnmarshalSQL(row Scanner) error
}

type PrimaryKey struct{}
