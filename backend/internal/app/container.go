package app

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// Authuser
	authuserAgg "github.com/Leviosa-care/leviosa/backend/internal/authuser/application/aggregator"
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
	couponPayment "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/stripe/coupon"
	pricePayment "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/stripe/price"
	productPayment "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/stripe/product"
	promotionCodePayment "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/stripe/promotion_code"

	// Booking
	allocationSvc "github.com/Leviosa-care/leviosa/backend/internal/booking/application/allocation"
	availabilitySvc "github.com/Leviosa-care/leviosa/backend/internal/booking/application/availability"
	bookingSvc "github.com/Leviosa-care/leviosa/backend/internal/booking/application/booking"
	buildingSvc "github.com/Leviosa-care/leviosa/backend/internal/booking/application/building"
	roomSvc "github.com/Leviosa-care/leviosa/backend/internal/booking/application/room"

	bookingPorts "github.com/Leviosa-care/leviosa/backend/internal/booking/ports"

	allocationRepo "github.com/Leviosa-care/leviosa/backend/internal/booking/infrastructure/postgres/allocation"
	availabilityRepo "github.com/Leviosa-care/leviosa/backend/internal/booking/infrastructure/postgres/availability"
	bookingRepo "github.com/Leviosa-care/leviosa/backend/internal/booking/infrastructure/postgres/booking"
	buildingRepo "github.com/Leviosa-care/leviosa/backend/internal/booking/infrastructure/postgres/building"
	roomRepo "github.com/Leviosa-care/leviosa/backend/internal/booking/infrastructure/postgres/room"

	// Common
	"github.com/Leviosa-care/leviosa/backend/internal/common/migrations"

	// Common auth session
	commonSession "github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"

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

	// Booking Repositories
	BuildingRepo     bookingPorts.BuildingRepository
	RoomRepo         bookingPorts.RoomRepository
	AllocationRepo   bookingPorts.RoomAllocationRepository
	AvailabilityRepo bookingPorts.AvailabilityRepository
	BookingRepo      bookingPorts.BookingRepository

	// Authuser Services
	UserService    authuserPorts.UserService
	PartnerService authuserPorts.PartnerService
	OTPService     authuserPorts.OTPService
	SessionService authuserPorts.SessionService
	AuthAggregator authuserPorts.AuthAggregatorService

	// Catalog Services
	CategoryService      catalogPorts.CategoryService
	ProductService       catalogPorts.ProductService
	PriceService         catalogPorts.PriceService
	ImageService         catalogPorts.ImageService
	CouponService        catalogPorts.CouponService
	PromotionCodeService catalogPorts.PromotionCodeService
	CatalogAggregator    catalogPorts.ProductAggregatorService
	CategoryAggregator   catalogPorts.CategoryImagesService

	// Booking Services
	BuildingService     bookingPorts.BuildingService
	RoomService         bookingPorts.RoomService
	AllocationService   bookingPorts.RoomAllocationService
	AvailabilityService bookingPorts.AvailabilityService
	BookingService      bookingPorts.BookingService
	MetricsService      bookingPorts.MetricsService
	PaymentService      bookingPorts.PaymentService
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
	if err := c.setupRepositories(ctx); err != nil {
		return nil, fmt.Errorf("setup repositories: %w", err)
	}

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

	// Vault client
	vaultConfig := api.DefaultConfig()
	vaultConfig.Address = c.Config.VaultAddr
	c.VaultClient, err = api.NewClient(vaultConfig)
	if err != nil {
		return fmt.Errorf("create vault client: %w", err)
	}
	c.VaultClient.SetToken(c.Config.VaultToken)

	// Encryption — the hashicorp providers read VAULT_ADDR and VAULT_TOKEN from env
	keyProvider, err := hashicorpkeys.NewTransitService()
	if err != nil {
		return fmt.Errorf("create hashicorp key provider: %w", err)
	}

	secretProvider, err := hashicorpsecrets.NewKVStore()
	if err != nil {
		return fmt.Errorf("create hashicorp secret provider: %w", err)
	}

	c.Crypto, err = encx.NewCrypto(ctx, keyProvider, secretProvider, encx.Config{
		KEKAlias:    c.Config.EncxKEKAlias,
		PepperAlias: c.Config.EncxPepperAlias,
	})
	if err != nil {
		return fmt.Errorf("create crypto service: %w", err)
	}

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

func (c *Container) setupRepositories(ctx context.Context) error {
	// Authuser repositories
	c.UserRepo = userRepo.New(ctx, c.DB)
	c.PartnerRepo = partnerRepo.New(ctx, c.DB)
	c.OTPRepo = otpRepo.New(c.RedisClient)
	c.SessionRepo = sessionRepo.New(c.RedisClient)
	c.StripeAdapter = stripeAdapter.NewService(c.Config.StripeSecretKey, "")

	// Catalog repositories
	c.CategoryRepo = categoryRepo.New(ctx, c.DB)
	c.ProductRepo = productRepo.New(ctx, c.DB)
	c.PriceRepo = priceRepo.New(ctx, c.DB)
	c.ImageRepo = imageRepo.New(ctx, c.DB)
	c.CouponRepo = couponRepo.New(ctx, c.DB)
	c.PromotionCodeRepo = promotionCodeRepo.New(ctx, c.DB)
	c.SharedRepo = sharedRepo.New(ctx, c.DB)
	c.ImageMedia = imageMedia.New(ctx, c.S3Client, c.Config.S3BucketName)

	// Booking repositories
	c.BuildingRepo = buildingRepo.New(ctx, c.DB)
	c.RoomRepo = roomRepo.New(ctx, c.DB)
	var err error
	c.AllocationRepo, err = allocationRepo.New(ctx, c.DB)
	if err != nil {
		return fmt.Errorf("create allocation repo: %w", err)
	}
	c.AvailabilityRepo = availabilityRepo.New(ctx, c.DB)
	c.BookingRepo = bookingRepo.New(ctx, c.DB)

	return nil
}

func (c *Container) setupServices(ctx context.Context) error {
	// Session service (needed early for auth middleware and other services)
	c.SessionService = sessionSvc.New(ctx, c.SessionRepo, c.Crypto)

	// Auth middleware (uses a separate minimal session repo interface for token lookups)
	authSessionRepo := commonSession.NewRedisSessionRepository(c.RedisClient)
	c.AuthMw = auth.NewSessionAuthMiddleware(authSessionRepo, c.Crypto, c.VaultClient)

	// OTP service
	var err error
	c.OTPService, err = otpSvc.New(ctx, c.OTPRepo, c.Crypto, c.RabbitMQ)
	if err != nil {
		return fmt.Errorf("create otp service: %w", err)
	}

	// Catalog stripe gateways
	stripeCoupon := couponPayment.NewCoupon(c.Config.StripeSecretKey, "")
	stripePrice := pricePayment.NewPrice(c.Config.StripeSecretKey, "")
	stripeProduct := productPayment.NewProduct(c.Config.StripeSecretKey, "")
	stripePromoCode := promotionCodePayment.NewPromotionCode(c.Config.StripeSecretKey, "")

	// Catalog services
	c.ImageService = imageSvc.New(c.ImageRepo, c.ImageMedia, c.SharedRepo)

	c.CategoryService = categorySvc.New(c.CategoryRepo, c.SharedRepo)

	c.CouponService = couponSvc.NewCouponService(c.CouponRepo)

	c.PromotionCodeService = promotionCodeSvc.New(c.PromotionCodeRepo, c.CouponRepo)

	c.PriceService = priceSvc.New(c.PriceRepo, c.SharedRepo, stripePrice)

	c.ProductService = productSvc.New(c.ProductRepo, c.SharedRepo, stripeProduct, stripePrice)

	// Catalog aggregators
	c.CatalogAggregator = catalogAgg.NewProductPricesAggregatorService(
		c.ProductService,
		c.PriceService,
		c.ImageService,
	)

	c.CategoryAggregator = catalogAgg.NewCategoryAggregatorService(
		c.CategoryService,
		c.ImageService,
	)

	// Authuser services
	c.UserService = userSvc.New(c.UserRepo, c.Crypto, c.StripeAdapter)

	c.PartnerService, err = partnerSvc.New(
		ctx,
		c.PartnerRepo,
		c.UserRepo,
		c.ProductService,
		c.CategoryService,
		c.Crypto,
		c.StripeAdapter,
	)
	if err != nil {
		return fmt.Errorf("create partner service: %w", err)
	}

	c.AuthAggregator = authuserAgg.New(
		c.OTPService,
		c.UserService,
		c.SessionService,
		c.PartnerService,
	)

	// Booking services (initialized but not fully wired yet)
	_ = allocationSvc.New
	_ = availabilitySvc.New
	_ = bookingSvc.New
	_ = buildingSvc.New
	_ = roomSvc.New

	// Suppress unused stripe gateway variables
	_ = stripeCoupon
	_ = stripePromoCode

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
