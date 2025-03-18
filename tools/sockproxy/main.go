package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"maps"
	"net"
	"net/http"
	"os"
)

func main() {
	port := flag.Int("p", 3000, "port to serve on")
	flag.Parse()

	path := flag.Arg(0)
	if path == "" {
		fmt.Fprintf(os.Stderr, "missing required path")
		os.Exit(1)
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, "unix", path)
			},
		},
	}

	err := http.ListenAndServe(fmt.Sprintf("localhost:%d", *port), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequest(r.Method, fmt.Sprintf("http://_%s?%s", r.URL.Path, r.URL.RawQuery), r.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to proxy request: %v", err)
			os.Exit(1)
		}

		req.Header = r.Header

		res, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to proxy request: %v", err)
			os.Exit(1)
		}
		defer res.Body.Close()

		header := w.Header()
		maps.Copy(header, res.Header)

		w.WriteHeader(res.StatusCode)
		if _, err := io.Copy(w, res.Body); err != nil {
			fmt.Fprintf(os.Stderr, "failed to proxy request: %v", err)
			os.Exit(1)
		}
	}))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to serve: %v", err)
		os.Exit(1)
	}
}
