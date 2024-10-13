package httputil

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cache"
)

type callbackCloser struct {
	io.Reader
	callback func() error
}

func (c callbackCloser) Close() error {
	return c.callback()
}

type Client struct {
	http.Client
	Cache       cache.Cache
	CacheMaxAge time.Duration
}

func NewClient(cache cache.Cache, maxAge time.Duration) *Client {
	return &Client{
		Client: http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout: 5 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout: 5 * time.Second,
			},
			Timeout: 10 * time.Second,
		},
		Cache:       cache,
		CacheMaxAge: maxAge,
	}
}

// DoCached returns a cached response for the request.
// If no cache entry exists, [http.Client.Do] is used. If the request succeeds,
// its response is cached if the response code is 2xx.
// It is the caller's responsibility to ensure that caching the request is
// sensible (i.e. only for GET requests).
// NOTE: Cached responses for URLs that were redirected will not have the
// correct request URL for the response - it will be the original request rather
// than the request to the final resource.
func (c *Client) DoCached(req *http.Request) (*http.Response, error) {
	if c.Cache == nil {
		panic("no cache configured")
	}

	ctx := req.Context()

	log := slog.With(slog.String("url", req.URL.String()))
	key := c.CacheKey(req)

	// Try to read from cache, only return on successful cache reads
	entry, err := c.Cache.Get(ctx, key)
	if err == nil {
		log.Debug("HTTP request cache hit")

		modTime := time.Time{}
		if entryInfo, ok := entry.(cache.EntryInfo); ok {
			modTime = entryInfo.ModTime()
		}

		outdated := !modTime.IsZero() && c.CacheMaxAge > 0 && time.Since(modTime) > c.CacheMaxAge
		if outdated {
			slog.Debug("HTTP request cache miss (entry found, but was outdated)")
			if err := entry.Close(); err != nil {
				slog.Warn("Failed to close HTTP response cache entry", slog.Any("error", err))
			}
		} else {
			slog.Debug("Reading cached response")
			res, err := http.ReadResponse(bufio.NewReader(entry), req)
			if err == nil {
				slog.Debug("HTTP response successfully read from cache")
				// TODO: is entry ever closed in this branch?
				return res, nil
			} else {
				slog.Warn("HTTP request cache read failure", slog.Any("error", err))
				if err := entry.Close(); err != nil {
					slog.Warn("Failed to close HTTP response cache entry", slog.Any("error", err))
				}
			}
		}
	} else if errors.Is(err, cache.ErrNotExist) {
		log.Debug("HTTP request cache miss")
	} else {
		log.Warn("HTTP request cache lookup failure", slog.Any("error", err))
	}

	// If no entry existed or reading from the cache failed, perform the request
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	// Cache 2xx
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		// Create a pipe, writing to cache as the body is read by the caller as to
		// not read the entire body into memory just for the sake of the cache
		originalBody := res.Body
		r, w := io.Pipe()
		wait := make(chan struct{})
		res.Body = callbackCloser{
			Reader: io.TeeReader(originalBody, w),
			callback: func() error {
				err := originalBody.Close()
				w.Close()
				select {
				case <-wait:
				case <-ctx.Done():
					return ctx.Err()
				}
				return err
			},
		}

		log.Debug("Storing HTTP cache entry once body is read")
		go func() {
			resCopy := http.Response{
				StatusCode:       res.StatusCode,
				ProtoMajor:       res.ProtoMajor,
				ProtoMinor:       res.ProtoMinor,
				Request:          res.Request,
				TransferEncoding: res.TransferEncoding,
				Trailer:          res.Trailer,
				Body:             io.NopCloser(r),
				ContentLength:    res.ContentLength,
				Header:           res.Header,
			}

			resReader, resWriter := io.Pipe()

			go func() {
				err := c.Cache.Set(ctx, key, resReader)
				if err == nil {
					log.Debug("HTTP request was cached successfully")
				} else {
					log.Warn("HTTP response cache failure", slog.Any("error", err))
				}
				close(wait)
			}()

			err := resCopy.Write(resWriter)
			resWriter.Close()
			if err == nil {
				log.Debug("HTTP response was successfully serialized")
			} else {
				log.Warn("HTTP response serialization failure", slog.Any("error", err))
			}
		}()
	} else {
		log.Debug("Skipping HTTP response cache as status code was not 2xx", slog.Int("statusCode", res.StatusCode))
	}

	return res, nil
}

func (c *Client) CacheKey(req *http.Request) string {
	return fmt.Sprintf("httputil/v1/%s/%s", req.Method, req.URL.String())
}
