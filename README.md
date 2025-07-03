# kubectl-fleet

**WARNING: this is a work-in-progress and unsupported prototype which is not for
production use.  Feedback is welcomed.**

## Introduction

kubectl-fleet is a kubectl plugin that allows users to navigate between
different hub and member clusters of [Azure Kubernetes Fleet
Manager](https://learn.microsoft.com/en-us/azure/kubernetes-fleet/) (Fleet)
instances.

## Installation

```sh
go install github.com/jim-minter/kubectl-fleet/cmd/kubectl-fleet@latest
export PATH=$PATH:$HOME/go/bin # if necessary
```

## Instructions

A Fleet ARM resource ID can be specified either via:

```sh
  --subscription 00000000-0000-0000-0000-000000000000 \
  --resource-group mygroup \
  --fleet-name myfleet
```

or via:

```sh
  --resource-id /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup/providers/Microsoft.ContainerService/fleets/myfleet
```

A Fleet member ARM resource ID can be specified either via

```sh
  --subscription 00000000-0000-0000-0000-000000000000 \
  --resource-group mygroup \
  --fleet-name myfleet \
  --member-name member1
```

or via:

```sh
  --resource-id /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup/providers/Microsoft.ContainerService/fleets/myfleet/members/member1
```

If `--subscription` is omitted, the current default `az` subscription will be
used (i.e. see `az account show`).

### 1. Use kubectl against the Fleet hub

Any kubectl command and options may be specified after writing `kubectl fleet
--resource-group mygroup --fleet-name myfleet`.

```sh
$ kubectl fleet --resource-group mygroup --fleet-name myfleet get namespaces
$ kubectl fleet --resource-group mygroup --fleet-name myfleet get crps
$ kubectl fleet --resource-group mygroup --fleet-name myfleet edit crp mycrp
```

etc.

### 2. Use kubectl against a Fleet member

To connect to a member instead of the hub, specify `--member-name member1`.

```sh
$ kubectl fleet --resource-group mygroup --fleet-name myfleet --member-name member1 get deployments -n mynamespace
$ kubectl fleet --resource-group mygroup --fleet-name myfleet --member-name member1 edit namespace default
```

etc.

### 3. Set kubectl context to the Fleet hub

As writing `fleet --resource-group mygroup --fleet-name myfleet` gets tiring,
you can use `kubectl fleet set-context`.

```sh
$ kubectl fleet set-context --resource-group mygroup --fleet-name myfleet
$ kubectl get namespaces
$ kubectl get crps
$ kubectl edit crp mycrp
```

etc.

### 4. Set kubectl context to a Fleet member

And `--member-name member1` works with `kubectl fleet set-context` too.

```sh
$ kubectl fleet set-context --resource-group mygroup --fleet-name myfleet --member-name member1
$ kubectl get deployments -n mynamespace
$ kubectl edit namespace default
```

etc.

### 5. Easily switch between the Fleet hub and members

You can also omit parameters like `--resource-group mygroup --fleet-name
myfleet` to switch easily between hub and member contexts.  Note that you need
to write `--member-name ''` to switch back to hub context from member context.

```sh
$ kubectl fleet set-context --resource-group mygroup --fleet-name myfleet  # on the hub
$ kubectl get crps
$ kubectl fleet set-context --member-name member1                          # now on member1
$ kubectl get deployments -n mynamespace
$ kubectl logs -n mynamespace -f -l app=myapp --all-containers
$ kubectl config current-context                                           # where was I again?
/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup/providers/Microsoft.ContainerService/fleets/myfleet/members/member1
$ kubectl fleet set-context --member-name member2                          # now on member2
$ kubectl get deployments -n mynamespace
$ kubectl fleet set-context --member-name ''                               # now back on the hub
```

etc.

### 6. Easily remember what the names of the Fleet members are

```sh
$ kubectl fleet members
NAME
member1
member2
```