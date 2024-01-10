# API

## Create the api secret

```bash
kubectl apply -f k8s/api/api-secret.yaml
```

## Create the api deployment

```bash
kubectl apply -f k8s/api/api-deployment.yaml
```

## Create the api service

```bash
kubectl apply -f k8s/api/api-service.yaml
```

## Create the api hpa

```bash
kubectl apply -f k8s/api/api-hpa.yaml
```
