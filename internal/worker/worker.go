package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/events"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/osv"
	"github.com/AlexGustafsson/cupdate/internal/platform/docker"
	"github.com/AlexGustafsson/cupdate/internal/platform/kubernetes"
	"github.com/AlexGustafsson/cupdate/internal/semver"
	"github.com/AlexGustafsson/cupdate/internal/store"
	"github.com/AlexGustafsson/cupdate/internal/workflow/imageworkflow"
	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = (*Worker)(nil)

// EventType is the type of an event.
type EventType string

const (
	// EventTypeUpdated is emitted whenever data of an image is updated.
	EventTypeUpdated EventType = "updated"
	// EventTypeProcessed is emitted whenever an image was processed.
	EventTypeProcessed EventType = "processed"
	// EventTypeNewVersionAvailable is emitted whenever the latest available
	// version of an image changes.
	EventTypeNewVersionAvailable EventType = "newVersionAvailable"
)

// Event describes a Worker event.
type Event struct {
	Reference string
	Type      EventType
}

// Worker processes raw container image entries, running the image workflow and
// storing the result to the state store.
// The worker produces events of the type [Event].
type Worker struct {
	*events.Hub[Event]

	httpClient   httputil.Requester
	store        *store.Store
	registryAuth *httputil.AuthMux

	processedCounter   prometheus.Counter
	processingDuration prometheus.Counter
	processingGauge    prometheus.Gauge
}

func New(httpClient httputil.Requester, store *store.Store, registryAuth *httputil.AuthMux) *Worker {
	return &Worker{
		Hub: events.NewHub[Event](),

		httpClient:   httpClient,
		store:        store,
		registryAuth: registryAuth,

		processedCounter: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "cupdate",
			Subsystem: "worker",
			Name:      "processed_images_total",
		}),
		processingDuration: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "cupdate",
			Subsystem: "worker",
			Name:      "processed_images_duration_seconds",
		}),
		processingGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "cupdate",
			Subsystem: "worker",
			Name:      "processing_images",
		}),
	}
}

// ProcessRawImage processes a raw image by the specified reference.
func (w *Worker) ProcessRawImage(ctx context.Context, reference oci.Reference) error {
	start := time.Now()
	w.processingGauge.Inc()
	defer w.processingGauge.Dec()

	image, err := w.store.GetRawImage(ctx, reference.String())
	if err != nil {
		return err
	}

	log := slog.With(slog.String("reference", reference.String()))
	log.DebugContext(ctx, "Processing reference")

	// Try to update the image's process time
	// NOTE: There's a race here if the entry has been modified or removed since
	// it was loaded from the store. It will eventually be corrent and consistent,
	// though. And it's unlikely to happen. So let's not keep a transaction during
	// processing for now. If it becomes important, we could keep an "etag" /
	// generation id in the document and throw an error if the expectation fails.
	// NOTE: Always update immediately as a failure to process or update the image
	// could be a reoccuring issue, so try to process other images before retrying
	// the failing image.
	image.LastProcessed = time.Now()
	if _, err := w.store.InsertRawImage(ctx, image); err != nil {
		return err
	}

	log.DebugContext(ctx, "Running workflow")
	data := &imageworkflow.Data{
		ImageReference:  reference,
		Image:           "",
		LatestReference: nil,
		Tags:            make([]string, 0),
		Description:     "",
		Links:           make([]models.ImageLink, 0),
		Vulnerabilities: make([]models.ImageVulnerability, 0),
		Graph:           image.Graph,
		RegistryAuth:    w.registryAuth,
	}

	for _, tag := range image.Tags {
		data.InsertTag(tag)
	}

	workflow := imageworkflow.New(w.httpClient, data)
	workflowRun, err := workflow.Run(ctx)
	if err != nil {
		ociJobFailed := workflowRun.Jobs[0].Result == models.JobRunResultFailed
		if ociJobFailed {
			log.ErrorContext(ctx, "Failed to run image workflow", slog.Any("error", err))
			data.InsertTag("status:failed")
		} else {
			log.WarnContext(ctx, "Failed to run a subset of image workflow", slog.Any("error", err))
			data.InsertTag("status:incomplete")
		}
		// Fallthrough - insert what we have
	}

	versionDiffSortable := semver.PackInt64(nil)

	// Add some basic tags
	if data.LatestReference != nil {
		if data.ImageReference.String() == data.LatestReference.String() {
			data.InsertTag("up-to-date")
		} else {
			data.InsertTag("outdated")
			// We know that the image is outdated, default to assuming the update is a
			// patch to handle cases like where digests diff but not the tag itself
			versionDiffSortable = semver.PackedSingleDigitPatchDiff
		}

		currentVersion, currentVersionErr := semver.ParseVersion(data.ImageReference.Version())
		newVersion, newVersionErr := semver.ParseVersion(data.LatestReference.Version())
		if currentVersion != nil && currentVersionErr == nil && newVersion != nil && newVersionErr == nil {
			diff := currentVersion.Diff(newVersion)
			if diff != "" {
				data.InsertTag("bump:" + diff)
				versionDiffSortable = semver.PackInt64(newVersion) - semver.PackInt64(currentVersion)
			}
		}
	}

	// Add Kubernetes namespace and Docker stack tags
	for _, node := range image.Graph.Nodes {
		switch node.Domain {
		case "kubernetes":
			if node.Type == kubernetes.ResourceKindCoreV1Namespace {
				data.InsertTag("namespace:" + node.Name)
			}
		case "docker":
			if node.Type == docker.ResourceKindSwarmNamespace {
				data.InsertTag("namespace:" + node.Name)
			} else if node.Type == docker.ResourceKindComposeProject {
				data.InsertTag("project:" + node.Name)
			}
		}
	}

	// Add risk task based on OpenSSF score
	if data.Scorecard.OK && data.Scorecard.Value != nil {
		// Don't add a label for low risk components
		if data.Scorecard.Value.Risk != models.ImageScorecardRiskLow {
			data.InsertTag("risk:" + string(data.Scorecard.Value.Risk))
		}
	}

	// Add vulnerable label if there are specified severities
	var maxSeverity osv.NormalizedSeverity
	for _, vulnerability := range data.Vulnerabilities {
		severity := vulnerability.NormalizedSeverity()
		if severity.Compare(maxSeverity) < 0 {
			maxSeverity = severity
		}
	}

	if maxSeverity != "" {
		switch maxSeverity {
		case osv.NormalizedSeverityCritical:
			data.InsertTag("vulnerability:critical")
		case osv.NormalizedSeverityHigh:
			data.InsertTag("vulnerability:high")
		case osv.NormalizedSeverityMedium:
			data.InsertTag("vulnerability:medium")
		case osv.NormalizedSeverityLow:
			data.InsertTag("vulnerability:low")
		default:
			data.InsertTag("vulnerability:unspecified")
		}
	}

	timeBeforeInsert := time.Now()

	result := models.Image{
		Reference:           data.ImageReference.String(),
		Created:             data.Created,
		LatestReference:     "",
		LatestCreated:       data.LatestCreated,
		VersionDiffSortable: versionDiffSortable,
		Description:         data.Description,
		Tags:                data.Tags,
		Image:               data.Image,
		Links:               data.Links,
		Vulnerabilities:     len(data.Vulnerabilities),
		LastModified:        time.Now(),
	}
	if data.LatestReference != nil {
		result.LatestReference = data.LatestReference.String()
	}
	if err := w.store.InsertImage(context.TODO(), &result); err != nil {
		log.ErrorContext(ctx, "Failed to insert image", slog.Any("error", err))
		// Fallthrough - try to insert what we have
	}

	if data.FullDescription.OK {
		if data.FullDescription.Value == nil {
			if err := w.store.DeleteImageDescription(ctx, reference.String()); err != nil {
				log.ErrorContext(ctx, "Failed to delete image description", slog.Any("error", err))
				// Fallthrough - try to insert what we have
			}
		} else {
			if err := w.store.InsertImageDescription(ctx, reference.String(), data.FullDescription.Value); err != nil {
				log.ErrorContext(ctx, "Failed to insert image description", slog.Any("error", err))
				// Fallthrough - try to insert what we have
			}
		}
	}

	if data.ReleaseNotes.OK {
		if data.ReleaseNotes.Value == nil {
			if err := w.store.DeleteImageReleaseNotes(ctx, reference.String()); err != nil {
				log.ErrorContext(ctx, "Failed to delete image release notes", slog.Any("error", err))
				// Fallthrough - try to insert what we have
			}
		} else {
			if err := w.store.InsertImageReleaseNotes(ctx, reference.String(), data.ReleaseNotes.Value); err != nil {
				log.ErrorContext(ctx, "Failed to insert image release notes", slog.Any("error", err))
				// Fallthrough - try to insert what we have
			}
		}
	}

	if err := w.store.InsertImageGraph(ctx, reference.String(), &data.Graph); err != nil {
		log.ErrorContext(ctx, "Failed to insert image graph", slog.Any("error", err))
		// Fallthrough - try to insert what we have
	}

	if data.Scorecard.OK {
		if data.Scorecard.Value == nil {
			if err := w.store.DeleteImageScorecard(ctx, reference.String()); err != nil {
				log.ErrorContext(ctx, "Failed to delete image scorecard", slog.Any("error", err))
				// Fallthrough - try to insert what we have
			}
		} else {
			if err := w.store.InsertImageScorecard(ctx, reference.String(), data.Scorecard.Value); err != nil {
				log.ErrorContext(ctx, "Failed to insert image scorecard", slog.Any("error", err))
				// Fallthrough - try to insert what we have
			}
		}
	}

	if data.Provenance.OK {
		if data.Provenance.Value == nil {
			if err := w.store.DeleteImageProvenance(ctx, reference.String()); err != nil {
				log.ErrorContext(ctx, "Failed to delete image provenance", slog.Any("error", err))
				// Fallthrough - try to insert what we have
			}
		} else {
			if err := w.store.InsertImageProvenance(ctx, reference.String(), data.Provenance.Value); err != nil {
				log.ErrorContext(ctx, "Failed to insert image provenance", slog.Any("error", err))
				// Fallthrough - try to insert what we have
			}
		}
	}

	if data.SBOM.OK {
		if data.SBOM.Value == nil {
			if err := w.store.DeleteImageSBOM(ctx, reference.String()); err != nil {
				log.ErrorContext(ctx, "Failed to delete image SBOM", slog.Any("error", err))
				// Fallthrough - try to insert what we have
			}
		} else {
			if err := w.store.InsertImageSBOM(ctx, reference.String(), data.SBOM.Value); err != nil {
				log.ErrorContext(ctx, "Failed to insert image SBOM", slog.Any("error", err))
				// Fallthrough - try to insert what we have
			}
		}
	}

	if err := w.store.InsertImageVulnerabilities(ctx, reference.String(), data.Vulnerabilities); err != nil {
		log.ErrorContext(ctx, "Failed to insert image vulnerabilities", slog.Any("error", err))
		// Fallthrough - try to insert what we have
	}

	if err := w.store.InsertWorkflowRun(ctx, reference.String(), workflowRun); err != nil {
		log.ErrorContext(ctx, "Failed to insert workflow run", slog.Any("error", err))
		// Fallthrough - try to insert what we have
	}

	timeAfterInsert := time.Now()

	log.DebugContext(ctx, "Processed image")
	w.processedCounter.Inc()
	w.processingDuration.Add(time.Since(start).Seconds())

	// Try to identify what changed
	changes, err := w.store.GetChanges(ctx, &store.GetChangesOptions{
		Reference: reference.String(),
		After:     timeBeforeInsert,
		Before:    timeAfterInsert,
	})
	if err != nil {
		log.ErrorContext(ctx, "Failed to identify changes", slog.Any("error", err))
	} else if len(changes) > 0 {
		log.DebugContext(ctx, "Updated image date", slog.Int("changes", len(changes)))
		// TODO: Group changes, create an event specifying the time. That way the
		// browser can ignore the event if it already updated after the time?
		w.Broadcast(ctx, Event{
			Reference: reference.String(),
			Type:      EventTypeUpdated,
		})

		// TODO: Have another readonly job for going over the changes made to
		// references to identify updates every now and then for third-party alerts.
		// For now, just do it on the RSS field? Perhaps try to use the change time
		// as the article time if the time of release is not found.
		// TODO: Instead of readonly job, just watch the events instead?
	}

	if result.LatestReference != "" && result.LatestReference != result.Reference {
		w.Broadcast(ctx, Event{
			Reference: reference.String(),
			Type:      EventTypeNewVersionAvailable,
		})
	}

	w.Broadcast(ctx, Event{
		Reference: reference.String(),
		Type:      EventTypeProcessed,
	})

	return nil
}

// Collect implements prometheus.Collector.
func (w *Worker) Collect(ch chan<- prometheus.Metric) {
	w.processedCounter.Collect(ch)
	w.processingDuration.Collect(ch)
	w.processingGauge.Collect(ch)
}

// Describe implements prometheus.Collector.
func (w *Worker) Describe(descs chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(w, descs)
}
