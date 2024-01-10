# Kubernetes Dashboard

## Install the Kubernetes Dashboard.

```bash
helm install kubernetes-dashboard k8s-dashboard/kubernetes-dashboard --version 7.0.0-alpha1 --set=cert-manager.enabled=false --set=app.ingress.enabled=false
```

### Upgrade the Kubernetes Dashboard.

```bash
helm upgrade --install kubernetes-dashboard kubernetes-dashboard/kubernetes-dashboard --create-namespace --namespace kubernetes-dashboard
```

### Uninstall the Kubernetes Dashboard.

```bash
helm delete kubernetes-dashboard
```

## Get the Kubernetes Dashboard URL.

```bash
export POD_NAME=$(kubectl get pods -n kubernetes-dashboard -l "app.kubernetes.io/name=kubernetes-dashboard,app.kubernetes.io/instance=kubernetes-dashboard" -o jsonpath="{.items[0].metadata.name}")
echo https://127.0.0.1:8443/
kubectl -n kubernetes-dashboard port-forward $POD_NAME 8443:8443
```

### Create a service account for the Kubernetes Dashboard.

```bash
kubectl create -f ./k8s/dashboard-account.yaml
```

### Get the token for the Kubernetes Dashboard.

```bash
kubectl -n kube-system create token admin-user
```

## Chaos Mesh

### Add Chaos Mesh repo.

```bash
helm repo add chaos-mesh https://charts.chaos-mesh.org
```

### Create namespace for Chaos Mesh.

```bash
kubectl create ns chaos-mesh
```

### Install Chaos Mesh.

```bash
helm install chaos-mesh chaos-mesh/chaos-mesh -n=chaos-mesh --version 2.6.2
```

### Run Chaos Mesh dashboard.

```bash
kubectl port-forward -n chaos-mesh svc/chaos-dashboard 2333:2333
```