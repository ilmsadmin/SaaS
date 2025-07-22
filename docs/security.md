# Security Guide - Zplus SaaS

## 1. Tổng quan Security Architecture

Zplus SaaS thực hiện bảo mật theo mô hình **Defense in Depth** với nhiều lớp bảo vệ:

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

## 2. Authentication System

### 2.1 JWT Token Implementation

**Token Structure:**
```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "sub": "user_id",
    "tenant_id": "tenant_slug",
    "email": "user@example.com",
    "roles": ["admin", "user"],
    "permissions": [
      "customers:read",
      "customers:write",
      "products:read"
    ],
    "session_id": "session_uuid",
    "exp": 1640995200,
    "iat": 1640908800
  }
}
```

**Go Implementation:**
```go
type JWTClaims struct {
    UserID      string   `json:"sub"`
    TenantID    string   `json:"tenant_id"`
    Email       string   `json:"email"`
    Roles       []string `json:"roles"`
    Permissions []string `json:"permissions"`
    SessionID   string   `json:"session_id"`
    jwt.RegisteredClaims
}

func GenerateToken(user *User, tenant *Tenant) (string, error) {
    claims := &JWTClaims{
        UserID:   user.ID,
        TenantID: tenant.ID,
        Email:    user.Email,
        Roles:    user.GetRoleNames(),
        Permissions: user.GetPermissions(),
        SessionID: uuid.New().String(),
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "zplus-saas",
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateToken(tokenString string) (*JWTClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(os.Getenv("JWT_SECRET")), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, fmt.Errorf("invalid token")
}
```

### 2.2 Refresh Token Strategy

**Refresh Token Implementation:**
```go
type RefreshToken struct {
    ID        string    `gorm:"primaryKey"`
    UserID    string    `gorm:"not null;index"`
    TenantID  string    `gorm:"not null;index"`
    Token     string    `gorm:"not null;unique"`
    ExpiresAt time.Time `gorm:"not null"`
    Used      bool      `gorm:"default:false"`
    CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (rt *RefreshToken) IsValid() bool {
    return !rt.Used && time.Now().Before(rt.ExpiresAt)
}

func CreateRefreshToken(userID, tenantID string) (*RefreshToken, error) {
    token := &RefreshToken{
        ID:        uuid.New().String(),
        UserID:    userID,
        TenantID:  tenantID,
        Token:     generateSecureToken(64),
        ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
    }
    
    if err := db.Create(token).Error; err != nil {
        return nil, err
    }
    
    return token, nil
}

func RefreshAccessToken(refreshToken string) (*AuthResponse, error) {
    var token RefreshToken
    if err := db.Where("token = ? AND used = false", refreshToken).First(&token).Error; err != nil {
        return nil, errors.New("invalid refresh token")
    }
    
    if !token.IsValid() {
        return nil, errors.New("refresh token expired")
    }
    
    // Mark as used
    token.Used = true
    db.Save(&token)
    
    // Generate new tokens
    user := getUserByID(token.UserID)
    tenant := getTenantByID(token.TenantID)
    
    newAccessToken, _ := GenerateToken(user, tenant)
    newRefreshToken, _ := CreateRefreshToken(user.ID, tenant.ID)
    
    return &AuthResponse{
        AccessToken:  newAccessToken,
        RefreshToken: newRefreshToken.Token,
        ExpiresIn:    86400, // 24 hours
    }, nil
}
```

### 2.3 Session Management

**Session Storage:**
```go
type Session struct {
    ID        string    `redis:"id"`
    UserID    string    `redis:"user_id"`
    TenantID  string    `redis:"tenant_id"`
    IPAddress string    `redis:"ip_address"`
    UserAgent string    `redis:"user_agent"`
    CreatedAt time.Time `redis:"created_at"`
    LastSeen  time.Time `redis:"last_seen"`
    Active    bool      `redis:"active"`
}

func CreateSession(userID, tenantID, ip, userAgent string) (*Session, error) {
    session := &Session{
        ID:        uuid.New().String(),
        UserID:    userID,
        TenantID:  tenantID,
        IPAddress: ip,
        UserAgent: userAgent,
        CreatedAt: time.Now(),
        LastSeen:  time.Now(),
        Active:    true,
    }
    
    // Store in Redis with TTL
    key := fmt.Sprintf("session:%s", session.ID)
    return session, redisClient.HMSet(ctx, key, session).Err()
}

func InvalidateSession(sessionID string) error {
    key := fmt.Sprintf("session:%s", sessionID)
    return redisClient.Del(ctx, key).Err()
}

func InvalidateAllUserSessions(userID string) error {
    pattern := fmt.Sprintf("session:*")
    keys, err := redisClient.Keys(ctx, pattern).Result()
    if err != nil {
        return err
    }
    
    for _, key := range keys {
        session := &Session{}
        if err := redisClient.HMGet(ctx, key, "user_id").Scan(session); err == nil {
            if session.UserID == userID {
                redisClient.Del(ctx, key)
            }
        }
    }
    
    return nil
}
```

## 3. Authorization System (RBAC)

### 3.1 Multi-tier RBAC Model

**3-tier Permission Structure:**

```
System Level:
├── Super Admin      → All system operations
├── System Admin     → Tenant management, billing
└── Support Staff    → Read-only access, support

Tenant Level:
├── Tenant Admin     → All tenant operations
├── Manager          → User management, module config
├── User             → Basic operations
└── Viewer           → Read-only access

Customer Level:
├── Student          → LMS module access
├── Teacher          → LMS content creation
├── Salesperson      → CRM lead management
└── Customer         → Portal access
```

**RBAC Implementation:**
```go
type Permission struct {
    ID       string `gorm:"primaryKey"`
    Name     string `gorm:"uniqueIndex;not null"`
    Resource string `gorm:"not null"` // users, customers, products
    Action   string `gorm:"not null"` // create, read, update, delete
    Level    string `gorm:"not null"` // system, tenant, customer
}

type Role struct {
    ID          string       `gorm:"primaryKey"`
    Name        string       `gorm:"not null"`
    DisplayName string       `gorm:"not null"`
    Level       string       `gorm:"not null"` // system, tenant, customer
    IsSystem    bool         `gorm:"default:false"`
    Permissions []Permission `gorm:"many2many:role_permissions;"`
}

type User struct {
    ID       string `gorm:"primaryKey"`
    Email    string `gorm:"uniqueIndex;not null"`
    TenantID string `gorm:"not null;index"`
    Roles    []Role `gorm:"many2many:user_roles;"`
}

func (u *User) HasPermission(resource, action string) bool {
    for _, role := range u.Roles {
        for _, permission := range role.Permissions {
            if permission.Resource == resource && permission.Action == action {
                return true
            }
        }
    }
    return false
}

func (u *User) HasRole(roleName string) bool {
    for _, role := range u.Roles {
        if role.Name == roleName {
            return true
        }
    }
    return false
}

func (u *User) GetPermissions() []string {
    var permissions []string
    seen := make(map[string]bool)
    
    for _, role := range u.Roles {
        for _, permission := range role.Permissions {
            key := fmt.Sprintf("%s:%s", permission.Resource, permission.Action)
            if !seen[key] {
                permissions = append(permissions, key)
                seen[key] = true
            }
        }
    }
    
    return permissions
}
```

### 3.2 Middleware Authorization

**Authentication Middleware:**
```go
func AuthMiddleware(requiredLevel string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        tokenString := extractToken(c)
        if tokenString == "" {
            return c.Status(401).JSON(fiber.Map{
                "error": "Authorization token required",
            })
        }
        
        claims, err := ValidateToken(tokenString)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{
                "error": "Invalid token",
            })
        }
        
        // Validate session
        session, err := GetSession(claims.SessionID)
        if err != nil || !session.Active {
            return c.Status(401).JSON(fiber.Map{
                "error": "Session expired",
            })
        }
        
        // Set context
        c.Locals("user_id", claims.UserID)
        c.Locals("tenant_id", claims.TenantID)
        c.Locals("permissions", claims.Permissions)
        
        return c.Next()
    }
}

func RequirePermission(resource, action string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        permissions := c.Locals("permissions").([]string)
        required := fmt.Sprintf("%s:%s", resource, action)
        
        for _, permission := range permissions {
            if permission == required {
                return c.Next()
            }
        }
        
        return c.Status(403).JSON(fiber.Map{
            "error": "Insufficient permissions",
        })
    }
}

func RequireRole(roleName string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userID := c.Locals("user_id").(string)
        user := getUserByID(userID)
        
        if !user.HasRole(roleName) {
            return c.Status(403).JSON(fiber.Map{
                "error": "Insufficient role",
            })
        }
        
        return c.Next()
    }
}
```

### 3.3 GraphQL Authorization

**Field-level Authorization:**
```go
func (r *userResolver) SensitiveData(ctx context.Context, obj *User) (*string, error) {
    reqCtx := getRequestContext(ctx)
    
    // Check if user can access sensitive data
    if !reqCtx.User.HasPermission("users", "read_sensitive") {
        return nil, nil // Return null for unauthorized access
    }
    
    // Admin can see all, users can only see their own
    if !reqCtx.User.HasRole("admin") && obj.ID != reqCtx.User.ID {
        return nil, nil
    }
    
    return &obj.SensitiveData, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
    reqCtx := getRequestContext(ctx)
    
    // Check permission
    if !reqCtx.User.HasPermission("users", "delete") {
        return false, errors.New("permission denied")
    }
    
    // Prevent self-deletion
    if id == reqCtx.User.ID {
        return false, errors.New("cannot delete yourself")
    }
    
    // Business logic...
    return true, nil
}
```

## 4. Tenant Isolation

### 4.1 Database-level Isolation

**Schema-per-Tenant Strategy:**
```go
func GetTenantDB(tenantID string) (*gorm.DB, error) {
    tenant, err := GetTenant(tenantID)
    if err != nil {
        return nil, err
    }
    
    schemaName := fmt.Sprintf("tenant_%s", tenant.Slug)
    
    // Create connection with schema
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s search_path=%s sslmode=require",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        schemaName,
    )
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    return db, nil
}

func TenantMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Extract tenant from subdomain or header
        tenantSlug := extractTenantSlug(c)
        
        if tenantSlug == "" {
            return c.Status(400).JSON(fiber.Map{
                "error": "Tenant not specified",
            })
        }
        
        // Validate tenant exists and is active
        tenant, err := GetTenantBySlug(tenantSlug)
        if err != nil {
            return c.Status(404).JSON(fiber.Map{
                "error": "Tenant not found",
            })
        }
        
        if tenant.Status != "active" {
            return c.Status(403).JSON(fiber.Map{
                "error": "Tenant suspended",
            })
        }
        
        // Set tenant context
        c.Locals("tenant", tenant)
        
        return c.Next()
    }
}
```

### 4.2 Row Level Security (RLS)

**PostgreSQL RLS Implementation:**
```sql
-- Enable RLS on tenant tables
ALTER TABLE customers ENABLE ROW LEVEL SECURITY;

-- Create policy for tenant isolation
CREATE POLICY tenant_isolation ON customers
    FOR ALL
    TO application_role
    USING (tenant_id = current_setting('app.tenant_id')::UUID);

-- Set tenant context in application
SET app.tenant_id = 'tenant_uuid_here';
```

**Go Implementation:**
```go
func SetTenantContext(db *gorm.DB, tenantID string) *gorm.DB {
    return db.Exec("SET app.tenant_id = ?", tenantID)
}

func (s *CustomerService) GetCustomers(tenantID string) ([]*Customer, error) {
    db := SetTenantContext(s.db, tenantID)
    
    var customers []*Customer
    err := db.Find(&customers).Error
    
    return customers, err
}
```

## 5. Data Encryption

### 5.1 Encryption at Rest

**Database Encryption:**
```sql
-- Enable transparent data encryption
ALTER SYSTEM SET ssl = on;
ALTER SYSTEM SET ssl_cert_file = '/etc/ssl/certs/server.crt';
ALTER SYSTEM SET ssl_key_file = '/etc/ssl/private/server.key';

-- Encrypt specific columns
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Store encrypted sensitive data
INSERT INTO users (email, password_hash, phone_encrypted) 
VALUES (
    'user@example.com',
    crypt('password123', gen_salt('bf')),
    pgp_sym_encrypt('+84901234567', 'encryption_key')
);

-- Decrypt data
SELECT 
    email,
    pgp_sym_decrypt(phone_encrypted, 'encryption_key') as phone
FROM users;
```

**Application-level Encryption:**
```go
import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
)

type Encryptor struct {
    gcm cipher.AEAD
}

func NewEncryptor(key []byte) (*Encryptor, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    return &Encryptor{gcm: gcm}, nil
}

func (e *Encryptor) Encrypt(plaintext string) (string, error) {
    nonce := make([]byte, e.gcm.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        return "", err
    }
    
    ciphertext := e.gcm.Seal(nonce, nonce, []byte(plaintext), nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (e *Encryptor) Decrypt(ciphertext string) (string, error) {
    data, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        return "", err
    }
    
    if len(data) < e.gcm.NonceSize() {
        return "", errors.New("ciphertext too short")
    }
    
    nonce := data[:e.gcm.NonceSize()]
    ciphertext_bytes := data[e.gcm.NonceSize():]
    
    plaintext, err := e.gcm.Open(nil, nonce, ciphertext_bytes, nil)
    if err != nil {
        return "", err
    }
    
    return string(plaintext), nil
}
```

### 5.2 Encryption in Transit

**TLS Configuration:**
```go
// TLS server configuration
func StartTLSServer() {
    cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
    if err != nil {
        log.Fatal(err)
    }
    
    tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{cert},
        MinVersion:   tls.VersionTLS12,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
            tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
        },
    }
    
    server := &http.Server{
        Addr:      ":443",
        TLSConfig: tlsConfig,
    }
    
    log.Fatal(server.ListenAndServeTLS("", ""))
}
```

## 6. Rate Limiting & DDoS Protection

### 6.1 Rate Limiting Implementation

**Redis-based Rate Limiting:**
```go
type RateLimiter struct {
    redis  *redis.Client
    window time.Duration
    limit  int
}

func NewRateLimiter(redis *redis.Client, window time.Duration, limit int) *RateLimiter {
    return &RateLimiter{
        redis:  redis,
        window: window,
        limit:  limit,
    }
}

func (rl *RateLimiter) Allow(key string) (bool, error) {
    ctx := context.Background()
    
    // Sliding window rate limiting
    now := time.Now().Unix()
    windowStart := now - int64(rl.window.Seconds())
    
    pipe := rl.redis.Pipeline()
    
    // Remove old entries
    pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))
    
    // Count current requests
    pipe.ZCard(ctx, key)
    
    // Add current request
    pipe.ZAdd(ctx, key, &redis.Z{Score: float64(now), Member: now})
    
    // Set expiration
    pipe.Expire(ctx, key, rl.window)
    
    results, err := pipe.Exec(ctx)
    if err != nil {
        return false, err
    }
    
    count := results[1].(*redis.IntCmd).Val()
    
    return count < int64(rl.limit), nil
}

func RateLimitMiddleware(limiter *RateLimiter) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Create rate limit key
        userID := c.Locals("user_id")
        tenantID := c.Locals("tenant_id")
        ip := c.IP()
        
        var key string
        if userID != nil {
            key = fmt.Sprintf("rate_limit:user:%s", userID)
        } else {
            key = fmt.Sprintf("rate_limit:ip:%s", ip)
        }
        
        if tenantID != nil {
            key = fmt.Sprintf("%s:tenant:%s", key, tenantID)
        }
        
        allowed, err := limiter.Allow(key)
        if err != nil {
            log.Printf("Rate limiter error: %v", err)
            return c.Next() // Allow on error
        }
        
        if !allowed {
            return c.Status(429).JSON(fiber.Map{
                "error": "Rate limit exceeded",
                "retry_after": 3600,
            })
        }
        
        return c.Next()
    }
}
```

### 6.2 DDoS Protection

**IP Whitelist/Blacklist:**
```go
type IPFilter struct {
    whitelist map[string]bool
    blacklist map[string]bool
    mu        sync.RWMutex
}

func (f *IPFilter) IsAllowed(ip string) bool {
    f.mu.RLock()
    defer f.mu.RUnlock()
    
    // Check blacklist first
    if f.blacklist[ip] {
        return false
    }
    
    // If whitelist exists, only allow whitelisted IPs
    if len(f.whitelist) > 0 {
        return f.whitelist[ip]
    }
    
    return true
}

func IPFilterMiddleware(filter *IPFilter) fiber.Handler {
    return func(c *fiber.Ctx) error {
        ip := c.IP()
        
        if !filter.IsAllowed(ip) {
            return c.Status(403).JSON(fiber.Map{
                "error": "Access denied",
            })
        }
        
        return c.Next()
    }
}
```

## 7. Audit Logging

### 7.1 Comprehensive Audit Trail

**Audit Log Structure:**
```go
type AuditLog struct {
    ID         string                 `gorm:"primaryKey"`
    TenantID   *string               `gorm:"index"`
    UserID     *string               `gorm:"index"`
    UserType   string                // system, tenant, customer
    Action     string                `gorm:"not null;index"`
    Resource   string                `gorm:"index"`
    ResourceID *string               `gorm:"index"`
    IPAddress  string
    UserAgent  string
    Details    map[string]interface{} `gorm:"type:jsonb"`
    CreatedAt  time.Time             `gorm:"default:CURRENT_TIMESTAMP;index"`
}

func LogAuditEvent(event AuditEvent) {
    auditLog := &AuditLog{
        ID:         uuid.New().String(),
        TenantID:   event.TenantID,
        UserID:     event.UserID,
        UserType:   event.UserType,
        Action:     event.Action,
        Resource:   event.Resource,
        ResourceID: event.ResourceID,
        IPAddress:  event.IPAddress,
        UserAgent:  event.UserAgent,
        Details:    event.Details,
        CreatedAt:  time.Now(),
    }
    
    // Async logging to avoid blocking
    go func() {
        if err := db.Create(auditLog).Error; err != nil {
            log.Printf("Failed to log audit event: %v", err)
        }
    }()
}
```

**Audit Middleware:**
```go
func AuditMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()
        
        // Capture request body for sensitive operations
        var requestBody []byte
        if c.Method() == "POST" || c.Method() == "PUT" || c.Method() == "DELETE" {
            requestBody = c.Body()
        }
        
        err := c.Next()
        
        // Log after request
        go func() {
            event := AuditEvent{
                TenantID:   getStringPtr(c.Locals("tenant_id")),
                UserID:     getStringPtr(c.Locals("user_id")),
                UserType:   getUserType(c),
                Action:     fmt.Sprintf("%s %s", c.Method(), c.Path()),
                Resource:   extractResource(c.Path()),
                ResourceID: extractResourceID(c),
                IPAddress:  c.IP(),
                UserAgent:  c.Get("User-Agent"),
                Details: map[string]interface{}{
                    "status_code":    c.Response().StatusCode(),
                    "response_time":  time.Since(start).Milliseconds(),
                    "request_body":   string(requestBody),
                    "content_length": len(c.Response().Body()),
                },
            }
            
            LogAuditEvent(event)
        }()
        
        return err
    }
}
```

## 8. Security Headers & CORS

### 8.1 Security Headers

**Security Headers Middleware:**
```go
func SecurityHeadersMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Prevent clickjacking
        c.Set("X-Frame-Options", "DENY")
        
        // Prevent MIME type sniffing
        c.Set("X-Content-Type-Options", "nosniff")
        
        // XSS protection
        c.Set("X-XSS-Protection", "1; mode=block")
        
        // Strict transport security
        c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        
        // Content security policy
        csp := "default-src 'self'; " +
               "script-src 'self' 'unsafe-inline' https://cdn.zplus.com; " +
               "style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; " +
               "font-src 'self' https://fonts.gstatic.com; " +
               "img-src 'self' data: https:; " +
               "connect-src 'self' wss: https:;"
        c.Set("Content-Security-Policy", csp)
        
        // Referrer policy
        c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
        
        // Permission policy
        c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
        
        return c.Next()
    }
}
```

### 8.2 CORS Configuration

**CORS Setup:**
```go
func SetupCORS() fiber.Handler {
    return cors.New(cors.Config{
        AllowOrigins: func() string {
            if os.Getenv("ENV") == "production" {
                return "https://app.zplus.com,https://*.zplus.com"
            }
            return "*" // Development only
        }(),
        AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
        AllowHeaders: "Origin,Content-Type,Accept,Authorization,X-Tenant-ID,X-Request-ID",
        AllowCredentials: true,
        MaxAge: 86400, // 24 hours
    })
}
```

## 9. Input Validation & Sanitization

### 9.1 Input Validation

**Validation Middleware:**
```go
import (
    "github.com/go-playground/validator/v10"
    "github.com/microcosm-cc/bluemonday"
)

var validate = validator.New()

type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=100"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8,max=128"`
    Phone    string `json:"phone" validate:"omitempty,e164"`
}

func ValidateRequest(req interface{}) error {
    if err := validate.Struct(req); err != nil {
        return fmt.Errorf("validation failed: %v", err)
    }
    return nil
}

func SanitizeHTML(input string) string {
    p := bluemonday.UGCPolicy()
    return p.Sanitize(input)
}

func ValidateAndSanitizeMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Skip for GET requests
        if c.Method() == "GET" {
            return c.Next()
        }
        
        // Parse and validate request body
        var body map[string]interface{}
        if err := c.BodyParser(&body); err != nil {
            return c.Status(400).JSON(fiber.Map{
                "error": "Invalid request body",
            })
        }
        
        // Sanitize string fields
        sanitizeMapValues(body)
        
        // Set sanitized body back
        c.Locals("validated_body", body)
        
        return c.Next()
    }
}

func sanitizeMapValues(m map[string]interface{}) {
    for k, v := range m {
        switch val := v.(type) {
        case string:
            m[k] = SanitizeHTML(val)
        case map[string]interface{}:
            sanitizeMapValues(val)
        case []interface{}:
            for i, item := range val {
                if itemMap, ok := item.(map[string]interface{}); ok {
                    sanitizeMapValues(itemMap)
                } else if itemStr, ok := item.(string); ok {
                    val[i] = SanitizeHTML(itemStr)
                }
            }
        }
    }
}
```

## 10. Security Monitoring & Alerting

### 10.1 Security Event Detection

**Anomaly Detection:**
```go
type SecurityMonitor struct {
    redis *redis.Client
}

func (sm *SecurityMonitor) DetectAnomalies(userID, ip string) {
    ctx := context.Background()
    
    // Check for multiple failed login attempts
    failedKey := fmt.Sprintf("failed_login:%s", ip)
    failed, _ := sm.redis.Get(ctx, failedKey).Int()
    
    if failed > 5 {
        sm.AlertSecurityEvent("MULTIPLE_FAILED_LOGINS", map[string]interface{}{
            "ip":       ip,
            "attempts": failed,
        })
    }
    
    // Check for unusual login times
    if sm.isUnusualLoginTime(userID) {
        sm.AlertSecurityEvent("UNUSUAL_LOGIN_TIME", map[string]interface{}{
            "user_id": userID,
            "time":    time.Now(),
        })
    }
    
    // Check for new device/location
    if sm.isNewDevice(userID, ip) {
        sm.AlertSecurityEvent("NEW_DEVICE_LOGIN", map[string]interface{}{
            "user_id": userID,
            "ip":      ip,
        })
    }
}

func (sm *SecurityMonitor) AlertSecurityEvent(eventType string, details map[string]interface{}) {
    alert := SecurityAlert{
        Type:      eventType,
        Severity:  getSeverity(eventType),
        Details:   details,
        Timestamp: time.Now(),
    }
    
    // Send to monitoring system
    go sm.sendAlert(alert)
    
    // Log for analysis
    log.Printf("Security alert: %s - %v", eventType, details)
}
```

### 10.2 Real-time Security Dashboard

**Security Metrics:**
```go
type SecurityMetrics struct {
    FailedLogins    int64     `json:"failed_logins"`
    BlockedIPs      int64     `json:"blocked_ips"`
    SuspiciousUsers int64     `json:"suspicious_users"`
    ActiveSessions  int64     `json:"active_sessions"`
    LastUpdated     time.Time `json:"last_updated"`
}

func GetSecurityMetrics() *SecurityMetrics {
    ctx := context.Background()
    
    // Get metrics from Redis
    metrics := &SecurityMetrics{}
    
    // Failed logins in last hour
    failedLogins, _ := redis.Get(ctx, "metrics:failed_logins:1h").Int64()
    metrics.FailedLogins = failedLogins
    
    // Blocked IPs
    blockedIPs, _ := redis.SCard(ctx, "blocked_ips").Result()
    metrics.BlockedIPs = blockedIPs
    
    // Active sessions
    activeSessions, _ := redis.DBSize(ctx).Result()
    metrics.ActiveSessions = activeSessions
    
    metrics.LastUpdated = time.Now()
    
    return metrics
}
```

## 11. Security Best Practices

### 11.1 Development Guidelines

**Secure Coding Practices:**

1. **Never store passwords in plain text**
2. **Always validate and sanitize input**
3. **Use parameterized queries to prevent SQL injection**
4. **Implement proper error handling without leaking sensitive information**
5. **Use HTTPS everywhere**
6. **Implement proper logging and monitoring**
7. **Regular security audits and penetration testing**

### 11.2 Deployment Security

**Production Checklist:**

- [ ] Enable HTTPS with strong TLS configuration
- [ ] Set up proper firewall rules
- [ ] Use environment variables for secrets
- [ ] Enable audit logging
- [ ] Set up monitoring and alerting
- [ ] Regular security updates
- [ ] Backup encryption
- [ ] Network segmentation
- [ ] Access control reviews
- [ ] Incident response plan

### 11.3 Regular Security Maintenance

**Monthly Tasks:**
- Review audit logs for anomalies
- Update dependencies with security patches
- Review user permissions and roles
- Check for unused accounts
- Verify backup integrity

**Quarterly Tasks:**
- Penetration testing
- Security training for development team
- Review and update security policies
- Access control audit
- Disaster recovery testing

**Annual Tasks:**
- Full security assessment
- Update incident response procedures
- Review compliance requirements
- Security architecture review
- Third-party security audit
