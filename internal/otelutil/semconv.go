package otelutil

import "go.opentelemetry.io/otel/attribute"

const (
	CupdateCacheStatusHit      string = "hit"
	CupdateCacheStatusMiss     string = "miss"
	CupdateCacheStatusError    string = "error"
	CupdateCacheStatusUncached string = "uncached"
)

const (
	CupdateCacheStatusKey = attribute.Key("cupdate.cache.status")
)

func CupdateCacheStatus(status string) attribute.KeyValue {
	return CupdateCacheStatusKey.String(status)
}
