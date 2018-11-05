# Cryptopus Secret Controller for Kubernetes / OpenShift

[![GitHub release](https://img.shields.io/github/release/puzzle/cryptopus-k8s-secretcontroller.svg)](https://github.com/puzzle/cryptopus-k8s-secretcontroller/releases/)
[![Docker Build Status](https://img.shields.io/docker/build/puzzleitc/cryptopus-k8s-secretcontroller.svg)](https://hub.docker.com/r/puzzleitc/cryptopus-k8s-secretcontroller)

This repository implements a Controller for watching SecretClaim Resources as
defined with a CustomResourceDefinition (CRD) and create Secrets bound to a Cryptopus Account


## SecretClaims

Sample Yaml File for a SecretClaim:

```yaml
apiVersion: cryptopussecretcontroller.puzzle.ch/v1alpha1
kind: SecretClaim
metadata:
  name: mysecretclaim
spec:
  secretName: mysecret
  id: [ 2256 ]
  refreshTime: 3600
  cryptopusSecret: my-namespace/cryptopusapi
```

* `secretName`: is the Name of the Secret that will be created
* `id`: Array of Ids for the Cryptopus Account to get
* `refreshTime`: if > 0, Time after which a Secrets gets an Update with Values from Cryptopus. Defaults to 0, and the `CONTROLLER_DEFAULT_REFRESH_TIME` (See Env Vars) is used
* `cryptopusSecret`: Secret with Cryptopus API Details (URL, Username, Token). For Details see below

## Secret with Cryptopus API Details

You have to create the following Secret with Details of your Cryptopus Instance.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: cryptopusapi
type: Opaque
data:
  CRYPTOPUS_API:
  CRYPTOPUS_API_USER:
  CRYPTOPUS_API_TOKEN:
```

* `CRYPTOPUS_API`: Cryptopus API URL (base64 encoded)
* `CRYPTOPUS_API_USER`: Cryptopus API User with Access to the Accounts you need (base64 encoded)
* `CRYPTOPUS_API_TOKEN`: Cryptopus Token for API User (base64 encoded)

## Resulting Secret

The Secret created by a SecretClaim will have a username, password key in the Data Section. username and password will have the account id as suffix. This allows to store more than one Cryptopus Account in the Secret.


```yaml
apiVersion: v1
data:
  password_2256: ...
  username_2256: ....
kind: Secret
metadata:
  name: mysecret
type: Opaque
```

## Troubleshooting

Look at the Event Section in the SecretClaim to get some Information if something is not working

```
oc describe SecretClaim mysecretclaim
Events:
  Type    Reason  Age    From                       Message
  ----    ------  ----   ----                       -------
  Normal  Synced  2m11s  cryptopussecretcontroller  SecretClaim synced successfully
```


## Remarks

In Openshift 3.7, Secrets are not automaticly deleted when the SecretClaim is deleted. So you have you remove it manually, otherwise you can't use the same Secret Name again with a SecretClaim


## Configuration

### Custom Resource Definition

This has to be created befor the Cryptopus Secret Controller is running

```yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: secretclaims.cryptopussecretcontroller.puzzle.ch
spec:
  group: cryptopussecretcontroller.puzzle.ch
  version: v1alpha1
  names:
    kind: SecretClaim
    plural: secretclaims
  scope: Namespaced

```

### Cluster Role for the ServiceAccount under which the Controller returns

We need the following ClusterRole

```yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: puzzle:app:cryptopussecretcontroller
rules:
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - '*'
- apiGroups:
  - cryptopussecretcontroller.puzzle.ch
  resources:
  - secretclaims
  - secretclaims/finalizers
  verbs:
  - '*'
- apiGroups:
  - ""
  resources:
  - events
  verbs:
  - create
  - patch
  - update
```

and a ClusterRoleBinding to the ServiceAccount which the SecretController Pod runs as:

```yaml
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: puzzle:app:cryptopussecretcontroller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: puzzle:app:cryptopussecretcontroller
subjects:
- kind: ServiceAccount
  name: default
  namespace: my-namespace
```

## Environment Variables

### CONTROLLER_DEFAULT_REFRESH_TIME

Default Refresh Time when not Set in SecretClaim. If not set, defaults to 300 [s]
