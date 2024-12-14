# Docker

Cupdate is made to run in Docker by mounting the Docker socket.

```shell
docker run --detach --volume /var/run/docker.sock:/var/run/docker.sock ghcr.io/alexgustafsson/cupdate:latest
```

See also the Docker compose example in [compose.yaml](compose.yaml).

```shell
docker compose up -f ./compose.yaml
```

## Config

For Docker, Cupdate requires the Docker host to be specified. For now, Cupdate
only supports using the Docker socket immediately. It should be configured using
the `CUPDATE_DOCKER_HOST` environment variable, setting it to
`unix:///var/run/docker.sock`, for example.

By default, only running containers are used by Cupdate. To use all containers,
set `CUPDATE_DOCKER_INCLUDE_ALL_CONTAINERS` to `true`.

See also [config.md](../config.md).
