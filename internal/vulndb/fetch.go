package vulndb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/oci"
)

func Fetch(ctx context.Context, httpClient httputil.Requester, destination string) error {
	workdir, err := os.MkdirTemp(os.TempDir(), "cupdate-vulndb-oci-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(workdir)

	client := oci.Client{
		Client: httpClient,
	}

	ref, err := oci.ParseReference("ghcr.io/alexgustafsson/cupdate/vulndb:latest")
	if err != nil {
		return err
	}

	manifestBlob, err := client.GetManifestBlob(ctx, ref)
	if err != nil {
		return err
	}

	manifestBlobContents, err := io.ReadAll(manifestBlob)
	manifestBlob.Close()
	if err != nil {
		return err
	}

	var baseManifest struct {
		SchemaVersion int    `json:"schemaVersion"`
		MediaType     string `json:"mediaType"`
		ArtifactType  string `json:"artifactType"`
	}
	if err := json.Unmarshal(manifestBlobContents, &baseManifest); err != nil {
		return err
	}

	if baseManifest.SchemaVersion != 2 ||
		baseManifest.MediaType != "application/vnd.oci.image.manifest.v1+json" ||
		baseManifest.ArtifactType != "application/vnd.cupdate.vulndb.v1+json" {
		return fmt.Errorf("unsupported manifest")
	}

	var manifest struct {
		Layers []struct {
			MediaType   string            `json:"mediaType"`
			Digest      string            `json:"digest"`
			Size        int               `json:"size"`
			Annotations map[string]string `json:"annotations"`
		} `json:"layers"`
		Annotations map[string]string `json:"annotations"`
	}
	if err := json.Unmarshal(manifestBlobContents, &manifest); err != nil {
		return err
	}

	var databaseBlobDigest string
	for _, layer := range manifest.Layers {
		if layer.MediaType == "application/x-sqlite3" {
			databaseBlobDigest = layer.Digest
			break
		}
	}

	if databaseBlobDigest == "" {
		return fmt.Errorf("artifact contains no database blob")
	}

	blob, err := client.GetBlob(ctx, ref, databaseBlobDigest, false)
	if err != nil {
		return err
	}
	defer blob.Close()

	file, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if _, err := io.Copy(file, blob); err != nil {
		return err
	}

	return nil
}
