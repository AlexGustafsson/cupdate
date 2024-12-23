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

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/rss"
	"github.com/AlexGustafsson/cupdate/internal/store"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrBadRequest = errors.New("bad request")
)

type Server struct {
	api *store.Store
	mux *http.ServeMux

	WebAddress string
}

func NewServer(api *store.Store, processQueue chan<- oci.Reference) *Server {
	s := &Server{
		api: api,
		mux: http.NewServeMux(),
	}

	s.mux.HandleFunc("GET /api/v1/tags", func(w http.ResponseWriter, r *http.Request) {
		tags, err := api.GetTags(r.Context())
		s.handleJSONResponse(w, r, tags, err)
	})

	s.mux.HandleFunc("GET /api/v1/images", func(w http.ResponseWriter, r *http.Request) {
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

		pageString := query.Get("page")
		var page int64 = 0
		if pageString != "" {
			var err error
			page, err = strconv.ParseInt(pageString, 10, 64)
			if err != nil {
				s.handleGenericResponse(w, r, err)
				return
			}
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
		}
		response, err := api.ListImages(r.Context(), listOptions)
		s.handleJSONResponse(w, r, response, err)
	})

	s.mux.HandleFunc("GET /api/v1/image", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		reference := query.Get("reference")

		response, err := api.GetImage(r.Context(), reference)
		if response == nil && err == nil {
			s.handleGenericResponse(w, r, ErrNotFound)
			return
		}

		s.handleJSONResponse(w, r, response, err)
	})

	s.mux.HandleFunc("GET /api/v1/image/description", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		reference := query.Get("reference")

		response, err := api.GetImageDescription(r.Context(), reference)
		if response == nil && err == nil {
			s.handleGenericResponse(w, r, ErrNotFound)
			return
		}

		s.handleJSONResponse(w, r, response, err)
	})

	s.mux.HandleFunc("GET /api/v1/image/release-notes", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		reference := query.Get("reference")

		response, err := api.GetImageReleaseNotes(r.Context(), reference)
		if response == nil && err == nil {
			s.handleGenericResponse(w, r, ErrNotFound)
			return
		}

		s.handleJSONResponse(w, r, response, err)
	})

	s.mux.HandleFunc("GET /api/v1/image/graph", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		reference := query.Get("reference")

		response, err := api.GetImageGraph(r.Context(), reference)
		s.handleJSONResponse(w, r, response, err)
	})

	s.mux.HandleFunc("POST /api/v1/image/scans", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		reference, err := oci.ParseReference(query.Get("reference"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		select {
		case <-r.Context().Done():
			w.WriteHeader(http.StatusRequestTimeout)
		case processQueue <- reference:
		}

		w.WriteHeader(http.StatusAccepted)
	})

	s.mux.HandleFunc("GET /api/v1/feed.rss", func(w http.ResponseWriter, r *http.Request) {
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

		page, err := api.ListImages(r.Context(), options)
		if err != nil {
			s.handleGenericResponse(w, r, err)
			return
		}

		items := make([]rss.Item, len(page.Images))
		for i, image := range page.Images {
			ref, err := oci.ParseReference(image.LatestReference)
			if err != nil {
				s.handleGenericResponse(w, r, err)
				return
			}

			items[i] = rss.Item{
				GUID: rss.NewDeterministicGUID(image.Reference),
				// TODO: Use image update time instead
				PubDate:     rss.Time(image.LastModified),
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

		slog.Error("Failed to handle request", slog.Any("error", err), slog.String("path", r.URL.Path))
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
		slog.Error("Failed to write response", slog.Any("error", err))
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	s.mux.ServeHTTP(w, r)
}
