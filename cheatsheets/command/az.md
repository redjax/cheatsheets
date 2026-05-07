---
description: "The [Azure CLI](https://github.com/Azure/azure-cli) is a CLI tool for interacting with Microsoft Azure."
last_updated: "{{last_update}}"
tags: ["command", "azure", "cloud"]
last_updated: "2026-05-07"
---
## Table of Contents <!-- omit in toc -->

- [About](#about)
- [Usage](#usage)
  - [Login](#login)
  - [User info](#user-info)
  - [Show active subscription](#show-active-subscription)
  - [Set/change subscription](#setchange-subscription)
  - [ACR](#acr)
  - [Resources](#resources)
    - [List](#list)
    - [Show](#show)
  - [Resource Groups](#resource-groups)
  - [Resource Graph](#resource-graph)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)
- [Links](#links)

# Azure CLI <!-- omit in toc -->

## About

The [Azure CLI](https://github.com/Azure/azure-cli) is a CLI tool for interacting with Microsoft Azure.

[Azure CLI docs](https://learn.microsoft.com/en-us/cli/azure/?view=azure-cli-latest)

## Usage

### Login

On most machines, running `az login` will start an authentication flow to get a token for future commands.

In some environments, like in WSL, you need to use `az login --use-device-code`. You may also need to pass a scope, like: `az login --scope https://management.core.windows.net//.default --use-device-code`.

### User info

- See more detailed information about the currently logged in user, use:

  ```shell
  az ad signed-in-user show
  ```

- List the user's roles and subscriptions:

  ```shell
  az role assignment list --assignee $(az ad signed-in-user show --query objectId -o tsv) --output table
  ```

- List user's object ID:

  ```shell
  az ad signed-in-user show --query objectId -o tsv
  ```

### Show active subscription

```shell
az account show --output table
```

### Set/change subscription

```shell
az account set --subscription "<subscription-name>"
```

### ACR

- List all registries in current subscription:

  ```shell
  az acr list --output table
  ```

- List repositories in a given container registry:

  ```shell
  az acr repository list --name <registry-name> --output table
  ```

### Resources

#### List

- List all Azure resources:

  ```shell
  az resource list --output table
  ```

- List resources by type (i.e. `vm`, `aks`, `acr`, etc):
  - List Azure VMs:

    ```
    az vm list -d --output table
    ```

  - List container registries:

    ```shell
    az acr list --output table
    ```

  - AKS clusters:

    ```shell
    az aks list --output table
    ```

  - App services (web apps):

    ```shell
    az webapp list --output table
    ```

  - Storage accounts:

    ```shell
    az storage account list --output table
    ```

  - SQL servers:

    ```shell
    az sql server list --output table
    ```

  - Key vaults:

    ```shell
    az keyvault list --output table
    ```

#### Show

- VM details
  - `az vm show --name <vm-name> --resource-group <group-name> --output json`

### Resource Groups

- List all resource groups:

  ```shell
  az group list --output table
  ```

- List resources in a resource group:

  ```shell
  az resource list --resource-group <group-name> --output table
  ```

### Resource Graph

You can use the Azure Resource Graph for cross-subscription queries. You may need to install the `resource-graph` extension first:

```shell
az extension add --name resource-graph
```

Queries generally look like:

```shell
az graph query -q "Resources | project name, type, location, resourceGroup" --output table
```

## Examples

## Troubleshooting

## Links

