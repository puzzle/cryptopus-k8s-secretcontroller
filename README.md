# Cryptopus Secret Controller

This repository implements a Controller for watching SecretClaim resources as
defined with a CustomResourceDefinition (CRD) and create Secrets bound to a Cryptopus Account


## SecretClaims

Sample Yaml File:

```yaml
apiVersion: cryptopussecretcontroller.puzzle.ch/v1alpha1
kind: SecretClaim
metadata:
  name: testclaim
spec:
  secretName: test
  id: 2256
  refreshTime: 3600
  cryptopusSecret: pitc-cryptopussecretcontroller-dev/cryptopuspuzzlesplattner
```


* `secretName`: is the Name of the Secret that will be created
* `id`: Id of the Cryptopus Account to get
* `refreshTime`: if > 0, Time after which a Secrets gets an Update with Values from Hashicorp Vault. Defaults to 0, and the `CONTROLLER_DEFAULT_REFRESH_TIME` (See Env Vars) is used
* `cryptopusSecret`: Secret with Cryptopus API Details (URL, Username, Token)

## Secret with Cryptopus API Details

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: cryptopuspuzzlesplattner
type: Opaque
data:
  CRYPTOPUS_API:
  CRYPTOPUS_API_USER:
  CRYPTOPUS_API_TOKEN:
```

* `CRYPTOPUS_API`: Cryptopus API URL (base64 encoded)
* `CRYPTOPUS_API_USER`: Cryptopus API User with Access to the Accounts you need (base64 encoded)
* `CRYPTOPUS_API_TOKEN`: Cryptopus Token for API User (base64 encoded)


## Configuration

### Custom Resource Definition

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


## Environment Variables

### CONTROLLER_DEFAULT_REFRESH_TIME

Default Refresh Time when not Set in SecretClaim. If not set, defaults to 300 [s]
