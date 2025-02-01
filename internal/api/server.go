package api

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

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

type Server struct {
	api *store.Store
	hub *events.Hub[worker.Event]
	mux *http.ServeMux

	WebAddress string
}

func NewServer(api *store.Store, hub *events.Hub[worker.Event], processQueue chan<- oci.Reference) *Server {
	s := &Server{
		api: api,
		hub: hub,
		mux: http.NewServeMux(),
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
		var limit int64 = 30
		if limitString != "" {
			var err error
			limit, err = strconv.ParseInt(limitString, 10, 64)
			if err != nil {
				s.handleGenericResponse(w, r, err)
				return
			}
		}

		listOptions := &store.ListImageOptions{
			Tags:  tags,
			Order: store.Order(order),
			Page:  int(page),
			Limit: int(limit),
			Sort:  store.Sort(sort),
			Query: query.Get("query"),
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

	s.mux.HandleFunc("POST /api/v1/image/scans", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := httputil.SpanFromRequest(r)
		span.SetAttributes(semconv.HTTPRoute("/api/v1/image/scans"))

		query := r.URL.Query()

		reference, err := oci.ParseReference(query.Get("reference"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		select {
		case <-ctx.Done():
			w.WriteHeader(http.StatusRequestTimeout)
		case processQueue <- reference:
			w.WriteHeader(http.StatusAccepted)
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
		options := &store.ListImageOptions{
			Tags: []string{"outdated"},
		}

		page, err := api.ListImages(ctx, options)
		if err != nil {
			s.handleGenericResponse(w, r, err)
			return
		}

		items := make([]rss.Item, len(page.Images))
		for i, image := range page.Images {
			if image.LatestReference == "" {
				continue
			}

			ref, err := oci.ParseReference(image.LatestReference)
			if err != nil {
				s.handleGenericResponse(w, r, err)
				return
			}

			pubDate := image.LastModified
			if image.LatestCreated != nil {
				pubDate = *image.LatestCreated
			}

			items[i] = rss.Item{
				GUID:        rss.NewDeterministicGUID(fmt.Sprintf("%s->%s", image.Reference, image.LatestReference)),
				PubDate:     rss.Time(pubDate),
				Title:       fmt.Sprintf("%s updated", ref.Name()),
				Link:        requestURL.Scheme + "://" + requestURL.Host + "/image?reference=" + url.QueryEscape(image.Reference),
				Description: fmt.Sprintf("%s updated to %s", ref.Name(), ref.Version()),
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

		for event := range s.hub.Subscribe(ctx) {
			var eventType models.EventType
			switch event.Type {
			case worker.EventTypeUpdated:
				eventType = models.EventTypeImageUpdated
			}

			data, err := json.Marshal(models.ImageEvent{
				Reference: event.Reference,
				Type:      eventType,
			})
			if err == nil {
				_, _ = fmt.Fprintf(w, "data:%s\n\n", data)
				if flusher, ok := w.(http.Flusher); ok {
					flusher.Flush()
				}
			}
		}
	})

	s.mux.HandleFunc("GET /api/v1/summary", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := httputil.SpanFromRequest(r)
		span.SetAttributes(semconv.HTTPRoute("/api/v1/summary"))

		response, err := api.Summary(ctx)
		s.handleJSONResponse(w, r, response, err)
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "traceresponse")

	httputil.InstrumentHandler(s.mux).ServeHTTP(w, r)
}
