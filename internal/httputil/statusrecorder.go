package httputil

import (
	"net/http"
)

var _ http.ResponseWriter = (*StatusRecorder)(nil)
var _ http.Flusher = (*StatusRecorder)(nil)

// StatusRecorder is a [http.ResponseWriter] and [http.Flusher] that records the
// status code.
type StatusRecorder struct {
	Writer http.ResponseWriter

	statusCode int
}

// Flush implements http.Flusher.
func (s *StatusRecorder) Flush() {
	if flusher, ok := s.Writer.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Header implements http.ResponseWriter.
func (s *StatusRecorder) Header() http.Header {
	return s.Writer.Header()
}

// Write implements http.ResponseWriter.
func (s *StatusRecorder) Write(b []byte) (int, error) {
	return s.Writer.Write(b)
}

// WriteHeader implements http.ResponseWriter.
func (s *StatusRecorder) WriteHeader(statusCode int) {
	s.statusCode = statusCode
	s.Writer.WriteHeader(statusCode)
}

// StatusCode returns the recorded status code.
func (s *StatusRecorder) StatusCode() int {
	return s.statusCode
}
