package app

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// Authuser
	authuserAgg "github.com/Leviosa-care/leviosa/backend/internal/authuser/application/aggregator"
	catalogSvc "github.com/Leviosa-care/leviosa/backend/internal/authuser/application/catalog"
	oauthSvc "github.com/Leviosa-care/leviosa/backend/internal/authuser/application/oauth"
	otpSvc "github.com/Leviosa-care/leviosa/backend/internal/authuser/application/otp"
	partnerSvc "github.com/Leviosa-care/leviosa/backend/internal/authuser/application/partner"
	sessionSvc "github.com/Leviosa-care/leviosa/backend/internal/authuser/application/session"
	userSvc "github.com/Leviosa-care/leviosa/backend/internal/authuser/application/user"

	authuserPorts "github.com/Leviosa-care/leviosa/backend/internal/authuser/ports"

	partnerRepo "github.com/Leviosa-care/leviosa/backend/internal/authuser/infrastructure/postgres/partner"
	userRepo "github.com/Leviosa-care/leviosa/backend/internal/authuser/infrastructure/postgres/user"
	otpRepo "github.com/Leviosa-care/leviosa/backend/internal/authuser/infrastructure/redis/otp"
	sessionRepo "github.com/Leviosa-care/leviosa/backend/internal/authuser/infrastructure/redis/session"
	stripeAdapter "github.com/Leviosa-care/leviosa/backend/internal/authuser/infrastructure/stripe"

	// Catalog
	catalogAgg "github.com/Leviosa-care/leviosa/backend/internal/catalog/application/aggregator"
	categorySvc "github.com/Leviosa-care/leviosa/backend/internal/catalog/application/category"
	couponSvc "github.com/Leviosa-care/leviosa/backend/internal/catalog/application/coupon"
	imageSvc "github.com/Leviosa-care/leviosa/backend/internal/catalog/application/image"
	priceSvc "github.com/Leviosa-care/leviosa/backend/internal/catalog/application/price"
	productSvc "github.com/Leviosa-care/leviosa/backend/internal/catalog/application/product"
	promotionCodeSvc "github.com/Leviosa-care/leviosa/backend/internal/catalog/application/promotion_code"

	catalogPorts "github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"

	categoryRepo "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/postgres/category"
	couponRepo "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/postgres/coupon"
	imageRepo "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/postgres/image"
	priceRepo "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/postgres/price"
	productRepo "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/postgres/product"
	promotionCodeRepo "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/postgres/promotion_code"
	sharedRepo "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/postgres/shared"
	imageMedia "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/s3/image"
	stripeCatalog "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/stripe"

	// Common
	"github.com/Leviosa-care/leviosa/backend/internal/common/migrations"

	// Middleware
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware/auth"

	// External
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/vault/api"
	"github.com/hengadev/encx"
	hashicorpkeys "github.com/hengadev/encx/providers/keys/hashicorp"
	hashicorpsecrets "github.com/hengadev/encx/providers/secrets/hashicorp"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/stripe/stripe-go/v82"
	stripeClient "github.com/stripe/stripe-go/v82/client"
)

// Container holds all application dependencies
type Container struct {
	Config *Config

	// Infrastructure
	DB          *pgxpool.Pool
	RedisClient *redis.Client
	RabbitMQ    *amqp.Connection
	S3Client    *s3.Client
	Crypto      encx.CryptoService
	StripeAPI   *stripeClient.API
	VaultClient *api.Client
	AuthMw      auth.AuthMiddleware

	// Authuser Repositories
	UserRepo      authuserPorts.UserRepository
	PartnerRepo   authuserPorts.PartnerRepository
	OTPRepo       authuserPorts.OTPRepository
	SessionRepo   authuserPorts.SessionRepository
	StripeAdapter authuserPorts.StripeService

	// Catalog Repositories
	CategoryRepo      catalogPorts.CategoryRepository
	ProductRepo       catalogPorts.ProductRepository
	PriceRepo         catalogPorts.PriceRepository
	ImageRepo         catalogPorts.ImageRepository
	CouponRepo        catalogPorts.CouponRepository
	PromotionCodeRepo catalogPorts.PromotionCodeRepository
	SharedRepo        catalogPorts.SharedRepository
	ImageMedia        catalogPorts.ImageMedia
	StripeCatalog     catalogPorts.StripeService

	// Authuser Services
	UserService    *userSvc.UserService
	PartnerService *partnerSvc.PartnerService
	OTPService     *otpSvc.OTPService
	SessionService *sessionSvc.SessionService
	OAuthService   *oauthSvc.Service
	CatalogService *catalogSvc.Service
	AuthAggregator *authuserAgg.AuthAggregatorService

	// Catalog Services
	CategoryService      *categorySvc.CategoryService
	ProductService       *productSvc.ProductService
	PriceService         *priceSvc.PriceService
	ImageService         *imageSvc.ImageService
	CouponService        *couponSvc.CouponService
	PromotionCodeService *promotionCodeSvc.PromotionCodeService
	CatalogAggregator    *catalogAgg.ProductAggregatorService
}

// NewContainer creates and wires all dependencies
func NewContainer(ctx context.Context, cfg *Config) (*Container, error) {
	c := &Container{Config: cfg}

	// Setup infrastructure
	if err := c.setupInfrastructure(ctx); err != nil {
		return nil, fmt.Errorf("setup infrastructure: %w", err)
	}

	// Run migrations
	if err := c.runMigrations(ctx); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	// Setup repositories
	c.setupRepositories(ctx)

	// Setup services
	if err := c.setupServices(ctx); err != nil {
		return nil, fmt.Errorf("setup services: %w", err)
	}

	return c, nil
}

func (c *Container) setupInfrastructure(ctx context.Context) error {
	// PostgreSQL
	pgCfg, err := pgxpool.ParseConfig(c.Config.PostgresURL)
	if err != nil {
		return fmt.Errorf("parse postgres config: %w", err)
	}
	pgCfg.MaxConns = 25
	pgCfg.MinConns = 5
	pgCfg.MaxConnLifetime = 30 * time.Minute
	pgCfg.MaxConnIdleTime = 5 * time.Minute

	c.DB, err = pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		return fmt.Errorf("create postgres pool: %w", err)
	}

	if err := c.DB.Ping(ctx); err != nil {
		return fmt.Errorf("ping postgres: %w", err)
	}

	// Redis
	c.RedisClient = redis.NewClient(&redis.Options{
		Addr:     c.Config.RedisAddr,
		Password: c.Config.RedisPassword,
		DB:       c.Config.RedisDB,
	})

	if err := c.RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping redis: %w", err)
	}

	// RabbitMQ (optional)
	if c.Config.RabbitMQURL != "" {
		c.RabbitMQ, err = amqp.Dial(c.Config.RabbitMQURL)
		if err != nil {
			// Don't fail - RabbitMQ is optional
			fmt.Printf("Warning: Failed to connect to RabbitMQ: %v\n", err)
		}
	}

	// S3
	awsCfg, err := awsConfig.LoadDefaultConfig(ctx,
		awsConfig.WithRegion(c.Config.S3Region),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			c.Config.S3AccessKeyID,
			c.Config.S3SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return fmt.Errorf("load aws config: %w", err)
	}

	s3Options := func(o *s3.Options) {
		if c.Config.S3Endpoint != "" {
			o.BaseEndpoint = aws.String(c.Config.S3Endpoint)
			o.UsePathStyle = true
		}
	}
	c.S3Client = s3.NewFromConfig(awsCfg, s3Options)

	// Vault & Encryption
	vaultConfig := api.DefaultConfig()
	vaultConfig.Address = c.Config.VaultAddr
	c.VaultClient, err = api.NewClient(vaultConfig)
	if err != nil {
		return fmt.Errorf("create vault client: %w", err)
	}
	c.VaultClient.SetToken(c.Config.VaultToken)

	keyProvider, err := hashicorpkeys.New(hashicorpkeys.Config{
		Address: c.Config.VaultAddr,
		Token:   c.Config.VaultToken,
	})
	if err != nil {
		return fmt.Errorf("create hashicorp key provider: %w", err)
	}

	secretProvider, err := hashicorpsecrets.New(hashicorpsecrets.Config{
		Address: c.Config.VaultAddr,
		Token:   c.Config.VaultToken,
	})
	if err != nil {
		return fmt.Errorf("create hashicorp secret provider: %w", err)
	}

	c.Crypto, err = encx.NewCryptoService(ctx, keyProvider, encx.WithSecretsProvider(secretProvider))
	if err != nil {
		return fmt.Errorf("create crypto service: %w", err)
	}

	// Stripe
	stripe.Key = c.Config.StripeSecretKey
	c.StripeAPI = &stripeClient.API{}
	c.StripeAPI.Init(c.Config.StripeSecretKey, nil)

	return nil
}

func (c *Container) runMigrations(ctx context.Context) error {
	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("pgx"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	gooseDB, err := sql.Open("pgx", c.DB.Config().ConnString())
	if err != nil {
		return fmt.Errorf("open sql.DB for migrations: %w", err)
	}
	defer gooseDB.Close()

	if err := goose.Up(gooseDB, "."); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}

func (c *Container) setupRepositories(ctx context.Context) {
	// Authuser repositories
	c.UserRepo = userRepo.New(ctx, c.DB)
	c.PartnerRepo = partnerRepo.New(ctx, c.DB)
	c.OTPRepo = otpRepo.New(c.RedisClient)
	c.SessionRepo = sessionRepo.New(c.RedisClient)
	c.StripeAdapter = stripeAdapter.NewService(c.StripeAPI)

	// Catalog repositories
	c.CategoryRepo = categoryRepo.New(ctx, c.DB)
	c.ProductRepo = productRepo.New(ctx, c.DB)
	c.PriceRepo = priceRepo.New(ctx, c.DB)
	c.ImageRepo = imageRepo.New(ctx, c.DB)
	c.CouponRepo = couponRepo.New(ctx, c.DB)
	c.PromotionCodeRepo = promotionCodeRepo.New(ctx, c.DB)
	c.SharedRepo = sharedRepo.New(ctx, c.DB)
	c.ImageMedia = imageMedia.New(ctx, c.S3Client, c.Config.S3BucketName)
	c.StripeCatalog = stripeCatalog.NewService(c.StripeAPI)
}

func (c *Container) setupServices(ctx context.Context) error {
	// Authuser services
	c.OTPService = otpSvc.New(c.OTPRepo, nil) // RabbitMQ consumer setup would go here if enabled

	c.SessionService = sessionSvc.New(c.SessionRepo)

	// Initialize AuthMiddleware (needs SessionRepo, Crypto, and VaultClient)
	c.AuthMw = auth.NewSessionAuthMiddleware(c.SessionRepo, c.Crypto, c.VaultClient)

	c.UserService = userSvc.New(
		c.UserRepo,
		c.SessionService,
		c.Crypto,
	)

	// Catalog service for authuser (cross-module dependency)
	c.CatalogService = catalogSvc.New(
		c.CategoryRepo,
		c.ProductRepo,
	)

	c.PartnerService = partnerSvc.New(
		c.PartnerRepo,
		c.UserRepo,
		c.CatalogService,
		c.Crypto,
		c.StripeAdapter,
		nil, // RabbitMQ connection (optional)
	)

	c.OAuthService = oauthSvc.New(&oauthSvc.OAuthConfig{
		GoogleClientID:     c.Config.GoogleClientID,
		GoogleClientSecret: c.Config.GoogleClientSecret,
		AppleClientID:      c.Config.AppleClientID,
		AppleClientSecret:  c.Config.AppleClientSecret,
		AppleTeamID:        c.Config.AppleTeamID,
		AppleKeyID:         c.Config.AppleKeyID,
		SessionSecret:      c.Config.SessionSecret,
	})

	c.AuthAggregator = authuserAgg.New(
		c.UserService,
		c.OTPService,
		c.SessionService,
		c.OAuthService,
		c.PartnerService,
		c.Crypto,
	)

	// Catalog services
	c.ImageService = imageSvc.New(
		c.ImageRepo,
		c.ImageMedia,
		c.SharedRepo,
	)

	c.CategoryService = categorySvc.New(
		c.CategoryRepo,
		c.ImageService,
		c.SharedRepo,
	)

	c.CouponService = couponSvc.New(
		c.CouponRepo,
		c.StripeCatalog,
	)

	c.PromotionCodeService = promotionCodeSvc.New(
		c.PromotionCodeRepo,
		c.StripeCatalog,
	)

	c.PriceService = priceSvc.New(
		c.PriceRepo,
		c.StripeCatalog,
	)

	c.ProductService = productSvc.New(
		c.ProductRepo,
		c.ImageService,
		c.StripeCatalog,
		c.SharedRepo,
	)

	c.CatalogAggregator = catalogAgg.New(
		c.CategoryService,
		c.ProductService,
		c.PriceService,
		c.CouponService,
		c.PromotionCodeService,
	)

	return nil
}

// Close gracefully closes all resources
func (c *Container) Close(ctx context.Context) error {
	var errs []error

	if c.DB != nil {
		c.DB.Close()
	}

	if c.RedisClient != nil {
		if err := c.RedisClient.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close redis: %w", err))
		}
	}

	if c.RabbitMQ != nil {
		if err := c.RabbitMQ.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close rabbitmq: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing resources: %v", errs)
	}

	return nil
}
