# Hướng dẫn Deployment - Zplus SaaS

## 1. Tổng quan Deployment

Zplus SaaS hỗ trợ nhiều phương thức deployment:
- **Docker Compose**: Cho development và small-scale production
- **Kubernetes**: Cho production scale với high availability
- **Docker Swarm**: Alternative cho Kubernetes
- **Manual**: Traditional server deployment

## 2. Environment Setup

### 2.1 Production Environment Variables

**Backend Environment (.env.prod):**
```env
# Application
ENV=production
PORT=8080
LOG_LEVEL=info

# Database
DB_HOST=postgres.internal
DB_PORT=5432
DB_USER=zplus_prod
DB_PASSWORD=your_secure_password
DB_NAME=zplus_system
DB_SSL_MODE=require
DB_MAX_CONNECTIONS=50

# Redis
REDIS_HOST=redis.internal
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password
REDIS_DB=0

# JWT
JWT_SECRET=your-super-secure-jwt-secret-minimum-32-characters
JWT_EXPIRES_IN=1h
JWT_REFRESH_EXPIRES_IN=7d

# File Storage
STORAGE_TYPE=s3
AWS_REGION=us-east-1
AWS_S3_BUCKET=zplus-files-prod
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key

# Email
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASS=your_sendgrid_api_key
FROM_EMAIL=noreply@zplus.com
FROM_NAME=Zplus SaaS

# Monitoring
SENTRY_DSN=https://your-sentry-dsn@sentry.io/project
NEW_RELIC_LICENSE_KEY=your_new_relic_key

# Rate Limiting
RATE_LIMIT_MAX=1000
RATE_LIMIT_WINDOW_MS=3600000

# Cors
CORS_ORIGIN=https://app.zplus.com,https://*.zplus.com
```

**Frontend Environment (.env.production):**
```env
# API Endpoints
NEXT_PUBLIC_API_URL=https://api.zplus.com
NEXT_PUBLIC_GRAPHQL_URL=https://api.zplus.com/graphql
NEXT_PUBLIC_WS_URL=wss://api.zplus.com/graphql

# Application
NEXT_PUBLIC_APP_NAME=Zplus SaaS
NEXT_PUBLIC_APP_URL=https://app.zplus.com

# Analytics
NEXT_PUBLIC_GA_TRACKING_ID=GA_MEASUREMENT_ID
NEXT_PUBLIC_HOTJAR_ID=your_hotjar_id

# Sentry
NEXT_PUBLIC_SENTRY_DSN=https://your-sentry-dsn@sentry.io/project

# Feature Flags
NEXT_PUBLIC_ENABLE_ANALYTICS=true
NEXT_PUBLIC_ENABLE_CHAT=true
```

## 3. Docker Compose Deployment

### 3.1 Production Docker Compose

**docker-compose.prod.yml:**
```yaml
version: '3.8'

services:
  # Load Balancer & Reverse Proxy
  traefik:
    image: traefik:v2.10
    command:
      - --api.dashboard=true
      - --api.insecure=false
      - --providers.docker=true
      - --providers.docker.exposedbydefault=false
      - --entrypoints.web.address=:80
      - --entrypoints.websecure.address=:443
      - --certificatesresolvers.letsencrypt.acme.tlschallenge=true
      - --certificatesresolvers.letsencrypt.acme.email=admin@zplus.com
      - --certificatesresolvers.letsencrypt.acme.storage=/letsencrypt/acme.json
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - ./letsencrypt:/letsencrypt
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.traefik.rule=Host(`traefik.zplus.com`)"
      - "traefik.http.routers.traefik.tls.certresolver=letsencrypt"
    networks:
      - zplus-network
    restart: unless-stopped

  # Database
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: zplus_system
      POSTGRES_USER: zplus_prod
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backups:/backups
    networks:
      - zplus-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U zplus_prod -d zplus_system"]
      interval: 30s
      timeout: 10s
      retries: 3

  # Redis
  redis:
    image: redis:7-alpine
    command: redis-server --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis_data:/data
    networks:
      - zplus-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 30s
      timeout: 10s
      retries: 3

  # Backend Gateway
  gateway:
    image: zplus/gateway:latest
    env_file:
      - .env.prod
    depends_on:
      - postgres
      - redis
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.api.rule=Host(`api.zplus.com`)"
      - "traefik.http.routers.api.tls.certresolver=letsencrypt"
      - "traefik.http.services.api.loadbalancer.server.port=8080"
    networks:
      - zplus-network
    restart: unless-stopped
    deploy:
      replicas: 2

  # Auth Service
  auth-service:
    image: zplus/auth-service:latest
    env_file:
      - .env.prod
    depends_on:
      - postgres
      - redis
    networks:
      - zplus-network
    restart: unless-stopped
    deploy:
      replicas: 2

  # File Service
  file-service:
    image: zplus/file-service:latest
    env_file:
      - .env.prod
    depends_on:
      - postgres
      - redis
    networks:
      - zplus-network
    restart: unless-stopped

  # CRM Service
  crm-service:
    image: zplus/crm-service:latest
    env_file:
      - .env.prod
    depends_on:
      - postgres
      - redis
    networks:
      - zplus-network
    restart: unless-stopped

  # Frontend
  frontend:
    image: zplus/frontend:latest
    environment:
      - NODE_ENV=production
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.app.rule=Host(`app.zplus.com`) || Host(`*.zplus.com`)"
      - "traefik.http.routers.app.tls.certresolver=letsencrypt"
      - "traefik.http.services.app.loadbalancer.server.port=3000"
    networks:
      - zplus-network
    restart: unless-stopped
    deploy:
      replicas: 2

  # Background Worker
  worker:
    image: zplus/worker:latest
    env_file:
      - .env.prod
    depends_on:
      - postgres
      - redis
    networks:
      - zplus-network
    restart: unless-stopped
    deploy:
      replicas: 2

volumes:
  postgres_data:
    driver: local
  redis_data:
    driver: local

networks:
  zplus-network:
    driver: bridge
```

### 3.2 Deployment Commands

```bash
# Deploy to production
docker-compose -f docker-compose.prod.yml up -d

# Update services
docker-compose -f docker-compose.prod.yml pull
docker-compose -f docker-compose.prod.yml up -d

# View logs
docker-compose -f docker-compose.prod.yml logs -f gateway

# Scale services
docker-compose -f docker-compose.prod.yml up -d --scale gateway=3
```

### 3.3 Health Checks

```bash
# Check all services
docker-compose -f docker-compose.prod.yml ps

# Test API health
curl -f https://api.zplus.com/health

# Test frontend
curl -f https://app.zplus.com
```

## 4. Kubernetes Deployment

### 4.1 Namespace và ConfigMaps

**namespace.yaml:**
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: zplus-prod
```

**configmap.yaml:**
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: zplus-config
  namespace: zplus-prod
data:
  DB_HOST: "postgres-service"
  REDIS_HOST: "redis-service"
  ENV: "production"
  LOG_LEVEL: "info"
```

**secrets.yaml:**
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: zplus-secrets
  namespace: zplus-prod
type: Opaque
data:
  DB_PASSWORD: <base64_encoded_password>
  JWT_SECRET: <base64_encoded_secret>
  REDIS_PASSWORD: <base64_encoded_password>
```

### 4.2 Database Deployment

**postgres.yaml:**
```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
  namespace: zplus-prod
spec:
  serviceName: postgres-service
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:15-alpine
        env:
        - name: POSTGRES_DB
          value: zplus_system
        - name: POSTGRES_USER
          value: zplus_prod
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: zplus-secrets
              key: DB_PASSWORD
        ports:
        - containerPort: 5432
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
        livenessProbe:
          exec:
            command:
            - pg_isready
            - -U
            - zplus_prod
            - -d
            - zplus_system
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          exec:
            command:
            - pg_isready
            - -U
            - zplus_prod
            - -d
            - zplus_system
          initialDelaySeconds: 5
          periodSeconds: 5
  volumeClaimTemplates:
  - metadata:
      name: postgres-storage
    spec:
      accessModes: ["ReadWriteOnce"]
      resources:
        requests:
          storage: 20Gi
---
apiVersion: v1
kind: Service
metadata:
  name: postgres-service
  namespace: zplus-prod
spec:
  selector:
    app: postgres
  ports:
  - port: 5432
    targetPort: 5432
  type: ClusterIP
```

### 4.3 Application Deployment

**gateway.yaml:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gateway
  namespace: zplus-prod
spec:
  replicas: 3
  selector:
    matchLabels:
      app: gateway
  template:
    metadata:
      labels:
        app: gateway
    spec:
      containers:
      - name: gateway
        image: zplus/gateway:latest
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: zplus-config
              key: DB_HOST
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: zplus-secrets
              key: DB_PASSWORD
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: zplus-secrets
              key: JWT_SECRET
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: gateway-service
  namespace: zplus-prod
spec:
  selector:
    app: gateway
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
```

### 4.4 Ingress Configuration

**ingress.yaml:**
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: zplus-ingress
  namespace: zplus-prod
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/websocket-services: "gateway-service"
    nginx.ingress.kubernetes.io/proxy-read-timeout: "3600"
    nginx.ingress.kubernetes.io/proxy-send-timeout: "3600"
spec:
  tls:
  - hosts:
    - api.zplus.com
    - app.zplus.com
    - "*.zplus.com"
    secretName: zplus-tls
  rules:
  - host: api.zplus.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: gateway-service
            port:
              number: 80
  - host: app.zplus.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: frontend-service
            port:
              number: 80
  - host: "*.zplus.com"
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: frontend-service
            port:
              number: 80
```

### 4.5 Horizontal Pod Autoscaler

**hpa.yaml:**
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: gateway-hpa
  namespace: zplus-prod
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: gateway
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

### 4.6 Deployment Commands

```bash
# Apply all manifests
kubectl apply -f infra/k8s/

# Check deployment status
kubectl get pods -n zplus-prod

# View logs
kubectl logs -f deployment/gateway -n zplus-prod

# Scale deployment
kubectl scale deployment gateway --replicas=5 -n zplus-prod

# Update deployment
kubectl set image deployment/gateway gateway=zplus/gateway:v1.2.0 -n zplus-prod

# Rollback deployment
kubectl rollout undo deployment/gateway -n zplus-prod
```

## 5. CI/CD Pipeline

### 5.1 GitHub Actions Workflow

**.github/workflows/deploy.yml:**
```yaml
name: Deploy to Production

on:
  push:
    branches: [main]
    tags: ['v*']

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Run tests
      run: |
        cd apps/backend
        go test ./...
    
    - name: Setup Node.js
      uses: actions/setup-node@v4
      with:
        node-version: '18'
    
    - name: Frontend tests
      run: |
        cd apps/frontend/web
        npm ci
        npm run test

  build:
    needs: test
    runs-on: ubuntu-latest
    strategy:
      matrix:
        service: [gateway, auth-service, crm-service, frontend]
    steps:
    - uses: actions/checkout@v4
    
    - name: Setup Docker Buildx
      uses: docker/setup-buildx-action@v3
    
    - name: Login to Registry
      uses: docker/login-action@v3
      with:
        registry: ${{ env.REGISTRY }}
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}
    
    - name: Build and push
      uses: docker/build-push-action@v5
      with:
        context: .
        file: ./apps/${{ matrix.service }}/Dockerfile
        push: true
        tags: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}/${{ matrix.service }}:${{ github.sha }}
        cache-from: type=gha
        cache-to: type=gha,mode=max

  deploy:
    needs: build
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
    - uses: actions/checkout@v4
    
    - name: Setup kubectl
      uses: azure/setup-kubectl@v3
      with:
        version: 'v1.28.0'
    
    - name: Setup kubeconfig
      run: |
        mkdir -p ~/.kube
        echo "${{ secrets.KUBECONFIG }}" | base64 -d > ~/.kube/config
    
    - name: Update image tags
      run: |
        sed -i 's|image: zplus/gateway:latest|image: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}/gateway:${{ github.sha }}|' infra/k8s/gateway.yaml
        sed -i 's|image: zplus/auth-service:latest|image: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}/auth-service:${{ github.sha }}|' infra/k8s/auth-service.yaml
    
    - name: Deploy to Kubernetes
      run: |
        kubectl apply -f infra/k8s/
        kubectl rollout status deployment/gateway -n zplus-prod
        kubectl rollout status deployment/frontend -n zplus-prod
    
    - name: Verify deployment
      run: |
        kubectl get pods -n zplus-prod
        curl -f https://api.zplus.com/health
```

### 5.2 Dockerfile Examples

**Backend Dockerfile (apps/gateway/Dockerfile):**
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./apps/gateway

# Production stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/apps/gateway/templates ./templates

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./main"]
```

**Frontend Dockerfile (apps/frontend/web/Dockerfile):**
```dockerfile
# Dependencies stage
FROM node:18-alpine AS deps
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

# Build stage
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Production stage
FROM node:18-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

ENV PORT=3000

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000/api/health || exit 1

CMD ["node", "server.js"]
```

## 6. Database Migration trong Production

### 6.1 Migration Strategy

```bash
# Create migration job
apiVersion: batch/v1
kind: Job
metadata:
  name: migrate-database
  namespace: zplus-prod
spec:
  template:
    spec:
      containers:
      - name: migrate
        image: zplus/migrate:latest
        env:
        - name: DB_URL
          value: "postgres://user:pass@postgres-service:5432/zplus_system?sslmode=require"
        command: ["migrate", "-path", "/migrations", "-database", "$(DB_URL)", "up"]
      restartPolicy: Never
  backoffLimit: 3
```

### 6.2 Rollback Strategy

```bash
# Rollback deployment
kubectl rollout undo deployment/gateway -n zplus-prod

# Rollback database (if needed)
kubectl run migrate-rollback --image=zplus/migrate:latest --rm -it -- \
  migrate -path /migrations -database "postgres://..." down 1
```

## 7. Monitoring và Logging

### 7.1 Prometheus & Grafana

**prometheus.yaml:**
```yaml
global:
  scrape_interval: 15s

scrape_configs:
- job_name: 'zplus-gateway'
  static_configs:
  - targets: ['gateway-service:8080']
  metrics_path: /metrics

- job_name: 'kubernetes-pods'
  kubernetes_sd_configs:
  - role: pod
  relabel_configs:
  - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
    action: keep
    regex: true
```

### 7.2 ELK Stack

**filebeat.yaml:**
```yaml
filebeat.inputs:
- type: container
  paths:
    - /var/log/containers/*zplus*.log

output.elasticsearch:
  hosts: ["elasticsearch:9200"]

setup.kibana:
  host: "kibana:5601"
```

### 7.3 Health Check Endpoints

```go
// Health check handler
func HealthHandler(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "status": "healthy",
        "timestamp": time.Now(),
        "uptime": time.Since(startTime).String(),
    })
}

// Readiness check
func ReadinessHandler(c *fiber.Ctx) error {
    // Check database connection
    if err := db.Ping(); err != nil {
        return c.Status(503).JSON(fiber.Map{
            "status": "not ready",
            "error": "database connection failed",
        })
    }
    
    // Check Redis connection
    if err := redis.Ping(); err != nil {
        return c.Status(503).JSON(fiber.Map{
            "status": "not ready",
            "error": "redis connection failed",
        })
    }
    
    return c.JSON(fiber.Map{
        "status": "ready",
        "checks": map[string]string{
            "database": "ok",
            "redis": "ok",
        },
    })
}
```

## 8. Backup và Disaster Recovery

### 8.1 Database Backup

**backup-cronjob.yaml:**
```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: postgres-backup
  namespace: zplus-prod
spec:
  schedule: "0 2 * * *" # Daily at 2 AM
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: postgres-backup
            image: postgres:15-alpine
            env:
            - name: PGPASSWORD
              valueFrom:
                secretKeyRef:
                  name: zplus-secrets
                  key: DB_PASSWORD
            command:
            - /bin/bash
            - -c
            - |
              DATE=$(date +%Y%m%d_%H%M%S)
              pg_dump -h postgres-service -U zplus_prod zplus_system > /backup/backup_$DATE.sql
              aws s3 cp /backup/backup_$DATE.sql s3://zplus-backups/
            volumeMounts:
            - name: backup-storage
              mountPath: /backup
          volumes:
          - name: backup-storage
            emptyDir: {}
          restartPolicy: OnFailure
```

### 8.2 File Storage Backup

```bash
# S3 to S3 backup
aws s3 sync s3://zplus-files-prod s3://zplus-backup-files --delete

# Backup to different region
aws s3 sync s3://zplus-files-prod s3://zplus-files-backup-eu --delete
```

## 9. Security trong Production

### 9.1 Network Policies

**network-policy.yaml:**
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: zplus-network-policy
  namespace: zplus-prod
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: postgres
    ports:
    - protocol: TCP
      port: 5432
  - to:
    - podSelector:
        matchLabels:
          app: redis
    ports:
    - protocol: TCP
      port: 6379
```

### 9.2 Pod Security Policies

**pod-security.yaml:**
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: gateway
spec:
  securityContext:
    runAsNonRoot: true
    runAsUser: 1001
    fsGroup: 1001
  containers:
  - name: gateway
    securityContext:
      allowPrivilegeEscalation: false
      readOnlyRootFilesystem: true
      capabilities:
        drop:
        - ALL
```

## 10. Performance Tuning

### 10.1 Resource Limits

```yaml
resources:
  requests:
    memory: "256Mi"
    cpu: "250m"
  limits:
    memory: "512Mi"
    cpu: "500m"
```

### 10.2 Database Tuning

```sql
-- PostgreSQL configuration for production
shared_buffers = 256MB
effective_cache_size = 1GB
maintenance_work_mem = 64MB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1
effective_io_concurrency = 200
```

### 10.3 CDN Configuration

```yaml
# CloudFront distribution
Resources:
  CloudFrontDistribution:
    Type: AWS::CloudFront::Distribution
    Properties:
      DistributionConfig:
        Origins:
        - DomainName: app.zplus.com
          Id: zplus-origin
          CustomOriginConfig:
            HTTPPort: 443
            OriginProtocolPolicy: https-only
        DefaultCacheBehavior:
          TargetOriginId: zplus-origin
          ViewerProtocolPolicy: redirect-to-https
          CachePolicyId: 4135ea2d-6df8-44a3-9df3-4b5a84be39ad
        Enabled: true
        HttpVersion: http2
```

## 11. Troubleshooting Production Issues

### 11.1 Common Issues

**Pod Crash Loop:**
```bash
# Check pod logs
kubectl logs -f pod/gateway-xxx -n zplus-prod

# Check pod events
kubectl describe pod gateway-xxx -n zplus-prod

# Check resource usage
kubectl top pod -n zplus-prod
```

**Database Connection Issues:**
```bash
# Test database connectivity
kubectl run pg-test --image=postgres:15-alpine --rm -it -- \
  psql -h postgres-service -U zplus_prod -d zplus_system

# Check database logs
kubectl logs -f statefulset/postgres -n zplus-prod
```

**Performance Issues:**
```bash
# Check resource usage
kubectl top nodes
kubectl top pods -n zplus-prod

# Check HPA status
kubectl get hpa -n zplus-prod

# Scale manually if needed
kubectl scale deployment gateway --replicas=5 -n zplus-prod
```

### 11.2 Emergency Procedures

**Service Outage:**
```bash
# Scale up immediately
kubectl scale deployment gateway --replicas=10 -n zplus-prod

# Switch to maintenance mode
kubectl patch ingress zplus-ingress -p '{"metadata":{"annotations":{"nginx.ingress.kubernetes.io/default-backend":"maintenance-service"}}}'

# Rollback to last known good version
kubectl rollout undo deployment/gateway -n zplus-prod
```

**Database Issues:**
```bash
# Switch to read-only mode
kubectl set env deployment/gateway READ_ONLY_MODE=true

# Restore from backup
kubectl create job --from=cronjob/postgres-backup restore-db
```

## 12. Maintenance Windows

### 12.1 Planned Maintenance

```bash
# 1. Notify users (maintenance mode)
kubectl apply -f maintenance-page.yaml

# 2. Scale down to minimum
kubectl scale deployment gateway --replicas=1 -n zplus-prod

# 3. Run migrations
kubectl apply -f migration-job.yaml

# 4. Update applications
kubectl set image deployment/gateway gateway=zplus/gateway:v1.2.0

# 5. Scale back up
kubectl scale deployment gateway --replicas=3 -n zplus-prod

# 6. Remove maintenance mode
kubectl delete -f maintenance-page.yaml
```

### 12.2 Zero-Downtime Deployment

```bash
# Rolling update with health checks
kubectl set image deployment/gateway gateway=zplus/gateway:v1.2.0 --record
kubectl rollout status deployment/gateway -n zplus-prod

# Monitor during rollout
watch kubectl get pods -n zplus-prod
```
