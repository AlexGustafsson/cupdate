# Cookbook

- [Update kustomizations using Cupdate's API](#update-kustomizations-using-cupdates-api)
- [Update compose files using Cupdate's API](#update-compose-files-using-cupdates-api)

## Update kustomizations using Cupdate's API

If you're using kustomizations to write and deploy your Kubernetes services you
can use Cupdate's API, curl, [yq](https://github.com/mikefarah/yq) and
[jq](https://github.com/jqlang/jq) to update the manifests.

Let's say you have a manifest like the one below.

```yaml
# kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: cupdate

images:
  - name: ghcr.io/alexgustafsson/cupdate
    newTag: 0.14.1

resources:
  - namespace.yml

  - cupdate
```

You can then use a script like
[update-kustomizations.sh](update-kustomizations.sh) to update such manifests
by using Cupdate's APIs.

The key here is to use yq to identify all image overrides and then curl to query
Cupdate. Finally, jq is used to parse the response and yq again to change the
manifest.

## Update compose files using Cupdate's API

If you're using Docker compose files to write and deploy your Docker services
you can use Cupdate's API, curl, [yq](https://github.com/mikefarah/yq) and
[jq](https://github.com/jqlang/jq) to update the manifests.

Let's say you have a compose file like the one below.

```yaml
# compose.yaml
services:
  cupdate:
    image: ghcr.io/alexgustafsson/cupdate:0.14.1
```

You can then use a script like
[update-compose-files.sh](update-compose-files.sh) to update such files by using
Cupdate's APIs.

The key here is to use yq to identify all image overrides and then curl to query
Cupdate. Finally, jq is used to parse the response and yq again to change the
manifest.
