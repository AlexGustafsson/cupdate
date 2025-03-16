# Contributing

The best way to contribute to Cupdate is to by using and testing it. If you face
issues, open a bug report and be ready to test fixes. If you have ideas for new
features, open a feature request.

If you want to poke around in Cupdate's code, the rest of this document contains
some basic info to get you started. If you end up wanting to contribute a
feature, please discuss the feature in a feature request issue first so that we
can make sure that the feature aligns with Cupdate's scope (see README) and that
no time is wasted developing features that might end up not being merged.

## Architecture

Cupdate is written in Go and comes with a frontend written using TypeScript,
React and Tailwind. For more information see [ARCHITECTURE.md](ARCHITECTURE.md).

## Building

Cupdate can be built on host using yarn and go, or inside of a container using
Docker.

### Building in Docker

Build Cupdate for running inside of a container.

```shell
docker build --tag ghcr.io/alexgustafsson/cupdate:latest .
```

Build Cupdate inside a container for running on host.

```shell
docker buildx build --target=export --output=. .
```

Build Cupdate inside a container for running on the specified platform.

```shell
docker buildx build --platform macos/arm64  --target=export --output=. .
```

### Building on host

```shell
yarn install
yarn build
go build -o cupdate cmd/cupdate/*.go
```

## Running

Cupdate supports both Kubernetes and Docker as the target platforms. Typically
Cupdate will run inside of these environments, but for development it can run
on a host and communicate the the platforms' APIs remotely.

### Preparing for Kubernetes

Proxy the Kubernetes API server.

```shell
kubectl proxy
```

Source the default dev config for the Kubernetes platform.

```shell
# Inspect
cat .env-kubernetes

# Bash etc.
source .env-kubernetes

# Fish
export (cat .env-kubernetes | xargs -L 1)
```

### Preparing for Docker

Symlink the Docker socket.

```shell
# NOTE: The path might be different on your machine
ln -s ~/.colima/default/docker.sock docker.sock
```

```shell
kubectl proxy
```

Source the default dev config for the Docker platform.

```shell
# Inspect
cat .env-docker

# Bash etc.
source .env-docker

# Fish
export (cat .env-docker | xargs -L 1)
```

### Running Cupdate

Start Cupdate.

```shell
go run cmd/cupdate/*.go
```

Optionally start the development web server for frontend development.

```shell
yarn run dev
```

Optionally use Jaeger for otel testing.

```shell
docker run --rm -it \
  -p 4317:4317 \
  -p 8081:16686 \
  jaegertracing/all-in-one

# NOTE: Start Cupdate with the required additional config
export CUPDATE_OTEL_TARGET=localhost:4317
export CUPDATE_OTEL_INSECURE=true
```

Optionally proxy a Docker socket to test Docker over TCP. Use the proxied port
as the Docker host rather then the one specified in `.env-docker`.

```shell
go run tools/sockproxy/*.go -p 3000 docker.sock
```

### Testing custom registry

To test custom registries and authentication, Zot can be used.

```shell
# Create a htpasswd for zot
htpasswd -bBn username password > integration/zot/htpasswd
```

```shell
docker run --rm -it -p 9090:9090 --volume "$PWD/integration/zot:/etc/zot:ro" ghcr.io/project-zot/zot-linux-arm64
```

Note that Zot's UI doesn't work on Safari ATM - you will just be logged out if
you log in.

Upload a simple image to Zot.

```shell
skopeo copy --override-os linux --dest-tls-verify=false docker://alpine:latest docker://localhost:9090/test/alpine:latest
```

Run the image from Zot.

```shell
docker run --rm -it localhost:9090/alpine
```

Start Cupdate targeting Docker, specifying the auth file.

```shell
export CUPDATE_REGISTRY_SECRETS="integration/zot/docker-basic-auth.json"
```

### Writing and running unit tests

Some tests directly use APIs on the internet, for "system tests". These tests
are by their nature flakey. As such, they don't run in the CI.

These tests are identified by their naming convention, `TestIntegration...` and
by the fact that they start by bailing if `-short` is specified when running the
tests.

As Cupdate has a lot of HTTP clients, there's a framework for writing table
tests for HTTP APIs used throughout tests. This framework should allow for a
near 100% test coverage in these clients. Additional tests may use the APIs
directly, as stated previously, and in these cases the tests may print data for
additional, manual, verification.

Tests are run by using go:

```shell
# Run all unit tests that are run in the CI
go test -race -short -v ./...

# Run all tests, even those using external APIs
go test -race -v ./...

# Run specific tests
go test -race -short -v ./internal/openssf/scorecard/...

# Collect coverage
go test -coverprofile coverage.out -race -v ./...

# Show coverage on a web page
go tool cover -html coverage.out
```

Some tests use containers to test integration with services. These should run
just fine, but on macOS, when using Colima, you'll have to specify the following
environment variables:

```shell
export TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE=/var/run/docker.sock
export DOCKER_HOST="unix://$HOME/.colima/docker.sock"
```
