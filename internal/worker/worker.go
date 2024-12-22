package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"github.com/AlexGustafsson/cupdate/internal/semver"
	"github.com/AlexGustafsson/cupdate/internal/store"
	"github.com/AlexGustafsson/cupdate/internal/workflow/imageworkflow"
	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = (*Worker)(nil)

type Worker struct {
	httpClient *httputil.Client
	store      *store.Store

	processedCounter   prometheus.Counter
	processingDuration prometheus.Counter
	processingGauge    prometheus.Gauge
}

func New(httpClient *httputil.Client, store *store.Store) *Worker {
	return &Worker{
		httpClient: httpClient,
		store:      store,

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

func (w *Worker) ProcessRawImage(ctx context.Context, reference oci.Reference) error {
	start := time.Now()
	w.processingGauge.Inc()
	defer w.processingGauge.Dec()

	image, err := w.store.GetRawImage(ctx, reference.String())
	if err != nil {
		return err
	}

	log := slog.With(slog.String("reference", reference.String()))
	log.Debug("Processing reference")

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
	if err := w.store.InsertRawImage(ctx, image); err != nil {
		return err
	}

	log.Debug("Running workflow")
	data := &imageworkflow.Data{
		ImageReference:  reference,
		Image:           "",
		LatestReference: nil,
		Tags:            make([]string, 0),
		Description:     "",
		FullDescription: nil,
		ReleaseNotes:    nil,
		Links:           make([]models.ImageLink, 0),
		Vulnerabilities: make([]models.ImageVulnerability, 0),
		Graph:           image.Graph,
	}

	for _, tag := range image.Tags {
		data.InsertTag(tag)
	}

	workflow := imageworkflow.New(w.httpClient, data)
	if err := workflow.Run(ctx); err != nil {
		log.Error("Failed to run pipeline for image", slog.Any("error", err))
		// Fallthrough - insert what we have
	}

	// If no new version was defined and the current version is using the "latest"
	// convention, the latest available reference is the current reference
	if reference.Version() == "latest" && data.LatestReference == nil {
		r := reference
		data.LatestReference = &r
	}

	// Add some basic tags
	if data.LatestReference != nil {
		if data.ImageReference.String() == data.LatestReference.String() {
			data.Tags = append(data.Tags, "up-to-date")
		} else {
			data.Tags = append(data.Tags, "outdated")
		}

		// Add tags based on version diff
		currentVersion, currentVersionErr := semver.ParseVersion(data.ImageReference.Version())
		newVersion, newVersionErr := semver.ParseVersion(data.LatestReference.Version())
		if currentVersion != nil && currentVersionErr == nil && newVersion != nil && newVersionErr == nil {
			diff := currentVersion.Diff(newVersion)
			if diff != "" {
				data.Tags = append(data.Tags, diff)
			}
		}
	}

	result := models.Image{
		Reference:       data.ImageReference.String(),
		LatestReference: "",
		Description:     data.Description,
		Tags:            data.Tags,
		Image:           data.Image,
		Links:           data.Links,
		Vulnerabilities: data.Vulnerabilities,
		LastModified:    time.Now(),
	}
	if data.LatestReference != nil {
		result.LatestReference = data.LatestReference.String()
	}
	if err := w.store.InsertImage(context.TODO(), &result); err != nil {
		log.Error("Failed to insert image", slog.Any("error", err))
		// Fallthrough - try to insert what we have
	}

	if data.FullDescription != nil {
		if err := w.store.InsertImageDescription(ctx, reference.String(), data.FullDescription); err != nil {
			log.Error("Failed to insert image description", slog.Any("error", err))
			// Fallthrough - try to insert what we have
		}
	}

	if data.ReleaseNotes != nil {
		if err := w.store.InsertImageReleaseNotes(ctx, reference.String(), data.ReleaseNotes); err != nil {
			log.Error("Failed to insert image description", slog.Any("error", err))
			// Fallthrough - try to insert what we have
		}
	}

	if err := w.store.InsertImageGraph(ctx, reference.String(), &data.Graph); err != nil {
		log.Error("Failed to insert image description", slog.Any("error", err))
		// Fallthrough - try to insert what we have
	}

	log.Debug("Updated data")
	w.processedCounter.Inc()
	w.processingDuration.Add(time.Since(start).Seconds())
	return nil
}

// Collect implements [prometheus.Collector].
func (w *Worker) Collect(ch chan<- prometheus.Metric) {
	w.processedCounter.Collect(ch)
	w.processingDuration.Collect(ch)
	w.processingGauge.Collect(ch)
}

// Describe implements [prometheus.Collector].
func (w *Worker) Describe(descs chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(w, descs)
}
