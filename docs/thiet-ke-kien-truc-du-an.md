# Thiết kế Kiến trúc Dự án - Zplus SaaS

## 1. Tổng quan Kiến trúc

Zplus SaaS được xây dựng trên kiến trúc **3-tier Multi-tenant** với sự phân tách rõ ràng giữa các tầng:

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

## 2. Kiến trúc Microservices

### 2.1 Sơ đồ Tổng quan

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

### 2.2 Technology Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Frontend** | Next.js | User Interface |
| **API Gateway** | GraphQL hoặc REST gateway | API Orchestration |
| **Load Balancer** | Traefik | Traffic Routing & Subdomain Management |
| **Backend Services** | Go Fiber + GORM | Business Logic |
| **Database** | PostgreSQL + MongoDB | Data Persistence |
| **Cache** | Redis | Session & Cache |
| **Message Queue** | Redis Queue | Async Processing & Background Jobs |
| **Auth** | JWT + RBAC | Authentication |

### 2.3 Vai trò chi tiết của Traefik

**Traefik** đóng vai trò quan trọng trong kiến trúc multi-tenant:

1. **Reverse Proxy & Load Balancer**: 
   - Phân phối tải giữa các instances của application
   - SSL/TLS termination
   - Health checks và failover

2. **Multi-tenant Routing**:
   - Phân tích subdomain từ request (`tenant-slug.domain.com`)
   - Inject tenant context vào header (`X-Tenant-ID`)
   - Route requests đến appropriate backend services

3. **Dynamic Configuration**:
   - Auto-discovery services trong Docker/Kubernetes
   - Dynamic SSL certificate management
   - Rate limiting và security middlewares

### 2.4 Message Queue Strategy

**Tại sao chọn Redis thay vì Kafka?**

| Aspect | Redis Queue | Kafka | Quyết định |
|--------|-------------|-------|------------|
| **Complexity** | Đơn giản, ít setup | Phức tạp, cần nhiều components | ✅ Redis cho MVP |
| **Latency** | Thấp (< 1ms) | Cao hơn (5-10ms) | ✅ Redis cho real-time |
| **Throughput** | Trung bình (100K msg/s) | Cao (1M+ msg/s) | Redis đủ cho scale hiện tại |
| **Persistence** | In-memory, optional persistence | Durable, replicated | Kafka cho critical events |
| **Use Cases** | Cache, sessions, simple queues | Event streaming, big data | Phù hợp với multi-tenant SaaS |

**Khi nào sẽ migrate sang Kafka?**
- Khi cần xử lý > 1M events/second
- Cần event streaming và real-time analytics
- Yêu cầu event sourcing patterns
- Audit trail và compliance requirements

## 3. Chiến lược API: GraphQL-First

### 3.1 Tại sao ưu tiên GraphQL?

Zplus SaaS áp dụng chiến lược **GraphQL-first** làm giao thức API chính, kết hợp với REST API cho các trường hợp cụ thể.

#### Lợi ích của GraphQL-First:

1. **Hiệu quả Frontend**: Một request duy nhất cho nhiều nguồn dữ liệu
2. **Type Safety**: Hệ thống type mạnh mẽ với auto-generated TypeScript
3. **Khả năng Real-time**: Hỗ trợ subscription tích hợp sẵn
4. **Multi-tenant Aware**: Schema-level tenant isolation
5. **Developer Experience**: GraphQL Playground và introspection

#### Thiết kế Schema GraphQL

```graphql
# Interface cơ bản cho multi-tenant
interface TenantEntity {
  id: ID!
  tenantId: TenantID!
  createdAt: DateTime!
  updatedAt: DateTime!
}

# Ví dụ: Entity khách hàng
type Customer implements TenantEntity {
  id: ID!
  tenantId: TenantID!
  name: String!
  email: String
  status: CustomerStatus!
  company: String
  createdAt: DateTime!
  updatedAt: DateTime!
}
```

#### Multi-tenant Context trong GraphQL

```go
// Request context với tenant isolation
type RequestContext struct {
    Tenant *TenantContext `json:"tenant"`
    User   *UserContext   `json:"user"`
}

// Authorization trong resolver
func (r *customerResolver) Customers(ctx context.Context) ([]*Customer, error) {
    reqCtx := getRequestContext(ctx)
    
    // Validate tenant access
    if err := r.requireTenantAuth(reqCtx); err != nil {
        return nil, err
    }
    
    // Query scoped to tenant
    return r.customerService.GetByTenant(reqCtx.Tenant.ID)
}
```

### 3.2 REST API Compatibility

REST endpoints được cung cấp cho:

1. **Simple CRUD operations**: Các thao tác đơn giản
2. **Third-party integrations**: Tích hợp bên thứ 3
3. **Health checks**: Monitoring và health check
4. **File uploads**: Upload file và binary data

#### Mẫu thiết kế REST API

```http
# REST endpoints với tenant scope
GET    /api/v1/customers
POST   /api/v1/customers
GET    /api/v1/customers/{id}
PUT    /api/v1/customers/{id}
DELETE /api/v1/customers/{id}

# Headers bắt buộc
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_slug}
```

### 3.3 Frontend-Backend Integration

#### GraphQL Code Generation

```typescript
// Auto-generated TypeScript types
export interface Customer {
  id: string;
  tenantId: string;
  name: string;
  email?: string;
  status: CustomerStatus;
  createdAt: string;
  updatedAt: string;
}

// Auto-generated React hooks
export const useGetCustomersQuery = (variables?: GetCustomersQueryVariables) => {
  return useQuery<GetCustomersQuery, GetCustomersQueryVariables>(
    GetCustomersDocument,
    variables
  );
};
```

#### Real-time Updates với Subscriptions

```typescript
// Frontend subscription cho live updates
const { data, loading } = useSubscription(NOTIFICATIONS_SUBSCRIPTION, {
  variables: { tenantId: currentTenant.id }
});

// GraphQL subscription
subscription NotificationUpdates($tenantId: TenantID!) {
  notifications(tenantId: $tenantId) {
    id
    type
    title
    message
    createdAt
  }
}
```

### 3.4 Authentication & Authorization

#### Cấu trúc JWT Token

```json
{
  "sub": "user_id",
  "tenant_id": "tenant_slug", 
  "email": "user@example.com",
  "roles": ["admin", "user"],
  "permissions": [
    "customers:read",
    "customers:write",
    "products:read"
  ],
  "exp": 1640995200,
  "iat": 1640908800
}
```

#### Authorization Patterns trong GraphQL

```go
// Resolver-level authorization
func (r *mutationResolver) CreateCustomer(ctx context.Context, input CreateCustomerInput) (*Customer, error) {
    reqCtx := getRequestContext(ctx)
    
    // Check specific permission
    if err := r.requirePermission(reqCtx, "customers", "write"); err != nil {
        return nil, err
    }
    
    // Business logic...
}

// Field-level authorization
func (r *customerResolver) SensitiveData(ctx context.Context, obj *Customer) (*string, error) {
    reqCtx := getRequestContext(ctx)
    
    if !reqCtx.User.HasRole("admin") {
        return nil, nil // Ẩn field cho non-admin
    }
    
    return &obj.SensitiveData, nil
}
```

## 4. Luồng Request

### 4.1 Request Flow Diagram

```
[Client] → [Traefik] → [GraphQL Gateway] → [Microservice] → [Database]
   │           │              │                  │              │
   │           │              │                  │              │
   │        Domain/         Middleware:       Business        Tenant
   │       Subdomain        - Auth             Logic          Schema
   │       Detection        - RBAC                            Access
   │                       - Rate Limit
   │                       - Tenant ID
```

### 4.2 Luồng xử lý chi tiết

1. **Client Request**: `student-zin100.myapp.com/graphql`

2. **Traefik Processing**:
   - **SSL Termination**: Xử lý HTTPS certificates
   - **Subdomain Analysis**: Parse `zin100` từ `student-zin100.myapp.com`
   - **Tenant Context Injection**: Thêm header `X-Tenant-ID: zin100`
   - **Load Balancing**: Chọn healthy instance của GraphQL Gateway
   - **Security Middlewares**: Rate limiting, CORS, security headers
   - **Route to Backend**: Forward request đến GraphQL Gateway

3. **GraphQL Gateway**:
   - Middleware authentication (JWT)
   - Middleware authorization (RBAC)
   - Middleware tenant validation
   - GraphQL query parsing và validation
   - Context injection cho resolvers
   - Route đến microservice tương ứng

4. **GraphQL Resolver Processing**:
   - Lấy tenant context: `getRequestContext(ctx)`  
   - Validate permissions cho từng field
   - Xử lý business logic với dataloader pattern
   - Truy cập database schema tương ứng
   - Trả về typed data

## 4. Database Architecture

### 4.1 Multi-tenant Strategy: Schema per Tenant

```
┌─────────────────────────────────────────────────────────────┐
│                      PostgreSQL Server                      │
├─────────────────────────────────────────────────────────────┤
│  Schema: public/system                                      │
│  ├── system_users       (System Admins)                    │
│  ├── tenants           (Tenant Registry)                   │
│  ├── plans             (Service Plans)                     │
│  ├── subscriptions     (Tenant Subscriptions)             │
│  ├── modules           (Available Modules)                 │
│  └── tenant_modules    (Enabled Modules per Tenant)       │
├─────────────────────────────────────────────────────────────┤
│  Schema: tenant_zin100                                      │
│  ├── users             (Tenant Users)                      │
│  ├── roles             (Tenant Roles)                      │
│  ├── permissions       (Tenant Permissions)                │
│  ├── customers         (End Customers)                     │
│  └── modules_config    (Module Configurations)             │
├─────────────────────────────────────────────────────────────┤
│  Schema: tenant_acme                                        │
│  ├── users, roles, permissions...                          │
│  └── (Same structure as above)                             │
└─────────────────────────────────────────────────────────────┘
```

### 4.2 Module-specific Tables

Mỗi module sẽ có bảng riêng trong schema của tenant:

**LMS Module**:
```sql
-- Schema: tenant_xyz
CREATE TABLE students (id, name, email, ...);
CREATE TABLE courses (id, title, description, ...);
CREATE TABLE enrollments (student_id, course_id, ...);
CREATE TABLE lessons (course_id, title, content, ...);
CREATE TABLE quizzes (course_id, questions, ...);
```

## 5. Security Architecture

### 5.1 Multi-layer Security

```
┌─────────────────────────────────────────────┐
│              Security Layers                │
├─────────────────────────────────────────────┤
│ 1. Network Security (SSL/TLS, Firewall)    │
├─────────────────────────────────────────────┤
│ 2. Application Security (JWT, Rate Limit)  │
├─────────────────────────────────────────────┤
│ 3. Authorization (RBAC Multi-tier)         │
├─────────────────────────────────────────────┤
│ 4. Data Security (Encryption, Audit)       │
├─────────────────────────────────────────────┤
│ 5. Tenant Isolation (Schema Separation)    │
└─────────────────────────────────────────────┘
```

### 5.2 RBAC Model

**System Level**:
- Super Admin
- System Admin  
- Support Staff

**Tenant Level**:
- Tenant Admin
- Manager
- User
- Viewer

**Customer Level**:
- Module-specific roles (Student, Teacher, Salesperson...)

## 6. Scalability & Performance

### 6.1 Horizontal Scaling

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  App Instance 1 │    │  App Instance 2 │    │  App Instance N │
│   (Container)   │    │   (Container)   │    │   (Container)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   Load Balancer │
                    │    (Traefik)    │
                    └─────────────────┘
```

### 6.2 Database Scaling

- **Read Replicas**: Cho read-heavy operations
- **Connection Pooling**: Optimize database connections  
- **Redis Caching**: Cache frequently accessed data
- **Database Sharding**: Phân chia tenant theo regions

### 6.3 Performance Optimization

- **CDN**: Static assets delivery
- **Lazy Loading**: Module loading on demand
- **Database Indexing**: Optimized queries
- **Background Jobs**: Async processing với Redis Queue
- **Event-driven Architecture**: Sử dụng Redis Pub/Sub cho real-time updates

### 6.4 Event-driven & Message Processing

**Redis Queue Use Cases**:
```go
// Background job processing
type BackgroundJob struct {
    Type    string                 `json:"type"`
    Payload map[string]interface{} `json:"payload"`
    TenantID string                `json:"tenant_id"`
}

// Examples:
// - Email notifications
// - PDF generation  
// - Data export/import
// - Webhook deliveries
// - Audit log processing
```

**Redis Pub/Sub Use Cases**:
```go
// Real-time notifications
channels := []string{
    "tenant:zin100:notifications",
    "tenant:zin100:chat",
    "system:maintenance",
}

// Examples:
// - Live chat messages
// - Real-time dashboard updates
// - System notifications
// - Module state changes
```

**Khi nào sẽ cần Kafka?**:
- Scale > 100K messages/second
- Event sourcing requirements
- Cross-service event streaming
- Audit trail compliance
- Analytics data pipeline

## 7. Module System Architecture

### 7.1 Module Registry

```go
type BaseModule interface {
    GetName() string
    GetVersion() string
    GetDependencies() []string
    Install(tenantID string) error
    Uninstall(tenantID string) error
    GetRoutes() []Route
    GetPermissions() []Permission
}
```

### 7.2 Module Lifecycle

1. **Module Registration**: Đăng ký module trong system
2. **Tenant Activation**: Kích hoạt module cho tenant
3. **Schema Migration**: Tạo tables cần thiết
4. **Route Registration**: Đăng ký API routes
5. **Permission Setup**: Thiết lập permissions

## 8. Project Structure - Monorepo Architecture

### 8.1 Cấu trúc dự án tách rời các module service

```
zplus-saas/
│
├── apps/
│   ├── backend/
│   │   ├── gateway/            # GraphQL hoặc REST gateway, xử lý auth/tenant
│   │   ├── auth/               # Xác thực + RBAC
│   │   ├── file/               # Quản lý file (upload/download)
│   │   ├── payment/            # Giao dịch & subscription
│   │   ├── crm/                # Khách hàng, liên hệ
│   │   ├── hrm/                # Nhân viên, chấm công, lương
│   │   ├── pos/                # Bán hàng: sản phẩm, hóa đơn
│   │   └── shared/             # Thư viện dùng chung
│   │
│   ├── frontend/
│   │   ├── web/                # Website chính (Next.js)
│   │   │   └── system/         # Giao diện quản lý của system admin
│   │   ├── admin/              # Giao diện quản lý (theo tenant)
│   │   └── ui/                 # UI components dùng chung
│
├── pkg/                        # SDK, lib dùng lại (go + js)
│
├── infra/
│   ├── db/                     # Migration cho system & tenants
│   ├── k8s/                    # Kubernetes manifests
│   ├── docker/                 # Dockerfiles
│   └── ci-cd/                  # GitHub Actions, ArgoCD
│
└── docs/
    └── architecture.md         # Tài liệu thiết kế hệ thống
```

### 8.2 Backend Services Architecture (Go Fiber + GORM)

**Gateway Service**:
```go
// apps/backend/gateway/
├── main.go                     # Entry point
├── handlers/                   # HTTP handlers
│   ├── graphql.go
│   └── rest.go
├── middleware/                 # Middleware layers
│   ├── auth.go
│   ├── tenant.go
│   └── ratelimit.go
└── config/                     # Configuration
    └── config.go
```

**Authentication Service**:
```go
// apps/backend/auth/
├── main.go
├── models/                     # GORM models
│   ├── user.go
│   ├── role.go
│   └── permission.go
├── services/                   # Business logic
│   ├── auth.go
│   └── rbac.go
└── handlers/                   # HTTP handlers
    └── auth.go
```

**Module Services** (CRM, HRM, POS):
```go
// apps/backend/{module}/
├── main.go
├── models/                     # GORM models specific to module
├── services/                   # Business logic
├── handlers/                   # HTTP handlers
└── migrations/                 # Database migrations
```

### 8.3 Frontend Architecture (Next.js)

**System Admin Frontend**:
```
// apps/frontend/web/system/
├── pages/                      # Next.js pages
│   ├── dashboard/
│   ├── tenants/
│   └── plans/
├── components/                 # React components
├── hooks/                      # Custom hooks
└── api/                        # API integration
```

**Tenant Admin Frontend**:
```
// apps/frontend/admin/
├── pages/                      # Next.js pages
├── components/                 # React components
├── contexts/                   # React contexts (tenant-specific)
└── api/                        # API integration
```

**Shared UI Components**:
```
// apps/frontend/ui/
├── components/                 # Reusable components
├── styles/                     # Shared styles
└── themes/                     # Multi-tenant theming
```

## 9. Deployment Architecture

### 9.1 Container Orchestration với Traefik

```yaml
# docker-compose.yml
version: '3.8'
services:
  traefik:
    image: traefik:v2.10
    container_name: zplus-traefik
    command:
      - "--api.insecure=true"
      - "--providers.docker=true"
      - "--providers.docker.exposedbydefault=false"
      - "--entrypoints.web.address=:80"
      - "--entrypoints.websecure.address=:443"
      - "--certificatesresolvers.letsencrypt.acme.tlschallenge=true"
      - "--certificatesresolvers.letsencrypt.acme.caserver=https://acme-v02.api.letsencrypt.org/directory"
      - "--certificatesresolvers.letsencrypt.acme.email=admin@zplus.com"
      - "--certificatesresolvers.letsencrypt.acme.storage=/letsencrypt/acme.json"
    ports:
      - "80:80"
      - "443:443"
      - "8080:8080"  # Traefik dashboard
    volumes:
      - "/var/run/docker.sock:/var/run/docker.sock:ro"
      - "./letsencrypt:/letsencrypt"
    networks:
      - zplus-network
  
  gateway:
    build: ./apps/backend/gateway
    container_name: zplus-gateway
    deploy:
      replicas: 3
    environment:
      - DATABASE_URL=${DATABASE_URL}
      - REDIS_URL=${REDIS_URL}
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.gateway.rule=Host(`*.zplus.com`)"
      - "traefik.http.routers.gateway.entrypoints=websecure"
      - "traefik.http.routers.gateway.tls.certresolver=letsencrypt"
      - "traefik.http.middlewares.tenant-headers.headers.customrequestheaders.X-Tenant-ID="
      - "traefik.http.routers.gateway.middlewares=tenant-headers"
    networks:
      - zplus-network
  
  postgres:
    image: postgres:15
    environment:
      - POSTGRES_DB=zplus_saas
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - zplus-network
  
  redis:
    image: redis:7-alpine
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data
    networks:
      - zplus-network

volumes:
  postgres_data:
  redis_data:

networks:
  zplus-network:
    driver: bridge
```

### 9.2 Kubernetes Deployment với Traefik

```yaml
# k8s/traefik-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: traefik
  namespace: zplus-saas
spec:
  replicas: 2
  selector:
    matchLabels:
      app: traefik
  template:
    metadata:
      labels:
        app: traefik
    spec:
      containers:
      - name: traefik
        image: traefik:v2.10
        args:
          - --api.insecure=true
          - --providers.kubernetes=true
          - --providers.kubernetes.ingressclass=traefik
          - --entrypoints.web.address=:80
          - --entrypoints.websecure.address=:443
        ports:
        - containerPort: 80
        - containerPort: 443
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: traefik
  namespace: zplus-saas
spec:
  selector:
    app: traefik
  ports:
  - name: web
    port: 80
    targetPort: 80
  - name: websecure
    port: 443
    targetPort: 443
  - name: admin
    port: 8080
    targetPort: 8080
  type: LoadBalancer
```

### 9.3 Environment Architecture

- **Development**: Single instance, local database
- **Staging**: Multi-instance, shared database  
- **Production**: Clustered, replicated database

## 11. Monitoring & Observability

### 10.1 Traefik Configuration Chi tiết

**File Configuration Structure**:
```
infra/docker/
├── traefik.yml          # Main Traefik configuration
├── dynamic.yml          # Dynamic routing rules
├── docker-compose.yml   # Container orchestration
└── letsencrypt/         # SSL certificates storage
```

**Multi-tenant Routing Logic**:
```yaml
# Dynamic routing for tenants
http:
  middlewares:
    tenant-headers:
      headers:
        customRequestHeaders:
          # tenant-slug.zplus.com → X-Tenant-ID: tenant-slug
          X-Tenant-ID: "{{ .Request.Host | regexReplaceAll `^([a-z0-9-]+)\\..*` `$1` }}"
  
  routers:
    api-tenant:
      rule: "Host(`{subdomain:[a-z0-9-]+}.zplus.com`) && PathPrefix(`/api`)"
      service: gateway-service
      middlewares:
        - tenant-headers
        - rate-limit
```

**Production Security**:
- SSL/TLS certificates via Let's Encrypt
- Rate limiting per tenant
- Security headers injection
- Access logs với tenant context
- Health checks và monitoring

### 10.2 Container Orchestration Strategy

**Development Environment**:
```bash
# Start all services
cd infra/docker
docker-compose up -d

# View Traefik dashboard
open http://localhost:8080

# Test tenant routing
curl -H "Host: tenant1.zplus.com" http://localhost/api/health
```

**Production Environment**:
- Kubernetes cluster với Traefik Ingress Controller
- Horizontal Pod Autoscaling dựa trên CPU/Memory
- Rolling deployments với zero downtime
- Centralized logging và monitoring

### 11.1 Monitoring Stack

- **Metrics**: Prometheus + Grafana
- **Logging**: ELK Stack (Elasticsearch, Logstash, Kibana)
- **Tracing**: Jaeger
- **Health Checks**: Built-in endpoints
- **Alerting**: Grafana Alerts

### 11.2 Key Metrics

- **System**: CPU, Memory, Disk usage
- **Application**: Response time, Error rate, Throughput
- **Business**: Tenant count, Active users, Revenue