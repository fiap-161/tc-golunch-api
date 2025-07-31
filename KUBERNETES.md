# Kubernetes Architecture - Tech Challenge FIAP 161

## Overview

This document describes the Kubernetes architecture implementation for the GoLunch application, meeting all specified requirements including scalability, security, and best practices.

## Architecture Components

### 1. Application Tier
- **Deployment**: `golunch-app-deployment`
- **Service**: `golunch-app-service` (LoadBalancer)
- **Replicas**: 3 initial, scalable 2-10 via HPA
- **Image**: Go application containerized with multi-stage Docker build

### 2. Database Tier
- **Deployment**: `postgres-deployment`
- **Service**: `postgres-service` (ClusterIP)
- **Storage**: EmptyDir volume (production should use PersistentVolume)
- **Database**: PostgreSQL 15-alpine

### 3. Configuration Management
- **ConfigMap**: `golunch-config` - Non-sensitive configuration
- **Secrets**: `golunch-secrets` - Database passwords, JWT keys, API tokens

### 4. Auto-scaling
- **HPA**: Horizontal Pod Autoscaler with CPU/Memory targets
- **Metrics**: 70% CPU, 80% Memory utilization thresholds

## Security Implementation

### 1. Secrets Management
- Database passwords stored in Kubernetes Secrets
- JWT secret keys in Secrets
- MercadoPago API tokens in Secrets
- All secrets base64 encoded

### 2. Network Security
- Database accessible only within cluster (ClusterIP)
- Application exposed via LoadBalancer
- No hardcoded credentials in code

### 3. Resource Security
- Resource limits prevent resource exhaustion
- Non-root container execution
- Minimal base image (scratch for final stage)

## Scalability Features

### 1. Horizontal Pod Autoscaler (HPA)
```yaml
minReplicas: 2
maxReplicas: 10
targetCPUUtilization: 70%
targetMemoryUtilization: 80%
```

### 2. Scaling Behavior
- **Scale Up**: Aggressive (100% increase, max 4 pods per 15s)
- **Scale Down**: Conservative (10% decrease per 60s, 5min stabilization)

### 3. Load Distribution
- LoadBalancer service distributes traffic across pods
- Session-less application design enables horizontal scaling

## Health Monitoring

### 1. Health Checks
- **Liveness Probe**: `/ping` endpoint, 30s initial delay
- **Readiness Probe**: `/ping` endpoint, 5s initial delay
- **Failure Threshold**: Automatic pod restart on failures

### 2. Resource Monitoring
- CPU and memory metrics for HPA decisions
- Container resource usage tracking

## Deployment Strategy

### 1. Rolling Updates
- Zero-downtime deployments
- Gradual pod replacement
- Automatic rollback on failure

### 2. Configuration Management
- Separate configs per environment
- Environment-specific secrets
- Kustomize for overlay management

## File Structure

```
k8s/
├── README.md                   # Detailed deployment guide
├── kustomization.yaml         # Kustomize configuration
├── namespace.yaml             # Namespace definition
├── configmap.yaml            # Non-sensitive config
├── secrets.yaml              # Sensitive configuration
├── postgres-deployment.yaml  # Database deployment
├── postgres-service.yaml     # Database service
├── app-deployment.yaml       # Application deployment
├── app-service.yaml          # Application service
└── hpa.yaml                  # Auto-scaling configuration
```

## Environment Variables

### Application Configuration
- `DATABASE_HOST`, `DATABASE_PORT`, `DATABASE_NAME`, `DATABASE_USER`
- `UPLOAD_DIR`, `GIN_MODE`

### Sensitive Configuration
- `SECRET_KEY` (JWT signing)
- `DATABASE_PASSWORD`
- `MERCADOPAGO_ACCESS_TOKEN`

## Production Considerations

### 1. Storage
- Replace EmptyDir with PersistentVolume for database
- Implement backup strategy
- Consider database clustering

### 2. Security Enhancements
- Implement Pod Security Standards
- Network policies for micro-segmentation
- Regular security scanning

### 3. Monitoring
- Implement comprehensive logging (ELK stack)
- Metrics collection (Prometheus/Grafana)
- Alerting for critical thresholds

### 4. High Availability
- Multi-zone deployment
- Database replication
- Ingress controller for advanced routing

## Compliance with Requirements

✅ **Functional Requirements**: All API endpoints maintained
✅ **Scalability**: HPA with CPU/Memory metrics
✅ **Security**: ConfigMaps and Secrets for configuration
✅ **Best Practices**: Deployments and Services architecture
✅ **GitHub Integration**: All manifests versioned in repository

This architecture provides a production-ready, scalable, and secure deployment of the GoLunch application on Kubernetes.