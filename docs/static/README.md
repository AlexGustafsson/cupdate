# Running Cupdate with a static set of images

> [!NOTE]
> Static file support is in beta and subject to change.

In cases where a user can't or doesn't want to have Cupdate automatically
discover images in use, a static file containing OCI references can be used
instead. Running Cupdate in a container is still the only supported way, so see
also the documentation for the other platforms.

Example use cases include but are not limited to:

- Images in use are based off of base images and you wish to track them instead
  (see [#425](https://github.com/AlexGustafsson/cupdate/issues/425)).
- Images are run in an unsupported platform or Cupdate can't reach the platform,
  such as when using Hashicorp's Nomad or when running Kubernetes in a locked
  down environment.
- You want to test Cupdate without running it in a production environment.

Specify the file Cupdate should read by using the `CUPDATE_STATIC_FILE_PATH`
environment variable. Cupdate will poll the file for changes.

The file should contain a list of OCI references, one per line. If possible,
include the digest in use like so:

```text
# Lines starting with '#' are ignored, as are empty lines
rhasspy/wyoming-whisper:2.5.0@sha256:0d78ad506e450fb113616650b7328233385905e2f2ed07fa59221012144500e3
victoriametrics/victoria-metrics:v1.128.0@sha256:c27e736a8aff888cf30c4f20ec648b767358694993d87e89afc6bf80f28991da
```

An example file is included in
[../../integration/static/references.txt](../../integration/static/references.txt).

- [Configuration](../config.md)
