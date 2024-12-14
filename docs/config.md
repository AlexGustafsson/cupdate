# Cupdate

## Configuration

Cupdate requires zero configuration, but is very configurable. Configuration is
done using environment variables.

| Environment variable                      | Description                                                                                                           | Default                     |
| ----------------------------------------- | --------------------------------------------------------------------------------------------------------------------- | --------------------------- |
| `CUPDATE_LOG_LEVEL`                       | `debug`, `info`, `warn`, `error`                                                                                      | `info`                      |
| `CUPDATE_API_ADDRESS`                     | The address to expose the API on.                                                                                     | `0.0.0.0`                   |
| `CUPDATE_API_PORT`                        | The port to expose the API on.                                                                                        | `8080`                      |
| `CUPDATE_WEB_DISABLED`                    | Whether or not to disable the web UI.                                                                                 | `false`                     |
| `CUPDATE_WEB_ADDRESS`                     | The URL at which the UI is available (such as `https://example.com`). Used for RSS feeds, should generally not be set | Automatically resolved      |
| `CUPDATE_CACHE_PATH`                      | A path to the boltdb file in which to store cache.                                                                    | `cachev1.boltdb`            |
| `CUPDATE_CACHE_MAX_AGE`                   | The maximum age of cache entries.                                                                                     | `24h`                       |
| `CUPDATE_DB_PATH`                         | A path to the sqlite file in which to store data.                                                                     | `dbv1.sqlite`               |
| `CUPDATE_PROCESSING_INTERVAL`             | The interval between worker runs.                                                                                     | `1h`                        |
| `CUPDATE_PROCESSING_ITEMS`                | The number of items (images) to process each worker run.                                                              | `10`                        |
| `CUPDATE_PROCESSING_MIN_AGE`              | The minimum age of an item (image) before being processed.                                                            | `72h`                       |
| `CUPDATE_PROCESSING_TIMEOUT`              | The maximum time one image may take to process before being terminated.                                               | `2m`                        |
| `CUPDATE_KUBERNETES_HOST`                 | The host of the Kubernetes API. For use with proxying.                                                                | Required to use Kubernetes. |
| `CUPDATE_KUBERNETES_INCLUDE_OLD_REPLICAS` | Whether or not to include old replica sets when scraping.                                                             | `false`                     |
| `CUPDATE_DOCKER_HOST`                     | Docker host address.                                                                                                  | Required to use Docker.     |
| `CUPDATE_DOCKER_INCLUDE_ALL_CONTAINERS`   | Whether or not to include containers in any state, not just running containers.                                       | `false`                     |
