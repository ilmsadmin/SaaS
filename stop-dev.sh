#!/bin/bash

# Zplus SaaS Development Stop Script
echo "🛑 Stopping Zplus SaaS Development Environment..."

# Stop backend services
for service in api-gateway auth-service tenant-service frontend; do
    if [ -f /tmp/${service}.pid ]; then
        PID=$(cat /tmp/${service}.pid)
        if ps -p $PID > /dev/null; then
            echo "🔄 Stopping ${service} (PID: $PID)..."
            kill $PID
            rm /tmp/${service}.pid
        else
            echo "⚠️  ${service} was not running"
            rm -f /tmp/${service}.pid
        fi
    else
        echo "⚠️  No PID file found for ${service}"
    fi
done

# Stop Docker services
echo "📦 Stopping infrastructure services..."
cd /Users/toan/Documents/project/SaaS
docker-compose down

# Clean up log files
echo "🧹 Cleaning up log files..."
rm -f /tmp/api-gateway.log /tmp/auth-service.log /tmp/tenant-service.log /tmp/frontend.log

echo "✅ All services stopped!"
