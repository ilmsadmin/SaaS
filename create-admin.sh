#!/bin/bash

# Tạo admin user mặc định cho Zplus SaaS
echo "👤 Tạo admin user mặc định..."

# Thông tin admin mặc định
ADMIN_EMAIL="admin@zplus.com"
ADMIN_PASSWORD="Admin@123456"
ADMIN_NAME="System Administrator"

# API endpoint
API_URL="http://localhost:8080/api/v1/admin/auth/create"

# Tạo admin user
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$ADMIN_EMAIL\",
    \"password\": \"$ADMIN_PASSWORD\",
    \"first_name\": \"System\",
    \"last_name\": \"Administrator\",
    \"role\": \"super_admin\"
  }"

echo ""
echo "✅ Admin user đã được tạo!"
echo ""
echo "📝 Thông tin đăng nhập admin:"
echo "   Email:    $ADMIN_EMAIL"
echo "   Password: $ADMIN_PASSWORD"
echo "   URL:      http://localhost:3001/login"
echo ""
echo "⚠️  Vui lòng đổi mật khẩu sau khi đăng nhập lần đầu!"
