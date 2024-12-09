# Cupdate

## Configuration

Cupdate requires zero configuration, but is very configurable. Configuration is
done using environment variables.

| Environment variable       | Description                                                             | Default          |
| -------------------------- | ----------------------------------------------------------------------- | ---------------- |
| `LOG_LEVEL`                | `debug`, `info`, `warn`, `error`                                        | `info`           |
| `API_ADDRESS`              | The address to expose the API on.                                       | `0.0.0.0`        |
| `API_PORT`                 | The port to expose the API on.                                          | `8080`           |
| `WEB_DISABLED`             | Whether or not to disable the web UI.                                   | `false`          |
| `CACHE_PATH`               | A path to the boltdb file in which to store cache.                      | `cachev1.boltdb` |
| `CACHE_MAX_AGE`            | The maximum age of cache entries.                                       | `24h`            |
| `DB_PATH`                  | A path to the sqlite file in which to store data.                       | `dbv1.sqlite`    |
| `PROCESSING_INTERVAL`      | The interval between worker runs.                                       | `1h`             |
| `PROCESSING_ITEMS`         | The number of items (images) to process each worker run.                | `10`             |
| `PROCESSING_MIN_AGE`       | The minimum age of an item (image) before being processed.              | `72h`            |
| `PROCESSING_TIMEOUT`       | The maximum time one image may take to process before being terminated. | `2m`             |
| `K8S_HOST`                 | The host of the Kubernetes API. For use with proxying.                  | Required.        |
| `K8S_INCLUDE_OLD_REPLICAS` | Whether or not to include old replica sets when scraping.               | `false`          |
