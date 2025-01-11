# Running Cupdate in Docker

Cupdate is made to run in Docker by mounting the Docker socket. Cupdate will
poll the socket for changes, reacting on changes made to containers.

To get started, run the following command to run Cupdate and expose its UI and
API on port 8080.

```shell
docker run --interactive --tty --rm \
  --volume "/var/run/docker.sock:/var/run/docker.sock:ro" \
  --mount type=tmpfs,destination=/tmp \
  --env CUPDATE_DOCKER_HOST=unix:///var/run/docker.sock \
  --publish 8080:8080 \
  ghcr.io/alexgustafsson/cupdate:0.14.1
```

To more easily configure Cupdate to your liking and persisting the configuration
it is recommended to use Docker Compose. To run Cupdate with Docker Compose, run
the following command.

```shell
docker compose -f ./docs/docker/compose.yaml up
```

The Compose file is configured using best practices, but can be adapted to suite
your needs.

If you do not want to mount the Docker socket, you can use a reverse proxy.
Cupdate uses the following API paths:

- `/version`
- `/containers/json`
- `/images/{id}/json`

If you wish to inspect the source code for the image, the relevant parts can be
found in `internal/platforms/docker/platform.go`.

## Config

When running Cupdate using Docker, the Docker host needs to be specified. For
now, Cupdate only supports using the Docker socket immediately. Its path should
be configured using the `CUPDATE_DOCKER_HOST` environment variable, setting it
to `unix:///var/run/docker.sock`, for example.

By default, only running containers are processed by Cupdate. To process all
containers, running or not, set `CUPDATE_DOCKER_INCLUDE_ALL_CONTAINERS` to
`true`.

Whilst the commands above are enough to get you started with Cupdate, you might
want to change some configuration to better suite your needs. Please see the
additional documentation in [../config.md](../config.md).
