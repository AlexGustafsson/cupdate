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
  ghcr.io/alexgustafsson/cupdate:0.22.2
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

If you wish to inspect the source code, the relevant parts can be found in
`internal/platforms/docker/platform.go`.

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

### TLS

By specifying a Docker host such as `https://docker.internal`, Cupdate will
automatically use TLS. However, many setups configure Docker using a self-signed
certificate chain. In such setups, Cupdate needs to know what certificates to
use. This is done by specifying the `CUPDATE_DOCKER_TLS_PATH` environment
variable to point to a directory containing the config.

Cupdate will look for the following files:

- `ca.pem` - If found, will control the trusted root CAs. May contain multiple
  certificates.
- `cert.pem` + `key.pem` - If found, will enable mTLS.

That is - if you don't want authentication, but still want TLS, it's enough to
have a `ca.pem`. Likewise, if you want mTLS but have configured trust in the
chain elsewhere, you can specify only `cert.pem` and `key.pem`.

In simple cases, like when you only have a single host, you can put the files
immediately in the specified directory.

```shell
${CUPDATE_DOCKER_TLS_PATH}
├── ca.pem
├── cert.pem
└── key.pem
```

If you want to specify certificates for each host, create a subdirectory for the
host's hostname (e.g. `docker.internal` or `192.168.0.116`) and put your files
there.

```shell
${CUPDATE_DOCKER_TLS_PATH}
├── ca.pem
├── cert.pem
├── key.pem
└── docker.internal # Overrides for the host docker.internal
    ├── ca.pem
    ├── cert.pem
    └── key.pem
```

If **no** applicable files are found in the subdirectory, the defaults, if any,
will be used.

The directory structure borrows from
[Uptime Kuma](https://github.com/louislam/uptime-kuma/wiki/How-to-Monitor-Docker-Containers)
to have a somewhat standard layout.

## Updating Cupdate

> [!NOTE]
> Before Cupdate hits v1.0.0, breaking changes can occur. Breaking changes could
> include API changes or changes to how the data is stored on disk. Breaking
> changes are communicated in release notes.

If you've installed Cupdate using the example compose file, please re-apply it
using the latest version to update Cupdate. If you've written custom manifests,
update the image version and refer to the release notes to learn if there are
additional changes required.
