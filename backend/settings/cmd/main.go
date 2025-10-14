package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	settingsHandler "github.com/Leviosa-care/settings/internal/adapters/http"
	settingsPostgres "github.com/Leviosa-care/settings/internal/adapters/postgres"
	settingsS3 "github.com/Leviosa-care/settings/internal/adapters/s3"
	settingsApp "github.com/Leviosa-care/settings/internal/application"

	"github.com/Leviosa-care/core/contracts/services"
	"github.com/Leviosa-care/core/logger"
	"github.com/Leviosa-care/core/middleware/auth"
	
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/vault/api"
	"github.com/hengadev/encx"
	hashicorpkeys "github.com/hengadev/encx/providers/keys/hashicorp"
	hashicorpsecrets "github.com/hengadev/encx/providers/secrets/hashicorp"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Settings service error: %s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// Setup graceful shutdown
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Setup logger
	logLevel := getEnvOrDefault("LOG_LEVEL", "info")
	logStyle := getEnvOrDefault("LOG_STYLE", "text")
	slogHandler, err := logger.SetHandler(logLevel, logStyle)
	if err != nil {
		return fmt.Errorf("failed to setup logger: %w", err)
	}
	slog.SetDefault(slog.New(slogHandler))

	slog.Info("Starting Settings Service", "service", services.Settings)

	// Load configuration
	cfg := loadConfig()
	
	// Setup Vault client
	vaultClient, err := setupVaultClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to setup Vault client: %w", err)
	}

	// Setup per-service crypto service
	crypto, err := setupCryptoService(ctx, vaultClient)
	if err != nil {
		return fmt.Errorf("failed to setup crypto service: %w", err)
	}

	// Setup database connection
	pgPool, err := setupPostgreSQL(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to setup PostgreSQL: %w", err)
	}
	defer pgPool.Close()

	// Setup Redis connection  
	redisClient, err := setupRedis(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to setup Redis: %w", err)
	}
	defer redisClient.Close()

	// Setup S3 client
	s3Client, err := setupS3Client(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to setup S3 client: %w", err)
	}

	// Setup RabbitMQ connection
	rabbitConn, err := setupRabbitMQ(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to setup RabbitMQ: %w", err)
	}
	defer rabbitConn.Close()

	// Create repositories
	settingsRepo := settingsPostgres.New(ctx, pgPool)
	settingsMedia := settingsS3.New(ctx, s3Client, cfg.S3.BucketName)

	// Create application service
	settingsService := settingsApp.New(settingsRepo, settingsMedia, crypto, rabbitConn)

	// Create authentication middleware with Vault client
	// Note: We need session repository for user auth, but settings service might not need it
	// For now, we'll pass nil and rely on service authentication
	authmw := auth.NewSessionAuthMiddleware(nil, crypto, vaultClient)

	// Create HTTP handler
	handler := settingsHandler.New(settingsService, authmw)

	// Setup HTTP server
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	
	// Add health check endpoint
	mux.HandleFunc("GET /health", healthCheckHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: mux,
	}

	// Start server in goroutine
	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("Settings HTTP server starting", "port", cfg.Server.Port)
		serverErrors <- server.ListenAndServe()
	}()

	// Wait for shutdown signal or server error
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case <-ctx.Done():
		slog.Info("Shutdown signal received, stopping server...")
		
		// Graceful shutdown with timeout
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()
		
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown error: %w", err)
		}
		
		slog.Info("Settings service shutdown complete")
		return nil
	}
}

// Configuration structure
type Config struct {
	Server struct {
		Port int
	}
	Database struct {
		Host     string
		Port     int
		Name     string
		User     string
		Password string
	}
	Redis struct {
		Host string
		Port int
	}
	RabbitMQ struct {
		Host     string
		Port     int
		User     string
		Password string
	}
	Vault struct {
		Address string
		Token   string
	}
	S3 struct {
		Endpoint        string
		Region          string
		AccessKeyID     string
		SecretAccessKey string
		BucketName      string
	}
}

func loadConfig() *Config {
	cfg := &Config{}
	
	// Server configuration
	cfg.Server.Port = getEnvAsIntOrDefault("SERVER_PORT", 8080)
	
	// Database configuration
	cfg.Database.Host = getEnvOrDefault("POSTGRES_HOST", "localhost")
	cfg.Database.Port = getEnvAsIntOrDefault("POSTGRES_PORT", 5432)
	cfg.Database.Name = getEnvOrDefault("POSTGRES_DB", "leviosa")
	cfg.Database.User = getEnvOrDefault("POSTGRES_USER", "postgres")
	cfg.Database.Password = getEnvOrDefault("POSTGRES_PASSWORD", "postgres")
	
	// Redis configuration
	cfg.Redis.Host = getEnvOrDefault("REDIS_HOST", "localhost")
	cfg.Redis.Port = getEnvAsIntOrDefault("REDIS_PORT", 6379)
	
	// RabbitMQ configuration
	cfg.RabbitMQ.Host = getEnvOrDefault("RABBITMQ_HOST", "localhost")
	cfg.RabbitMQ.Port = getEnvAsIntOrDefault("RABBITMQ_PORT", 5672)
	cfg.RabbitMQ.User = getEnvOrDefault("RABBITMQ_USER", "guest")
	cfg.RabbitMQ.Password = getEnvOrDefault("RABBITMQ_PASSWORD", "guest")
	
	// Vault configuration
	cfg.Vault.Address = getEnvOrDefault("VAULT_ADDR", "http://localhost:8200")
	cfg.Vault.Token = getEnvOrDefault("VAULT_TOKEN", "")
	
	// S3 configuration
	cfg.S3.Endpoint = getEnvOrDefault("AWS_ENDPOINT_URL", "")
	cfg.S3.Region = getEnvOrDefault("AWS_DEFAULT_REGION", "us-east-1")
	cfg.S3.AccessKeyID = getEnvOrDefault("AWS_ACCESS_KEY_ID", "")
	cfg.S3.SecretAccessKey = getEnvOrDefault("AWS_SECRET_ACCESS_KEY", "")
	cfg.S3.BucketName = getEnvOrDefault("S3_BUCKET_NAME", "leviosa-settings")
	
	return cfg
}

// Setup functions

func setupVaultClient(cfg *Config) (*api.Client, error) {
	vaultConfig := api.DefaultConfig()
	vaultConfig.Address = cfg.Vault.Address
	
	client, err := api.NewClient(vaultConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Vault client: %w", err)
	}
	
	if cfg.Vault.Token != "" {
		client.SetToken(cfg.Vault.Token)
	}
	
	return client, nil
}

func setupCryptoService(ctx context.Context, vaultClient *api.Client) (encx.CryptoService, error) {
	// Create KMS provider (KeyManagementService) for cryptographic operations
	kms, err := hashicorpkeys.NewTransitService()
	if err != nil {
		return nil, fmt.Errorf("failed to create KMS provider: %w", err)
	}

	// Create secrets provider (SecretManagementService) for pepper storage
	secrets, err := hashicorpsecrets.NewKVStore()
	if err != nil {
		return nil, fmt.Errorf("failed to create secrets provider: %w", err)
	}

	// Use service-specific encryption key and pepper for GDPR compliance
	serviceKeyName := fmt.Sprintf("%s-encryption-key", services.Settings)
	servicePepperAlias := services.Settings

	// Create explicit Config struct with service-specific values
	cfg := encx.Config{
		KEKAlias:    serviceKeyName,  // Use key name, not full path
		PepperAlias: servicePepperAlias, // Use service name, not full Vault path
	}

	// Create crypto service with new v0.6.0 API signature
	crypto, err := encx.NewCrypto(ctx, kms, secrets, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create crypto service: %w", err)
	}

	slog.Info("Crypto service initialized", "service", services.Settings, "key", serviceKeyName)
	return crypto, nil
}

func setupPostgreSQL(ctx context.Context, cfg *Config) (*pgxpool.Pool, error) {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)
	
	pgCfg, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PostgreSQL config: %w", err)
	}
	
	// Configure pool settings
	pgCfg.MaxConns = 10
	pgCfg.MinConns = 2
	pgCfg.MaxConnLifetime = time.Hour
	pgCfg.MaxConnIdleTime = 30 * time.Minute
	
	pool, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create PostgreSQL pool: %w", err)
	}
	
	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}
	
	slog.Info("PostgreSQL connected", "host", cfg.Database.Host, "port", cfg.Database.Port)
	return pool, nil
}

func setupRedis(ctx context.Context, cfg *Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
	})
	
	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}
	
	slog.Info("Redis connected", "host", cfg.Redis.Host, "port", cfg.Redis.Port)
	return client, nil
}

func setupS3Client(ctx context.Context, cfg *Config) (*s3.Client, error) {
	// Build AWS config
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.S3.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.S3.AccessKeyID,
			cfg.S3.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	
	// Create S3 client
	var client *s3.Client
	if cfg.S3.Endpoint != "" {
		// Use custom endpoint (e.g., Localstack)
		client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.S3.Endpoint)
			o.UsePathStyle = true
		})
	} else {
		// Use AWS S3
		client = s3.NewFromConfig(awsCfg)
	}
	
	slog.Info("S3 client initialized", "region", cfg.S3.Region, "bucket", cfg.S3.BucketName)
	return client, nil
}

func setupRabbitMQ(ctx context.Context, cfg *Config) (*amqp091.Connection, error) {
	connectionString := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
	)
	
	conn, err := amqp091.Dial(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	
	slog.Info("RabbitMQ connected", "host", cfg.RabbitMQ.Host, "port", cfg.RabbitMQ.Port)
	return conn, nil
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy", "service": "settings"}`))
}

// Utility functions for environment variables

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue := parseIntOrDefault(value, defaultValue); intValue != defaultValue {
			return intValue
		}
	}
	return defaultValue
}

func parseIntOrDefault(value string, defaultValue int) int {
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	return defaultValue
}
