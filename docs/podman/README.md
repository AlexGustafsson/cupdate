# Running Cupdate in Podman

> [!NOTE]
> Podman support is in beta and subject to change.

Cupdate is made to run in Podman by mounting the compatibility mode Docker
socket. Cupdate will poll the socket for changes, reacting on changes made to
containers. Pods are not supported for now.

If you haven't already, you'll need to mount the Podman socket. Depending on
your setup, this requires root privileges.

```shell
systemctl enable --now podman.socket
```

Next, run the following command to run Cupdate and expose its UI and
API on port 8080.

```shell
podman run --interactive --tty --rm \
  --volume "/var/run/user/1000/podman.sock:/var/run/podman.sock:ro" \
  --mount type=tmpfs,destination=/tmp \
  --env CUPDATE_DOCKER_HOST=unix:///var/run/podman.sock \
  --publish 8080:8080 \
  ghcr.io/alexgustafsson/cupdate:0.22.2
```

As the Podman support is based on the Docker-compatible socket, the rest of the
Cupdate config is identical to when you're using Docker.

- [Configuration](../config.md)
- [Docker setup](../docker/README.md)
