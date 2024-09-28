package models

type Store struct {
	Tags   []*Tag
	Images []*Image
	// Descriptions is mapped by OCI reference.
	Descriptions map[string]*ImageDescription
	// ReleaseNotes is mapped by OCI reference.
	ReleaseNotes map[string]*ImageReleaseNotes
	// Graphs is mapped by OCI reference.
	Graphs map[string]Graph
}
