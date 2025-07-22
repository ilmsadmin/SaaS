# Thiết kế Tổng quan Dự án - Zplus SaaS

## 1. Giới thiệu

Zplus SaaS là một nền tảng SaaS (Software as a Service) đa tenant được thiết kế với kiến trúc 3 tầng phân quyền rõ ràng: System → Tenant → Customer. Hệ thống được xây dựng với khả năng mở rộng cao, có thể tích hợp nhiều module khác nhau như CRM, POS, LMS, HRM, Checkin...

## 2. Mục tiêu dự án

### 2.1 Mục tiêu chính
- Xây dựng nền tảng SaaS đa tenant hiện đại
- Hỗ trợ quản lý nhiều khách hàng (tenant) với dữ liệu tách biệt
- Cung cấp hệ thống module linh hoạt, có thể bật/tắt theo nhu cầu
- Đảm bảo bảo mật cao với hệ thống RBAC đa tầng

### 2.2 Mục tiêu kỹ thuật
- Kiến trúc microservices dễ mở rộng
- Hiệu suất cao với Redis cache
- API GraphQL linh hoạt
- Hỗ trợ custom domain/subdomain cho từng tenant

## 3. Đối tượng sử dụng

### 3.1 System Admin
- Quản trị viên hệ thống
- Quản lý tenant, gói dịch vụ, thanh toán
- Cấu hình module cho từng tenant

### 3.2 Tenant Admin
- Quản trị viên của từng tổ chức/công ty
- Quản lý người dùng, phân quyền trong phạm vi tenant
- Cấu hình tính năng, tích hợp dịch vụ

### 3.3 End User (Customer)
- Người dùng cuối sử dụng các module cụ thể
- Học viên (LMS), khách hàng (CRM), nhân viên (HRM)...

## 4. Phạm vi dự án

### 4.1 Chức năng cốt lõi
- **System Layer**: Quản lý tenant, gói cước, module, thanh toán
- **Tenant Layer**: RBAC, quản lý người dùng, khách hàng, cấu hình
- **Customer Layer**: Giao diện sử dụng các module chuyên biệt

### 4.2 Module được hỗ trợ
- **CRM**: Quản lý khách hàng, bán hàng
- **LMS**: Học tập trực tuyến
- **POS**: Bán hàng tại điểm
- **HRM**: Quản lý nhân sự
- **Checkin**: Chấm công, điểm danh

## 5. Lợi ích

### 5.1 Cho nhà cung cấp dịch vụ
- Tiết kiệm chi phí vận hành
- Dễ dàng mở rộng khách hàng
- Quản lý tập trung, hiệu quả

### 5.2 Cho khách hàng (tenant)
- Chi phí thấp, thanh toán theo sử dụng
- Không cần đầu tư hạ tầng IT
- Cập nhật tính năng tự động

### 5.3 Cho người dùng cuối
- Truy cập mọi lúc, mọi nơi
- Giao diện thân thiện, dễ sử dụng
- Tích hợp đa nền tảng

## 6. Yêu cầu kỹ thuật tổng quan

### 6.1 Performance
- Hỗ trợ hàng nghìn tenant đồng thời
- Thời gian phản hồi < 200ms
- Uptime 99.9%

### 6.2 Security
- Mã hóa dữ liệu end-to-end
- JWT authentication
- RBAC đa tầng
- Audit log đầy đủ

### 6.3 Scalability
- Horizontal scaling
- Database sharding theo tenant
- CDN cho static assets
- Load balancing

