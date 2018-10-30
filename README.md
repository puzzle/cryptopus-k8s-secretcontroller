# vaultsecretcontroller

This repository implements a controller for watching SecretClaim resources as
defined with a CustomResourceDefinition (CRD) and create Secrets bound to a Hashicorp Vault Sercret


# SecretClaims

Sample Yaml File:

```yaml
apiVersion: vaultsecretcontroller.mobi.ch/v1alpha1
kind: SecretClaim
metadata:
  name: example-secretclaim
spec:
  secretName: mysecret
  path: secret/test
  role: my-app-role
  refreshTime: -1
  serviceAccount: ""
```


* `secretName`: is the Name of the Secret that will be created
* `path`: path where the secret is found in Hashicorp Vault
* `role`: Role to be used for kubernetes auth in Hashicorp Vault `vault write auth/kubernetes/login role=? jwt=?`
* `refreshTime`: if > 0, Time after which a Secrets gets an Update with Values from Hashicorp Vault. Defaults to 0, and the `CONTROLLER_DEFAULT_REFRESH_TIME` (See Env Vars) is used
* `serviceAccount`: Service Account to be used for kubernetes auth in Hashicorp Vault. If not set "", the default Namespace ServiceAccount is used

## Configuration

### Custom Resource Definition

```yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: secretclaims.vaultsecretcontroller.mobi.ch
spec:
  group: vaultsecretcontroller.mobi.ch
  version: v1alpha1
  names:
    kind: SecretClaim
    plural: secretclaims
  scope: Namespaced
```

### Vault Secret Controller

* `--kubeconfig`: Path to a kubeconfig. Only required if out-of-cluster.
* `--master`: The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.

### Hashicorp Vault Configuratoin

For Hashicorp Vault Client use the official Environment Variables as listed in https://www.vaultproject.io/docs/commands/index.html

### CA Cert for Vault Configuration

You can create a Config Map with the CA:

```yaml
apiVersion: v1
data:
  ca.pem: |-
    -----BEGIN CERTIFICATE-----
    ....
    -----END CERTIFICATE-----
kind: ConfigMap
metadata:

  name: vault-ca
```

Mount Config Map and set `VAULT_CACERT` to the mounted File


## Other Environment Variables

### CONTROLLER_DEFAULT_REFRESH_TIME

Default Refresh Time when not Set in SecretClaim. If not set, defaults to 300 [s]
