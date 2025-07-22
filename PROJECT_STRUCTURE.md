# Zplus SaaS - Project Structure

```
zplus-saas/
├── README.md                           # Main project documentation
├── PROJECT_STRUCTURE.md               # This file - project structure overview
├── docker-compose.yml                 # Local development setup
├── docker-compose.prod.yml            # Production docker compose
├── Makefile                           # Build and development scripts
├── .gitignore                         # Git ignore patterns
├── .env.example                       # Environment variables template
├── go.work                            # Go workspace configuration
├── package.json                       # Root package.json for workspace
├── turbo.json                         # Turborepo configuration
│
├── apps/                              # Main applications
│   ├── backend/                       # Go microservices
│   │   ├── api-gateway/               # GraphQL/REST API Gateway
│   │   │   ├── cmd/                   # Application entry points
│   │   │   │   └── main.go            # Gateway main application
│   │   │   ├── internal/              # Private application code
│   │   │   │   ├── handlers/          # HTTP/GraphQL handlers
│   │   │   │   │   ├── graphql.go     # GraphQL resolvers
│   │   │   │   │   ├── rest.go        # REST endpoints
│   │   │   │   │   └── websocket.go   # WebSocket handlers
│   │   │   │   ├── middleware/        # HTTP middleware
│   │   │   │   │   ├── auth.go        # Authentication middleware
│   │   │   │   │   ├── cors.go        # CORS middleware
│   │   │   │   │   ├── ratelimit.go   # Rate limiting
│   │   │   │   │   └── tenant.go      # Tenant resolution
│   │   │   │   ├── graphql/           # GraphQL specific code
│   │   │   │   │   ├── schema.graphql # GraphQL schema
│   │   │   │   │   ├── resolvers.go   # Generated resolvers
│   │   │   │   │   └── directives.go  # Custom directives
│   │   │   │   └── proxy/             # Service proxy logic
│   │   │   │       ├── client.go      # HTTP client for microservices
│   │   │   │       └── router.go      # Request routing logic
│   │   │   ├── configs/               # Configuration files
│   │   │   │   ├── config.go          # Configuration struct
│   │   │   │   └── config.yaml        # Default configuration
│   │   │   ├── docs/                  # API documentation
│   │   │   │   └── swagger.yaml       # OpenAPI/Swagger spec
│   │   │   ├── Dockerfile             # Container build file
│   │   │   ├── go.mod                 # Go module dependencies
│   │   │   └── go.sum                 # Go module checksums
│   │   │
│   │   ├── auth-service/              # Authentication & Authorization Service
│   │   │   ├── cmd/
│   │   │   │   └── main.go
│   │   │   ├── internal/
│   │   │   │   ├── models/            # Database models
│   │   │   │   │   ├── user.go        # User entity
│   │   │   │   │   ├── tenant.go      # Tenant entity
│   │   │   │   │   ├── role.go        # Role entity
│   │   │   │   │   ├── permission.go  # Permission entity
│   │   │   │   │   └── session.go     # Session entity
│   │   │   │   ├── services/          # Business logic
│   │   │   │   │   ├── auth_service.go        # Authentication logic
│   │   │   │   │   ├── user_service.go        # User management
│   │   │   │   │   ├── tenant_service.go      # Tenant management
│   │   │   │   │   ├── rbac_service.go        # RBAC logic
│   │   │   │   │   └── jwt_service.go         # JWT token handling
│   │   │   │   ├── handlers/          # HTTP handlers
│   │   │   │   │   ├── auth_handler.go        # Auth endpoints
│   │   │   │   │   ├── user_handler.go        # User CRUD
│   │   │   │   │   └── tenant_handler.go      # Tenant CRUD
│   │   │   │   ├── repositories/      # Data access layer
│   │   │   │   │   ├── user_repo.go           # User repository
│   │   │   │   │   ├── tenant_repo.go         # Tenant repository
│   │   │   │   │   └── session_repo.go        # Session repository
│   │   │   │   └── middleware/        # Service-specific middleware
│   │   │   │       ├── jwt.go                 # JWT validation
│   │   │   │       └── rbac.go                # RBAC enforcement
│   │   │   ├── migrations/            # Database migrations
│   │   │   │   ├── 001_init_users.sql
│   │   │   │   ├── 002_init_tenants.sql
│   │   │   │   └── 003_init_rbac.sql
│   │   │   ├── tests/                 # Unit & integration tests
│   │   │   │   ├── auth_test.go
│   │   │   │   ├── user_test.go
│   │   │   │   └── integration_test.go
│   │   │   ├── Dockerfile
│   │   │   ├── go.mod
│   │   │   └── go.sum
│   │   │
│   │   ├── file-service/              # File Storage Service (MinIO/S3)
│   │   │   ├── cmd/
│   │   │   ├── internal/
│   │   │   │   ├── models/
│   │   │   │   ├── services/
│   │   │   │   ├── handlers/
│   │   │   │   └── repositories/
│   │   │   ├── migrations/
│   │   │   ├── tests/
│   │   │   ├── Dockerfile
│   │   │   ├── go.mod
│   │   │   └── go.sum
│   │   │
│   │   ├── payment-service/           # Payment & Billing Service
│   │   │   ├── cmd/
│   │   │   ├── internal/
│   │   │   │   ├── models/
│   │   │   │   │   ├── subscription.go
│   │   │   │   │   ├── payment.go
│   │   │   │   │   ├── invoice.go
│   │   │   │   │   └── plan.go
│   │   │   │   ├── services/
│   │   │   │   │   ├── payment_service.go
│   │   │   │   │   ├── subscription_service.go
│   │   │   │   │   ├── billing_service.go
│   │   │   │   │   └── webhook_service.go
│   │   │   │   ├── handlers/
│   │   │   │   └── repositories/
│   │   │   ├── migrations/
│   │   │   ├── tests/
│   │   │   ├── Dockerfile
│   │   │   ├── go.mod
│   │   │   └── go.sum
│   │   │
│   │   ├── crm-service/               # Customer Relationship Management
│   │   │   ├── cmd/
│   │   │   ├── internal/
│   │   │   │   ├── models/
│   │   │   │   │   ├── customer.go
│   │   │   │   │   ├── lead.go
│   │   │   │   │   ├── opportunity.go
│   │   │   │   │   ├── contact.go
│   │   │   │   │   └── activity.go
│   │   │   │   ├── services/
│   │   │   │   │   ├── customer_service.go
│   │   │   │   │   ├── lead_service.go
│   │   │   │   │   ├── sales_service.go
│   │   │   │   │   └── analytics_service.go
│   │   │   │   ├── handlers/
│   │   │   │   └── repositories/
│   │   │   ├── migrations/
│   │   │   ├── tests/
│   │   │   ├── Dockerfile
│   │   │   ├── go.mod
│   │   │   └── go.sum
│   │   │
│   │   ├── lms-service/               # Learning Management System
│   │   │   ├── cmd/
│   │   │   ├── internal/
│   │   │   │   ├── models/
│   │   │   │   │   ├── course.go
│   │   │   │   │   ├── lesson.go
│   │   │   │   │   ├── enrollment.go
│   │   │   │   │   ├── quiz.go
│   │   │   │   │   └── progress.go
│   │   │   │   ├── services/
│   │   │   │   │   ├── course_service.go
│   │   │   │   │   ├── enrollment_service.go
│   │   │   │   │   ├── progress_service.go
│   │   │   │   │   └── video_service.go
│   │   │   │   ├── handlers/
│   │   │   │   └── repositories/
│   │   │   ├── migrations/
│   │   │   ├── tests/
│   │   │   ├── Dockerfile
│   │   │   ├── go.mod
│   │   │   └── go.sum
│   │   │
│   │   ├── pos-service/               # Point of Sale System
│   │   │   ├── cmd/
│   │   │   ├── internal/
│   │   │   │   ├── models/
│   │   │   │   │   ├── product.go
│   │   │   │   │   ├── inventory.go
│   │   │   │   │   ├── order.go
│   │   │   │   │   ├── transaction.go
│   │   │   │   │   └── receipt.go
│   │   │   │   ├── services/
│   │   │   │   │   ├── product_service.go
│   │   │   │   │   ├── inventory_service.go
│   │   │   │   │   ├── order_service.go
│   │   │   │   │   └── payment_service.go
│   │   │   │   ├── handlers/
│   │   │   │   └── repositories/
│   │   │   ├── migrations/
│   │   │   ├── tests/
│   │   │   ├── Dockerfile
│   │   │   ├── go.mod
│   │   │   └── go.sum
│   │   │
│   │   ├── hrm-service/               # Human Resource Management
│   │   │   ├── cmd/
│   │   │   ├── internal/
│   │   │   │   ├── models/
│   │   │   │   │   ├── employee.go
│   │   │   │   │   ├── department.go
│   │   │   │   │   ├── payroll.go
│   │   │   │   │   ├── attendance.go
│   │   │   │   │   └── performance.go
│   │   │   │   ├── services/
│   │   │   │   │   ├── employee_service.go
│   │   │   │   │   ├── payroll_service.go
│   │   │   │   │   ├── attendance_service.go
│   │   │   │   │   └── performance_service.go
│   │   │   │   ├── handlers/
│   │   │   │   └── repositories/
│   │   │   ├── migrations/
│   │   │   ├── tests/
│   │   │   ├── Dockerfile
│   │   │   ├── go.mod
│   │   │   └── go.sum
│   │   │
│   │   ├── checkin-service/           # Check-in & Attendance System
│   │   │   ├── cmd/
│   │   │   ├── internal/
│   │   │   │   ├── models/
│   │   │   │   │   ├── checkin.go
│   │   │   │   │   ├── location.go
│   │   │   │   │   ├── qrcode.go
│   │   │   │   │   └── face_recognition.go
│   │   │   │   ├── services/
│   │   │   │   │   ├── checkin_service.go
│   │   │   │   │   ├── location_service.go
│   │   │   │   │   ├── qr_service.go
│   │   │   │   │   └── face_service.go
│   │   │   │   ├── handlers/
│   │   │   │   └── repositories/
│   │   │   ├── migrations/
│   │   │   ├── tests/
│   │   │   ├── Dockerfile
│   │   │   ├── go.mod
│   │   │   └── go.sum
│   │   │
│   │   └── shared/                    # Shared utilities across services
│   │       ├── database/              # Database utilities
│   │       │   ├── postgres.go        # PostgreSQL connection
│   │       │   ├── mongodb.go         # MongoDB connection
│   │       │   ├── redis.go           # Redis connection
│   │       │   └── migrations.go      # Migration utilities
│   │       ├── utils/                 # Common utilities
│   │       │   ├── logger.go          # Structured logging
│   │       │   ├── validator.go       # Input validation
│   │       │   ├── encryption.go      # Encryption utilities
│   │       │   ├── email.go           # Email service
│   │       │   └── queue.go           # Background job queue
│   │       ├── middleware/            # Shared middleware
│   │       │   ├── auth.go            # Authentication middleware
│   │       │   ├── tenant.go          # Tenant isolation
│   │       │   ├── logging.go         # Request logging
│   │       │   └── metrics.go         # Metrics collection
│   │       ├── types/                 # Shared types
│   │       │   ├── tenant.go          # Tenant types
│   │       │   ├── user.go            # User types
│   │       │   └── response.go        # API response types
│   │       └── config/                # Shared configuration
│   │           ├── database.go        # Database config
│   │           ├── redis.go           # Redis config
│   │           └── services.go        # Service discovery config
│   │
│   └── frontend/                      # Frontend applications
│       ├── web/                       # Next.js Web Application
│       │   ├── src/
│       │   │   ├── app/               # Next.js 14 App Router
│       │   │   │   ├── (system)/      # System admin routes
│       │   │   │   │   ├── dashboard/
│       │   │   │   │   ├── tenants/
│       │   │   │   │   ├── users/
│       │   │   │   │   └── settings/
│       │   │   │   ├── (tenant)/      # Tenant-specific routes
│       │   │   │   │   ├── [tenant]/  # Dynamic tenant routing
│       │   │   │   │   │   ├── dashboard/
│       │   │   │   │   │   ├── crm/
│       │   │   │   │   │   ├── lms/
│       │   │   │   │   │   ├── pos/
│       │   │   │   │   │   ├── hrm/
│       │   │   │   │   │   └── checkin/
│       │   │   │   ├── auth/          # Authentication pages
│       │   │   │   │   ├── login/
│       │   │   │   │   ├── register/
│       │   │   │   │   └── forgot-password/
│       │   │   │   ├── api/           # API routes
│       │   │   │   │   ├── auth/
│       │   │   │   │   ├── graphql/
│       │   │   │   │   └── webhooks/
│       │   │   │   ├── globals.css    # Global styles
│       │   │   │   ├── layout.tsx     # Root layout
│       │   │   │   ├── loading.tsx    # Loading UI
│       │   │   │   ├── error.tsx      # Error UI
│       │   │   │   └── not-found.tsx  # 404 page
│       │   │   ├── components/        # Reusable components
│       │   │   │   ├── ui/            # shadcn/ui components
│       │   │   │   │   ├── button.tsx
│       │   │   │   │   ├── input.tsx
│       │   │   │   │   ├── modal.tsx
│       │   │   │   │   └── ...
│       │   │   │   ├── layout/        # Layout components
│       │   │   │   │   ├── header.tsx
│       │   │   │   │   ├── sidebar.tsx
│       │   │   │   │   ├── navigation.tsx
│       │   │   │   │   └── footer.tsx
│       │   │   │   ├── modules/       # Module-specific components
│       │   │   │   │   ├── crm/
│       │   │   │   │   │   ├── customer-list.tsx
│       │   │   │   │   │   ├── lead-form.tsx
│       │   │   │   │   │   └── sales-dashboard.tsx
│       │   │   │   │   ├── lms/
│       │   │   │   │   │   ├── course-list.tsx
│       │   │   │   │   │   ├── video-player.tsx
│       │   │   │   │   │   └── quiz-component.tsx
│       │   │   │   │   └── ...
│       │   │   │   └── common/        # Common components
│       │   │   │       ├── data-table.tsx
│       │   │   │       ├── form-builder.tsx
│       │   │   │       ├── chart-widgets.tsx
│       │   │   │       └── tenant-switcher.tsx
│       │   │   ├── lib/               # Utility libraries
│       │   │   │   ├── utils.ts       # Utility functions
│       │   │   │   ├── auth.ts        # Authentication logic
│       │   │   │   ├── graphql.ts     # GraphQL client
│       │   │   │   ├── tenant.ts      # Tenant resolution
│       │   │   │   └── validations.ts # Form validations
│       │   │   ├── hooks/             # Custom React hooks
│       │   │   │   ├── useAuth.ts     # Authentication hook
│       │   │   │   ├── useTenant.ts   # Tenant context hook
│       │   │   │   ├── useLocalStorage.ts
│       │   │   │   └── useDebounce.ts
│       │   │   ├── store/             # Zustand stores
│       │   │   │   ├── auth-store.ts  # Authentication state
│       │   │   │   ├── tenant-store.ts # Tenant state
│       │   │   │   ├── ui-store.ts    # UI state
│       │   │   │   └── module-stores/ # Module-specific stores
│       │   │   │       ├── crm-store.ts
│       │   │   │       ├── lms-store.ts
│       │   │   │       └── ...
│       │   │   └── types/             # TypeScript types
│       │   │       ├── auth.ts        # Authentication types
│       │   │       ├── tenant.ts      # Tenant types
│       │   │       ├── user.ts        # User types
│       │   │       └── modules/       # Module-specific types
│       │   │           ├── crm.ts
│       │   │           ├── lms.ts
│       │   │           └── ...
│       │   ├── public/                # Static assets
│       │   │   ├── images/
│       │   │   ├── icons/
│       │   │   ├── favicon.ico
│       │   │   └── manifest.json
│       │   ├── .env.example           # Environment variables
│       │   ├── .env.local             # Local environment
│       │   ├── next.config.js         # Next.js configuration
│       │   ├── tailwind.config.js     # Tailwind CSS config
│       │   ├── tsconfig.json          # TypeScript config
│       │   ├── package.json           # Dependencies
│       │   ├── Dockerfile             # Container build
│       │   └── README.md              # Frontend documentation
│       │
│       ├── mobile/                    # React Native Mobile App
│       │   ├── src/
│       │   │   ├── screens/           # App screens
│       │   │   │   ├── auth/
│       │   │   │   ├── dashboard/
│       │   │   │   ├── modules/
│       │   │   │   └── settings/
│       │   │   ├── components/        # Reusable components
│       │   │   ├── navigation/        # Navigation setup
│       │   │   ├── store/             # State management
│       │   │   ├── services/          # API services
│       │   │   ├── utils/             # Utilities
│       │   │   └── types/             # TypeScript types
│       │   ├── assets/                # Static assets
│       │   ├── app.json               # Expo configuration
│       │   ├── package.json
│       │   ├── tsconfig.json
│       │   └── README.md
│       │
│       └── admin/                     # System Admin Dashboard
│           ├── src/
│           │   ├── pages/             # Admin pages
│           │   │   ├── system/        # System management
│           │   │   ├── tenants/       # Tenant management
│           │   │   ├── billing/       # Billing management
│           │   │   └── monitoring/    # System monitoring
│           │   ├── components/        # Admin components
│           │   ├── lib/               # Admin utilities
│           │   ├── hooks/             # Admin hooks
│           │   ├── store/             # Admin state
│           │   └── types/             # Admin types
│           ├── public/
│           ├── package.json
│           ├── tsconfig.json
│           └── README.md
│
├── packages/                          # Shared packages
│   ├── ui/                           # Shared UI components
│   │   ├── src/
│   │   │   ├── components/           # Common UI components
│   │   │   │   ├── button/
│   │   │   │   ├── input/
│   │   │   │   ├── modal/
│   │   │   │   ├── table/
│   │   │   │   └── chart/
│   │   │   ├── hooks/                # UI-related hooks
│   │   │   ├── utils/                # UI utilities
│   │   │   └── index.ts              # Package exports
│   │   ├── package.json
│   │   ├── tsconfig.json
│   │   └── README.md
│   │
│   ├── types/                        # Shared TypeScript types
│   │   ├── src/
│   │   │   ├── auth.ts               # Authentication types
│   │   │   ├── tenant.ts             # Tenant types
│   │   │   ├── user.ts               # User types
│   │   │   ├── api.ts                # API response types
│   │   │   ├── modules/              # Module types
│   │   │   │   ├── crm.ts
│   │   │   │   ├── lms.ts
│   │   │   │   ├── pos.ts
│   │   │   │   ├── hrm.ts
│   │   │   │   └── checkin.ts
│   │   │   └── index.ts
│   │   ├── package.json
│   │   └── tsconfig.json
│   │
│   ├── utils/                        # Shared utilities
│   │   ├── src/
│   │   │   ├── date.ts               # Date utilities
│   │   │   ├── format.ts             # Formatting utilities
│   │   │   ├── validation.ts         # Validation utilities
│   │   │   ├── encryption.ts         # Encryption utilities
│   │   │   ├── api.ts                # API utilities
│   │   │   └── index.ts
│   │   ├── package.json
│   │   └── tsconfig.json
│   │
│   └── config/                       # Shared configuration
│       ├── src/
│       │   ├── eslint/               # ESLint configurations
│       │   │   ├── base.js
│       │   │   ├── react.js
│       │   │   └── node.js
│       │   ├── prettier/             # Prettier configurations
│       │   │   └── index.js
│       │   ├── tailwind/             # Tailwind configurations
│       │   │   └── base.js
│       │   └── typescript/           # TypeScript configurations
│       │       ├── base.json
│       │       ├── react.json
│       │       └── node.json
│       ├── package.json
│       └── README.md
│
├── infra/                            # Infrastructure as Code
│   ├── docker/                       # Docker configurations
│   │   ├── api-gateway/
│   │   │   └── Dockerfile
│   │   ├── auth-service/
│   │   │   └── Dockerfile
│   │   ├── nginx/
│   │   │   ├── Dockerfile
│   │   │   └── nginx.conf
│   │   ├── postgres/
│   │   │   ├── Dockerfile
│   │   │   └── init.sql
│   │   ├── redis/
│   │   │   ├── Dockerfile
│   │   │   └── redis.conf
│   │   └── docker-compose.override.yml
│   │
│   ├── k8s/                          # Kubernetes manifests
│   │   ├── namespaces/
│   │   │   ├── production.yaml
│   │   │   └── staging.yaml
│   │   ├── deployments/
│   │   │   ├── api-gateway.yaml
│   │   │   ├── auth-service.yaml
│   │   │   ├── crm-service.yaml
│   │   │   ├── lms-service.yaml
│   │   │   ├── pos-service.yaml
│   │   │   ├── hrm-service.yaml
│   │   │   ├── checkin-service.yaml
│   │   │   ├── file-service.yaml
│   │   │   └── payment-service.yaml
│   │   ├── services/
│   │   │   ├── api-gateway-svc.yaml
│   │   │   ├── auth-service-svc.yaml
│   │   │   └── ...
│   │   ├── configmaps/
│   │   │   ├── api-gateway-config.yaml
│   │   │   ├── auth-service-config.yaml
│   │   │   └── ...
│   │   ├── secrets/
│   │   │   ├── database-secrets.yaml
│   │   │   ├── jwt-secrets.yaml
│   │   │   └── ...
│   │   ├── ingress/
│   │   │   ├── api-gateway-ingress.yaml
│   │   │   └── tenant-ingress.yaml
│   │   └── storage/
│   │       ├── postgres-pvc.yaml
│   │       ├── redis-pvc.yaml
│   │       └── minio-pvc.yaml
│   │
│   ├── terraform/                    # Infrastructure provisioning
│   │   ├── modules/
│   │   │   ├── vpc/
│   │   │   │   ├── main.tf
│   │   │   │   ├── variables.tf
│   │   │   │   └── outputs.tf
│   │   │   ├── eks/
│   │   │   │   ├── main.tf
│   │   │   │   ├── variables.tf
│   │   │   │   └── outputs.tf
│   │   │   ├── rds/
│   │   │   │   ├── main.tf
│   │   │   │   ├── variables.tf
│   │   │   │   └── outputs.tf
│   │   │   └── redis/
│   │   │       ├── main.tf
│   │   │       ├── variables.tf
│   │   │       └── outputs.tf
│   │   ├── environments/
│   │   │   ├── production/
│   │   │   │   ├── main.tf
│   │   │   │   ├── variables.tf
│   │   │   │   └── terraform.tfvars
│   │   │   └── staging/
│   │   │       ├── main.tf
│   │   │       ├── variables.tf
│   │   │       └── terraform.tfvars
│   │   ├── provider.tf
│   │   ├── versions.tf
│   │   └── README.md
│   │
│   └── monitoring/                   # Monitoring stack
│       ├── prometheus/
│       │   ├── prometheus.yml
│       │   ├── alert.rules.yml
│       │   └── docker-compose.yml
│       ├── grafana/
│       │   ├── dashboards/
│       │   │   ├── api-gateway.json
│       │   │   ├── microservices.json
│       │   │   ├── business-metrics.json
│       │   │   └── infrastructure.json
│       │   ├── provisioning/
│       │   │   ├── datasources/
│       │   │   └── dashboards/
│       │   └── grafana.ini
│       ├── jaeger/
│       │   ├── jaeger.yml
│       │   └── docker-compose.yml
│       └── elk/
│           ├── elasticsearch/
│           │   └── elasticsearch.yml
│           ├── logstash/
│           │   ├── logstash.conf
│           │   └── pipelines.yml
│           ├── kibana/
│           │   └── kibana.yml
│           └── docker-compose.yml
│
├── scripts/                          # Development and deployment scripts
│   ├── build.sh                      # Build all services
│   ├── deploy.sh                     # Deployment script
│   ├── setup-dev.sh                  # Development environment setup
│   ├── test.sh                       # Run all tests
│   ├── migrate.sh                    # Database migrations
│   ├── seed.sh                       # Database seeding
│   ├── backup.sh                     # Database backup
│   ├── restore.sh                    # Database restore
│   └── monitoring/
│       ├── setup-prometheus.sh
│       ├── setup-grafana.sh
│       └── setup-elk.sh
│
├── tests/                            # Integration and E2E tests
│   ├── integration/                  # Integration tests
│   │   ├── api/
│   │   │   ├── auth_test.go
│   │   │   ├── crm_test.go
│   │   │   ├── lms_test.go
│   │   │   └── ...
│   │   ├── database/
│   │   │   ├── postgres_test.go
│   │   │   ├── mongodb_test.go
│   │   │   └── redis_test.go
│   │   └── fixtures/
│   │       ├── users.json
│   │       ├── tenants.json
│   │       └── ...
│   ├── e2e/                          # End-to-end tests
│   │   ├── playwright/
│   │   │   ├── auth.spec.ts
│   │   │   ├── crm.spec.ts
│   │   │   ├── lms.spec.ts
│   │   │   └── tenant-switching.spec.ts
│   │   ├── cypress/
│   │   │   ├── integration/
│   │   │   ├── fixtures/
│   │   │   └── support/
│   │   └── config/
│   │       ├── playwright.config.ts
│   │       └── cypress.config.ts
│   ├── load/                         # Load testing
│   │   ├── k6/
│   │   │   ├── api-load-test.js
│   │   │   ├── user-journey.js
│   │   │   └── stress-test.js
│   │   └── artillery/
│   │       ├── load-test.yml
│   │       └── scenarios/
│   └── security/                     # Security testing
│       ├── zap/
│       │   └── security-test.py
│       └── bandit/
│           └── security-config.yml
│
├── tools/                            # Development tools
│   ├── codegen/                      # Code generation tools
│   │   ├── graphql-codegen.yml       # GraphQL code generation
│   │   ├── gorm-gen/                 # GORM model generation
│   │   └── openapi-gen/              # OpenAPI client generation
│   ├── migration/                    # Database migration tools
│   │   ├── migrate-tool.go
│   │   └── seed-data/
│   ├── monitoring/                   # Monitoring tools
│   │   ├── health-check.go
│   │   └── metrics-collector.go
│   └── deployment/                   # Deployment tools
│       ├── release.sh
│       ├── rollback.sh
│       └── health-check.sh
│
└── docs/                             # Documentation (existing)
    ├── api-documentation.md
    ├── database-schema.md
    ├── deployment.md
    ├── installation.md
    ├── module-development.md
    ├── security.md
    ├── thiet-ke-kien-truc-database.md
    ├── thiet-ke-kien-truc-du-an.md
    └── thiet-ke-tong-quan-du-an.md
```

## Key Design Principles

### 1. **Microservices Architecture**
- Each service is independent with its own database
- Service communication through API Gateway
- Shared utilities in `/shared` package

### 2. **Multi-tenancy**
- Database schema isolation per tenant
- Tenant-aware routing at API Gateway level
- Tenant context propagated through all services

### 3. **Modular Frontend**
- Module-specific components and stores
- Dynamic routing based on tenant configuration
- Shared UI components in packages

### 4. **Infrastructure as Code**
- Complete Kubernetes manifests
- Terraform for cloud infrastructure
- Docker containers for all services

### 5. **Monitoring & Observability**
- Prometheus metrics collection
- Grafana dashboards
- Distributed tracing with Jaeger
- Centralized logging with ELK stack

### 6. **Testing Strategy**
- Unit tests per service
- Integration tests for API endpoints
- E2E tests with Playwright/Cypress
- Load testing with k6/Artillery
- Security testing with OWASP ZAP

### 7. **Development Experience**
- Monorepo with Turborepo
- Shared packages for common functionality
- Docker Compose for local development
- Automated testing and deployment pipelines

This structure supports the complete SaaS platform with all modules, proper separation of concerns, scalability, and maintainability.
