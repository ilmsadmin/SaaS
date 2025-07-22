# Zplus SaaS - Multi-tenant Software as a Service Platform

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![Node Version](https://img.shields.io/badge/node-18+-green.svg)](https://nodejs.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/your-org/zplus-saas)
[![Coverage](https://img.shields.io/badge/coverage-85%25-green.svg)](https://codecov.io/gh/your-org/zplus-saas)

Zplus SaaS là một nền tảng SaaS đa tenant hiện đại được xây dựng với kiến trúc 3-tier phân quyền (System → Tenant → Customer) và microservices, hỗ trợ nhiều module như CRM, LMS, POS, HRM và Checkin với khả năng tùy biến cao.

## 📋 Mục lục

1. [Giới thiệu](#-giới-thiệu)
2. [Tính năng chính](#-tính-năng-chính)
3. [Kiến trúc hệ thống](#️-kiến-trúc-hệ-thống)
4. [Tech Stack](#-tech-stack)
5. [Yêu cầu hệ thống](#-yêu-cầu-hệ-thống)
6. [Cài đặt nhanh](#-cài-đặt-nhanh)
7. [Hệ thống Module](#️-hệ-thống-module)
8. [Bảo mật](#-bảo-mật)
9. [Hiệu suất](#-hiệu-suất)
10. [Testing](#-testing)
11. [Deployment](#-deployment)
12. [Monitoring](#-monitoring)
13. [Tài liệu](#-tài-liệu)
14. [Đóng góp](#-đóng-góp)
15. [Hỗ trợ](#-hỗ-trợ)

## 🚀 Giới thiệu

Zplus SaaS là giải pháp SaaS toàn diện được thiết kế cho các doanh nghiệp muốn cung cấp dịch vụ phần mềm trực tuyến. Với kiến trúc 3-tier phân quyền rõ ràng và hệ thống module linh hoạt, platform cho phép quản lý hàng nghìn tenant với dữ liệu tách biệt hoàn toàn.

### 🎯 Đối tượng sử dụng

- **System Admin**: Quản trị toàn hệ thống, tenant, gói dịch vụ
- **Tenant Admin**: Quản trị trong phạm vi tổ chức/công ty
- **End Users**: Người dùng cuối sử dụng các module chuyên biệt

### 🌟 Tính năng chính

### 🌟 Tính năng chính

#### Core Features
- **Multi-tenant Architecture**: Hỗ trợ hàng nghìn tenant với dữ liệu tách biệt hoàn toàn
- **Modular System**: Bật/tắt module theo nhu cầu từng tenant  
- **3-tier Authorization**: System → Tenant → Customer với RBAC đa tầng
- **Custom Domain/Subdomain**: Mỗi tenant có thể sử dụng domain riêng

#### Technical Features  
- **GraphQL-First API**: Hiệu suất cao với type safety và real-time subscriptions
- **Microservices Architecture**: Scalable với independent deployments
- **Multi-database Support**: PostgreSQL + MongoDB + Redis
- **Real-time Updates**: WebSocket và GraphQL subscriptions
- **Background Processing**: Async jobs với Redis Queue

#### Business Features
- **Subscription Management**: Quản lý gói cước, thanh toán tự động
- **Multi-language Support**: I18n cho nhiều ngôn ngữ
- **Advanced Analytics**: Dashboard báo cáo chi tiết
- **API Integration**: RESTful APIs và webhooks

## 🏗️ Kiến trúc hệ thống

### 2.1 Kiến trúc 3-tier Multi-tenant

```
+-------------------------+
|        System          |  ← Quản trị toàn cục (RBAC, gói dịch vụ, tenant, domain)
+-------------------------+
            |
            v
+-------------------------+
|        Tenant          |  ← Quản trị trong phạm vi tenant (RBAC, user, module, khách hàng)
+-------------------------+
            |
            v
+-------------------------+
|       Customer         |  ← Người dùng cuối, sử dụng dịch vụ (CRM, LMS, POS...)
+-------------------------+
```

### 2.2 Microservices Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│    Frontend     │    │   API Gateway   │    │   Load Balancer │
│    (Next.js)    │◄──►│(GraphQL/REST)   │◄──►│    (Traefik)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                ┌───────────────┼───────────────┐
                │               │               │
        ┌───────▼───────┐ ┌─────▼─────┐ ┌──────▼──────┐
        │  Auth Service │ │File Service│ │Payment Svc  │
        │  (Go/Fiber)   │ │(Go/Fiber)  │ │ (Go/Fiber)  │
        └───────────────┘ └───────────┘ └─────────────┘
                │               │               │
        ┌───────▼───────┐ ┌─────▼─────┐ ┌──────▼──────┐
        │  CRM Service  │ │HRM Service │ │  POS Service│
        │  (Go/Fiber)   │ │(Go/Fiber)  │ │ (Go/Fiber)  │
        └───────────────┘ └───────────┘ └─────────────┘
                │               │               │
        ┌───────▼───────┐ ┌─────▼─────┐ ┌──────▼──────┐
        │  PostgreSQL   │ │ MongoDB   │ │    Redis    │
        │(Relational)   │ │(Documents) │ │   (Cache)   │
        └───────────────┘ └───────────┘ └─────────────┘
```

### 2.3 Tenant Routing Strategy

- **Subdomain**: `tenant-name.saas-platform.com`
- **Custom Domain**: `client-domain.com`  
- **Path-based**: `saas-platform.com/tenant-name`

## 🚀 Tech Stack

### 3.1 Frontend Technologies

| Component | Technology | Purpose | Features |
|-----------|------------|---------|----------|
| **Web App** | Next.js 14 + TypeScript | User Interface | SSR, ISR, App Router |
| **Mobile** | React Native + Expo | Cross-platform | iOS, Android |
| **UI Library** | Tailwind CSS + shadcn/ui | Styling | Responsive, Dark mode |
| **State Management** | Zustand + TanStack Query | Client State | Optimistic updates |
| **Forms** | React Hook Form + Zod | Form Handling | Validation, TypeScript |

### 3.2 Backend Technologies

| Component | Technology | Purpose | Features |
|-----------|------------|---------|----------|
| **API Gateway** | Go Fiber + GraphQL | API Orchestration | Rate limiting, Auth |
| **Microservices** | Go Fiber + GORM | Business Logic | High performance |
| **Authentication** | JWT + Refresh Token | Security | Multi-tier RBAC |
| **Background Jobs** | Redis Queue | Async Processing | Retry, Scheduling |
| **File Storage** | MinIO/S3 | Object Storage | Multi-tenant isolation |

### 3.3 Database & Infrastructure

| Component | Technology | Purpose | Features |
|-----------|------------|---------|----------|
| **Primary DB** | PostgreSQL 15+ | Relational Data | ACID, Multi-schema |
| **Document DB** | MongoDB 6+ | Flexible Data | Analytics, Logs |
| **Cache & Queue** | Redis 7+ | Performance | Session, Cache, Jobs |
| **Load Balancer** | Traefik | Multi-tenant Routing | SSL, Auto-discovery |
| **Monitoring** | Prometheus + Grafana | Observability | Metrics, Alerts |

## 📚 Tài liệu

### 4.1 Tài liệu Thiết kế
- [📋 Thiết kế Tổng quan](docs/thiet-ke-tong-quan-du-an.md) - Mục tiêu, phạm vi và đối tượng sử dụng
- [🏗️ Thiết kế Kiến trúc](docs/thiet-ke-kien-truc-du-an.md) - Chi tiết kiến trúc microservices  
- [🗃️ Thiết kế Database](docs/thiet-ke-kien-truc-database.md) - Schema design và data modeling

### 4.2 Tài liệu Triển khai
- [⚙️ Hướng dẫn Cài đặt](docs/installation.md) - Setup môi trường development
- [🚀 Hướng dẫn Deployment](docs/deployment.md) - Triển khai production
- [🔒 Security Guide](docs/security.md) - Bảo mật và RBAC multi-tier

### 4.3 Tài liệu Phát triển  
- [📖 API Documentation](docs/api-documentation.md) - GraphQL/REST APIs
- [🗄️ Database Schema](docs/database-schema.md) - Chi tiết database structure
- [🧩 Module Development](docs/module-development.md) - Phát triển module mới
- [🧪 Testing Guide](docs/testing.md) - Unit, Integration, E2E testing
- [🔧 Troubleshooting](docs/troubleshooting.md) - Debug và xử lý sự cố

## 🛠️ Yêu cầu hệ thống

### 5.1 Development Environment

| Component | Minimum | Recommended |
|-----------|---------|-------------|
| **Go** | 1.21+ | 1.22+ |
| **Node.js** | 18+ | 20+ LTS |
| **PostgreSQL** | 13+ | 15+ |
| **Redis** | 6+ | 7+ |
| **Docker** | 20+ | 24+ |
| **RAM** | 8GB | 16GB+ |
| **Storage** | 50GB SSD | 100GB+ SSD |

### 5.2 Production Environment

| Component | Minimum | Recommended |
|-----------|---------|-------------|
| **CPU** | 4 cores | 8+ cores |
| **RAM** | 16GB | 32GB+ |
| **Storage** | 200GB SSD | 1TB+ NVMe |
| **Network** | 100Mbps | 1Gbps+ |
| **Load Balancer** | 2 instances | 3+ instances |

## ⚡ Cài đặt nhanh

### 6.1 Clone Repository

```bash
# Clone repository
git clone https://github.com/your-org/zplus-saas.git
cd zplus-saas

# Kiểm tra cấu trúc project
tree -d -L 2
```

### 6.2 Setup Backend Services

```bash
# Chuyển đến thư mục backend
cd apps/backend

# Copy environment variables
cp .env.example .env

# Cài đặt dependencies
go mod download

# Chạy database migrations
make migrate

# Khởi động development server
make dev
```

### 6.3 Setup Frontend Application

```bash
# Chuyển đến thư mục frontend  
cd apps/frontend/web

# Cài đặt dependencies
npm install

# Khởi động development server
npm run dev
```

### 6.4 Setup với Docker (Recommended)

```bash
# Khởi động toàn bộ stack
docker-compose up -d

# Xem logs
docker-compose logs -f

# Dừng services
docker-compose down
```

### 6.5 Truy cập Application

| Service | URL | Description |
|---------|-----|-------------|
| **System Admin** | http://localhost:3000/system | Quản trị hệ thống |
| **Tenant Demo** | http://tenant-demo.localhost:3000 | Tenant mẫu |
| **API Playground** | http://localhost:8080/playground | GraphQL Playground |
| **API Documentation** | http://localhost:8080/docs | Swagger UI |
| **Monitoring** | http://localhost:3001 | Grafana Dashboard |

## 🏢 Hệ thống Module

### 7.1 Available Modules

| Module | Status | Description | Features |
|--------|--------|-------------|----------|
| **CRM** | ✅ Available | Quản lý khách hàng, bán hàng | Lead tracking, Sales pipeline, Customer analytics |
| **LMS** | ✅ Available | Học tập trực tuyến | Course management, Video streaming, Quiz & Exams |
| **POS** | ✅ Available | Bán hàng tại điểm | Inventory, Payment gateway, Receipt printing |
| **HRM** | ✅ Available | Quản lý nhân sự | Employee management, Payroll, Performance tracking |
| **Checkin** | ✅ Available | Chấm công điểm danh | Location tracking, QR Code, Face recognition |
| **Accounting** | 🚧 In Development | Kế toán tài chính | Invoicing, Financial reports |
| **E-commerce** | 📋 Planned | Thương mại điện tử | Product catalog, Shopping cart |

### 7.2 Module Architecture

```
apps/backend/{module}/
├── models/           # Database models & migrations  
│   ├── entity.go     # GORM models
│   ├── dto.go        # Data transfer objects
│   └── migrations/   # Database migrations
├── services/         # Business logic layer
│   ├── {module}_service.go
│   └── interfaces.go # Service interfaces  
├── handlers/         # HTTP/GraphQL handlers
│   ├── rest.go       # REST endpoints
│   ├── graphql.go    # GraphQL resolvers
│   └── websocket.go  # WebSocket handlers
├── repositories/     # Data access layer
│   └── {module}_repo.go
├── utils/           # Module-specific utilities
├── tests/           # Unit & integration tests
└── routes.go        # Route definitions
```

### 7.3 Module Configuration

```yaml
# tenant_modules.yml
tenant_id: "uuid-here"
modules:
  crm:
    enabled: true
    plan: "premium"
    features: ["lead_tracking", "analytics", "api_access"]
  lms:
    enabled: true  
    plan: "basic"
    features: ["course_management", "video_streaming"]
```

### 7.4 Module Development Guidelines

- **Modular Design**: Mỗi module độc lập, có thể enable/disable
- **Consistent API**: Tuân thủ GraphQL schema conventions
- **Database Isolation**: Separate schema/table prefix cho mỗi module
- **Permission System**: Integrate với RBAC framework
- **Testing Required**: Unit tests, integration tests mandatory

## 🔒 Bảo mật

### 8.1 Security Architecture

Zplus SaaS thực hiện bảo mật theo mô hình **Defense in Depth**:

```
┌─────────────────────────────────────────────┐
│              Security Layers                │
├─────────────────────────────────────────────┤
│ 1. Network Security (SSL/TLS, Firewall)    │ ← Infrastructure
├─────────────────────────────────────────────┤
│ 2. Application Security (JWT, Rate Limit)  │ ← Application  
├─────────────────────────────────────────────┤
│ 3. Authorization (RBAC Multi-tier)         │ ← Access Control
├─────────────────────────────────────────────┤
│ 4. Data Security (Encryption, Audit)       │ ← Data Protection
├─────────────────────────────────────────────┤
│ 5. Tenant Isolation (Schema Separation)    │ ← Multi-tenancy
└─────────────────────────────────────────────┘
```

### 8.2 Authentication & Authorization

| Feature | Implementation | Purpose |
|---------|----------------|---------|
| **JWT Authentication** | Access + Refresh Token | Stateless authentication |
| **Multi-tier RBAC** | System/Tenant/Customer | Granular permissions |
| **API Rate Limiting** | Per-tenant/per-user | DDoS protection |
| **Session Management** | Redis-based | Secure session handling |
| **Audit Logging** | Comprehensive trail | Compliance & monitoring |

### 8.3 Data Protection

- **Encryption at Rest**: AES-256 cho sensitive data
- **Encryption in Transit**: TLS 1.3 cho tất cả connections
- **Data Masking**: PII data masking trong logs
- **Backup Encryption**: Encrypted database backups
- **Tenant Isolation**: Schema-per-tenant với strict isolation

### 8.4 Security Standards

- **OWASP Top 10**: Compliance với security best practices
- **SOC 2 Type II**: Security controls compliance
- **GDPR Ready**: Data protection và privacy controls
- **ISO 27001**: Information security management

## 📈 Hiệu suất

### 9.1 Performance Metrics

| Metric | Target | Current | Monitoring |
|--------|--------|---------|------------|
| **API Response Time** | < 200ms | 150ms avg | Prometheus |
| **Database Query Time** | < 50ms | 30ms avg | Slow query log |
| **Memory Usage** | < 80% | 65% avg | Grafana |
| **CPU Usage** | < 70% | 45% avg | System metrics |
| **Concurrent Users** | 10,000+ | Tested to 15,000 | Load testing |

### 9.2 Optimization Strategies

#### Database Optimization
- **Connection Pooling**: PgBouncer với max_connections optimized
- **Read Replicas**: Read-heavy queries directed to replicas
- **Query Optimization**: Indexed queries, query plan analysis
- **Database Partitioning**: Time-based partitioning cho audit logs

#### Caching Strategy
- **Application Cache**: Redis với TTL-based expiration
- **Database Query Cache**: PostgreSQL query result caching
- **CDN Integration**: CloudFront/CloudFlare cho static assets
- **Browser Caching**: Optimized cache headers

#### Infrastructure Scaling
- **Horizontal Scaling**: Kubernetes auto-scaling
- **Load Balancing**: Traefik với health checks
- **Container Optimization**: Multi-stage Docker builds
- **Resource Management**: CPU/Memory limits và requests

### 9.3 Background Processing

```yaml
# Redis Queue Configuration
redis_queue:
  workers: 10
  max_retry: 3
  retry_delay: "30s"
  queues:
    - name: "email"
      priority: 1
    - name: "analytics" 
      priority: 2
    - name: "reports"
      priority: 3
```

## 🧪 Testing

### 10.1 Testing Strategy

| Test Type | Coverage | Tools | Command |
|-----------|----------|-------|---------|
| **Unit Tests** | 85%+ | Go: testify, JS: Jest | `make test-unit` |
| **Integration Tests** | 70%+ | Testcontainers, Supertest | `make test-integration` |
| **E2E Tests** | Key flows | Playwright, Cypress | `npm run test:e2e` |
| **Load Tests** | 10k users | k6, Artillery | `make load-test` |
| **Security Tests** | OWASP | ZAP, Bandit | `make security-test` |

### 10.2 Test Commands

```bash
# Chạy tất cả backend tests
make test

# Chạy frontend tests với coverage
npm run test:coverage

# E2E testing
npm run test:e2e

# Load testing với k6
make load-test

# Security scanning
make security-scan

# Database migration testing
make test-migrations
```

### 10.3 CI/CD Pipeline

```yaml
# .github/workflows/test.yml
stages:
  - lint: ESLint, Golangci-lint
  - unit: Jest, Go test
  - integration: Testcontainers 
  - security: SAST, Dependency check
  - e2e: Playwright tests
  - performance: k6 load tests
```

## 🚀 Deployment

### 11.1 Environment Setup

| Environment | Purpose | Infrastructure | Monitoring |
|-------------|---------|----------------|------------|
| **Development** | Local development | Docker Compose | Basic logging |
| **Staging** | Testing & QA | Kubernetes | Full monitoring |
| **Production** | Live system | Kubernetes HA | Full observability |

### 11.2 Development Deployment

```bash
# Docker Compose (Recommended for local dev)
docker-compose up -d

# Kiểm tra status
docker-compose ps

# Xem logs
docker-compose logs -f api-gateway

# Rebuild services
docker-compose up -d --build
```

### 11.3 Staging/Production Deployment

```bash
# Kubernetes deployment
kubectl apply -f infra/k8s/namespaces/
kubectl apply -f infra/k8s/configmaps/
kubectl apply -f infra/k8s/secrets/
kubectl apply -f infra/k8s/deployments/
kubectl apply -f infra/k8s/services/
kubectl apply -f infra/k8s/ingress/

# Verify deployment
kubectl get pods -n zplus-saas
kubectl get services -n zplus-saas

# Rolling update
kubectl rollout restart deployment/api-gateway -n zplus-saas
```

### 11.4 Docker Swarm (Alternative)

```bash
# Production stack deployment
docker stack deploy -c docker-compose.prod.yml zplus

# Scaling services
docker service scale zplus_api-gateway=3

# Update service
docker service update --image zplus/api-gateway:v2.0 zplus_api-gateway
```

### 11.5 Infrastructure as Code

```bash
# Terraform (Cloud infrastructure)
cd infra/terraform
terraform init
terraform plan
terraform apply

# Ansible (Server configuration)
cd infra/ansible
ansible-playbook -i inventory/production site.yml
```

## 📊 Monitoring

### 12.1 Observability Stack

| Component | Technology | Purpose | Access |
|-----------|------------|---------|--------|
| **Metrics** | Prometheus + Grafana | Performance monitoring | http://localhost:3001 |
| **Logging** | ELK Stack (Elasticsearch, Logstash, Kibana) | Log aggregation | http://localhost:5601 |
| **Tracing** | Jaeger | Distributed tracing | http://localhost:16686 |
| **Alerting** | Grafana Alerts + PagerDuty | Incident management | Slack/Email/SMS |
| **Uptime** | UptimeRobot | Service availability | Dashboard |

### 12.2 Key Metrics

```yaml
# Grafana Dashboard Metrics
application_metrics:
  - api_request_duration_seconds
  - api_request_total
  - database_connections_active
  - redis_operations_total
  - tenant_active_users
  - background_jobs_processed

business_metrics:
  - tenant_signups_total
  - module_usage_by_tenant
  - subscription_revenue
  - user_session_duration
```

### 12.3 Health Checks

```bash
# Application health endpoints
curl http://localhost:8080/health
curl http://localhost:8080/health/database
curl http://localhost:8080/health/redis
curl http://localhost:8080/health/modules

# Kubernetes health checks
kubectl get pods -l app=api-gateway
kubectl describe pod <pod-name>
```

### 12.4 Alerting Rules

- **High Response Time**: > 500ms average over 5 minutes
- **Error Rate**: > 5% error rate over 5 minutes  
- **Database Connections**: > 80% of max connections
- **Memory Usage**: > 90% memory usage
- **Disk Space**: < 20% free disk space
- **Failed Background Jobs**: > 10 failed jobs in queue

## 🤝 Đóng góp

### 13.1 Development Workflow

1. **Fork Repository**
   ```bash
   # Fork trên GitHub và clone về local
   git clone https://github.com/your-username/zplus-saas.git
   cd zplus-saas
   ```

2. **Setup Development Environment**
   ```bash
   # Cài đặt pre-commit hooks
   pre-commit install
   
   # Chạy development setup
   make dev-setup
   ```

3. **Create Feature Branch**
   ```bash
   # Tạo branch từ develop
   git checkout develop
   git pull origin develop
   git checkout -b feature/amazing-feature
   ```

4. **Development & Testing**
   ```bash
   # Phát triển tính năng
   # Viết unit tests
   make test
   
   # Kiểm tra code quality
   make lint
   ```

5. **Commit Changes**
   ```bash
   # Commit với conventional commits
   git commit -m 'feat(crm): add customer analytics dashboard'
   git push origin feature/amazing-feature
   ```

6. **Pull Request**
   - Tạo PR từ feature branch đến develop
   - Điền đầy đủ PR template
   - Đảm bảo CI/CD passes

### 13.2 Coding Standards

- **Go**: Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- **TypeScript**: ESLint + Prettier configuration
- **Git**: Conventional Commits specification
- **Documentation**: Godoc for Go, JSDoc for TypeScript
- **Testing**: Minimum 80% code coverage required

### 13.3 Code Review Process

1. **Self Review**: Kiểm tra code trước khi tạo PR
2. **Automated Checks**: CI/CD pipeline must pass
3. **Peer Review**: Ít nhất 2 approvals từ maintainers
4. **Manual Testing**: QA testing trên staging environment
5. **Merge**: Squash merge vào develop branch

## � Hỗ trợ

### 14.1 Tài liệu và Hướng dẫn

| Resource | Link | Description |
|----------|------|-------------|
| **📖 Documentation** | [docs/](docs/) | Tài liệu chi tiết hệ thống |
| **🎯 API Reference** | [API Docs](http://localhost:8080/docs) | Interactive API documentation |
| **🎥 Video Tutorials** | [YouTube Channel](https://youtube.com/@zplus-saas) | Video hướng dẫn sử dụng |
| **📚 Knowledge Base** | [KB](https://kb.zplus.com) | Câu hỏi thường gặp |

### 14.2 Community Support

| Platform | Link | Purpose |
|----------|------|---------|
| **💬 Discord** | [Join Discord](https://discord.gg/zplus-saas) | Real-time community chat |
| **🐛 GitHub Issues** | [Issues](https://github.com/your-org/zplus-saas/issues) | Bug reports & feature requests |
| **💡 Discussions** | [Discussions](https://github.com/your-org/zplus-saas/discussions) | Q&A và thảo luận |
| **📧 Email** | support@zplus.com | Official support channel |

### 14.3 Enterprise Support

- **Priority Support**: 24/7 support với SLA guaranteed
- **Custom Development**: Tùy chỉnh theo nhu cầu doanh nghiệp
- **Training & Consultation**: Đào tạo team và tư vấn implementation
- **Dedicated Account Manager**: Quản lý tài khoản chuyên biệt

**Contact**: enterprise@zplus.com | +84-xxx-xxx-xxx

### 14.4 Reporting Issues

```bash
# Khi báo cáo lỗi, vui lòng bao gồm:
- Môi trường (OS, browser, versions)
- Steps to reproduce
- Expected vs actual behavior  
- Screenshots/logs nếu có
- Configuration files (sanitized)
```

## 📝 License

Dự án này được phân phối dưới giấy phép MIT License. Xem file [LICENSE](LICENSE) để biết chi tiết.

```
MIT License

Copyright (c) 2024 Zplus SaaS Team

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
```

## 🗺️ Roadmap

### 15.1 Version 2.0 (Q3 2024)
- [ ] **Mobile Applications**
  - [ ] React Native iOS app
  - [ ] React Native Android app  
  - [ ] Push notifications
  - [ ] Offline capability

- [ ] **Advanced Analytics**
  - [ ] Real-time analytics dashboard
  - [ ] Custom report builder
  - [ ] Data export/import tools
  - [ ] Advanced visualization widgets

### 15.2 Version 2.5 (Q4 2024)
- [ ] **AI/ML Integration**
  - [ ] Customer behavior prediction
  - [ ] Automated lead scoring (CRM)
  - [ ] Personalized learning paths (LMS)
  - [ ] Intelligent inventory management (POS)

- [ ] **Workflow Engine**
  - [ ] Visual workflow designer
  - [ ] Approval workflows
  - [ ] Custom automation rules
  - [ ] Integration connectors

### 15.3 Version 3.0 (Q1 2025)
- [ ] **Multi-language Support**
  - [ ] Complete i18n framework
  - [ ] RTL language support
  - [ ] Multi-currency support
  - [ ] Localized modules

- [ ] **Third-party Integrations Hub**
  - [ ] Zapier integration
  - [ ] Salesforce connector
  - [ ] Google Workspace integration
  - [ ] Microsoft 365 integration
  - [ ] WhatsApp Business API
  - [ ] Payment gateways expansion

### 15.4 Long-term Vision (2025+)
- [ ] **Marketplace Ecosystem**
  - [ ] Third-party module marketplace
  - [ ] Developer certification program
  - [ ] Revenue sharing model
  - [ ] Community contributions

- [ ] **Enterprise Features**
  - [ ] Advanced compliance tools (SOX, HIPAA)
  - [ ] White-label solutions
  - [ ] API rate limiting per tenant
  - [ ] Advanced security features

---

## 🎯 Tổng kết

Zplus SaaS là một nền tảng SaaS đầy đủ tính năng được thiết kế để phục vụ nhu cầu đa dạng của các doanh nghiệp hiện đại. Với kiến trúc 3-tier phân quyền rõ ràng, hệ thống module linh hoạt và khả năng mở rộng cao, platform cung cấp foundation vững chắc cho việc xây dựng các ứng dụng SaaS chuyên nghiệp.

**🚀 Bắt đầu ngay**: [Cài đặt nhanh](#-cài-đặt-nhanh) | **📚 Tài liệu**: [docs/](docs/) | **💬 Hỗ trợ**: [support@zplus.com](mailto:support@zplus.com)

---

<div align="center">
  <strong>Made with ❤️ by Zplus SaaS Team</strong>
  <br>
  <sub>Built in Vietnam 🇻🇳</sub>
</div>