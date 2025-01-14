# Contributing

We are open for contributions, big or small. The best way to contribute is to
open a bug report or feature request and testing fixes. If you want to poke
around in Cupdate's code, the rest of this document contains some basic info to
get you started.

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

Run an image using zot instead.

```shell
docker run --rm -it localhost:9090/alpine
```

Start Cupdate targeting Docker, specifying the auth file.

```shell
export CUPDATE_REGISTRY_SECRETS="integration/zot/docker-basic-auth.json"
```
