# Cupdate

## Configuration

Cupdate requires zero configuration, but is very configurable. Configuration is
done using environment variables.

| Environment variable                    | Description                                                                                                           | Default                         |
| --------------------------------------- | --------------------------------------------------------------------------------------------------------------------- | ------------------------------- |
| `CUPDATE_LOG_LEVEL`                     | `debug`, `info`, `warn`, `error`                                                                                      | `info`                          |
| `CUPDATE_API_ADDRESS`                   | The address to expose the API on.                                                                                     | `0.0.0.0`                       |
| `CUPDATE_API_PORT`                      | The port to expose the API on.                                                                                        | `8080`                          |
| `CUPDATE_WEB_DISABLED`                  | Whether or not to disable the web UI.                                                                                 | `false`                         |
| `CUPDATE_WEB_ADDRESS`                   | The URL at which the UI is available (such as `https://example.com`). Used for RSS feeds, should generally not be set | Automatically resolved          |
| `CUPDATE_HTTP_USER_AGENT`               | The User Agent string to use for HTTP requests.                                                                       | `Cupdate/1.0`                   |
| `CUPDATE_CACHE_PATH`                    | A path to the boltdb file in which to store cache.                                                                    | `cachev1.boltdb`                |
| `CUPDATE_CACHE_MAX_AGE`                 | The maximum age of cache entries.                                                                                     | `24h`                           |
| `CUPDATE_DB_PATH`                       | A path to the sqlite file in which to store data.                                                                     | `dbv1.sqlite`                   |
| `CUPDATE_PROCESSING_INTERVAL`           | The interval between worker runs.                                                                                     | `1h`                            |
| `CUPDATE_PROCESSING_ITEMS`              | The number of items (images) to process each worker run.                                                              | `10`                            |
| `CUPDATE_PROCESSING_MIN_AGE`            | The minimum age of an item (image) before being processed.                                                            | `72h`                           |
| `CUPDATE_PROCESSING_TIMEOUT`            | The maximum time one image may take to process before being terminated.                                               | `2m`                            |
| `CUPDATE_PROCESSING_QUEUE_BURST`        | Number of items that can be processed in a short burst.                                                               | `10`                            |
| `CUPDATE_PROCESSING_QUEUE_RATE`         | The desired processing rate under normal circumstances.                                                               | `1m`                            |
| `CUPDATE_WORKFLOW_CLEANUP_MAX_AGE`      | The maximum age of a workflow run before it's removed.                                                                | `48h`                           |
| `CUPDATE_WORKFLOW_CLEANUP_INTERVAL`     | The time between workflow run cleanup iterations.                                                                     | `1h`                            |
| `CUPDATE_KUBERNETES_HOST`               | The host of the Kubernetes API. For use with proxying.                                                                | Required to use Kubernetes.     |
| `CUPDATE_DOCKER_HOST`                   | One or more comma-separated Docker host URIs. Supports unix://path and tcp://host:port URIs.                          | Required to use Docker.         |
| `CUPDATE_DOCKER_INCLUDE_ALL_CONTAINERS` | Whether or not to include containers in any state, not just running containers.                                       | `false`                         |
| `CUPDATE_OTEL_TARGET`                   | Target URL to an Open Telemetry GRPC ingest endpoint.                                                                 | Required to use Open Telemetry. |
| `CUPDATE_OTEL_INSECURE`                 | Disable client transport security for the Open Telemetry GRPC connection.                                             | `false`                         |
| `CUPDATE_REGISTRY_SECRETS`              | Path to a JSON file containing registry secrets. See Docker's config.json and Kubernetes' `imagePullSecrets`.         | None                            |

### Labels

Cupdate can take additional resource-specific configuration via the use of
labels. For Docker, this means annotating the container image or the container /
service itself. In Kubernetes, any resource in the image's tree can be annotated
to configure Cupdate.

Each label has two aliases to follow both the Docker and Kubernetes conventions.

The web UI can be used to see what labels are used for a given image on the
image's page, in the graph view. Nodes that have labels configured will have an
icon / tooltip indicating that the behavior might differ from the defaults. The
exact labels used can be seen by clicking on a node, which shows a dialog.

#### `ignore`

- Kubernetes: `config.cupdate/ignore`
- Docker: `cupdate.config.ignore`

Set to `true` to ignore the resource subtree (e.g. deployment, pod, or container
). Defaults to `false`.

#### `stay-on-current-major`

- Kubernetes: `config.cupdate/stay-on-current-major`
- Docker: `cupdate.config.stay-on-current-major`

Set to `true` to stay on the current major track specified by images' semantic
tags.

Examples of updates made by Cupdate by default:

- `alpine:3.21.2` -> `alpine:3.21.3` (patch on current major track)
- `node:22.14.0` -> `node:23.8.0` (end of current major track, new major available)

With `stay-on-current-major` set to `true`, Cupdate wouldn't recommend node to
be updated as the newer version is no longer on the same major.

This is useful for databases where updating to a newer major version can take a
lot of time and require changes to applications.

#### `pin`

- Kubernetes: `config.cupdate/pin`
- Docker: `cupdate.config.pin`

Set to `true` to pin images' tags, meaning Cupdate will only recommend updates
to the underlying manifest identified by the tag, as opposed to semantic
updates.

Examples of updates made by Cupdate by default:

- `alpine:3.21.2` -> `alpine:3.21.3` (patch on current major track)
- `node:22.14.0` -> `node:23.8.0` (end of current major track, new major available)

With `pin` set to `true`, neither update would be recommended. However, Cupdate
would still recommend updates like the following, if their tags have been
overwritten, pointing to a newer manifest.

- `latest` -> `latest`
- `alpine:3` -> `alpine:3`

This is useful for databases where updating to a newer version can take a lot of
time and require changes to applications. It is recommended to only use `pin` if
necessary, preferring the use of `stay-on-current-major` alongside a tag that
specifies an as granular version as possible (i.e. `alpine:3.21.2` as opposed to
`alpine:3`). But in cases where the tag is not necessarily semantic, such as
where `12.0.0` -> `12.1.0` would mean a major change, `pin` can be used to force
the tag to not change.

#### Examples

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: cupdate
    config.cupdate/ignore: "true"
# ...
```

```yaml
# compose.yaml
services:
  cupdate:
    labels:
      - cupdate.config.ignore: "true"
```
