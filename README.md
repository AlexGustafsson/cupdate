<p align="center">
  <img src=".github/logo.png" alt="Logo">
</p>

# Cupdate

Cupdate is a zero-config service that helps you keep your container images
up-to-date. It automatically identifies container images in use in your
Kubernetes cluster or on your Docker host. Cupdate then identifies the latest
available version and makes this data and more available to you via a UI, API or
through an RSS feed.

Cupdate is for those who like the process of keeping their services up-to-date,
looking through what's outdated and what features new updates bring. Cupdate
will not help you deploy the updates. If you deploy your services using things
like [flux](https://github.com/fluxcd/flux2), then there are great services that
will modify your manifests for you, such as Dependabot or
[Renovate](https://github.com/renovatebot/renovate). Cupdate is not about that,
nor will it ever be. That's not to say that Cupdate won't integrate well with
such services. Cupdate can still act as a dashboard of your deployed services,
visualizing their graphs and versions. Cupdate's APIs can also be used to write
such services/scripts with ease. There's an example script in the
[cookbook](docs/cookbook/README.md).

Features:

- Zero configuration required
- Performant and lightweight - uses virtually zero CPU and roughly 14MiB RAM
- Auto-detect container images in Kubernetes and Docker
- Auto-detect the latest available container image versions
- UI for discovering updates
- Subscribe to updates via an RSS feed
- Graphs image versions' dependants explaining why they're in use
- Vulnerability scanning via Docker Scout, Quay and the
  GitHub Advisory Database through [vulndb](#vulndb)
- APIs for custom integrations

Supported registries:

- docker.io
- ghcr.io
- quay.io
- lscr.io
- registry.k8s.io, k8s.gcr.io
- registry.gitlab.com

Supported data sources:

- Docker Hub, Docker Scout
- GitHub, GitHub Container Registry
- GitLab
- Quay

## Getting started

Cupdate can be deployed using Kubernetes or Docker. It's designed to run well
with minimal required configuration. Please refer to the platform-specific
documentation for more information on how to get started with Cupdate:

- Running Cupdate using Kubernetes:
  [docs/kubernetes/README.md](docs/kubernetes/README.md)
- Running Cupdate using Docker:
  [docs/docker/README.md](docs/docker/README.md)

Cupdate can expose metrics and traces. For more information on how to use them,
see [docs/observability/README.md](docs/observability/README.md).

Although not recommended or intended, Cupdate can be run directly on host. In
that case, please build Cupdate and run it using the instructions in
[CONTRIBUTING.md](CONTRIBUTING.md).

## Screenshots

| Light mode                                                                                            | Dark mode                                                                                           |
| ----------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- |
| ![Dashboard screenshot in light mode](./docs/screenshots/dashboard-light.png)                         | ![Dashboard screenshot in dark mode](./docs/screenshots/dashboard-dark.png)                         |
| ![Dashboard screenshot on small screen in light mode](./docs/screenshots/dashboard-small-light.png)   | ![Dashboard screenshot on small screen in dark mode](./docs/screenshots/dashboard-small-dark.png)   |
| ![Image page screenshot in light mode](./docs/screenshots/image-page-light.png)                       | ![Image page screenshot in dark mode](./docs/screenshots/image-page-dark.png)                       |
| ![Image page release screenshot page in light mode](./docs/screenshots/image-page-release-light.png)  | ![Image page release screenshot in dark mode](./docs/screenshots/image-page-release-dark.png)       |
| ![Image page graph screenshot page in light mode](./docs/screenshots/image-page-graph-light.png)      | ![Image page graph screenshot in dark mode](./docs/screenshots/image-page-graph-dark.png)           |
| ![Vulnerable image page screenshot in light mode](./docs/screenshots/image-page-vulnerable-light.png) | ![Vulnerable image page screenshot in dark mode](./docs/screenshots/image-page-vulnerable-dark.png) |

## Vulndb

Vulndb is a tiny sqlite file that contains information useful to statically look
up known vulnerabilities in container images based on their source repositories.
For now it uses GitHub's advisory database.

For more information see [tools/vulndb/README.md](tools/vulndb/README.md).

The database is updated daily and published as an OCI artifact used by Cupdate.
The artifact is available here:
<https://github.com/AlexGustafsson/cupdate/pkgs/container/cupdate%2Fvulndb>.

For more advanced scanning requirements, use something like
[Trivy](https://github.com/aquasecurity/trivy) or
[Grype](https://github.com/anchore/grype).
