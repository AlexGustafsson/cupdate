package vulndb

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
)

func AutoFetchAndOpen(ctx context.Context, path string, httpClient *httputil.Client, maxAge time.Duration) (*Conn, error) {
	stat, err := os.Stat(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	if stat != nil && stat.IsDir() {
		return nil, fmt.Errorf("cannot set up vulndb - path exists and is a directory")
	}

	if stat == nil || time.Since(stat.ModTime()) > maxAge {
		if err := Fetch(ctx, httpClient, path); err != nil {
			return nil, err
		}
	}

	return Open(path)
}
