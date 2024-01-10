# Database

## Create the database secret

```bash
kubectl apply -f k8s/database/database-secret.yaml
```

## Create the database persistent volume

```bash
kubectl apply -f k8s/database/database-volume.yaml
```

## Create the database persistent volume claim

```bash
kubectl apply -f k8s/database/database-pvc.yaml
```

## Create the database stateful set

```bash
kubectl apply -f k8s/database/database-stateful.yaml
```

## Create the database service

```bash
kubectl apply -f k8s/database/database-service.yaml
```
