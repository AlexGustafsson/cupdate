# Kubernetes

Cupdate is made to run well in Kubernetes. It is intended to be deployed as a
single node, optionally persisting its state to a persistent volume. Cupdate
will automatically react to changes to resources and update its data
accordingly.

Cupdate is intended to be run using a service account.

Please refer to [`rbac.yaml`](./rbac.yaml) and [`service.yaml`](./service.yaml)
for examples on how Cupdate can be configured to run in Kubernetes.

## Graph

The diagram below shows the graphing supported by the Kubernetes Cupdate
platform.

```mermaid
flowchart TD
    Namespace --> ResourceKindAppsV1Deployment
    Namespace --> ResourceKindAppsV1DaemonSet
    Namespace --> ResourceKindAppsV1ReplicaSet
    Namespace --> ResourceKindAppsV1StatefulSet
    Namespace --> ResourceKindBatchV1CronJob
    Namespace --> ResourceKindBatchV1Job
    Namespace --> ResourceKindCoreV1Pod
    Namespace --> ResourceKindCoreV1Pod

    ResourceKindAppsV1Deployment --> ResourceKindCoreV1Pod
    ResourceKindAppsV1DaemonSet --> ResourceKindCoreV1Pod
    ResourceKindAppsV1ReplicaSet --> ResourceKindCoreV1Pod
    ResourceKindAppsV1StatefulSet --> ResourceKindCoreV1Pod
    ResourceKindBatchV1CronJob --> ResourceKindCoreV1Pod
    ResourceKindBatchV1Job --> ResourceKindCoreV1Pod

    ResourceKindCoreV1Pod --> ResourceKindCoreV1Container

    ResourceKindCoreV1Container --> Image
```
