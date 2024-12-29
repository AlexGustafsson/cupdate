package oci

import (
	"context"
	"os"
	"path/filepath"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	oras "oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"
)

func PushArtifact(ctx context.Context, path string, username string, password string) error {
	workdir, err := os.MkdirTemp(os.TempDir(), "cupdate-vulndb-oci-*")
	if err != nil {
		return err
	}

	fs, err := file.New(workdir)
	if err != nil {
		return err
	}
	defer fs.Close()

	dbPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	mediaType := "application/x-sqlite3"
	fileDescriptor, err := fs.Add(ctx, path, mediaType, dbPath)
	if err != nil {
		return err
	}
	fileDescriptors := []v1.Descriptor{fileDescriptor}

	artifactType := "application/vnd.cupdate.vulndb.v1+json"
	opts := oras.PackManifestOptions{
		Layers: fileDescriptors,
	}
	manifestDescriptor, err := oras.PackManifest(ctx, fs, oras.PackManifestVersion1_1, artifactType, opts)
	if err != nil {
		return err
	}

	manifestDescriptor.Annotations["org.opencontainers.image.source"] = "https://github.com/AlexGustafsson/cupdate"
	manifestDescriptor.Annotations["org.opencontainers.image.description"] = `Cupdate's vulnerability database.`

	tag := "latest"
	if err = fs.Tag(ctx, manifestDescriptor, tag); err != nil {
		return err
	}

	reg := "ghcr.io"
	repo, err := remote.NewRepository(reg + "/alexgustafsson/cupdate-vulndb")
	if err != nil {
		return err
	}

	repo.Client = &auth.Client{
		Client: retry.DefaultClient,
		Cache:  auth.NewCache(),
		Credential: auth.StaticCredential(reg, auth.Credential{
			Username: username,
			Password: password,
		}),
	}

	_, err = oras.Copy(ctx, fs, tag, repo, tag, oras.DefaultCopyOptions)
	return err
}
