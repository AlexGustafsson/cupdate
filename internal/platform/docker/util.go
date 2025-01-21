package docker

import (
	"github.com/AlexGustafsson/cupdate/internal/oci"
)

func getImageReference(image string, repoDigests []string) (oci.Reference, error) {
	// Parse the image as defined by the container
	containerRef, err := oci.ParseReference(image)
	if err != nil {
		return oci.Reference{}, err
	}

	// Try to parse the image as defined by the runtime, which can contain so
	// called "repo tags" and "repo digests".
	//
	// Example from a container:
	// "Image": "ghcr.io/project-zot/zot",
	// "ImageID": "sha256:f34ecd4430b1b512b8183173e31dd0f6f542b574a54cfc4e5de49df180682e57",
	//
	// Example from the resolved image:
	// "RepoDigests": [
	//   "ghcr.io/project-zot/zot-linux-arm64@sha256:5106e25775d64a613655660085f92db2ea56cc08c1fc6b840a8235dd3329c7fe",
	//   "ghcr.io/project-zot/zot@sha256:115a5eec6d9391912cd5d0b750b6b3f3886c2984e1ca5d51c4d9f430dc3c7b2e"
	// ],
	// "RepoTags": [
	// [
	//   "ghcr.io/project-zot/zot-linux-arm64:latest",
	//   "ghcr.io/project-zot/zot:latest"
	// ]
	//
	// Note that the digest is not always to the manifest list, but might as well
	// be to the actual manifest in use. It is not well-defined.
	// Note that, like in the above example, the resolved domain, path or tag
	// does not have to match the one defined in the container if multiple
	// versions have been pulled by the Docker engine

	// Try to find a repo digest that matches the way the image is specified by
	// the container
	for _, repoDigest := range repoDigests {
		digestRef, err := oci.ParseReference(repoDigest)
		if err != nil {
			continue
		}

		if digestRef.Domain == containerRef.Domain && digestRef.Path == containerRef.Path {
			// Digests lack the tag - keep it the way it was defined in the container
			if containerRef.HasTag {
				digestRef.HasTag = true
				digestRef.Tag = containerRef.Tag
			}
			return digestRef, nil
		}
	}

	// If we can't find a resolved reference, use the one from the container
	return containerRef, nil
}
