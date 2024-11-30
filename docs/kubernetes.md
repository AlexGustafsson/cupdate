# Kubernetes

> [!WARNING]
> WIP

## Origin

```mermaid
flowchart TD
    Namespace --> ResourceKindAppsV1Deployment
    Namespace --> ResourceKindAppsV1DaemonSet
    Namespace --> ResourceKindAppsV1ReplicaSet
    Namespace --> ResourceKindAppsV1StatefulSet
    Namespace --> ResourceKindBatchV1CronJob
    Namespace --> ResourceKindBatchV1Job
    Namespace --> ResourceKindCoreV1Pod
    Namespace --> Pod

    ResourceKindAppsV1Deployment --> Pod
    ResourceKindAppsV1DaemonSet --> Pod
    ResourceKindAppsV1ReplicaSet --> Pod
    ResourceKindAppsV1StatefulSet --> Pod
    ResourceKindBatchV1CronJob --> Pod
    ResourceKindBatchV1Job --> Pod
    ResourceKindCoreV1Pod --> Pod

    Pod --> Container

    Container --> Image
```
