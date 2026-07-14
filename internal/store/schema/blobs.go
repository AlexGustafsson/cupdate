package schema

import "time"

type GraphBlob struct {
	Edges map[string]map[string]bool `json:"edges"`
	Nodes map[string]GraphNode       `json:"nodes"`
}

type GraphNode struct {
	Domain         string            `json:"domain"`
	Type           string            `json:"type"`
	Name           string            `json:"name"`
	Labels         map[string]string `json:"labels,omitempty"`
	InternalLabels map[string]string `json:"internalLabels,omitempty"`
}

type AnnotationsBlob map[string]string

type TagsBlob []string
type LinksBlob []string

type WorkflowRunResult string

const (
	WorkflowRunResultSucceeded WorkflowRunResult = "succeeded"
	WorkflowRunResultFailed    WorkflowRunResult = "failed"
)

type WorkflowRunBlob struct {
	TraceID         string            `json:"traceId"`
	Started         time.Time         `json:"started"`
	DurationSeconds float64           `json:"duration"`
	Result          WorkflowRunResult `json:"result"`
	Jobs            []JobRun          `json:"jobs"`
}

type JobRunResult string

const (
	JobRunResultSucceeded JobRunResult = "succeeded"
	JobRunResultSkipped   JobRunResult = "skipped"
	JobRunResultFailed    JobRunResult = "failed"
)

type JobRun struct {
	Result          JobRunResult `json:"result"`
	Steps           []StepRun    `json:"steps"`
	DependsOn       []string     `json:"dependsOn"`
	JobID           string       `json:"jobId,omitempty"`
	JobName         string       `json:"jobName,omitempty"`
	Started         time.Time    `json:"started,omitempty"`
	DurationSeconds float64      `json:"duration,omitempty"`
}

type StepRunResult string

const (
	StepRunResultSucceeded StepRunResult = "succeeded"
	StepRunResultSkipped   StepRunResult = "skipped"
	StepRunResultFailed    StepRunResult = "failed"
)

type StepRun struct {
	Result          StepRunResult `json:"result"`
	StepName        string        `json:"stepName,omitempty"`
	Started         time.Time     `json:"started,omitempty"`
	DurationSeconds float64       `json:"duration,omitempty"`
	Error           string        `json:"error,omitempty"`
}
