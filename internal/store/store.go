package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "embed" // Embed SQL files

	"github.com/AlexGustafsson/cupdate/internal/events"
	"github.com/AlexGustafsson/cupdate/internal/models"
	_ "modernc.org/sqlite"
)

//go:embed createTablesIfNotExist.sql
var createTablesIfNotExist string

//go:embed trackChanges.sql
var trackChanges string

type EventType string

const (
	EventTypeUpdated EventType = "updated"
)

type Event struct {
	Reference string
	Type      EventType
}

type Store struct {
	*events.Hub[Event]

	// mutex must be held when performing write operations
	mutex sync.Mutex
	db    *sql.DB
}

// TODO: For single rows use QueryRowContext instead of QueryContext

func New(uri string, readonly bool) (*Store, error) {
	// Use WAL to allow multiple readers
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

		_, err = db.Exec(trackChanges)
		if err != nil {
			db.Close()
			return nil, err
		}
	}

	return &Store{
		db:  db,
		Hub: events.NewHub[Event](),
	}, nil
}

func (s *Store) InsertRawImage(ctx context.Context, image *models.RawImage) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	tags, err := json.Marshal(image.Tags)
	if err != nil {
		return false, err
	}

	graph, err := json.Marshal(image.Graph)
	if err != nil {
		return false, err
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
		RETURNING lastProcessed;`)
	if err != nil {
		return false, err
	}

	res, err := statement.QueryContext(ctx, image.Reference, tags, graph, lastProcessed)
	statement.Close()
	if err != nil {
		return false, err
	}

	res.Next()
	err = res.Scan(&lastProcessed)
	res.Close()
	if err != nil {
		return false, err
	}

	_ = s.Broadcast(ctx, Event{
		Reference: image.Reference,
		Type:      EventTypeUpdated,
	})

	return lastProcessed == nil, nil
}

func (s *Store) GetRawImage(ctx context.Context, reference string) (*models.RawImage, error) {
	statement, err := s.db.PrepareContext(ctx, `SELECT
	reference, tags, graph, lastProcessed
	FROM raw_images WHERE reference = ?;`)
	if err != nil {
		return nil, err
	}

	res, err := statement.QueryContext(ctx, reference)
	statement.Close()
	if err != nil {
		return nil, err
	}

	if !res.Next() {
		return nil, fmt.Errorf("raw image not found")
	}

	var rawImage models.RawImage
	var tags []byte
	var graph []byte
	var lastProcessed *time.Time
	err = res.Scan(&rawImage.Reference, &tags, &graph, &lastProcessed)
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

	return &rawImage, nil
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
	s.mutex.Lock()
	defer s.mutex.Unlock()

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	statement, err := tx.PrepareContext(ctx, `INSERT INTO images
	(reference, created, latestReference, latestCreated, versionDiffSortable, description, lastModified, imageUrl)
	VALUES
	(?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(reference) DO UPDATE SET
		created=excluded.created,
		latestReference=excluded.latestReference,
		latestCreated=excluded.latestCreated,
		versionDiffSortable=excluded.versionDiffSortable,
		description=excluded.description,
		lastModified=excluded.lastModified,
		imageUrl=excluded.imageUrl
	;`)
	if err != nil {
		tx.Rollback()
		return err
	}

	// TODO: Implement the scan interface for models
	var latestReference *string
	if image.LatestReference != "" {
		latestReference = &image.LatestReference
	}
	_, err = statement.ExecContext(ctx, image.Reference, image.Created, latestReference, image.LatestCreated, image.VersionDiffSortable, image.Description, image.LastModified, image.Image)
	statement.Close()
	if err != nil {
		tx.Rollback()
		return err
	}

	// First clear out tags for an easy way of removing those that are no longer
	// referenced
	statement, err = tx.PrepareContext(ctx, `DELETE FROM images_tags WHERE reference = ?;`)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = statement.ExecContext(ctx, image.Reference)
	statement.Close()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Add tags
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

	// Add links
	serializedLinks, err := json.Marshal(image.Links)
	if err != nil {
		tx.Rollback()
		return err
	}

	statement, err = tx.PrepareContext(ctx, `INSERT INTO images_linksv2
		(reference, links)
		VALUES
		(?, ?)
		ON CONFLICT(reference) DO UPDATE SET
			links=excluded.links
		;`)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = statement.ExecContext(ctx, image.Reference, serializedLinks)
	statement.Close()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Add vulnerabilities
	serializedVulnerabilities, err := json.Marshal(image.Vulnerabilities)
	if err != nil {
		tx.Rollback()
		return err
	}

	statement, err = tx.PrepareContext(ctx, `INSERT INTO images_vulnerabilitiesv2
		(reference, count, vulnerabilities)
		VALUES
		(?, ?, ?)
		ON CONFLICT(reference) DO UPDATE SET
			count=excluded.count,
			vulnerabilities=vulnerabilities
		;`)
	if err != nil {
		return err
	}

	_, err = statement.ExecContext(ctx, image.Reference, len(image.Vulnerabilities), serializedVulnerabilities)
	statement.Close()
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	_ = s.Broadcast(ctx, Event{
		Reference: image.Reference,
		Type:      EventTypeUpdated,
	})

	return nil
}

func (s *Store) GetImage(ctx context.Context, reference string) (*models.Image, error) {
	statement, err := s.db.PrepareContext(ctx, `SELECT
	reference, created, latestReference, latestCreated, versionDiffSortable, description, imageUrl, lastModified
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

	var latestReference *string
	err = res.Scan(&image.Reference, &image.Created, &latestReference, &image.LatestCreated, &image.VersionDiffSortable, &image.Description, &image.Image, &image.LastModified)
	res.Close()
	if err != nil {
		return nil, err
	}

	if latestReference != nil {
		image.LatestReference = *latestReference
	}

	image.Tags, err = s.GetImagesTags(ctx, reference)
	if err != nil {
		return nil, err
	}

	image.Links, err = s.GetImagesLinks(ctx, reference)
	if err != nil {
		return nil, err
	}

	image.Vulnerabilities, err = s.GetImageVulnerabilities(ctx, reference)
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
	statement, err := s.db.PrepareContext(ctx, `SELECT links FROM images_linksv2 WHERE reference = ?;`)
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

	var serializedLinks []byte
	err = res.Scan(&serializedLinks)
	res.Close()
	if err != nil {
		return nil, err
	}

	var links []models.ImageLink
	if err := json.Unmarshal(serializedLinks, &links); err != nil {
		return nil, err
	}

	return links, nil
}

func (s *Store) GetImageVulnerabilities(ctx context.Context, reference string) ([]models.ImageVulnerability, error) {
	statement, err := s.db.PrepareContext(ctx, `SELECT vulnerabilities FROM images_vulnerabilitiesv2 WHERE reference = ?;`)
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

	var serializedVulnerabilities []byte
	err = res.Scan(&serializedVulnerabilities)
	res.Close()
	if err != nil {
		return nil, err
	}

	var vulnerabilities []models.ImageVulnerability
	if err := json.Unmarshal(serializedVulnerabilities, &vulnerabilities); err != nil {
		return nil, err
	}

	return vulnerabilities, nil
}

func (s *Store) GetTags(ctx context.Context) ([]string, error) {
	statement, err := s.db.PrepareContext(ctx, `SELECT DISTINCT tag FROM images_tags;`)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	res, err := statement.QueryContext(ctx)
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

func (s *Store) InsertImageDescription(ctx context.Context, reference string, description *models.ImageDescription) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

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
	if err != nil {
		return err
	}

	_ = s.Broadcast(ctx, Event{
		Reference: reference,
		Type:      EventTypeUpdated,
	})

	return nil
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
	s.mutex.Lock()
	defer s.mutex.Unlock()

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
	if err != nil {
		return err
	}

	_ = s.Broadcast(ctx, Event{
		Reference: reference,
		Type:      EventTypeUpdated,
	})

	return nil
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
	s.mutex.Lock()
	defer s.mutex.Unlock()

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
	if err != nil {
		return err
	}

	_ = s.Broadcast(ctx, Event{
		Reference: reference,
		Type:      EventTypeUpdated,
	})

	return nil
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
	OrderAcending   Order = "asc"
	OrderDescending Order = "desc"
)

type Sort string

const (
	SortReference Sort = "reference"
	SortBump      Sort = "bump"
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
	// Sort defaults to SortBump.
	Sort Sort
	// Query is an Sqlite full text search query.
	Query string
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
	sort, ok := map[Sort]string{
		"":            "bump",
		SortReference: "reference",
		SortBump:      "bump",
	}[options.Sort]
	if !ok {
		return nil, fmt.Errorf("invalid sort")
	}

	// NOTE: This mapping is done to hard code strings used in SQL queries to
	// prevent injection attacks
	order, ok := map[Order]string{
		"":              "",
		OrderAcending:   "ASC",
		OrderDescending: "DESC",
	}[options.Order]
	if !ok {
		return nil, fmt.Errorf("invalid order property")
	}
	if order == "" {
		// Different sorts have different orders that make sense to use as a default
		switch sort {
		case "bump":
			order = "DESC"
		default:
			order = "ASC"
		}
	}

	page := max(options.Page, 0)

	offset := page * limit

	var result models.ImagePage
	result.Images = make([]models.Image, 0)

	summary, err := s.Summary(ctx)
	if err != nil {
		return nil, err
	}
	result.Summary = *summary

	orderClause := ""
	switch sort {
	case "reference":
		orderClause = "ORDER BY images.reference " + order
	case "bump":
		orderClause = "ORDER BY images.versionDiffSortable " + order + ", images.reference"
	}

	limitClause := "LIMIT ? OFFSET ?"

	whereClause := ""
	if len(options.Tags) > 0 {
		whereClause += fmt.Sprintf("WHERE images_tags.tag IN (%s)", "?"+strings.Repeat(", ?", len(options.Tags)-1))
	}
	if options.Query != "" {
		if len(options.Tags) > 0 {
			whereClause += " AND "
		} else {
			whereClause += "WHERE "
		}
		whereClause += "images.reference IN (SELECT reference from images_fts WHERE images_fts MATCH ?)"
	}

	groupByClause := "GROUP BY images.reference"

	havingClause := ""
	if len(options.Tags) > 0 {
		havingClause = "HAVING COUNT(*) = ?"
	}

	statement, err := s.db.PrepareContext(ctx, `SELECT COUNT(1) OVER () FROM images LEFT OUTER JOIN images_tags ON images_tags.reference = images.reference `+whereClause+" "+groupByClause+" "+havingClause+";")
	if err != nil {
		return nil, err
	}

	args := make([]any, 0)
	if len(options.Tags) > 0 {
		for _, tag := range options.Tags {
			args = append(args, tag)
		}
		if options.Query != "" {
			args = append(args, ftsEscape(options.Query))
		}
		args = append(args, len(options.Tags))
	} else if options.Query != "" {
		args = append(args, ftsEscape(options.Query))
	}
	res, err := statement.QueryContext(ctx, args...)
	statement.Close()
	if err != nil {
		return nil, err
	}

	var totalMatches int
	if res.Next() {
		if err := res.Scan(&totalMatches); err != nil {
			res.Close()
			return nil, err
		}
	} else {
		if err := res.Err(); err != nil {
			res.Close()
			return nil, err
		}

		totalMatches = 0
	}
	res.Close()
	result.Pagination.Total = totalMatches

	statement, err = s.db.PrepareContext(ctx, `SELECT images.reference FROM images LEFT OUTER JOIN images_tags ON images_tags.reference = images.reference `+whereClause+" "+groupByClause+" "+havingClause+" "+orderClause+" "+limitClause+";")
	if err != nil {
		return nil, err
	}

	args = make([]any, 0)
	if len(options.Tags) > 0 {
		for _, tag := range options.Tags {
			args = append(args, tag)
		}
		if options.Query != "" {
			args = append(args, strconv.Quote(options.Query))
		}
		args = append(args, len(options.Tags))
	} else if options.Query != "" {
		args = append(args, ftsEscape(options.Query))
	}
	args = append(args, limit)
	args = append(args, offset)
	res, err = statement.QueryContext(ctx, args...)
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

	result.Pagination.Size = limit
	// Page index starts at 1
	result.Pagination.Page = page + 1

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

func (s *Store) Summary(ctx context.Context) (*models.ImagePageSummary, error) {
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
	res, err = s.db.QueryContext(ctx, `SELECT COUNT(1) FROM images WHERE latestReference IS NOT NULL AND reference != latestReference;`)
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

	// Total vulnerable images
	res, err = s.db.QueryContext(ctx, `SELECT SUM(count) FROM images_vulnerabilitiesv2;`)
	if err != nil {
		return nil, err
	}

	if !res.Next() {
		return nil, res.Err()
	}

	var totalVulnerableImages int
	if err := res.Scan(&totalVulnerableImages); err != nil {
		res.Close()
		return nil, err
	}
	res.Close()

	// Total raw images
	res, err = s.db.QueryContext(ctx, `SELECT COUNT(1) FROM raw_images;`)
	if err != nil {
		return nil, err
	}

	if !res.Next() {
		return nil, res.Err()
	}

	var totalRawImages int
	if err := res.Scan(&totalRawImages); err != nil {
		res.Close()
		return nil, err
	}
	res.Close()

	return &models.ImagePageSummary{
		Images:     totalImages,
		Outdated:   totalOutdatedImages,
		Vulnerable: totalVulnerableImages,
		Processing: totalRawImages - totalImages,
	}, nil
}

type ImageReferenceUpdate struct {
	Reference string
	Time      time.Time

	OldLatestReference string
	NewLatestReference string
}

func (s *Store) GetReferenceUpdates(ctx context.Context) ([]ImageReferenceUpdate, error) {
	res, err := s.db.QueryContext(ctx, `SELECT reference, time, oldLatestReference, newLatestReference FROM images_reference_updates;`)
	if err != nil {
		return nil, err
	}

	updates := make([]ImageReferenceUpdate, 0)
	for res.Next() {
		var update ImageReferenceUpdate
		err := res.Scan(
			&update.Reference,
			&update.Time,
			&update.OldLatestReference,
			&update.NewLatestReference,
		)
		if err != nil {
			res.Close()
			return nil, err
		}
		updates = append(updates, update)
	}
	res.Close()
	if err := res.Err(); err != nil {
		return nil, err
	}

	return updates, nil
}

type Change struct {
	Reference string
	Time      time.Time
	Type      string

	ChangedBasic           bool
	ChangedLinks           bool
	ChangedReleaseNotes    bool
	ChangedDescription     bool
	ChangedGraph           bool
	ChangedVulnerabilities bool
}

type GetChangesOptions struct {
	Reference string
	After     time.Time
	Before    time.Time
}

func (s *Store) GetChanges(ctx context.Context, options *GetChangesOptions) ([]Change, error) {
	whereClauses := make([]string, 0)
	parameters := make([]any, 0)

	if options != nil && options.Reference != "" {
		whereClauses = append(whereClauses, "reference = ?")
		parameters = append(parameters, options.Reference)
	}

	if options != nil && !options.After.IsZero() {
		whereClauses = append(whereClauses, "time >= ?")
		parameters = append(parameters, options.After.Round(time.Second))
	}

	if options != nil && !options.Before.IsZero() {
		whereClauses = append(whereClauses, "time <= ?")
		parameters = append(parameters, options.Before.Round(time.Second))
	}

	whereClause := strings.Join(whereClauses, " AND ")

	query := `SELECT reference, time, type, changedBasic, changedLinks, changedReleaseNotes, changedDescription, changedGraph, changedVulnerabilities FROM images_changes`
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	statement, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	res, err := statement.QueryContext(ctx, parameters...)
	statement.Close()
	if err != nil {
		return nil, err
	}

	updates := make([]Change, 0)
	for res.Next() {
		var update Change
		err := res.Scan(
			&update.Reference,
			&update.Time,
			&update.Type,
			&update.ChangedBasic,
			&update.ChangedLinks,
			&update.ChangedReleaseNotes,
			&update.ChangedDescription,
			&update.ChangedGraph,
			&update.ChangedVulnerabilities,
		)
		if err != nil {
			res.Close()
			return nil, err
		}
		updates = append(updates, update)
	}
	res.Close()
	if err := res.Err(); err != nil {
		return nil, err
	}

	return updates, nil
}

type DeleteChangesOptions struct {
	After  time.Time
	Before time.Time
}

func (s *Store) DeleteChanges(ctx context.Context, options *DeleteChangesOptions) error {
	whereClauses := make([]string, 0)
	parameters := make([]any, 0)

	if options != nil && !options.After.IsZero() {
		whereClauses = append(whereClauses, "time > ?")
		parameters = append(parameters, options.After)
	}

	if options != nil && !options.Before.IsZero() {
		whereClauses = append(whereClauses, "time < ?")
		parameters = append(parameters, options.Before)
	}

	whereClause := strings.Join(whereClauses, " AND ")

	query := `DELETE FROM images_changes`
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	statement, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = statement.ExecContext(ctx, parameters...)
	statement.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

// ftsEscape escapes a string for use with sqlite's full text search.
// It is not a security feature, it just ensures that all searches are full text
// and not using fts' query syntax.
func ftsEscape(s string) string {
	// The trailing * makes this a prefix search which will allow more natural
	// matches on smaller queries
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"*`
}
