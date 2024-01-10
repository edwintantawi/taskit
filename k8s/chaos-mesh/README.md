# Chaos Mesh

## Add Chaos Mesh repo

```bash
helm repo add chaos-mesh https://charts.chaos-mesh.org
```

## Create namespace for Chaos Mesh

```bash
kubectl create ns chaos-mesh
```

## Install Chaos Mesh

```bash
helm install chaos-mesh chaos-mesh/chaos-mesh -n=chaos-mesh --version 2.6.2
```

## Run Chaos Mesh dashboard

```bash
kubectl port-forward -n chaos-mesh svc/chaos-dashboard 2333:2333
```

## Get the token for the Chaos Mesh dashboard

```bash
kubectl -n kube-system create token admin-user
```
