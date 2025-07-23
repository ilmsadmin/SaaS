# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY apps/backend/hrm-service ./apps/backend/hrm-service
COPY apps/backend/shared ./apps/backend/shared

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./apps/backend/hrm-service/cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS and wget for healthcheck
RUN apk --no-cache add ca-certificates tzdata wget

# Set timezone
ENV TZ=Asia/Ho_Chi_Minh

# Create non-root user
RUN adduser -D -s /bin/sh hrm-user

# Set working directory
WORKDIR /home/hrm-user/

# Copy binary from builder stage
COPY --from=builder /app/main ./

# Change ownership
RUN chown hrm-user:hrm-user main

# Switch to non-root user
USER hrm-user

# Expose port
EXPOSE 8089

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8089/health || exit 1

# Run the application
CMD ["./main"]
