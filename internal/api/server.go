package api

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/events"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/rss"
	"github.com/AlexGustafsson/cupdate/internal/store"
	"github.com/AlexGustafsson/cupdate/internal/worker"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrBadRequest = errors.New("bad request")
)

type Platform interface {
	Graph(context.Context) error
}

type Server struct {
	api         *store.Store
	platform    Platform
	workerHub   *events.Hub[worker.Event]
	platformHub *events.Hub[models.PlatformEvent]
	logoProxy   LogoProxy
	mux         *http.ServeMux

	WebAddress string
}

func NewServer(
	api *store.Store,
	workerHub *events.Hub[worker.Event],
	platformHub *events.Hub[models.PlatformEvent],
	processQueue *worker.Queue[oci.Reference],
	logoProxy LogoProxy,
	platform Platform,
) *Server {
	s := &Server{
		api:         api,
		platform:    platform,
		workerHub:   workerHub,
		platformHub: platformHub,
		logoProxy:   logoProxy,
		mux:         http.NewServeMux(),
	}

	s.mux.HandleFunc("GET /api/v1/tags", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := httputil.SpanFromRequest(r)
		span.SetAttributes(semconv.HTTPRoute("/api/v1/tags"))

		tags, err := api.GetTags(ctx)
		s.handleJSONResponse(w, r, tags, err)
	})

	s.mux.HandleFunc("GET /api/v1/images", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := httputil.SpanFromRequest(r)
		span.SetAttributes(semconv.HTTPRoute("/api/v1/images"))

		query := r.URL.Query()

		tags, ok := query["tag"]
		if !ok {
			tags = make([]string, 0)
		}

		tagOperator := store.TagOperatorAnd
		switch query.Get("tagop") {
		case "and":
			tagOperator = store.TagOperatorAnd
		case "or":
			tagOperator = store.TagOperatorOr
		}

		sort := query.Get("sort")
		if sort != "" && sort != "reference" && sort != "bump" {
			s.handleGenericResponse(w, r, ErrBadRequest)
			return
		}

		order := query.Get("order")
		if order != "" && order != "desc" && order != "asc" {
			s.handleGenericResponse(w, r, ErrBadRequest)
			return
		}

		// Parse the page index, if given
		pageString := query.Get("page")
		var page int64 = 0
		if pageString != "" {
			var err error
			page, err = strconv.ParseInt(pageString, 10, 64)
			if err != nil {
				s.handleGenericResponse(w, r, err)
				return
			}

			// Page index starts at 1
			if page < 1 {
				s.handleGenericResponse(w, r, ErrBadRequest)
				return
			}
			page -= 1
		}

		limitString := query.Get("limit")
		var limit int = 30
		if limitString != "" {
			var err error
			l, err := strconv.ParseInt(limitString, 10, 32)
			if err != nil {
				s.handleGenericResponse(w, r, err)
				return
			}
			limit = int(l)
		}

		listOptions := &store.ListImageOptions{
			Tags:        tags,
			TagOperator: tagOperator,
			Order:       store.Order(order),
			Page:        int(page),
			Limit:       int(limit),
			Sort:        store.Sort(sort),
			Query:       query.Get("query"),
		}

		response, err := api.ListImages(ctx, listOptions)
		s.handleJSONResponse(w, r, response, err)
	})

	s.mux.HandleFunc("GET /api/v1/image", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := httputil.SpanFromRequest(r)
		span.SetAttributes(semconv.HTTPRoute("/api/v1/image"))

		query := r.URL.Query()

		reference := query.Get("reference")

		response, err := api.GetImage(ctx, reference)
		if response == nil && err == nil {
			s.handleGenericResponse(w, r, ErrNotFound)
			return
		}

		s.handleJSONResponse(w, r, response, err)
	})

	s.mux.HandleFunc("GET /api/v1/image/description", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := httputil.SpanFromRequest(r)
		span.SetAttributes(semconv.HTTPRoute("/api/v1/image/description"))

		query := r.URL.Query()

		reference := query.Get("reference")

		response, err := api.GetImageDescription(ctx, reference)
		if response == nil && err == nil {
			s.handleGenericResponse(w, r, ErrNotFound)
			return
		}

		s.handleJSONResponse(w, r, response, err)
	})

	s.mux.HandleFunc("GET /api/v1/image/release-notes", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := httputil.SpanFromRequest(r)
		span.SetAttributes(semconv.HTTPRoute("/api/v1/image/release-notes"))

		query := r.URL.Query()

		reference := query.Get("reference")

		response, err := api.GetImageReleaseNotes(ctx, reference)
		if response == nil && err == nil {
			s.handleGenericResponse(w, r, ErrNotFound)
			return
		}

		s.handleJSONResponse(w, r, response, err)
	})

	s.mux.HandleFunc("GET /api/v1/image/graph", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := httputil.SpanFromRequest(r)
		span.SetAttributes(semconv.HTTPRoute("/api/v1/image/graph"))

		query := r.URL.Query()

		reference := query.Get("reference")

		response, err := api.GetImageGraph(ctx, reference)
		s.handleJSONResponse(w, r, response, err)
	})

	s.mux.HandleFunc("GET /api/v1/image/scorecard", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := httputil.SpanFromRequest(r)
		span.SetAttributes(semconv.HTTPRoute("/api/v1/image/scorecard"))

		query := r.URL.Query()

		reference := query.Get("reference")

		response, err := api.GetImageScorecard(ctx, reference)
		s.handleJSONResponse(w, r, response, err)
	})

	s.mux.HandleFunc("GET /api/v1/image/provenance", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := httputil.SpanFromRequest(r)
		span.SetAttributes(semconv.HTTPRoute("/api/v1/image/provenance"))

		query := r.URL.Query()

		reference := query.Get("reference")

		response, err := api.GetImageProvenance(ctx, reference)
		s.handleJSONResponse(w, r, response, err)
	})

	s.mux.HandleFunc("GET /api/v1/image/sbom", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := httputil.SpanFromRequest(r)
		span.SetAttributes(semconv.HTTPRoute("/api/v1/image/sbom"))

		query := r.URL.Query()

		reference := query.Get("reference")

		response, err := api.GetImageSBOM(ctx, reference)
		s.handleJSONResponse(w, r, response, err)
	})

	s.mux.HandleFunc("GET /api/v1/image/vulnerabilities", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := httputil.SpanFromRequest(r)
		span.SetAttributes(semconv.HTTPRoute("/api/v1/image/vulnerabilities"))

		query := r.URL.Query()

		reference := query.Get("reference")

		vulnerabilities, err := api.GetImageVulnerabilities(ctx, reference)
		if err != nil {
			s.handleJSONResponse(w, r, nil, err)
			return
		}

		response := struct {
			Vulnerabilities []models.ImageVulnerability `json:"vulnerabilities"`
		}{
			Vulnerabilities: vulnerabilities,
		}
		s.handleJSONResponse(w, r, response, err)
	})

	s.mux.HandleFunc("POST /api/v1/image/scans", func(w http.ResponseWriter, r *http.Request) {
		_, span := httputil.SpanFromRequest(r)
		span.SetAttributes(semconv.HTTPRoute("/api/v1/image/scans"))

		query := r.URL.Query()

		reference, err := oci.ParseReference(query.Get("reference"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		processQueue.PushFront(reference)
		w.WriteHeader(http.StatusAccepted)
	})

	// NOTE: For now, there's no use case of exposing multiple workflows, but
	// let's have room for it in the APIs
	s.mux.HandleFunc("GET /api/v1/image/workflows/latest", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := httputil.SpanFromRequest(r)
		span.SetAttributes(semconv.HTTPRoute("/api/v1/image/workflows/latest"))

		query := r.URL.Query()

		reference := query.Get("reference")

		response, err := api.GetLatestWorkflowRun(ctx, reference)
		s.handleJSONResponse(w, r, response, err)
	})

	s.mux.HandleFunc("GET /api/v1/image/logo", func(w http.ResponseWriter, r *http.Request) {
		_, span := httputil.SpanFromRequest(r)
		span.SetAttributes(semconv.HTTPRoute("/api/v1/image/logo"))

		query := r.URL.Query()

		reference := query.Get("reference")

		ref, err := oci.ParseReference(reference)
		if err != nil {
			s.handleGenericResponse(w, r, ErrBadRequest)
			return
		}

		err = s.logoProxy.ServeLogo(w, r, ref)
		if err == ErrNotFound {
			// Ask user agents to cache that the image was not found for a few minutes
			w.Header().Set("Cache-Control", "max-age=300")
			s.handleGenericResponse(w, r, ErrNotFound)
		} else if err != nil {
			s.handleGenericResponse(w, r, err)
			return
		}
	})

	s.mux.HandleFunc("GET /api/v1/feed.rss", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := httputil.SpanFromRequest(r)
		span.SetAttributes(semconv.HTTPRoute("/api/v1/feed.rss"))

		var requestURL *url.URL
		var err error
		if s.WebAddress == "" {
			requestURL, err = httputil.ResolveRequestURL(r)
		} else {
			requestURL, err = url.Parse(s.WebAddress)
		}
		if err != nil {
			s.handleGenericResponse(w, r, ErrBadRequest)
			return
		}

		// TODO: When we support other sort properties (like latest release), sort
		// by that
		// TODO: We currently use the default count. IIRC, it's good practice in RSS
		// to return just the latest ~20 items.
		options := &store.GetUpdateOptions{
			Limit: 20,
		}

		updates, err := api.GetUpdates(ctx, options)
		if err != nil {
			s.handleGenericResponse(w, r, err)
			return
		}

		items := make([]rss.Item, len(updates))
		for i, update := range updates {
			// TODO: Use annotations to identify some human readable version like the
			// UI, as opposed to using the reference immediately if it's just a shaid
			newRef, err := oci.ParseReference(update.NewReference)
			if err != nil {
				s.handleGenericResponse(w, r, err)
				return
			}

			oldRef, err := oci.ParseReference(update.OldReference)
			if err != nil {
				s.handleGenericResponse(w, r, err)
				return
			}

			pubDate := update.Identified
			if update.Released != nil {
				pubDate = *update.Released
			}

			items[i] = rss.Item{
				GUID:        rss.NewDeterministicGUID(update.NewReference),
				PubDate:     rss.Time(pubDate),
				Title:       fmt.Sprintf("%s updated", newRef.Name()),
				Link:        requestURL.Scheme + "://" + requestURL.Host + "/image?reference=" + url.QueryEscape(update.OldReference),
				Description: fmt.Sprintf("%s updated to %s from %s", newRef.Name(), newRef.Version(), oldRef.Version()),
			}
		}

		feed := rss.Feed{
			Version: "2.0",
			Channels: []rss.Channel{
				{
					Title:       "Cupdate",
					Link:        requestURL.Scheme + "://" + requestURL.Host,
					Description: "Container images discovered by Cupdate",
					Items:       items,
				},
			},
		}

		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)

		encoder := xml.NewEncoder(w)
		encoder.Indent("", "\t")

		if _, err := w.Write([]byte(xml.Header)); err != nil {
			return
		}
		if err := encoder.Encode(&feed); err != nil {
			return
		}
	})

	s.mux.HandleFunc("GET /api/v1/events", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := httputil.SpanFromRequest(r)
		span.SetAttributes(semconv.HTTPRoute("/api/v1/events"))

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)

		workerEvents := s.workerHub.Subscribe(ctx)
		platformEvents := s.platformHub.Subscribe(ctx)
		for {
			var data []byte
			var err error
			select {
			case <-r.Context().Done():
				return
			case event := <-workerEvents:
				var eventType models.EventType
				switch event.Type {
				case worker.EventTypeUpdated:
					eventType = models.EventTypeImageUpdated
				case worker.EventTypeProcessed:
					eventType = models.EventTypeImageProcessed
				case worker.EventTypeNewVersionAvailable:
					eventType = models.EventTypeImageNewVersionAvailable
				}

				data, err = json.Marshal(models.ImageEvent{
					Reference: event.Reference,
					Type:      eventType,
				})
			case event := <-platformEvents:
				data, err = json.Marshal(event)
			}

			if err == nil {
				_, _ = fmt.Fprintf(w, "data:%s\n\n", data)
				if flusher, ok := w.(http.Flusher); ok {
					flusher.Flush()
				}
			} else {
				slog.Warn("Failed to marshal event", slog.Any("error", err))
			}
		}
	})

	s.mux.HandleFunc("GET /api/v1/summary", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := httputil.SpanFromRequest(r)
		span.SetAttributes(semconv.HTTPRoute("/api/v1/summary"))

		response, err := api.Summary(ctx)
		s.handleJSONResponse(w, r, response, err)
	})

	s.mux.HandleFunc("POST /api/v1/images/poll", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := httputil.SpanFromRequest(r)
		span.SetAttributes(semconv.HTTPRoute("/api/v1/images/poll"))

		err := platform.Graph(ctx)
		s.handleGenericResponse(w, r, err)
	})

	return s
}

func (s *Server) handleGenericResponse(w http.ResponseWriter, r *http.Request, err error) bool {
	if err == ErrBadRequest {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return false
	} else if err == ErrNotFound {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return false
	} else if err != nil {
		// The request was likely just aborted
		if err == r.Context().Err() {
			http.Error(w, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
			return false
		}

		slog.ErrorContext(r.Context(), "Failed to handle request", slog.Any("error", err), slog.String("path", r.URL.Path))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return false
	}

	return true
}

func (s *Server) handleJSONResponse(w http.ResponseWriter, r *http.Request, response any, err error) {
	ok := s.handleGenericResponse(w, r, err)
	if !ok {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.ErrorContext(r.Context(), "Failed to write response", slog.Any("error", err))
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") && !strings.Contains(r.Header.Get("Accept"), "text/event-stream") {
		w.Header().Set("Content-Encoding", "gzip")
		gzip := &httputil.GzipWriter{ResponseWriter: w}
		defer gzip.Close()

		w = gzip
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "traceresponse")

	httputil.InstrumentHandler(s.mux).ServeHTTP(w, r)
}
