package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// FileRecord represents a file record in database
type FileRecord struct {
	ID               string     `json:"id" db:"id"`
	Filename         string     `json:"filename" db:"filename"`
	OriginalFilename string     `json:"original_filename" db:"original_filename"`
	ContentType      string     `json:"content_type" db:"content_type"`
	FileSize         int64      `json:"file_size" db:"file_size"`
	FilePath         string     `json:"file_path" db:"file_path"`
	UploadUserID     string     `json:"upload_user_id" db:"upload_user_id"`
	TenantID         string     `json:"tenant_id" db:"tenant_id"`
	Metadata         string     `json:"metadata,omitempty" db:"metadata"`
	IsPublic         bool       `json:"is_public" db:"is_public"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// FileUploadResponse represents file upload response
type FileUploadResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Files   []FileRecord `json:"files,omitempty"`
}

// FileHandler handles file-related requests
type FileHandler struct {
	db        *sqlx.DB
	uploadDir string
}

func NewFileHandler(db *sqlx.DB) *FileHandler {
	uploadDir := os.Getenv("UPLOAD_DIRECTORY")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("Failed to create upload directory: %v", err)
	}

	return &FileHandler{
		db:        db,
		uploadDir: uploadDir,
	}
}

// generateUniqueFilename generates a unique filename
func (h *FileHandler) generateUniqueFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	name := strings.TrimSuffix(originalName, ext)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s_%d%s", name, timestamp, ext)
}

// getMimeType returns the MIME type of the file
func (h *FileHandler) getMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".txt":
		return "text/plain"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	default:
		return "application/octet-stream"
	}
}

// UploadFiles handles file upload
func (h *FileHandler) UploadFiles(c *fiber.Ctx) error {
	tenantID := c.Get("X-Tenant-ID", "demo-tenant")
	userID := c.Get("X-User-ID", "demo-user")

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(FileUploadResponse{
			Success: false,
			Message: "Failed to parse multipart form",
		})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(FileUploadResponse{
			Success: false,
			Message: "No files provided",
		})
	}

	var uploadedFiles []FileRecord
	for _, file := range files {
		// Generate unique filename
		uniqueFileName := h.generateUniqueFilename(file.Filename)
		filePath := filepath.Join(h.uploadDir, uniqueFileName)

		// Save file to disk
		if err := c.SaveFile(file, filePath); err != nil {
			log.Printf("Failed to save file %s: %v", file.Filename, err)
			continue
		}

		// Create file record
		fileRecord := FileRecord{
			Filename:         uniqueFileName,
			OriginalFilename: file.Filename,
			ContentType:      h.getMimeType(file.Filename),
			FileSize:         file.Size,
			FilePath:         filePath,
			UploadUserID:     userID,
			TenantID:         tenantID,
			Metadata:         "{}",
			IsPublic:         false,
		}

		// Insert into database
		query := `INSERT INTO files (filename, original_filename, content_type, file_size, file_path, upload_user_id, tenant_id, metadata, is_public)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id, created_at, updated_at`

		err = h.db.QueryRow(query, fileRecord.Filename, fileRecord.OriginalFilename,
			fileRecord.ContentType, fileRecord.FileSize, fileRecord.FilePath, fileRecord.UploadUserID,
			fileRecord.TenantID, fileRecord.Metadata, fileRecord.IsPublic).Scan(&fileRecord.ID, &fileRecord.CreatedAt, &fileRecord.UpdatedAt)

		if err != nil {
			log.Printf("Failed to insert file record: %v", err)
			os.Remove(filePath)
			continue
		}

		uploadedFiles = append(uploadedFiles, fileRecord)
	}

	if len(uploadedFiles) == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(FileUploadResponse{
			Success: false,
			Message: "Failed to upload any files",
		})
	}

	return c.JSON(FileUploadResponse{
		Success: true,
		Message: fmt.Sprintf("Successfully uploaded %d file(s)", len(uploadedFiles)),
		Files:   uploadedFiles,
	})
}

// GetFiles handles listing files
func (h *FileHandler) GetFiles(c *fiber.Ctx) error {
	tenantID := c.Get("X-Tenant-ID", "demo-tenant")

	query := `SELECT id, filename, original_filename, content_type, file_size, file_path, upload_user_id, tenant_id, metadata, is_public, created_at, updated_at 
			  FROM files WHERE tenant_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`

	var files []FileRecord
	err := h.db.Select(&files, query, tenantID)
	if err != nil {
		log.Printf("Failed to fetch files: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch files",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    files,
	})
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Database connection
	dbHost := os.Getenv("DATABASE_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DATABASE_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DATABASE_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres123"
	}
	dbName := os.Getenv("DATABASE_NAME")
	if dbName == "" {
		dbName = "zplus_saas"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Connected to database successfully")

	// Initialize handler
	fileHandler := NewFileHandler(db)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024, // 100MB limit
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service":   "file-service",
			"status":    "healthy",
			"timestamp": time.Now(),
		})
	})

	// File routes
	api := app.Group("/api/v1")
	api.Post("/files/upload", fileHandler.UploadFiles)
	api.Get("/files", fileHandler.GetFiles)

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8088"
	}

	log.Printf("File Service starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
