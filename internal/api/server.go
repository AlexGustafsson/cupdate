package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/store"
	"k8s.io/utils/strings/slices"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrBadRequest = errors.New("bad request")
)

type Server struct {
	api *store.Store
	mux *http.ServeMux
}

func NewServer(api *store.Store) *Server {
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

		tags := slices.Filter(nil, strings.Split(query.Get("tags"), ","), func(s string) bool {
			return s != ""
		})

		sort := query.Get("sort")
		if sort != "" && sort != "reference" {
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
			Tags:         tags,
			Order:        store.Order(order),
			Page:         int(page),
			Limit:        int(limit),
			SortProperty: store.SortProperty(sort),
		}
		response, err := api.ListImages(r.Context(), listOptions)
		s.handleJSONResponse(w, r, response, err)
	})

	s.mux.HandleFunc("GET /api/v1/image", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		reference := query.Get("reference")

		response, err := api.GetImage(r.Context(), reference)
		s.handleJSONResponse(w, r, response, err)
	})

	s.mux.HandleFunc("GET /api/v1/image/description", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		reference := query.Get("reference")

		response, err := api.GetImageDescription(r.Context(), reference)
		s.handleJSONResponse(w, r, response, err)
	})

	s.mux.HandleFunc("GET /api/v1/image/release-notes", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		reference := query.Get("reference")

		response, err := api.GetImageReleaseNotes(r.Context(), reference)
		s.handleJSONResponse(w, r, response, err)
	})

	s.mux.HandleFunc("GET /api/v1/image/graph", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		reference := query.Get("reference")

		response, err := api.GetImageGraph(r.Context(), reference)
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

		slog.Error("Failed to handle request", slog.Any("error", err), slog.String("path", "/tags"))
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
