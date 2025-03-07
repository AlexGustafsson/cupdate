package otelutil

import "go.opentelemetry.io/otel/attribute"

const (
	CupdateCacheStatusHit      string = "hit"
	CupdateCacheStatusMiss     string = "miss"
	CupdateCacheStatusError    string = "error"
	CupdateCacheStatusUncached string = "uncached"
)

const CupdateCacheStatusKey = attribute.Key("cupdate.cache.status")

// CupdateWorkflowStepName returns an attribute for a cached entry's status.
func CupdateCacheStatus(status string) attribute.KeyValue {
	return CupdateCacheStatusKey.String(status)
}

const (
	CupdateWorkflowRunSpanName = "cupdate.workflow.run"
	CupdateWorkflowNameKey     = attribute.Key("cupdate.workflow.name")
)

// CupdateWorkflowStepName returns an attribute for a workflow's name.
func CupdateWorkflowName(name string) attribute.KeyValue {
	return CupdateWorkflowNameKey.String(name)
}

const (
	CupdateWorkflowJobRunSpanName = "cupdate.workflow.job.run"
	CupdateWorkflowJobNameKey     = attribute.Key("cupdate.workflow.job.name")
)

// CupdateWorkflowStepName returns an attribute for a workflow job's name.
func CupdateWorkflowJobName(name string) attribute.KeyValue {
	return CupdateWorkflowJobNameKey.String(name)
}

const (
	CupdateWorkflowStepRunSpanName     = "cupdate.workflow.step.run"
	CupdateWorkflowStepPostRunSpanName = "cupdate.workflow.step.post-run"
	CupdateWorkflowStepNameKey         = attribute.Key("cupdate.workflow.step.name")
)

// CupdateWorkflowStepName returns an attribute for a workflow step's name.
func CupdateWorkflowStepName(name string) attribute.KeyValue {
	return CupdateWorkflowStepNameKey.String(name)
}
