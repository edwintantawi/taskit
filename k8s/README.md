# Kubernetes

## Setup kube system

```bash
kubectl apply -f k8s/kube-system
```

## Setup all application service manifests

### Create the database resource

```bash
kubectl apply -f k8s/database
```

### Create the api resource

```bash
kubectl apply -f k8s/api
```

### Create the web resource

```bash
kubectl apply -f k8s/web
```

## Test the application

### Simulate pod failure

```bash
kubectl apply -f k8s/chaos-mesh/pod-failure.yaml
```

### Simulate stress test

```bash
kubectl apply -f k8s/chaos-mesh/stress-test.yaml
```
