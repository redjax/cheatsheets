---
description: "Kubectl is a tool for running commands against Kubernetes clusters."
last_updated: "{{last_update}}"
tags:
  - "kubectl"
  - "command"
  - "kubernetes"
  - "cli"
  - "container"
  - "orchestration"
---

# Kubectl <!-- omit in toc -->

[Kubectl](https://kubernetes.io/docs/reference/kubectl/) is a command line utility for controlling a Kubernetes cluster. It uses a `~/.kube/config` file copied from the cluster for connection, and lets you deploy things to the cluster, view its running state and configuration, read logs, start/stop/restart containers, and bind cluster ports to local ports to access webUIs securely.

## Usage

General syntax:

```shell
kubectl [command] [TYPE] [NAME] [flags]
```

### Cluster and context commands

| Command                                                        | Description                                         |
| -------------------------------------------------------------- | --------------------------------------------------- |
| `kubectl cluster-info`                                         | Display cluster endpoint information                |
| `kubectl version`                                              | Display client and server Kubernetes versions       |
| `kubectl get nodes`                                            | List cluster nodes                                  |
| `kubectl get nodes -o wide`                                    | List nodes with additional information              |
| `kubectl config current-context`                               | Display the active context                          |
| `kubectl config get-contexts`                                  | List available contexts                             |
| `kubectl config use-context [context]`                         | Switch to a Kubernetes context                      |
| `kubectl config view`                                          | Display the merged kubeconfig                       |
| `kubectl config view --minify`                                 | Display configuration for the current context       |
| `kubectl config set-context --current --namespace=[namespace]` | Set the default namespace for the current context   |
| `kubectl explain [resource]`                                   | Display documentation for a resource                |
| `kubectl explain [resource].[field]`                           | Display documentation for a specific resource field |

### Namespace commands

| Command                                | Description                                   |
| -------------------------------------- | --------------------------------------------- |
| `kubectl get namespaces`               | List namespaces                               |
| `kubectl get ns`                       | List namespaces using the short resource name |
| `kubectl create namespace [namespace]` | Create a namespace                            |
| `kubectl delete namespace [namespace]` | Delete a namespace                            |
| `kubectl get pods -A`                  | List pods in every namespace                  |
| `kubectl get all -n [namespace]`       | List common resources in a namespace          |

### Resource commands


| Command                                            | Description                                   |
| -------------------------------------------------- | --------------------------------------------- |
| `kubectl get [resource]`                           | List resources                                |
| `kubectl get [resource] [name]`                    | Display a specific resource                   |
| `kubectl get [resource] -o wide`                   | Display resources with additional information |
| `kubectl get [resource] -o yaml`                   | Display a resource as YAML                    |
| `kubectl get [resource] -o json`                   | Display a resource as JSON                    |
| `kubectl get [resource] -o name`                   | Display only resource names                   |
| `kubectl get [resource] --show-labels`             | Display resource labels                       |
| `kubectl get [resource] -l [key]=[value]`          | Filter resources by label                     |
| `kubectl get [resource] --sort-by=[jsonpath]`      | Sort resources by a field                     |
| `kubectl describe [resource] [name]`               | Display detailed information about a resource |
| `kubectl label [resource] [name] [key]=[value]`    | Add or update a resource label                |
| `kubectl annotate [resource] [name] [key]=[value]` | Add or update a resource annotation           |

### Pod commands

| Command                                                                    | Description                                       |
| -------------------------------------------------------------------------- | ------------------------------------------------- |
| `kubectl get pods`                                                         | List pods in the current namespace                |
| `kubectl get pods -o wide`                                                 | List pods with node and IP information            |
| `kubectl get pods --watch`                                                 | Watch pod changes                                 |
| `kubectl describe pod [pod]`                                               | Display detailed pod information                  |
| `kubectl logs [pod]`                                                       | Display pod logs                                  |
| `kubectl logs -f [pod]`                                                    | Follow pod logs                                   |
| `kubectl logs [pod] --previous`                                            | Display logs from the previous container instance |
| `kubectl logs [pod] -c [container]`                                        | Display logs for a specific container             |
| `kubectl logs [pod] --all-containers=true`                                 | Display logs from every container                 |
| `kubectl exec [pod] -- [command]`                                          | Execute a command in a pod                        |
| `kubectl exec -it [pod] -- /bin/sh`                                        | Open an interactive shell in a pod                |
| `kubectl exec -it [pod] -c [container] -- /bin/sh`                         | Open a shell in a specific container              |
| `kubectl cp [pod]:[path] [local-path]`                                     | Copy a file from a pod                            |
| `kubectl cp [local-path] [pod]:[path]`                                     | Copy a file to a pod                              |
| `kubectl port-forward pod/[pod] [local-port]:[pod-port]`                   | Forward a local port to a pod                     |
| `kubectl port-forward service/[service] [local-port]:[service-port]`       | Forward a local port to a service                 |
| `kubectl run [name] --image=[image]`                                       | Create and run a pod                              |
| `kubectl run [name] --rm -it --restart=Never --image=[image] -- [command]` | Run a temporary interactive pod                   |

### Deployment commands

| Command                                                         | Description                         |
| --------------------------------------------------------------- | ----------------------------------- |
| `kubectl get deployments`                                       | List deployments                    |
| `kubectl create deployment [name] --image=[image]`              | Create a deployment                 |
| `kubectl set image deployment/[name] [container]=[image]`       | Update a deployment image           |
| `kubectl scale deployment/[name] --replicas=[number]`           | Scale a deployment                  |
| `kubectl rollout status deployment/[name]`                      | Watch deployment rollout status     |
| `kubectl rollout history deployment/[name]`                     | Display deployment rollout history  |
| `kubectl rollout history deployment/[name] --revision=[number]` | Display a specific rollout revision |
| `kubectl rollout undo deployment/[name]`                        | Undo the most recent rollout        |
| `kubectl rollout undo deployment/[name] --to-revision=[number]` | Undo to a specific revision         |
| `kubectl rollout restart deployment/[name]`                     | Restart all pods in a deployment    |
| `kubectl rollout pause deployment/[name]`                       | Pause a deployment rollout          |
| `kubectl rollout resume deployment/[name]`                      | Resume a deployment rollout         |

### Service and networking commands

| Command                                                              | Description                               |
| -------------------------------------------------------------------- | ----------------------------------------- |
| `kubectl get services`                                               | List services                             |
| `kubectl describe service [service]`                                 | Display detailed service information      |
| `kubectl get endpoints [service]`                                    | Display service endpoints                 |
| `kubectl get endpointslices`                                         | List EndpointSlices                       |
| `kubectl proxy`                                                      | Start a local proxy to the Kubernetes API |
| `kubectl port-forward service/[service] [local-port]:[service-port]` | Forward local traffic to a service        |

### Apply and delete commands

| Command                                            | Description                                          |
| -------------------------------------------------- | ---------------------------------------------------- |
| `kubectl apply -f [file]`                          | Create or update resources from a manifest           |
| `kubectl apply -f [directory]`                     | Create or update resources from files in a directory |
| `kubectl apply -R -f [directory]`                  | Recursively apply manifests                          |
| `kubectl apply -f -`                               | Apply a manifest from standard input                 |
| `kubectl diff -f [file]`                           | Display changes before applying a manifest           |
| `kubectl apply --dry-run=client -f [file] -o yaml` | Generate a manifest without contacting the cluster   |
| `kubectl apply --dry-run=server -f [file]`         | Validate a manifest on the server without saving it  |
| `kubectl delete -f [file]`                         | Delete resources defined in a manifest               |
| `kubectl delete [resource] [name]`                 | Delete a named resource                              |
| `kubectl delete [resource] -l [key]=[value]`       | Delete resources matching a label                    |

### Events and troubleshooting commands

| Command                                                                      | Description                                   |
| ---------------------------------------------------------------------------- | --------------------------------------------- |
| `kubectl get events`                                                         | List events in the current namespace          |
| `kubectl get events -A`                                                      | List events in every namespace                |
| `kubectl get events --sort-by='.metadata.creationTimestamp'`                 | Sort events chronologically                   |
| `kubectl get events --field-selector=type=Warning`                           | List warning events                           |
| `kubectl get pods --field-selector=status.phase!=Running`                    | List pods that are not running                |
| `kubectl top nodes`                                                          | Display node resource usage                   |
| `kubectl top pods`                                                           | Display pod resource usage                    |
| `kubectl top pods -A`                                                        | Display pod resource usage in every namespace |
| `kubectl debug -it [pod] --image=[image] --target=[container] -- [command]`  | Add an ephemeral debugging container          |
| `kubectl debug [pod] -it --copy-to=[debug-pod] --image=[image] -- [command]` | Create a debugging copy of a pod              |
| `kubectl debug node/[node] -it --image=[image]`                              | Create a debugging pod on a node              |

### Node commands

| Command                                                           | Description                                   |
| ----------------------------------------------------------------- | --------------------------------------------- |
| `kubectl get nodes`                                               | List nodes                                    |
| `kubectl describe node [node]`                                    | Display detailed node information             |
| `kubectl cordon [node]`                                           | Mark a node as unschedulable                  |
| `kubectl drain [node] --ignore-daemonsets`                        | Evict pods and prepare a node for maintenance |
| `kubectl drain [node] --ignore-daemonsets --delete-emptydir-data` | Drain a node and remove `emptyDir` data       |
| `kubectl uncordon [node]`                                         | Mark a node as schedulable                    |

## Examples

### Inspect a workload

```shell
kubectl get deployment/web -o wide
kubectl describe deployment/web
kubectl rollout status deployment/web
kubectl get pods -l app=web -o wide
kubectl logs -l app=web --all-containers=true --prefix=true
```

### Debug a failing pod

List non-running pods and recent warning events:

```shell
kubectl get pods -A \
  --field-selector=status.phase!=Running

kubectl get events -A \
  --field-selector=type=Warning \
  --sort-by='.metadata.creationTimestamp'
```

Inspect the pod and its logs:

```shell
kubectl describe pod [pod]
kubectl logs [pod] --all-containers=true
kubectl logs [pod] --previous
```

Open a temporary debugging shell:

```shell
kubectl debug -it [pod] \
  --image=busybox:1.37 \
  --target=[container] \
  -- sh
```

### Drain a node for maintenance

```shell
kubectl cordon worker-01

kubectl drain worker-01 \
  --ignore-daemonsets \
  --delete-emptydir-data

# Perform maintenance.

kubectl uncordon worker-01
```

## Troubleshooting

Check the active context and namespace:

```shell
kubectl config current-context
kubectl config view --minify
```

Check cluster connectivity:

```shell
kubectl cluster-info
kubectl get nodes
```

Inspect recent events:

```shell
kubectl get events -A \
  --sort-by='.metadata.creationTimestamp'
```

Inspect a resource:

```shell
kubectl describe pod [pod]
kubectl describe deployment [deployment]
kubectl describe service [service]
```

Check logs from a restarting container:

```shell
kubectl logs [pod] --previous
```

Check resource usage:

```shell
kubectl top nodes
kubectl top pods -A
```

> [!WARNING]
> `kubectl delete`, `kubectl drain`, and rollout commands can affect running workloads. Use `kubectl describe`, `kubectl diff`, and dry-run options before making destructive or large-scale changes.

## Links

- [kubectl Quick Reference](https://kubernetes.io/docs/reference/kubectl/quick-reference/)
- [kubectl Overview](https://kubernetes.io/docs/reference/kubectl/)
- [kubectl Command Reference](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands)
- [kubectl Rollout Reference](https://kubernetes.io/docs/reference/kubectl/generated/kubectl_rollout/)
- [kubectl Exec Reference](https://kubernetes.io/docs/reference/kubectl/generated/kubectl_exec/)
- [kubectl Debug Reference](https://kubernetes.io/docs/reference/kubectl/generated/kubectl_debug/)
