package models

type UnprocessedStore struct {
	Tags   []*Tag
	Images []*Image
	// Descriptions is mapped by name:version
	Descriptions map[string]*ImageDescription
	// ReleaseNotes is mapped by name:version
	ReleaseNotes map[string]*ImageReleaseNotes
	// Graphs is mapped by name:version
	Graphs map[string][]*Graph
}

type Store struct {
	Tags   []*Tag
	Images []*Image
	// Descriptions is mapped by name:version
	Descriptions map[string]*ImageDescription
	// ReleaseNotes is mapped by name:version
	ReleaseNotes map[string]*ImageReleaseNotes
	// Graphs is mapped by name:version
	Graphs map[string]*Graph
}
