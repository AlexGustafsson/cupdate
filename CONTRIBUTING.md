# Contributing

We are open for contributions, big or small. If you can't code it yourself,
please feel free to open an issue or discussion on GitHub.

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
