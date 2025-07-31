# Kubernetes Deployment for GoLunch Application

This directory contains all the Kubernetes manifest files required to deploy the GoLunch application in a Kubernetes cluster.

## Architecture Overview

The application consists of:
- **GoLunch API Application**: A Go-based REST API for order management
- **PostgreSQL Database**: Database for persistent storage
- **Horizontal Pod Autoscaler**: Automatic scaling based on CPU and memory usage

## Security Best Practices

- **Secrets**: Sensitive values like database passwords and API keys are stored in Kubernetes Secrets
- **ConfigMaps**: Non-sensitive configuration is stored in ConfigMaps
- **Resource Limits**: All containers have resource requests and limits defined
- **Health Checks**: Liveness and readiness probes are configured
- **Non-root User**: Application runs with minimal privileges

## Files Description

| File | Description |
|------|-------------|
| `namespace.yaml` | Creates the tech-challenge-fiap161 namespace |
| `configmap.yaml` | Non-sensitive configuration values |
| `secrets.yaml` | Sensitive configuration values (base64 encoded) |
| `postgres-deployment.yaml` | PostgreSQL database deployment |
| `postgres-service.yaml` | PostgreSQL service for internal communication |
| `app-deployment.yaml` | GoLunch application deployment |
| `app-service.yaml` | GoLunch service for external access |
| `hpa.yaml` | Horizontal Pod Autoscaler configuration |
| `kustomization.yaml` | Kustomize configuration for easy deployment |

## Deployment Instructions

### Prerequisites

1. Kubernetes cluster (v1.20+)
2. `kubectl` configured to access your cluster
3. Metrics server installed for HPA functionality

### Before Deployment

1. **Update Secrets**: Edit `secrets.yaml` and replace the base64 encoded values with your actual secrets:
   ```bash
   echo -n "your-actual-secret-key" | base64
   echo -n "your-actual-database-password" | base64
   echo -n "your-mercadopago-access-token" | base64
   ```

2. **Build and Push Docker Image**: 
   ```bash
   docker build -t your-registry/golunch-app:latest .
   docker push your-registry/golunch-app:latest
   ```

3. **Update Image Reference**: Edit `app-deployment.yaml` and update the image reference to point to your registry.

### Deploy Using kubectl

```bash
# Deploy all resources
kubectl apply -f k8s/

# Or deploy in order
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/postgres-deployment.yaml
kubectl apply -f k8s/postgres-service.yaml
kubectl apply -f k8s/app-deployment.yaml
kubectl apply -f k8s/app-service.yaml
kubectl apply -f k8s/hpa.yaml
```

### Deploy Using Kustomize

```bash
kubectl apply -k k8s/
```

## Scaling Configuration

The HPA is configured with:
- **Min Replicas**: 2
- **Max Replicas**: 10
- **CPU Target**: 70% utilization
- **Memory Target**: 80% utilization

Scaling behavior:
- **Scale Up**: Fast scaling with up to 100% increase every 15 seconds
- **Scale Down**: Conservative scaling with 10% decrease every 60 seconds, with 5-minute stabilization window

## Monitoring and Health Checks

- **Liveness Probe**: Checks `/ping` endpoint every 10 seconds after 30-second initial delay
- **Readiness Probe**: Checks `/ping` endpoint every 5 seconds after 5-second initial delay

## Accessing the Application

After deployment, the application will be available:
- **Internal Access**: `http://golunch-app-service.tech-challenge-fiap161.svc.cluster.local`
- **External Access**: Depends on your LoadBalancer implementation
- **API Documentation**: `http://<external-ip>/swagger/index.html`

## Troubleshooting

### Check Pod Status
```bash
kubectl get pods -n tech-challenge-fiap161
```

### View Application Logs
```bash
# View logs from specific deployment
kubectl logs -f deployment/golunch-app-deployment -n tech-challenge-fiap161

# View logs from all replicas with pod identification
kubectl logs -f -l app=golunch-app -n tech-challenge-fiap161 --prefix=true --all-containers=true
```

### Check Database Connection
```bash
kubectl exec -it deployment/postgres-deployment -n tech-challenge-fiap161 -- psql -U golunch_user -d golunch
```

### Monitor HPA Status
```bash
kubectl get hpa -n tech-challenge-fiap161
kubectl describe hpa golunch-app-hpa -n tech-challenge-fiap161
```

## Environment Variables

The application uses the following environment variables:

### From ConfigMap
- `DATABASE_HOST`: PostgreSQL host
- `DATABASE_PORT`: PostgreSQL port
- `DATABASE_NAME`: Database name
- `DATABASE_USER`: Database user
- `UPLOAD_DIR`: Directory for file uploads
- `GIN_MODE`: Gin framework mode

### From Secrets
- `SECRET_KEY`: JWT secret key
- `DATABASE_PASSWORD`: Database password
- `MERCADOPAGO_ACCESS_TOKEN`: MercadoPago API token

## Resource Requirements

### Application Pods
- **Requests**: 128Mi memory, 100m CPU
- **Limits**: 256Mi memory, 500m CPU

### Database Pod
- **Requests**: 256Mi memory, 250m CPU
- **Limits**: 512Mi memory, 500m CPU