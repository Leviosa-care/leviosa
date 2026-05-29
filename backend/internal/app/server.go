package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	// Authuser HTTP handlers
	authHandler "github.com/Leviosa-care/leviosa/backend/internal/authuser/interface/auth"
	partnerHandler "github.com/Leviosa-care/leviosa/backend/internal/authuser/interface/partner"
	userHandler "github.com/Leviosa-care/leviosa/backend/internal/authuser/interface/user"

	// Booking HTTP handlers
	allocationHandler "github.com/Leviosa-care/leviosa/backend/internal/booking/interface/allocation"
	availabilityHandler "github.com/Leviosa-care/leviosa/backend/internal/booking/interface/availability"
	bookingHandler "github.com/Leviosa-care/leviosa/backend/internal/booking/interface/booking"
	buildingHandler "github.com/Leviosa-care/leviosa/backend/internal/booking/interface/building"
	metricsHandler "github.com/Leviosa-care/leviosa/backend/internal/booking/interface/metrics"
	roomHandler "github.com/Leviosa-care/leviosa/backend/internal/booking/interface/room"

	// Catalog HTTP handlers
	categoryHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/category"
	couponHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/coupon"
	imageHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/image"
	priceHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/price"
	productHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/product"
	promotionCodeHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/promotion_code"

	// Messaging HTTP handlers
	messagingHandler "github.com/Leviosa-care/leviosa/backend/internal/messaging/interface"

	// Common
	"github.com/Leviosa-care/leviosa/backend/internal/common/envmode"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

// Server represents the HTTP server
type Server struct {
	container  *Container
	httpServer *http.Server
	logger     *slog.Logger
}

// NewServer creates a new HTTP server with all routes.
// It initialises package-level middleware state (e.g. CORS origin) before
// the server can accept connections.
func NewServer(container *Container, logger *slog.Logger) *Server {
	middleware.SetAllowedOrigin(container.Config.FrontendOrigin)

	return &Server{
		container: container,
		logger:    logger,
	}
}

// Start starts the HTTP server
func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	// Setup routes
	s.setupRoutes(mux)

	// Create HTTP server
	// WriteTimeout is 0 (disabled) because SSE connections stay open
	// indefinitely. Individual handlers are protected by context cancellation
	// and the server's graceful shutdown.
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.container.Config.ServerPort),
		Handler:      s.applyMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 0, // disabled for SSE
		IdleTimeout:  120 * time.Second,
	}

	s.logger.InfoContext(ctx, "Starting server", "port", s.container.Config.ServerPort)

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.InfoContext(ctx, "Shutting down server...")

	// Shutdown HTTP server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}

	// Close container resources
	if err := s.container.Close(ctx); err != nil {
		return fmt.Errorf("close container: %w", err)
	}

	s.logger.InfoContext(ctx, "Server shutdown complete")
	return nil
}

func (s *Server) setupRoutes(mux *http.ServeMux) {
	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		httpx.RespondWithJSON(w, map[string]string{"status": "ok"}, http.StatusOK)
	})

	// Authuser routes
	s.setupAuthuserRoutes(mux)

	// Catalog routes
	s.setupCatalogRoutes(mux)

	// Booking routes
	s.setupBookingRoutes(mux)

	// Messaging routes
	s.setupMessagingRoutes(mux)
}

func (s *Server) setupAuthuserRoutes(router *http.ServeMux) {
	// Auth handler
	authH := authHandler.New(
		s.container.AuthAggregator,
		s.container.AuthMw,
	)
	authH.RegisterRoutes(router)

	// User handler
	userH := userHandler.New(
		s.container.UserService,
		s.container.AuthMw,
	)
	userH.RegisterRoutes(router)

	partnerH := partnerHandler.New(
		s.container.PartnerService,
		s.container.AuthMw,
	)
	partnerH.RegisterRoutes(router)
}

func (s *Server) setupCatalogRoutes(router *http.ServeMux) {
	// Category handler
	categoryH := categoryHandler.New(
		s.container.CategoryService,
		s.container.ImageService,
		s.container.CategoryAggregator,
		s.container.AuthMw,
	)
	categoryH.RegisterRoutes(router)

	// Product handler
	productH := productHandler.New(
		s.container.ProductService,
		s.container.CatalogAggregator,
		s.container.AuthMw,
	)
	productH.RegisterRoutes(router)

	// Price handler
	priceH := priceHandler.New(
		s.container.PriceService,
		s.container.AuthMw,
	)
	priceH.RegisterRoutes(router)

	// Image handler
	imageH := imageHandler.New(
		s.container.ImageService,
		s.container.AuthMw,
	)
	imageH.RegisterRoutes(router)

	// Coupon handler
	couponH := couponHandler.New(
		s.container.CouponService,
		s.container.AuthMw,
	)
	couponH.RegisterRoutes(router)

	// Promotion code handler
	promotionCodeH := promotionCodeHandler.New(
		s.container.PromotionCodeService,
		s.container.AuthMw,
	)
	promotionCodeH.RegisterRoutes(router)
}

func (s *Server) setupBookingRoutes(router *http.ServeMux) {
	// Building handler
	buildingH := buildingHandler.New(
		s.container.BuildingService,
		s.container.AuthMw,
	)
	buildingH.RegisterRoutes(router)

	// Room handler
	roomH := roomHandler.New(
		s.container.RoomService,
		s.container.AuthMw,
	)
	roomH.RegisterRoutes(router)

	// Allocation handler
	allocationH := allocationHandler.New(
		s.container.AllocationService,
		s.container.AuthMw,
	)
	allocationH.RegisterRoutes(router)

	// Availability handler
	availabilityH := availabilityHandler.New(
		s.container.AvailabilityService,
		s.container.AuthMw,
	)
	availabilityH.RegisterRoutes(router)

	// Booking handler
	bookingH := bookingHandler.New(
		s.container.BookingService,
		s.container.PaymentService,
		s.container.AuthMw,
	)
	bookingH.RegisterRoutes(router)

	// Metrics handler (analytics endpoints)
	metricsH := metricsHandler.New(
		s.container.MetricsService,
		s.container.AuthMw,
	)
	metricsH.RegisterRoutes(router)
}

func (s *Server) setupMessagingRoutes(router *http.ServeMux) {
	messagingH := messagingHandler.New(
		s.container.MessagingService,
		s.container.AuthMw,
		s.container.MessagingBroker,
	)
	messagingH.RegisterRoutes(router)
}

func (s *Server) applyMiddleware(handler http.Handler) http.Handler {
	// CORS middleware
	corsHandler := middleware.EnableCORS(handler.ServeHTTP)
	handler = http.HandlerFunc(corsHandler)

	// Logging middleware
	handler = middleware.AttachLogger(envmode.Prod, s.logger)(handler)

	return handler
}
