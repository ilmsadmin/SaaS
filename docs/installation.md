# Hướng dẫn Cài đặt - Zplus SaaS

## 1. Yêu cầu hệ thống

### Development Environment

| Component | Version | Purpose |
|-----------|---------|---------|
| **Go** | 1.21+ | Backend services |
| **Node.js** | 18+ | Frontend applications |
| **PostgreSQL** | 13+ | Primary database |
| **Redis** | 6+ | Cache & message queue |
| **Docker** | 20+ | Containerization (optional) |
| **Git** | 2.30+ | Version control |

### Production Environment

| Component | Minimum | Recommended |
|-----------|---------|-------------|
| **CPU** | 4 cores | 8+ cores |
| **RAM** | 8GB | 16GB+ |
| **Storage** | 100GB SSD | 500GB+ SSD |
| **Network** | 100Mbps | 1Gbps+ |

## 2. Cài đặt Dependencies

### 2.1 Go Installation

**macOS:**
```bash
brew install go
```

**Ubuntu/Debian:**
```bash
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

**Windows:**
Download từ https://golang.org/dl/

### 2.2 Node.js Installation

**Using Node Version Manager (Recommended):**
```bash
# Install nvm
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash

# Install Node.js
nvm install 18
nvm use 18
nvm alias default 18
```

### 2.3 PostgreSQL Installation

**macOS:**
```bash
brew install postgresql@15
brew services start postgresql@15
```

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

**Docker (Recommended for Development):**
```bash
docker run --name zplus-postgres \
  -e POSTGRES_DB=zplus_system \
  -e POSTGRES_USER=zplus \
  -e POSTGRES_PASSWORD=your_password \
  -p 5432:5432 \
  -d postgres:15
```

### 2.4 Redis Installation

**macOS:**
```bash
brew install redis
brew services start redis
```

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install redis-server
sudo systemctl start redis-server
sudo systemctl enable redis-server
```

**Docker:**
```bash
docker run --name zplus-redis \
  -p 6379:6379 \
  -d redis:7-alpine
```

## 3. Clone và Setup Project

### 3.1 Clone Repository

```bash
git clone https://github.com/your-org/zplus-saas.git
cd zplus-saas
```

### 3.2 Environment Configuration

**Backend Configuration:**
```bash
# Copy environment template
cp apps/backend/.env.example apps/backend/.env

# Edit environment variables
nano apps/backend/.env
```

**Sample .env file:**
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=zplus
DB_PASSWORD=your_password
DB_NAME=zplus_system

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRES_IN=24h

# Server
PORT=8080
ENV=development

# File Storage
STORAGE_TYPE=local
STORAGE_PATH=./uploads

# Email (Optional)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password
```

**Frontend Configuration:**
```bash
# Copy environment template
cp apps/frontend/web/.env.example apps/frontend/web/.env.local

# Edit environment variables
nano apps/frontend/web/.env.local
```

**Sample .env.local file:**
```env
# API Endpoints
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_GRAPHQL_URL=http://localhost:8080/graphql
NEXT_PUBLIC_WS_URL=ws://localhost:8080/graphql

# Application
NEXT_PUBLIC_APP_NAME=Zplus SaaS
NEXT_PUBLIC_APP_URL=http://localhost:3000

# Feature Flags
NEXT_PUBLIC_ENABLE_ANALYTICS=false
NEXT_PUBLIC_ENABLE_CHAT=true
```

## 4. Database Setup

### 4.1 Tạo Database

```bash
# Connect to PostgreSQL
psql -U postgres

# Create system database
CREATE DATABASE zplus_system;
CREATE USER zplus WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE zplus_system TO zplus;

# Exit psql
\q
```

### 4.2 Run Migrations

```bash
cd apps/backend

# Install Go dependencies
go mod download

# Run system migrations
make migrate

# Create sample tenant
make seed
```

### 4.3 Sample Migration Commands

```bash
# Backend Makefile commands
make migrate         # Run all migrations
make migrate-up      # Apply pending migrations
make migrate-down    # Rollback last migration
make migrate-reset   # Reset all migrations
make seed           # Seed sample data
```

## 5. Backend Setup

### 5.1 Install Dependencies

```bash
cd apps/backend

# Download Go modules
go mod download
go mod tidy

# Install development tools
go install github.com/cosmtrek/air@latest
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### 5.2 Build và Run

**Development Mode (với hot reload):**
```bash
# Start with Air (hot reload)
make dev

# Or manual
air -c .air.toml
```

**Production Mode:**
```bash
# Build binary
make build

# Run binary
./bin/zplus-backend
```

### 5.3 Verify Backend

```bash
# Health check
curl http://localhost:8080/health

# GraphQL playground
open http://localhost:8080/playground
```

## 6. Frontend Setup

### 6.1 Install Dependencies

```bash
cd apps/frontend/web

# Install npm dependencies
npm install

# Install UI component dependencies
cd ../ui
npm install
```

### 6.2 Development Server

```bash
cd apps/frontend/web

# Start development server
npm run dev

# Start with specific port
npm run dev -- -p 3001
```

### 6.3 Build Production

```bash
# Build for production
npm run build

# Start production server
npm start
```

### 6.4 Verify Frontend

```bash
# Open browser
open http://localhost:3000

# System admin panel
open http://localhost:3000/system

# Test tenant subdomain (requires hosts file setup)
open http://demo.localhost:3000
```

## 7. Docker Setup (Alternative)

### 7.1 Docker Compose Development

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### 7.2 Individual Service Commands

```bash
# Start only database services
docker-compose up -d postgres redis

# Start backend
docker-compose up -d backend

# Start frontend
docker-compose up -d frontend
```

### 7.3 Docker Health Checks

```bash
# Check container status
docker-compose ps

# Check service health
docker-compose exec backend /bin/sh -c "curl -f http://localhost:8080/health"
```

## 8. Development Tools Setup

### 8.1 IDE Configuration

**VS Code Extensions:**
```json
{
  "recommendations": [
    "golang.go",
    "bradlc.vscode-tailwindcss",
    "ms-vscode.vscode-typescript-next",
    "graphql.vscode-graphql",
    "ms-vscode.vscode-json"
  ]
}
```

**VS Code Settings:**
```json
{
  "go.toolsManagement.autoUpdate": true,
  "go.useLanguageServer": true,
  "go.formatTool": "goimports",
  "go.lintTool": "golangci-lint"
}
```

### 8.2 Git Hooks Setup

```bash
# Install pre-commit hooks
npm install -g @commitlint/cli @commitlint/config-conventional
npm install -g husky

# Setup husky
npx husky install
npx husky add .husky/pre-commit "make lint"
npx husky add .husky/commit-msg "npx commitlint --edit $1"
```

### 8.3 Code Quality Tools

```bash
# Install Go linting tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install frontend linting
cd apps/frontend/web
npm install -D eslint prettier @typescript-eslint/parser

# Run linting
make lint        # Backend
npm run lint     # Frontend
```

## 9. Troubleshooting

### 9.1 Common Issues

**Database Connection Issues:**
```bash
# Check PostgreSQL status
pg_isready -h localhost -p 5432

# Check database exists
psql -U zplus -d zplus_system -c "\l"

# Reset connections
sudo systemctl restart postgresql
```

**Redis Connection Issues:**
```bash
# Check Redis status
redis-cli ping

# Check Redis info
redis-cli info

# Restart Redis
sudo systemctl restart redis-server
```

**Port Conflicts:**
```bash
# Check what's using port 8080
lsof -i :8080

# Kill process using port
kill -9 $(lsof -t -i:8080)
```

### 9.2 Environment Issues

**Go Environment:**
```bash
# Check Go version
go version

# Check Go environment
go env

# Clear module cache
go clean -modcache
```

**Node Environment:**
```bash
# Check Node version
node --version
npm --version

# Clear npm cache
npm cache clean --force

# Reinstall node_modules
rm -rf node_modules package-lock.json
npm install
```

### 9.3 Performance Issues

**Memory Issues:**
```bash
# Check system memory
free -h

# Check Docker memory usage
docker stats

# Increase Node.js memory limit
export NODE_OPTIONS="--max-old-space-size=4096"
```

## 10. Next Steps

Sau khi cài đặt thành công:

1. **Đọc tài liệu API**: [API Documentation](api-documentation.md)
2. **Tìm hiểu Database Schema**: [Database Schema](database-schema.md)
3. **Học cách phát triển Module**: [Module Development](module-development.md)
4. **Setup Testing**: [Testing Guide](testing.md)
5. **Chuẩn bị Deployment**: [Deployment Guide](deployment.md)

## 11. Support

Nếu gặp vấn đề trong quá trình cài đặt:

- Kiểm tra [Troubleshooting Guide](troubleshooting.md)
- Tạo issue trên GitHub
- Liên hệ team development
