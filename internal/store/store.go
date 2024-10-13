package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "embed" // Embed SQL files

	"github.com/AlexGustafsson/cupdate/internal/models"
	_ "modernc.org/sqlite"
)

//go:embed createTablesIfNotExist.sql
var createTablesIfNotExist string

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
		_, err = db.Exec(createTablesIfNotExist)
		if err != nil {
			db.Close()
			return nil, err
		}
	}

	return &Store{db: db}, nil
}

func (s *Store) InsertRawImage(ctx context.Context, image *models.RawImage) error {
	tags, err := json.Marshal(image.Tags)
	if err != nil {
		return err
	}

	graph, err := json.Marshal(image.Graph)
	if err != nil {
		return err
	}

	var lastProcessed *time.Time
	if !image.LastProcessed.IsZero() {
		lastProcessed = &image.LastProcessed
	}

	statement, err := s.db.PrepareContext(ctx, `INSERT INTO raw_images
		(reference, tags, graph, lastProcessed)
		VALUES
		(?, ?, ?, ?)
		ON CONFLICT(reference) DO UPDATE SET
			tags=excluded.tags,
			graph=excluded.graph,
			lastProcessed=coalesce(excluded.lastProcessed, lastProcessed)
		;`)
	if err != nil {
		return err
	}

	_, err = statement.ExecContext(ctx, image.Reference, tags, graph, lastProcessed)
	statement.Close()
	if err != nil {
		return err
	}

	return nil
}

type ListRawImagesOptions struct {
	NotUpdatedSince time.Time
	Limit           int
}

func (s *Store) ListRawImages(ctx context.Context, options *ListRawImagesOptions) ([]models.RawImage, error) {
	if options == nil {
		options = &ListRawImagesOptions{}
	}

	limit := 30
	if options.Limit > 0 {
		limit = min(options.Limit, 30)
	}

	notUpdatedSince := time.Now()
	if !options.NotUpdatedSince.IsZero() {
		notUpdatedSince = options.NotUpdatedSince
	}

	statement, err := s.db.PrepareContext(ctx, `SELECT
	reference, tags, graph, lastProcessed
	FROM raw_images WHERE lastProcessed IS NULL OR lastProcessed < ? ORDER BY lastProcessed ASC LIMIT ?;`)
	if err != nil {
		return nil, err
	}

	// TODO: Implement the scan interface for models
	res, err := statement.QueryContext(ctx, notUpdatedSince, limit)
	statement.Close()
	if err != nil {
		return nil, err
	}

	rawImages := make([]models.RawImage, 0)
	for res.Next() {
		var rawImage models.RawImage
		var tags []byte
		var graph []byte
		var lastProcessed *time.Time
		err := res.Scan(&rawImage.Reference, &tags, &graph, &lastProcessed)
		if err != nil {
			res.Close()
			return nil, err
		}

		if err := json.Unmarshal(tags, &rawImage.Tags); err != nil {
			res.Close()
			return nil, err
		}

		if err := json.Unmarshal(graph, &rawImage.Graph); err != nil {
			res.Close()
			return nil, err
		}

		if lastProcessed != nil {
			rawImage.LastProcessed = *lastProcessed
		}

		rawImages = append(rawImages, rawImage)
	}
	res.Close()
	if err := res.Err(); err != nil {
		return nil, err
	}

	return rawImages, nil
}

func (s *Store) InsertImage(ctx context.Context, image *models.Image) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	statement, err := tx.PrepareContext(ctx, `INSERT INTO images
	(reference, latestReference, description, lastModified, imageUrl)
	VALUES
	(?, ?, ?, ?, ?)
	ON CONFLICT(reference) DO UPDATE SET
		latestReference=excluded.latestReference,
		description=excluded.description,
		lastModified=excluded.lastModified,
		imageUrl=excluded.imageUrl
	;`)
	if err != nil {
		tx.Rollback()
		return err
	}

	// TODO: Implement the scan interface for models
	_, err = statement.ExecContext(ctx, image.Reference, image.LatestReference, image.Description, image.LastModified, image.Image)
	statement.Close()
	if err != nil {
		tx.Rollback()
		return err
	}

	// TODO: Removed tags are not removed from db
	for _, tag := range image.Tags {
		statement, err := tx.PrepareContext(ctx, `INSERT INTO images_tags
		(reference, tag)
		VALUES
		(?, ?)
		ON CONFLICT(reference, tag) DO NOTHING
		;`)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = statement.ExecContext(ctx, image.Reference, tag)
		statement.Close()
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// TODO: Removed links are not removed from db
	for _, link := range image.Links {
		statement, err := tx.PrepareContext(ctx, `INSERT INTO images_links
		(reference, type, url)
		VALUES
		(?, ?, ?)
		ON CONFLICT(reference, url) DO UPDATE SET
			type=excluded.type
		;`)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = statement.ExecContext(ctx, image.Reference, link.Type, link.URL)
		statement.Close()
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *Store) InsertTag(ctx context.Context, tag *models.Tag) error {
	statement, err := s.db.PrepareContext(ctx, `INSERT INTO tags
		(name, color, description)
		VALUES
		(?, ?, ?)
		ON CONFLICT(name) DO UPDATE SET
			color=excluded.color,
			description=excluded.description
		;`)
	if err != nil {
		return err
	}

	_, err = statement.ExecContext(ctx, tag.Name, tag.Color, tag.Description)
	statement.Close()
	return err
}

func (s *Store) GetImage(ctx context.Context, reference string) (*models.Image, error) {
	statement, err := s.db.PrepareContext(ctx, `SELECT
	reference, latestReference, description, imageUrl, lastModified
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
	statement, err := s.db.PrepareContext(ctx, `INSERT INTO images_descriptions
		(reference, html, markdown)
		VALUES
		(?, ?, ?)
		ON CONFLICT(reference) DO UPDATE SET
			html=excluded.html,
			markdown=excluded.markdown
		;`)
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
	statement, err := s.db.PrepareContext(ctx, `INSERT INTO images_release_notes
		(reference, title, html, markdown, released)
		VALUES
		(?, ?, ?, ?, ?)
		ON CONFLICT(reference) DO UPDATE SET
			title=excluded.title,
			html=excluded.html,
			markdown=excluded.markdown,
			released=excluded.released
		;`)
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
	statement, err := s.db.PrepareContext(ctx, `INSERT INTO images_graphs
		(reference, graph)
		VALUES
		(?, ?)
		ON CONFLICT(reference) DO UPDATE SET
			graph=excluded.graph
		;`)
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
// Returns the number of affected rows.
func (s *Store) DeleteNonPresent(ctx context.Context, references []string) (int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	// Create a temperate table that is kept throughout the transaction. This in
	// order to be able to handle any number of references, without limitation
	_, err = tx.ExecContext(ctx, "CREATE TABLE temp.present_references (reference TEXT);")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	statement, err := tx.PrepareContext(ctx, "INSERT INTO temp.present_references (reference) VALUES (?);")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, reference := range references {
		_, err := statement.ExecContext(ctx, reference)
		if err != nil {
			_ = statement.Close()
			_ = tx.Rollback()
			return 0, err
		}
	}
	if err := statement.Close(); err != nil {
		tx.Rollback()
		return 0, err
	}

	res, err := tx.ExecContext(ctx, "DELETE FROM raw_images WHERE reference NOT IN (SELECT reference FROM temp.present_references);")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	_, err = tx.ExecContext(ctx, "DROP TABLE temp.present_references;")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
