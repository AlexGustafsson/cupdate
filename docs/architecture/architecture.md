# Architecture

![A simplified overview of the architecture](overview.excalidraw.png)

Cupdate discovers container images that are in use in a _platform_. Next,
Cupdate discovers new versions for these container images in their respective
OCI registry. Lastly data is enriched from sources like Docker Hub and GitHub,
depending on the information gathered about the image from the registry.

## Cupdate

![An overview of the parts that constitute Cupdate](cupdate.excalidraw.png)

## Platforms

### Kubernetes

![An overview of how Cupdate uses Kubernetes](kubernetes.excalidraw.png)

When running in Kubernetes, Cupdate lists and then watches all resources that
references an image. Resources such as pods directly refences an image that is
in use. Resources like deployments reference images through pod templates.

### Docker

![An overview of how Cupdate uses Docker](docker.excalidraw.png)

When using Docker, Cupdate uses `docker.sock` directly to identify images and
containers using those images.
