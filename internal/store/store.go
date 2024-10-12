package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/models"
	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

// TODO: For single rows use QueryRowContext instead of QueryContext

func New(uri string, readonly bool) (*Store, error) {
	uri += "?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(1000)&_time_format=sqlite"
	if readonly {
		uri += "&_pragma=query_only(true)"
	}

	db, err := sql.Open("sqlite", uri)
	if err != nil {
		return nil, err
	}

	if !readonly {
		// SEE: docs/architecture/database.md
		_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS raw_images (
		reference TEXT PRIMARY KEY NOT NULL,
		tags BLOB,
		graph BLOB
	)

	CREATE TABLE IF NOT EXISTS images (
		reference TEXT PRIMARY KEY NOT NULL,
		latestReference TEXT NOT NULL,
		description TEXT NOT NULL,
		lastModified DATETIME NOT NULL,
		imageUrl TEXT NOT NULL
		FOREIGN KEY(reference) REFERENCES raw_images(reference)
		ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS tags (
		name TEXT PRIMARY KEY NOT NULL,
		color TEXT NOT NULL,
		description TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS images_tags (
		reference TEXT NOT NULL,
		tag TEXT NOT NULL,
		PRIMARY KEY (reference, tag),
		FOREIGN KEY(reference) REFERENCES images(reference)
		FOREIGN KEY(tag) REFERENCES tags(name)
		ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS images_links (
		reference TEXT NOT NULL,
		url TEXT NOT NULL,
		type TEXT NOT NULL,
		PRIMARY KEY (reference, url),
		FOREIGN KEY(reference) REFERENCES images(reference)
		ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS images_release_notes (
		reference TEXT NOT NULL,
		title TEXT NOT NULL,
		html TEXT NOT NULL,
		markdown TEXT NOT NULL,
		released DATETIME NOT NULL,
		PRIMARY KEY (reference),
		FOREIGN KEY(reference) REFERENCES images(reference)
		ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS images_descriptions (
		reference TEXT NOT NULL,
		html TEXT NOT NULL,
		markdown TEXT NOT NULL,
		PRIMARY KEY (reference),
		FOREIGN KEY(reference) REFERENCES images(reference)
		ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS images_graphs (
		reference TEXT NOT NULL,
		graph BLOB NOT NULL,
		PRIMARY KEY (reference),
		FOREIGN KEY(reference) REFERENCES images(reference)
		ON DELETE CASCADE
	);
	`)
		if err != nil {
			db.Close()
			return nil, err
		}
	}

	return &Store{db: db}, nil
}

func (s *Store) InsertRawImage(ctx context.Context, reference string, tags []string, graph models.Graph) error {
	statement, err := s.db.PrepareContext(ctx, `INSERT OR REPLACE INTO raw_images
		(reference, tags, graph)
		VALUES
		(?, ?, ?);`)
	if err != nil {
		return err
	}

	_, err = statement.ExecContext(ctx, tags, graph)
	statement.Close()
	if err != nil {
		return err
	}

	return nil
}

// TODO: Get all raw images that haven't been updated since ...
func (s *Store) GetRawImages(ctx context.Context, lastModified time.Duration) (*models.Image, error) {
	statement, err := s.db.PrepareContext(ctx, `SELECT
	reference, latestReference, description, image, lastModified
	FROM images WHERE reference = ?;`)
	if err != nil {
		return nil, err
	}

	var image models.Image

	// TODO: Implement the scan interface for models
	res, err := statement.QueryContext(ctx, reference)
	statement.Close()
	if err != nil {
		return nil, err
	}

	if !res.Next() {
		res.Close()
		return nil, res.Err()
	}

	err = res.Scan(&image.Reference, &image.LatestReference, &image.Description, &image.Image, &image.LastModified)
	res.Close()
	if err != nil {
		return nil, err
	}

	image.Tags, err = s.GetImagesTags(ctx, reference)
	if err != nil {
		return nil, err
	}

	image.Links, err = s.GetImagesLinks(ctx, reference)
	if err != nil {
		return nil, err
	}

	return &image, nil
}

func (s *Store) InsertImage(ctx context.Context, image *models.Image) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	statement, err := tx.PrepareContext(ctx, `INSERT OR REPLACE INTO images
	(reference, latestReference, description, lastModified, image)
	VALUES
	(?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}

	// TODO: Implement the scan interface for models
	_, err = statement.ExecContext(ctx, image.Reference, image.LatestReference, image.Description, image.LastModified, image.Image)
	statement.Close()
	if err != nil {
		return err
	}

	for _, tag := range image.Tags {
		statement, err := tx.PrepareContext(ctx, `INSERT OR REPLACE INTO images_tags
		(reference, tag)
		VALUES
		(?, ?);`)
		if err != nil {
			return err
		}

		_, err = statement.ExecContext(ctx, image.Reference, tag)
		statement.Close()
		if err != nil {
			return err
		}
	}

	for _, link := range image.Links {
		statement, err := tx.PrepareContext(ctx, `INSERT OR REPLACE INTO images_links
		(reference, type, url)
		VALUES
		(?, ?, ?);`)
		if err != nil {
			return err
		}

		_, err = statement.ExecContext(ctx, image.Reference, link.Type, link.URL)
		statement.Close()
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Store) InsertTag(ctx context.Context, tag *models.Tag) error {
	statement, err := s.db.PrepareContext(ctx, `INSERT OR REPLACE INTO tags
		(name, color, description)
		VALUES
		(?, ?, ?);`)
	if err != nil {
		return err
	}

	_, err = statement.ExecContext(ctx, tag.Name, tag.Color, tag.Description)
	statement.Close()
	return err
}

func (s *Store) GetImage(ctx context.Context, reference string) (*models.Image, error) {
	statement, err := s.db.PrepareContext(ctx, `SELECT
	reference, latestReference, description, image, lastModified
	FROM images WHERE reference = ?;`)
	if err != nil {
		return nil, err
	}

	var image models.Image

	// TODO: Implement the scan interface for models
	res, err := statement.QueryContext(ctx, reference)
	statement.Close()
	if err != nil {
		return nil, err
	}

	if !res.Next() {
		res.Close()
		return nil, res.Err()
	}

	err = res.Scan(&image.Reference, &image.LatestReference, &image.Description, &image.Image, &image.LastModified)
	res.Close()
	if err != nil {
		return nil, err
	}

	image.Tags, err = s.GetImagesTags(ctx, reference)
	if err != nil {
		return nil, err
	}

	image.Links, err = s.GetImagesLinks(ctx, reference)
	if err != nil {
		return nil, err
	}

	return &image, nil
}

func (s *Store) GetImagesTags(ctx context.Context, reference string) ([]string, error) {
	statement, err := s.db.PrepareContext(ctx, `SELECT tag FROM images_tags WHERE reference = ?;`)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	res, err := statement.QueryContext(ctx, reference)
	if err != nil {
		return nil, err
	}

	tags := make([]string, 0)
	for res.Next() {
		var tag string
		err := res.Scan(&tag)
		if err != nil {
			res.Close()
			return nil, err
		}
		tags = append(tags, tag)
	}
	res.Close()
	if err := res.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (s *Store) GetImagesLinks(ctx context.Context, reference string) ([]models.ImageLink, error) {
	statement, err := s.db.PrepareContext(ctx, `SELECT type, url FROM images_links WHERE reference = ?;`)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	res, err := statement.QueryContext(ctx, reference)
	if err != nil {
		return nil, err
	}

	links := make([]models.ImageLink, 0)
	for res.Next() {
		var link models.ImageLink
		err := res.Scan(&link.Type, &link.URL)
		if err != nil {
			res.Close()
			return nil, err
		}
		links = append(links, link)
	}
	res.Close()
	if err := res.Err(); err != nil {
		return nil, err
	}

	return links, nil
}

func (s *Store) GetTags(ctx context.Context) ([]models.Tag, error) {
	statement, err := s.db.PrepareContext(ctx, `SELECT name, color, description FROM tags;`)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	res, err := statement.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	tags := make([]models.Tag, 0)
	for res.Next() {
		var tag models.Tag
		err := res.Scan(&tag.Name, &tag.Color, &tag.Description)
		if err != nil {
			res.Close()
			return nil, err
		}
		tags = append(tags, tag)
	}
	res.Close()
	if err := res.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (s *Store) InsertImageDescription(ctx context.Context, reference string, description *models.ImageDescription) error {
	statement, err := s.db.PrepareContext(ctx, `INSERT OR REPLACE INTO images_descriptions
		(reference, html, markdown)
		VALUES
		(?, ?, ?);`)
	if err != nil {
		return err
	}

	_, err = statement.ExecContext(ctx, reference, description.HTML, description.Markdown)
	statement.Close()
	return err
}

func (s *Store) GetImageDescription(ctx context.Context, reference string) (*models.ImageDescription, error) {
	statement, err := s.db.PrepareContext(ctx, `SELECT html, markdown FROM images_descriptions WHERE reference = ?;`)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	res, err := statement.QueryContext(ctx, reference)
	if err != nil {
		return nil, err
	}

	if !res.Next() {
		res.Close()
		return nil, res.Err()
	}

	var description models.ImageDescription
	err = res.Scan(&description.HTML, &description.Markdown)
	res.Close()
	if err != nil {
		return nil, err
	}

	return &description, nil
}

func (s *Store) InsertImageReleaseNotes(ctx context.Context, reference string, releaseNotes *models.ImageReleaseNotes) error {
	statement, err := s.db.PrepareContext(ctx, `INSERT OR REPLACE INTO images_release_notes
		(reference, title, html, markdown, released)
		VALUES
		(?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}

	_, err = statement.ExecContext(ctx, reference, releaseNotes.Title, releaseNotes.HTML, releaseNotes.Markdown, releaseNotes.Released)
	statement.Close()
	return err
}

func (s *Store) GetImageReleaseNotes(ctx context.Context, reference string) (*models.ImageReleaseNotes, error) {
	statement, err := s.db.PrepareContext(ctx, `SELECT title, html, markdown, released FROM images_release_notes WHERE reference = ?;`)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	res, err := statement.QueryContext(ctx, reference)
	if err != nil {
		return nil, err
	}

	if !res.Next() {
		res.Close()
		return nil, res.Err()
	}

	var releaseNotes models.ImageReleaseNotes
	err = res.Scan(&releaseNotes.Title, &releaseNotes.HTML, &releaseNotes.Markdown, &releaseNotes.Released)
	res.Close()
	if err != nil {
		return nil, err
	}

	return &releaseNotes, nil
}

func (s *Store) InsertImageGraph(ctx context.Context, reference string, graph *models.Graph) error {
	statement, err := s.db.PrepareContext(ctx, `INSERT OR REPLACE INTO images_graphs
		(reference, graph)
		VALUES
		(?, ?);`)
	if err != nil {
		return err
	}

	serializedGraph, err := json.Marshal(graph)
	if err != nil {
		statement.Close()
		return err
	}

	_, err = statement.ExecContext(ctx, reference, serializedGraph)
	statement.Close()
	return err
}

func (s *Store) GetImageGraph(ctx context.Context, reference string) (*models.Graph, error) {
	statement, err := s.db.PrepareContext(ctx, `SELECT graph FROM images_graphs WHERE reference = ?;`)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	res, err := statement.QueryContext(ctx, reference)
	if err != nil {
		return nil, err
	}

	if !res.Next() {
		res.Close()
		return nil, res.Err()
	}

	var serializedGraph []byte
	err = res.Scan(&serializedGraph)
	res.Close()
	if err != nil {
		return nil, err
	}

	var graph *models.Graph
	if err := json.Unmarshal(serializedGraph, &graph); err != nil {
		return nil, err
	}

	return graph, nil
}

type Order string

const (
	OrderAcending   Order = "ASC"
	OrderDescending Order = "DESC"
)

type SortProperty string

const (
	SortPropertyReference    SortProperty = "reference"
	SortPropertyLastModified SortProperty = "last_modified"
)

type ListImageOptions struct {
	// Tags defaults to nil (don't filter by tags).
	Tags []string
	// Order defaults to OrderAscending.
	Order Order
	// Page defaults to 0.
	Page int
	// Limit defaults to 30.
	Limit int
	// SortProperty defaults to SortPropertyReference.
	SortProperty SortProperty
}

func (s *Store) ListImages(ctx context.Context, options *ListImageOptions) (*models.ImagePage, error) {
	if options == nil {
		options = &ListImageOptions{}
	}

	limit := 30
	if options.Limit > 0 {
		limit = min(options.Limit, 30)
	}

	// NOTE: This mapping is done to hard code strings used in SQL queries to
	// prevent injection attacks
	sortProperty, ok := map[SortProperty]string{
		"":                       "reference",
		SortPropertyReference:    "reference",
		SortPropertyLastModified: "lastModified",
	}[options.SortProperty]
	if !ok {
		return nil, fmt.Errorf("invalid sort property")
	}

	// NOTE: This mapping is done to hard code strings used in SQL queries to
	// prevent injection attacks
	order, ok := map[Order]string{
		"":              "ASC",
		OrderAcending:   "ASC",
		OrderDescending: "DESC",
	}[options.Order]
	if !ok {
		return nil, fmt.Errorf("invalid sort property")
	}

	page := max(options.Page, 0)

	offset := page * limit

	// Total images
	res, err := s.db.QueryContext(ctx, `SELECT COUNT(1) FROM images;`)
	if err != nil {
		return nil, err
	}

	if !res.Next() {
		return nil, res.Err()
	}

	var totalImages int
	if err := res.Scan(&totalImages); err != nil {
		res.Close()
		return nil, err
	}
	res.Close()

	// Total outdated images
	res, err = s.db.QueryContext(ctx, `SELECT COUNT(1) FROM images WHERE reference != latestReference;`)
	if err != nil {
		return nil, err
	}

	if !res.Next() {
		return nil, res.Err()
	}

	var totalOutdatedImages int
	if err := res.Scan(&totalOutdatedImages); err != nil {
		res.Close()
		return nil, err
	}
	res.Close()

	var result models.ImagePage
	result.Images = make([]models.Image, 0)
	result.Summary.Images = totalImages
	result.Summary.Outdated = totalOutdatedImages
	// TODO:
	// result.Summary.Pods

	orderClause := "ORDER BY " + sortProperty + " " + order

	limitClause := "LIMIT ? OFFSET ?"

	// TODO: Support tag filter
	res, err = s.db.QueryContext(ctx, `SELECT COUNT(1) FROM images;`)
	if err != nil {
		return nil, err
	}

	if !res.Next() {
		return nil, res.Err()
	}

	var totalMatches int
	if err := res.Scan(&totalMatches); err != nil {
		res.Close()
		return nil, err
	}
	res.Close()
	result.Pagination.Total = totalMatches

	// TODO: Support tag filter
	statement, err := s.db.PrepareContext(ctx, `SELECT reference FROM images `+orderClause+" "+limitClause+";")
	if err != nil {
		return nil, err
	}

	res, err = statement.QueryContext(ctx, options.Limit, offset)
	statement.Close()
	if err != nil {
		return nil, err
	}

	for res.Next() {
		var image models.Image
		err := res.Scan(&image.Reference)
		if err != nil {
			res.Close()
			return nil, err
		}
		result.Images = append(result.Images, image)
	}
	res.Close()
	if err := res.Err(); err != nil {
		return nil, err
	}

	result.Pagination.Size = options.Limit
	result.Pagination.Page = options.Page

	for i := range result.Images {
		image, err := s.GetImage(ctx, result.Images[i].Reference)
		if err != nil {
			return nil, err
		}
		result.Images[i] = *image
	}

	return &result, nil
}

// DeleteNonPresent deletes all images that are not referenced.
func (s *Store) DeleteNonPresent(ctx context.Context, references []string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Create a temperate table that is kept throughout the transaction. This in
	// order to be able to handle any number of references, without limitation
	_, err = tx.ExecContext(ctx, "CREATE TABLE temp.present_references (reference TEXT);")
	if err != nil {
		tx.Rollback()
		return err
	}

	statement, err := tx.PrepareContext(ctx, "INSERT INTO temp.present_references (reference) VALUES (?);")
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, reference := range references {
		_, err := statement.ExecContext(ctx, reference)
		if err != nil {
			_ = statement.Close()
			_ = tx.Rollback()
			return err
		}
	}
	if err := statement.Close(); err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM raw_images WHERE reference NOT IN (SELECT reference FROM temp.present_references);")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, "DROP TABLE temp.present_references;")
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
