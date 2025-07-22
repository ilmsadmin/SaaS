package config

import (
	"os"
	"strconv"
)

type Config struct {
	// Application
	Environment string
	AppName     string
	AppVersion  string
	Port        string
	Debug       bool

	// Database
	DatabaseURL string
	RedisURL    string
	MongoURL    string

	// JWT
	JWTSecret           string
	JWTExpiresIn        string
	JWTRefreshExpiresIn string
	JWTIssuer           string
	JWTAudience         string

	// Services URLs
	AuthServiceURL    string
	CRMServiceURL     string
	HRMServiceURL     string
	POSServiceURL     string
	LMSServiceURL     string
	CheckinServiceURL string
	PaymentServiceURL string
	FileServiceURL    string

	// CORS
	CORSAllowOrigins string
	CORSAllowMethods string
	CORSAllowHeaders string

	// MinIO/S3
	MinIOEndpoint  string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOBucket    string
	MinIOUseSSL    bool
}

func Load() *Config {
	return &Config{
		// Application
		Environment: getEnv("ENVIRONMENT", "development"),
		AppName:     getEnv("APP_NAME", "zplus-saas"),
		AppVersion:  getEnv("APP_VERSION", "1.0.0"),
		Port:        getEnv("PORT", "8080"),
		Debug:       getEnvBool("DEBUG", true),

		// Database
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres123@localhost:5432/zplus_saas?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		MongoURL:    getEnv("MONGODB_URL", "mongodb://admin:admin123@localhost:27017"),

		// JWT
		JWTSecret:           getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
		JWTExpiresIn:        getEnv("JWT_EXPIRES_IN", "24h"),
		JWTRefreshExpiresIn: getEnv("JWT_REFRESH_EXPIRES_IN", "7d"),
		JWTIssuer:           getEnv("JWT_ISSUER", "zplus-saas"),
		JWTAudience:         getEnv("JWT_AUDIENCE", "zplus-users"),

		// Services URLs
		AuthServiceURL:    getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
		CRMServiceURL:     getEnv("CRM_SERVICE_URL", "http://localhost:8082"),
		HRMServiceURL:     getEnv("HRM_SERVICE_URL", "http://localhost:8083"),
		POSServiceURL:     getEnv("POS_SERVICE_URL", "http://localhost:8084"),
		LMSServiceURL:     getEnv("LMS_SERVICE_URL", "http://localhost:8085"),
		CheckinServiceURL: getEnv("CHECKIN_SERVICE_URL", "http://localhost:8086"),
		PaymentServiceURL: getEnv("PAYMENT_SERVICE_URL", "http://localhost:8087"),
		FileServiceURL:    getEnv("FILE_SERVICE_URL", "http://localhost:8088"),

		// CORS
		CORSAllowOrigins: getEnv("CORS_ALLOW_ORIGINS", "*"),
		CORSAllowMethods: getEnv("CORS_ALLOW_METHODS", "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS"),
		CORSAllowHeaders: getEnv("CORS_ALLOW_HEADERS", "Origin, Content-Type, Accept, Authorization"),

		// MinIO/S3
		MinIOEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinIOAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinIOSecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin123"),
		MinIOBucket:    getEnv("MINIO_BUCKET", "zplus-files"),
		MinIOUseSSL:    getEnvBool("MINIO_USE_SSL", false),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
