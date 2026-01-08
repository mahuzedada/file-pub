package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"file-pub/image"
	"file-pub/internal/common"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	_ "github.com/go-sql-driver/mysql"
)

// Config holds application configuration
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	S3Bucket   string
	S3Region   string
	Port       string
}

func main() {
	config := loadConfig()

	app, err := initApp(config)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer app.DB.Close()

	// Setup routes
	http.HandleFunc("/", app.ImageHandler.HandleHome)
	http.HandleFunc("/upload", app.ImageHandler.HandleUpload)
	http.HandleFunc("/health", app.handleHealth)

	log.Printf("Server starting on port %s", config.Port)
	log.Printf("Database: %s@%s:%s/%s", config.DBUser, config.DBHost, config.DBPort, config.DBName)
	log.Printf("S3 Bucket: %s (Region: %s)", config.S3Bucket, config.S3Region)

	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// App holds application dependencies
type App struct {
	DB           *sql.DB
	S3Client     *s3.S3
	ImageHandler *image.ImageHandler
	Config       Config
}

func loadConfig() Config {
	config := Config{
		DBHost:     common.GetEnv("DB_HOST", "localhost"),
		DBPort:     common.GetEnv("DB_PORT", "3306"),
		DBUser:     common.GetEnv("DB_USER", "root"),
		DBPassword: common.GetEnv("DB_PASSWORD", "password"),
		DBName:     common.GetEnv("DB_NAME", "filepub"),
		S3Bucket:   common.GetEnv("S3_BUCKET", ""),
		S3Region:   common.GetEnv("S3_REGION", "us-east-1"),
		Port:       common.GetEnv("PORT", "8080"),
	}

	if config.S3Bucket == "" {
		log.Fatal("S3_BUCKET environment variable is required")
	}

	return config
}

func initApp(config Config) (*App, error) {
	// Initialize database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test database connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create table if not exists
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	// Initialize AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.S3Region),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	s3Client := s3.New(sess)
	uploader := s3manager.NewUploader(sess)

	// Parse templates
	templates, err := template.ParseGlob("templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	// Initialize domain services
	imageRepo := image.NewImageRepository(db)
	imageService := image.NewImageService(imageRepo, uploader, config.S3Bucket)
	imageHandler := image.NewImageHandler(imageService, templates)

	return &App{
		DB:           db,
		S3Client:     s3Client,
		ImageHandler: imageHandler,
		Config:       config,
	}, nil
}

func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS images (
		id VARCHAR(36) PRIMARY KEY,
		filename VARCHAR(255) NOT NULL,
		original_name VARCHAR(255) NOT NULL,
		s3_key VARCHAR(512) NOT NULL,
		s3_url VARCHAR(1024) NOT NULL,
		content_type VARCHAR(100) NOT NULL,
		size BIGINT NOT NULL,
		uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_uploaded_at (uploaded_at DESC)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	_, err := db.Exec(query)
	return err
}

func (app *App) handleHealth(w http.ResponseWriter, r *http.Request) {
	// Check database connection
	if err := app.DB.Ping(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "Database unhealthy: %v\n", err)
		return
	}

	// Check S3 access
	_, err := app.S3Client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(app.Config.S3Bucket),
	})
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "S3 bucket unhealthy: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
	fmt.Fprintf(w, "Database: Connected\n")
	fmt.Fprintf(w, "S3 Bucket: Accessible\n")
}
