package httputil

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cache"
	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = (*Client)(nil)

type Client struct {
	http.Client
	cache       cache.Cache
	cacheMaxAge time.Duration

	requestsCounter  *prometheus.CounterVec
	cacheHitsCounter prometheus.Counter
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

		cache:       cache,
		cacheMaxAge: maxAge,

		requestsCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "cupdate",
			Subsystem: "http",
			Name:      "requests_total",
		}, []string{"hostname", "method", "code"}),

		cacheHitsCounter: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "cupdate",
			Subsystem: "http",
			Name:      "cache_hits_total",
		}),
	}
}

// See [http.Client.Do].
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	res, err := c.Client.Do(req)
	if err == nil {
		c.requestsCounter.WithLabelValues(req.URL.Host, req.Method, strconv.FormatInt(int64(res.StatusCode), 10)).Inc()
	}
	return res, err
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
	ctx := req.Context()

	log := slog.With(slog.String("url", req.URL.String()))
	key := c.CacheKey(req)

	// Try to read from cache, only return on successful cache reads
	entry, err := c.cache.Get(ctx, key)
	if err == nil {
		slog.Debug("HTTP response cache hit")
		c.cacheHitsCounter.Inc()
		res, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(entry)), req)
		if err == nil {
			slog.Debug("HTTP response successfully read from cache")
			c.requestsCounter.WithLabelValues(req.URL.Host, req.Method, strconv.FormatInt(int64(res.StatusCode), 10)).Inc()
			return res, nil
		} else {
			slog.Warn("HTTP request cache parse failure", slog.Any("error", err))
		}
	} else if errors.Is(err, cache.ErrNotExist) {
		log.Debug("HTTP request cache miss")
	} else {
		log.Warn("HTTP request cache lookup failure", slog.Any("error", err))
	}

	// If no entry existed or reading from the cache failed, perform the request
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	// Cache 2xx
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		log.Debug("Caching HTTP response")

		// Let's try to not be smart about streaming the result to / from cache,
		// something that was previously done. Just read the body to memory and
		// cache the response as a blob. The requests that are cached in cupdate's
		// use cases are small enough that it shouldn't be an issue
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		if err := res.Body.Close(); err != nil {
			return nil, err
		}

		// Serialize the response
		var buffer bytes.Buffer
		res.Body = io.NopCloser(bytes.NewReader(body))
		res.Write(&buffer)

		// Restore the response body
		res.Body = io.NopCloser(bytes.NewReader(body))

		err = c.cache.Set(ctx, key, buffer.Bytes(), &cache.SetEntryOptions{Expires: time.Now().Add(c.cacheMaxAge)})
		if err == nil {
			log.Debug("HTTP request was cached successfully")
		} else {
			log.Warn("HTTP response cache failure", slog.Any("error", err))
		}
	} else {
		log.Debug("Skipping HTTP response cache as status code was not 2xx", slog.Int("statusCode", res.StatusCode))
	}

	return res, nil
}

func (c *Client) CacheKey(req *http.Request) string {
	return fmt.Sprintf("httputil/v1/%s/%s", req.Method, req.URL.String())
}

// Collect implements [prometheus.Collector].
func (c *Client) Collect(ch chan<- prometheus.Metric) {
	c.requestsCounter.Collect(ch)
	c.cacheHitsCounter.Collect(ch)
}

// Describe implements [prometheus.Collector].
func (c *Client) Describe(descs chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, descs)
}
