# Kube System

## Install the Kubernetes Dashboard

```bash
helm install kubernetes-dashboard k8s-dashboard/kubernetes-dashboard --version 7.0.0-alpha1 --set=cert-manager.enabled=false --set=app.ingress.enabled=false
```

### Upgrade the Kubernetes Dashboard

```bash
helm upgrade --install kubernetes-dashboard kubernetes-dashboard/kubernetes-dashboard --create-namespace --namespace kubernetes-dashboard
```

### Uninstall the Kubernetes Dashboard

```bash
helm delete kubernetes-dashboard
```

## Get the Kubernetes Dashboard URL

```bash
export POD_NAME=$(kubectl get pods -n kubernetes-dashboard -l "app.kubernetes.io/name=kubernetes-dashboard,app.kubernetes.io/instance=kubernetes-dashboard" -o jsonpath="{.items[0].metadata.name}")
echo https://127.0.0.1:8443/
kubectl -n kubernetes-dashboard port-forward $POD_NAME 8443:8443
```

### Create a service account for the Kubernetes Dashboard

```bash
kubectl create -f ./k8s/kube-system/dashboard-account.yaml
```

### Get the token for the Kubernetes dashboard

```bash
kubectl -n kube-system create token admin-user
```
