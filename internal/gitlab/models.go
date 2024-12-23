package gitlab

import "time"

type ContainerRepository struct {
	ID       string `json:"id"`
	Location string `json:"location"`
}

type ContainerRepositoryTag struct {
	Digest      string    `json:"digest"`
	Location    string    `json:"location"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"createdAt"`
	PublishedAt time.Time `json:"publishedAt"`

	// ... unused fields
}

type Blob struct {
	ID            string `json:"id"`
	Path          string `json:"path"`
	Name          string `json:"name"`
	Extension     string `json:"extension"`
	Size          int    `json:"size"`
	MimeType      string `json:"mime_type"`
	Binary        bool   `json:"binary"`
	RawPath       string `json:"raw_path"`
	BlamePath     string `json:"blame_path"`
	CommitsPath   string `json:"commits_path"`
	TreePath      string `json:"tree_path"`
	Permalink     string `json:"permalink"`
	LastCommitSHA string `json:"last_commit_sha"`
	HTML          string `json:"html"`
	Raw           []byte `json:"-"`

	// ... some available fields are ignored
}
