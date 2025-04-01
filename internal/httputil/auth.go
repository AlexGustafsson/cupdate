package httputil

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// AuthHandler implements request authentication.
type AuthHandler interface {
	// HandleAuth authenticates a request.
	HandleAuth(*http.Request) error
}

type AuthHandlerFunc func(*http.Request) error

func (f AuthHandlerFunc) HandleAuth(r *http.Request) error {
	return f(r)
}

var _ AuthHandler = (*BasicAuthHandler)(nil)

// BasicAuthHandler auths requests using a username/password via the Basic
// authorization scheme.
type BasicAuthHandler struct {
	Username string
	Password string
}

func (h BasicAuthHandler) HandleAuth(r *http.Request) error {
	credentials := base64.StdEncoding.EncodeToString([]byte(h.Username + ":" + h.Password))
	r.Header.Set("Authorization", "Basic "+credentials)

	return nil
}

// AuthMux is an HTTP auth multiplexer. It matches URLs of auth requests against
// a list of registered patterns and calls the handler for the pattern that most
// closely matches the request.
type AuthMux struct {
	mutex    sync.RWMutex
	patterns map[string]AuthHandler
	header   http.Header
}

func NewAuthMux() *AuthMux {
	return &AuthMux{
		patterns: make(map[string]AuthHandler),
		header:   make(http.Header),
	}
}

// Handle registers handler for pattern.
func (a *AuthMux) Handle(pattern string, handler AuthHandler) {
	a.register(pattern, handler)
}

// Handle registers handler for pattern.
func (a *AuthMux) HandleFunc(pattern string, handler func(*http.Request) error) {
	a.register(pattern, AuthHandlerFunc(handler))
}

// HandleAuth implements [AuthHandler.HandleAuth].
func (a *AuthMux) HandleAuth(r *http.Request) error {
	a.mutex.RLock()

	handler := a.match(r.URL)
	if handler == nil {
		handler = a.patterns[""]
	}

	for k, v := range a.header {
		r.Header[k] = v
	}

	a.mutex.RUnlock()

	if handler == nil {
		return nil
	}

	return handler.HandleAuth(r)
}

// SetHeader sets a header to set on all requests.
func (a *AuthMux) SetHeader(key string, value string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.header.Set(key, value)
}

func (a *AuthMux) register(pattern string, handler AuthHandler) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.patterns[pattern] = handler
}

func (a *AuthMux) match(url *url.URL) AuthHandler {
	// TODO: This code is not especially nice or efficient. For example, it parses
	// the URLs every iteration
	for pattern, handler := range a.patterns {
		if pattern == "" {
			continue
		}

		// Compare scheme, if set in pattern - otherwise allow any scheme
		if !strings.HasPrefix(pattern, "https://") && !strings.HasPrefix(pattern, "http://") {
			pattern = url.Scheme + "://" + pattern
		}

		u, err := url.Parse(pattern)
		if err != nil {
			continue
		}

		if url.Scheme != u.Scheme {
			continue
		}

		// The Docker client matches either the image or hostname, where the API can
		// have a /v2/ or /v1/ prefix
		if strings.HasPrefix(u.Path, "/v2/") || strings.HasPrefix(u.Path, "/v1/") {
			u.Path = u.Path[3:]
		}

		// If the pattern has a path specified, match it
		if u.Path != "" && u.Path != "/" {
			p := url.Path
			if !strings.HasSuffix(p, "/") {
				p += "/"
			}

			// Make sure paths in patterns match full segments
			if !strings.HasSuffix(u.Path, "/") {
				u.Path += "/"
			}

			if !strings.HasPrefix(p, u.Path) {
				continue
			}
		}

		urlParts := strings.Split(url.Host, ".")
		patternParts := strings.Split(u.Host, ".")

		// The pattern has more parts of its hostname than the URL
		if len(urlParts) < len(patternParts) {
			continue
		}

		matched := true
		for i := 0; i < len(patternParts); i++ {
			if strings.HasPrefix(patternParts[i], "*") {
				if !strings.HasSuffix(urlParts[i], patternParts[i][1:]) {
					matched = false
					break
				}
			} else if strings.HasSuffix(patternParts[i], "*") {
				// TODO:
			} else if urlParts[i] != patternParts[i] {
				matched = false
				break
			}
		}

		if matched {
			return handler
		}
	}

	return nil
}

// Copy copies all patterns from another [AuthMux] to this one.
func (a *AuthMux) Copy(other *AuthMux) {
	if other == nil {
		return
	}

	a.mutex.Lock()
	defer a.mutex.Unlock()

	other.mutex.RLock()
	defer other.mutex.RUnlock()

	for pattern, handler := range other.patterns {
		a.patterns[pattern] = handler
	}

	for k, v := range other.header {
		a.header[k] = v
	}
}
