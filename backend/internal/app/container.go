package app

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// Authuser
	authuserAgg "github.com/Leviosa-care/leviosa/backend/internal/authuser/application/aggregator"
	bookingClientAdapter "github.com/Leviosa-care/leviosa/backend/internal/authuser/infrastructure/booking"
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
	pricePayment "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/stripe/price"
	productPayment "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/stripe/product"

	// Messaging
	messagingSvc         "github.com/Leviosa-care/leviosa/backend/internal/messaging/application"
	messagingBooking     "github.com/Leviosa-care/leviosa/backend/internal/messaging/infrastructure/booking"
	messagingAuthuser    "github.com/Leviosa-care/leviosa/backend/internal/messaging/infrastructure/authuser"
	messagingRepo        "github.com/Leviosa-care/leviosa/backend/internal/messaging/infrastructure/postgres"
	messagingSSE         "github.com/Leviosa-care/leviosa/backend/internal/messaging/infrastructure/sse"
	messagingPorts       "github.com/Leviosa-care/leviosa/backend/internal/messaging/ports"

	// Booking
	allocationSvc "github.com/Leviosa-care/leviosa/backend/internal/booking/application/allocation"
	availabilitySvc "github.com/Leviosa-care/leviosa/backend/internal/booking/application/availability"
	bookingSvc "github.com/Leviosa-care/leviosa/backend/internal/booking/application/booking"
	buildingSvc "github.com/Leviosa-care/leviosa/backend/internal/booking/application/building"
	metricsSvc "github.com/Leviosa-care/leviosa/backend/internal/booking/application/metrics"
	roomSvc "github.com/Leviosa-care/leviosa/backend/internal/booking/application/room"

	bookingPorts "github.com/Leviosa-care/leviosa/backend/internal/booking/ports"

	allocationRepo "github.com/Leviosa-care/leviosa/backend/internal/booking/infrastructure/postgres/allocation"
	availabilityRepo "github.com/Leviosa-care/leviosa/backend/internal/booking/infrastructure/postgres/availability"
	bookingRepo "github.com/Leviosa-care/leviosa/backend/internal/booking/infrastructure/postgres/booking"
	buildingRepo "github.com/Leviosa-care/leviosa/backend/internal/booking/infrastructure/postgres/building"
	metricsRepo "github.com/Leviosa-care/leviosa/backend/internal/booking/infrastructure/postgres/metrics"
	roomRepo "github.com/Leviosa-care/leviosa/backend/internal/booking/infrastructure/postgres/room"
	roomScheduleRepo "github.com/Leviosa-care/leviosa/backend/internal/booking/adapters/postgres"
	bookingAuthuser "github.com/Leviosa-care/leviosa/backend/internal/booking/infrastructure/authuser"
	bookingNotification "github.com/Leviosa-care/leviosa/backend/internal/booking/infrastructure/notification"
	bookingStripe "github.com/Leviosa-care/leviosa/backend/internal/booking/infrastructure/stripe"

	// Notification
	notificationApp "github.com/Leviosa-care/leviosa/backend/internal/notification/application"
	smtpClient "github.com/Leviosa-care/leviosa/backend/internal/notification/infrastructure/smtp"
	twilioClient "github.com/Leviosa-care/leviosa/backend/internal/notification/infrastructure/twilio"
	notificationPorts "github.com/Leviosa-care/leviosa/backend/internal/notification/ports"

	// Settings
	settingsApp "github.com/Leviosa-care/leviosa/backend/internal/settings/application"
	settingsMedia "github.com/Leviosa-care/leviosa/backend/internal/settings/infrastructure/s3"
	noopPublisher "github.com/Leviosa-care/leviosa/backend/internal/settings/infrastructure/noop"
	settingsRepo "github.com/Leviosa-care/leviosa/backend/internal/settings/infrastructure/postgres"
	settingsRabbitMQ "github.com/Leviosa-care/leviosa/backend/internal/settings/infrastructure/rabbitmq"
	settingsPorts "github.com/Leviosa-care/leviosa/backend/internal/settings/ports"

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
	BuildingRepo       bookingPorts.BuildingRepository
	RoomRepo           bookingPorts.RoomRepository
	AllocationRepo     bookingPorts.RoomAllocationRepository
	AvailabilityRepo   bookingPorts.AvailabilityRepository
	BookingRepo        bookingPorts.BookingRepository
	MetricsRepo        bookingPorts.MetricsRepository
	RoomScheduleRepo   bookingPorts.RoomScheduleRepository
	BookingAuthuserCLi bookingPorts.AuthUserClient

	// Authuser Services
	UserService     authuserPorts.UserService
	PartnerService  authuserPorts.PartnerService
	OTPService      authuserPorts.OTPService
	SessionService  authuserPorts.SessionService
	AuthAggregator  authuserPorts.AuthAggregatorService
	BookingClient   authuserPorts.BookingClient

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

	// Messaging Repositories
	MessageRepo messagingPorts.MessageRepository

	// Messaging Infrastructure
	MessagingBroker *messagingSSE.Broker

	// Messaging Services
	MessagingService messagingPorts.MessagingService

	// Notification
	BookingNotificationService bookingPorts.BookingNotificationService
	NotificationService       notificationPorts.NotificationService

	// Settings
	SettingsService       settingsPorts.SettingsService
	SettingsRepo          settingsPorts.SettingsRepository
	SettingsMedia         settingsPorts.SettingsMedia
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
	c.MetricsRepo = metricsRepo.New(c.DB)
	c.RoomScheduleRepo = roomScheduleRepo.NewRoomScheduleRepository(c.DB)

	// Messaging repositories
	c.MessageRepo = messagingRepo.New(ctx, c.DB)

	// Settings repositories
	c.SettingsRepo = settingsRepo.New(ctx, c.DB)
	c.SettingsMedia = settingsMedia.New(ctx, c.S3Client, c.Config.S3BucketName)

	return nil
}

func (c *Container) setupServices(ctx context.Context) error {
	// Session service (needed early for auth middleware and other services)
	c.SessionService = sessionSvc.New(ctx, c.SessionRepo, c.Crypto)

	// Auth middleware (uses a separate minimal session repo interface for token lookups)
	authSessionRepo := commonSession.NewRedisSessionRepository(c.RedisClient)
	c.AuthMw = auth.NewSessionAuthMiddleware(authSessionRepo, c.Crypto, c.VaultClient,
		auth.WithServiceKeyCacheTTL(time.Duration(c.Config.ServiceKeyCacheTTLSeconds)*time.Second),
	)

	// Shared notification clients — built once, shared between booking notification
	// adapter and the canonical notification service (which is also used by OTP delivery).
	sharedEmailClient := smtpClient.NewSMTPClient(smtpClient.SMTPConfig{
		Host:     c.Config.SMTPHost,
		Port:     c.Config.SMTPPort,
		Username: c.Config.SMTPUsername,
		Password: c.Config.SMTPPassword,
	})

	var sharedSMSClient notificationPorts.SMSService
	if c.Config.TwilioAccountSID != "" && c.Config.TwilioAuthToken != "" && c.Config.TwilioPhoneNumber != "" {
		sharedSMSClient = twilioClient.NewTwilioClient(
			c.Config.TwilioAccountSID,
			c.Config.TwilioAuthToken,
			c.Config.TwilioPhoneNumber,
		)
	}

	// Catalog stripe gateways
	stripePrice := pricePayment.NewPrice(c.Config.StripeSecretKey, "")
	stripeProduct := productPayment.NewProduct(c.Config.StripeSecretKey, "")

	// Catalog services
	c.ImageService = imageSvc.New(c.ImageRepo, c.ImageMedia, c.SharedRepo)
	c.CategoryService = categorySvc.New(c.CategoryRepo, c.SharedRepo)
	c.CouponService = couponSvc.NewCouponService(c.CouponRepo)
	c.PromotionCodeService = promotionCodeSvc.New(c.PromotionCodeRepo, c.CouponRepo)
	c.PriceService = priceSvc.New(c.PriceRepo, c.SharedRepo, stripePrice)
	c.ProductService = productSvc.New(c.ProductRepo, c.SharedRepo, stripeProduct, stripePrice, c.ImageRepo)
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

	var err error

	// Partner service depends on catalog (ProductService, CategoryService)
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

	// Settings service — wired before OTP so the notification service can read settings.
	var settingsPublisher settingsPorts.EventPublisher
	if c.RabbitMQ != nil {
		settingsPublisher = settingsRabbitMQ.NewPublisher(c.RabbitMQ)
	} else {
		settingsPublisher = noopPublisher.NewPublisher()
	}
	c.SettingsService = settingsApp.New(c.SettingsRepo, c.SettingsMedia, c.Crypto, settingsPublisher)

	// Canonical notification service
	notificationSettingsProvider := notificationApp.NewSettingsProvider(c.SettingsService)
	notificationMailService := notificationApp.NewMailService(sharedEmailClient, notificationSettingsProvider)

	var notificationSmsService *notificationApp.SMSService
	if sharedSMSClient != nil {
		notificationSmsService = notificationApp.NewSMSService(sharedSMSClient)
	}

	c.NotificationService = notificationApp.NewNotificationService(
		notificationMailService,
		notificationSmsService,
		notificationSettingsProvider,
	)

	// OTP service — requires notification service for in-process email delivery
	c.OTPService, err = otpSvc.New(ctx, c.OTPRepo, c.Crypto, c.NotificationService)
	if err != nil {
		return fmt.Errorf("create otp service: %w", err)
	}

	// Auth aggregator — held as concrete type so SetBookingClient can be called
	// without a type assertion after BookingService is wired.
	authAgg := authuserAgg.New(
		c.OTPService,
		c.UserService,
		c.SessionService,
		c.PartnerService,
		nil, // BookingClient injected below
	)
	c.AuthAggregator = authAgg

	// Booking services
	// Stripe payment gateway for bookings
	bookingStripe := bookingStripe.NewService(c.Config.StripeSecretKey, "", c.Config.StripeWebhookSecret)

	// AuthUser client for booking module
	c.BookingAuthuserCLi = bookingAuthuser.NewInProcessClient(c.PartnerService, c.UserService)

	c.MetricsService = metricsSvc.New(c.MetricsRepo, c.Crypto)

	c.BuildingService = buildingSvc.New(c.BuildingRepo, c.Crypto)

	c.RoomService = roomSvc.New(c.RoomRepo, c.BuildingRepo, c.Crypto)

	c.AllocationService = allocationSvc.New(c.AllocationRepo, c.RoomRepo, c.BookingAuthuserCLi, c.Crypto)

	c.AvailabilityService = availabilitySvc.New(c.AvailabilityRepo, c.AllocationRepo, c.RoomRepo, c.RoomScheduleRepo, c.ProductService, c.Crypto)

	notificationAdapter := bookingNotification.NewBookingNotificationAdapter(
		sharedEmailClient,
		sharedSMSClient,
		c.Config.FrontendOrigin,
		bookingNotification.NewInProcessUserFetcher(c.UserService),
		bookingNotification.NewInProcessRoomFetcher(c.RoomService),
		bookingNotification.NewInProcessBuildingFetcher(c.BuildingService),
		c.ProductService,
	)
	c.BookingNotificationService = notificationAdapter

	c.BookingService = bookingSvc.New(
		c.BookingRepo,
		c.AvailabilityRepo,
		bookingStripe,
		c.ProductService,
		c.PriceService,
		notificationAdapter,
		c.Crypto,
		bookingSvc.WithRoomService(c.RoomService),
		bookingSvc.WithAuthUserClient(c.BookingAuthuserCLi),
		bookingSvc.WithTokenSecret([]byte(c.Config.BookingTokenSecret)),
	)

	c.PaymentService = bookingStripe

	// Wire BookingClient into AuthAggregator (now that BookingService is created)
	c.BookingClient = bookingClientAdapter.NewInProcessClient(c.BookingService)
	authAgg.SetBookingClient(c.BookingClient)

	// Messaging services
	bookChecker := messagingBooking.New(c.BookingService)
	nameFetcher := messagingAuthuser.New(c.UserService)
	c.MessagingBroker = messagingSSE.NewBroker(nil) // logger set later in server
	c.MessagingService = messagingSvc.New(c.MessageRepo, c.Crypto, bookChecker, nameFetcher, c.MessagingBroker)

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
