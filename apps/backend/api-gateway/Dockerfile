# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files from root
COPY go.mod go.sum ./
RUN go mod download

# Copy all necessary source code
COPY apps/backend/api-gateway ./apps/backend/api-gateway
COPY apps/backend/shared ./apps/backend/shared

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./apps/backend/api-gateway/cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS and wget for healthcheck
RUN apk --no-cache add ca-certificates tzdata wget

# Set timezone
ENV TZ=Asia/Ho_Chi_Minh

# Create non-root user
RUN adduser -D -s /bin/sh api-user

# Set working directory
WORKDIR /home/api-user/

# Copy binary from builder stage
COPY --from=builder /app/main ./

# Change ownership
RUN chown api-user:api-user main

# Switch to non-root user
USER api-user

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run command
CMD ["./main"]
