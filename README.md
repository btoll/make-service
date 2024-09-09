# make-service

This is a simple tool for anyone to create a new single or multi-deployment service to quickly get started testing their application.

It will build the service and put it in the `./build/` directory.

## Example

Create a single deployment service in the `development`, `beta` and `production` environments:

`recipe.yaml`

```yaml
name: aion-payments-micro
deployments:
  - name: aionpaymentsconsumer
    image:
      name: foo
      tag: development
    environments:
      - name: development
        replicas: 2
        image:
          name: foo
          tag: development
        envs:
          - FOO=foo
          - ION.Host__ApplicationName=aionpaymentsconsumer
          - ASPNETCORE_ENVIRONMENT=dev
      - name: beta
        replicas: 3
        image:
          name: foo
          tag: beta
        envs:
          - BAR=bar
          - ION.Host__ApplicationName=aionpaymentsconsumer
          - ASPNETCORE_ENVIRONMENT=beta
      - name: production
        replicas: 4
        image:
          name: foo
          tag: production
        envs:
          - QUUX=quux
          - ION.Host__ApplicationName=aionpaymentsconsumer
          - ASPNETCORE_ENVIRONMENT=production
```

```bash
$ ./make-service --filename recipe.yaml
$ tree build
build/
└── aion-payments-micro/
    └── aionpaymentsconsumer/
        ├── base/
        │   ├── deployment.yaml
        │   ├── kustomization.yaml
        │   └── service.yaml
        └── overlays/
            ├── beta/
            │   ├── env
            │   └── kustomization.yaml
            ├── development/
            │   ├── env
            │   └── kustomization.yaml
            └── production/
                ├── env
                └── kustomization.yaml

8 directories, 9 files
```

Test it by specifying the `development` environment:

```bash
$ kubectl kustomize build/aion-payments-micro/aionpaymentsconsumer/overlays/development/
```

This should produce the following manifests:

```yaml
apiVersion: v1
data:
  ASPNETCORE_ENVIRONMENT: dev
  FOO: foo
  ION.Host__ApplicationName: aionpaymentsconsumer
kind: ConfigMap
metadata:
  name: env-aionpaymentsconsumer-gg59c79c78
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: aionpaymentsconsumer
    kafka: "true"
  name: aionpaymentsconsumer
  namespace: default
spec:
  replicas: 2
  selector:
    matchLabels:
      app: aionpaymentsconsumer
  template:
    metadata:
      labels:
        app: aionpaymentsconsumer
    spec:
      containers:
      - envFrom:
        - configMapRef:
            name: env-aionpaymentsconsumer-gg59c79c78
        - configMapRef:
            name: kubernetes-container-user
        image: 451310829282.dkr.ecr.us-east-1.amazonaws.com/aion/foo:development
        imagePullPolicy: Always
        name: aionpaymentsconsumer
      nodeSelector:
        node_type: default

```

---

Create a multi-deployment service in the `development` environment with Kubernetes services:

`recipe.yaml`

```yaml
name: aion-payments-micro
deployments:
  - name: aionpaymentsconsumer
    image:
      name: foo
      tag: development
    environments:
      - name: development
        replicas: 2
        image:
          name: foo
          tag: development
        envs:
          - FOO=foo
          - ION.Host__ApplicationName=aionpaymentsconsumer
          - ASPNETCORE_ENVIRONMENT=dev
      - name: beta
        replicas: 3
        image:
          name: foo
          tag: beta
        envs:
          - BAR=bar
          - ION.Host__ApplicationName=aionpaymentsconsumer
          - ASPNETCORE_ENVIRONMENT=beta
      - name: production
        replicas: 4
        image:
          name: foo
          tag: production
        envs:
          - QUUX=quux
          - ION.Host__ApplicationName=aionpaymentsconsumer
          - ASPNETCORE_ENVIRONMENT=production
  - name: aionpaymentsmicro
    image:
      name: bar
      tag: latest
    environments:
      - name: development
        replicas: 4
        image:
          name: bar
          tag: development
        envs:
          - bar=bar
  - name: aionpaymentsorganizationconsumer
    image:
      name: quux
      tag: latest
    environments:
      - name: development
        replicas: 7
        image:
          name: quux
          tag: development
        envs:
          - QUUX=quux

services:
  - name: aionpaymentsconsumer
    port: 90
    targetPort: 9191
  - name: aionpaymentsorganizationconsumer
    port: 80
    targetPort: 8080
```

```bash
$ ./make-service --filename recipe.yaml
```

