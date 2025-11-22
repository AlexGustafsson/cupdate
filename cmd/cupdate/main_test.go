package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunLifecycle(t *testing.T) {
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	basePath := t.TempDir()

	err := os.WriteFile(filepath.Join(basePath, "references.txt"), []byte{}, 0644)
	require.NoError(t, err)

	environ := []string{
		"CUPDATE_LOG_LEVEL=debug",

		"CUPDATE_API_ADDR=127.0.0.1",
		"CUPADTE_API_PORT=",

		"CUPDATE_CACHE_PATH=" + filepath.Join(basePath, "cachev1.boltdb"),

		"CUPDATE_DB_PATH=" + filepath.Join(basePath, "dbv1.boltdb"),

		"CUPDATE_STATIC_FILE_PATH=" + filepath.Join(basePath, "references.txt"),
	}

	signals := make(chan os.Signal, 1)
	signals <- os.Interrupt

	exitCode := run(environ, signals)
	assert.Zero(t, exitCode)
}

func TestRunForcefulExit(t *testing.T) {
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	basePath := t.TempDir()

	err := os.WriteFile(filepath.Join(basePath, "references.txt"), []byte{}, 0644)
	require.NoError(t, err)

	environ := []string{
		"CUPDATE_LOG_LEVEL=debug",

		"CUPDATE_API_ADDR=127.0.0.1",
		"CUPADTE_API_PORT=",

		"CUPDATE_CACHE_PATH=" + filepath.Join(basePath, "cachev1.boltdb"),

		"CUPDATE_DB_PATH=" + filepath.Join(basePath, "dbv1.boltdb"),

		"CUPDATE_STATIC_FILE_PATH=" + filepath.Join(basePath, "references.txt"),
	}

	signals := make(chan os.Signal, 2)
	signals <- os.Interrupt
	signals <- os.Interrupt

	exitCode := run(environ, signals)
	assert.NotZero(t, exitCode)
}
