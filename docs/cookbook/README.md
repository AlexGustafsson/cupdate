# Cookbook

- [Update kustomizations using Cupdate's API](#update-kustomizations-using-cupdates-api)
- [Update compose files using Cupdate's API](#update-compose-files-using-cupdates-api)
- [Showing summary in Homepage](#showing-summary-in-homepage)

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
    newTag: 0.16.0

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
    image: ghcr.io/alexgustafsson/cupdate:0.16.0
```

You can then use a script like
[update-compose-files.sh](update-compose-files.sh) to update such files by using
Cupdate's APIs.

The key here is to use yq to identify all image overrides and then curl to query
Cupdate. Finally, jq is used to parse the response and yq again to change the
manifest.

## Showing summary in Homepage

If you're using [Homepage](https://github.com/gethomepage/homepage) you can use
Cupdate's API and Homepage's
[custom API integration](https://gethomepage.dev/widgets/services/customapi/) to
show the outdated and vulnerable image counts directly in Homepage.

Additionally, using
[Kubernetes service annotations](https://gethomepage.dev/configs/kubernetes/#services)
it can be as easy as adding the following labels to your ingress:

```yaml
# Annotations on your Kubernetes service
gethomepage.dev/widget.type: customapi
# Update the widget URL to where you're actually hosting Cupdate, but use the
# /api/v1/summary endpoint
gethomepage.dev/widget.url: http://cupdate.cupdate.svc.cluster.local/api/v1/summary
gethomepage.dev/widget.method: GET
gethomepage.dev/widget.mappings.0.label: outdated
gethomepage.dev/widget.mappings.0.field: outdated
gethomepage.dev/widget.mappings.1.label: images
gethomepage.dev/widget.mappings.1.field: images
gethomepage.dev/widget.mappings.2.label: vulnerable
gethomepage.dev/widget.mappings.2.field: vulnerable
gethomepage.dev/widget.mappings.3.label: processing
gethomepage.dev/widget.mappings.3.field: processing
gethomepage.dev/widget.mappings.4.label: failed
gethomepage.dev/widget.mappings.4.field: failed
```

The same support exists in Docker, but looks a little bit different:

```yaml
# Labels on your Docker container / Docker Compose service
labels:
  - homepage.group=Docker
  - homepage.name=Cupdate
  - homepage.icon=cupdate.png
  # Update the widget URL to where you're actually hosting Cupdate
  - homepage.href=https://cupdate.internal
  - homepage.description=Cupdate keeps track of image updates.
  - homepage.widget.type=customapi
  # Update the widget URL to where you're actually hosting Cupdate, but use the
  # /api/v1/summary endpoint
  - homepage.widget.url=http://cupdate.internal/api/v1/summary
  - homepage.widget.method=GET
  - homepage.widget.mappings[0].label=outdated
  - homepage.widget.mappings[0].field=outdated
  - homepage.widget.mappings[1].label=images
  - homepage.widget.mappings[1].field=images
  - homepage.widget.mappings[2].label=vulnerable
  - homepage.widget.mappings[2].field=vulnerable
  - homepage.widget.mappings[3].label=processing
  - homepage.widget.mappings[3].field=processing
  - homepage.widget.mappings[4].label=failed
  - homepage.widget.mappings[4].field=failed
``
```
