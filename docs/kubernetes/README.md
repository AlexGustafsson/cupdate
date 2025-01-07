# Running Cupdate in Kubernetes

Cupdate is made for running in Kubernetes. It is intended to be deployed as a
single instance, using the Kubernetes APIs to react on changes made to
deployments, containers, replica sets and more.

To get started, run the command below to inspect the manifests about to be
applied.

```shell
kubectl apply --dry-run=client -o yaml 'https://github.com/AlexGustafsson/cupdate/deploy?ref=v0.14.1'
```

Next, run the following command to apply the manifests.

```shell
kubectl apply -k 'https://github.com/AlexGustafsson/cupdate/deploy?ref=v0.14.1'
```

If you're running Kubernetes with RBAC, Cupdate needs additional configuration.
To install Cupdate with support for RBAC, run the following command.

```shell
kubectl apply -k 'https://github.com/AlexGustafsson/cupdate/deploy/overlays/rbac?ref=v0.14.1'
```

## Config

> [!NOTE]
> As there are a lot of different ways to expose services, Cupdate is not
> deployed without any ingress.

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
kustomize build 'https://github.com/AlexGustafsson/cupdate/deploy/overlays/rbac?ref=v0.14.1' > cupdate.yaml
```

By default, Cupdate will ignore old replica sets kept around by Kubernetes to
enable rollback of services. To include them, set
`CUPDATE_KUBERNETES_INCLUDE_OLD_REPLICAS` to `true`.

Whilst the commands above are enough to get you started with Cupdate, you might
want to change some configuration to better suite your needs. Please see the
additional documentation in [../config.md](../config.md).
