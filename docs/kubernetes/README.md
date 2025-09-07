# Running Cupdate in Kubernetes

Cupdate is made for running in Kubernetes. It is intended to be deployed as a
single instance, using the Kubernetes APIs to react on changes made to
deployments, containers, replica sets and more.

To get started, run the command below to inspect the manifests about to be
applied.

```shell
kubectl apply --dry-run=client -o yaml 'https://github.com/AlexGustafsson/cupdate/deploy?ref=v0.22.0'
```

Next, run the following command to apply the manifests.

```shell
kubectl apply -k 'https://github.com/AlexGustafsson/cupdate/deploy?ref=v0.22.0'
```

If you're running Kubernetes with RBAC, Cupdate needs additional configuration.
To install Cupdate with support for RBAC, run the following command.

```shell
kubectl apply -k 'https://github.com/AlexGustafsson/cupdate/deploy/overlays/rbac?ref=v0.22.0'
```

## Config

> [!NOTE]
> As there are a lot of different ways to expose services, Cupdate is deployed
> without any ingress.

> [!NOTE]
> Without additional configuration, Cupdate is deployed without any persistent
> state. This will work, but may require additional time after startup for all
> images to be processed.

To more easily configure Cupdate, it's recommended to use a
`kustomization.yaml` file. You can copy [kustomization.yaml](kustomization.yaml)
and then run `kubectl apply -k kustomization.yaml` to deploy Cupdate.

For even more configurability, build the complete manifests and modify them to
your liking.

```shell
kustomize build 'https://github.com/AlexGustafsson/cupdate/deploy/overlays/rbac?ref=v0.22.0' > cupdate.yaml
```

By default, Cupdate will ignore old replica sets kept around by Kubernetes to
enable rollback of services. To include them, set
`CUPDATE_KUBERNETES_INCLUDE_OLD_REPLICAS` to `true`.

Cupdate uses an event-driven architecture to keep its state up-to-date and
properly reflecting Kubernetes. That means that when applicable resources are
changed, Cupdate will produce an initial graph of relationships and images in
use and make it available to the user. Though the processing is cheap and the
events are debounced and API calls cached by default, it could result in
additional API calls made to the Kubernetes APIs. These should cheap as well. If
resources are frequently changed, (more than once every minute or so) you can
control the delay by using the `CUPDATE_KUBERNETES_DEBOUNCE_INTERVAL`
environment variable.

Whilst the commands above are enough to get you started with Cupdate, you might
want to change some configuration to better suite your needs. Please see the
additional documentation in [../config.md](../config.md).

## Updating Cupdate

> [!NOTE]
> Before Cupdate hits v1.0.0, breaking changes can occur. Breaking changes could
> include API changes or changes to how the data is stored on disk. Breaking
> changes are communicated in release notes.

If you've installed Cupdate using the official kustomization, please re-apply it
using the latest version to update Cupdate. If you've written custom manifests,
update the image version and refer to the release notes to learn if there are
additional changes required.
