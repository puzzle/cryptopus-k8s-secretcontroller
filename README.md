# Cryptopus Secret Controller for Kubernetes / OpenShift

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
  cryptopusSecret: pitc-cryptopussecretcontroller-dev/cryptopuspuzzlesplattner
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

## Resulting Secret

The Secret created by a SecretClaim will have a username, password key in the Data Section:


```yaml
apiVersion: v1
data:
  password: ...
  username: ....
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
  namespace: pitc-cryptopussecretcontroller-dev
```

### ImageStream

```yaml
apiVersion: image.openshift.io/v1
kind: ImageStream
metadata:
  labels:
    app: secretcontroller
  name: secretcontroller
spec:
  lookupPolicy:
    local: false
  tags:
  - annotations:
      openshift.io/imported-from: quay.io/splattner/cryptopussecretcontroller
    from:
      kind: DockerImage
      name: quay.io/splattner/cryptopussecretcontroller
    importPolicy: {}
    name: latest
    referencePolicy:
      type: Source
```

### Deployment


```yaml
apiVersion: v1
items:
- apiVersion: apps.openshift.io/v1
  kind: DeploymentConfig
  metadata:
    labels:
      app: secretcontroller
    name: secretcontroller
  spec:
    replicas: 1
    selector:
      app: secretcontroller
      deploymentconfig: secretcontroller
    strategy:
      activeDeadlineSeconds: 21600
      resources: {}
      rollingParams:
        intervalSeconds: 1
        maxSurge: 25%
        maxUnavailable: 25%
        timeoutSeconds: 600
        updatePeriodSeconds: 1
      type: Rolling
    template:
      metadata:
        labels:
          app: secretcontroller
          deploymentconfig: secretcontroller
      spec:
        containers:
        - image: quay.io/splattner/cryptopussecretcontroller@sha256:40b1393e1d1c1c1b3c43f79bfd9014bf49adb75f4a1d827143c4da966df4254d
          imagePullPolicy: Always
          name: secretcontroller
          resources: {}
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
        dnsPolicy: ClusterFirst
        restartPolicy: Always
        schedulerName: default-scheduler
        securityContext: {}
        terminationGracePeriodSeconds: 30
    test: false
    triggers:
    - type: ConfigChange
    - imageChangeParams:
        automatic: true
        containerNames:
        - secretcontroller
        from:
          kind: ImageStreamTag
          name: secretcontroller:latest
          namespace: pitc-cryptopussecretcontroller-dev
        lastTriggeredImage: quay.io/splattner/cryptopussecretcontroller@sha256:40b1393e1d1c1c1b3c43f79bfd9014bf49adb75f4a1d827143c4da966df4254d
      type: ImageChange
```


## Environment Variables

### CONTROLLER_DEFAULT_REFRESH_TIME

Default Refresh Time when not Set in SecretClaim. If not set, defaults to 300 [s]
