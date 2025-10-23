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
	specializationHandler "github.com/Leviosa-care/leviosa/backend/internal/authuser/interface/specialization"
	userHandler "github.com/Leviosa-care/leviosa/backend/internal/authuser/interface/user"

	// Catalog HTTP handlers
	categoryHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/category"
	couponHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/coupon"
	imageHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/image"
	priceHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/price"
	productHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/product"
	promotionCodeHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/promotion_code"

	// Common
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

// Server represents the HTTP server
type Server struct {
	container  *Container
	httpServer *http.Server
	logger     *slog.Logger
}

// NewServer creates a new HTTP server with all routes
func NewServer(container *Container, logger *slog.Logger) *Server {
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
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.container.Config.ServerPort),
		Handler:      s.applyMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
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
}

func (s *Server) setupAuthuserRoutes(mux *http.ServeMux) {
	// Auth handler
	authH := authHandler.NewHandler(
		s.container.AuthAggregator,
		s.container.Crypto,
	)
	authHandler.RegisterRoutes(mux, authH)

	// User handler
	userH := userHandler.NewHandler(
		s.container.UserService,
		s.container.Crypto,
	)
	userHandler.RegisterRoutes(mux, userH)

	// Partner handler
	partnerH := partnerHandler.NewHandler(
		s.container.PartnerService,
		s.container.Crypto,
	)
	partnerHandler.RegisterRoutes(mux, partnerH)

	// Specialization handler
	specializationH := specializationHandler.NewHandler(
		s.container.SpecializationService,
		s.container.Crypto,
	)
	specializationHandler.RegisterRoutes(mux, specializationH)
}

func (s *Server) setupCatalogRoutes(mux *http.ServeMux) {
	// Category handler
	categoryH := categoryHandler.NewHandler(
		s.container.CatalogAggregator,
		s.container.Crypto,
	)
	categoryHandler.RegisterRoutes(mux, categoryH)

	// Product handler
	productH := productHandler.NewHandler(
		s.container.CatalogAggregator,
		s.container.Crypto,
	)
	productHandler.RegisterRoutes(mux, productH)

	// Price handler
	priceH := priceHandler.NewHandler(
		s.container.CatalogAggregator,
		s.container.Crypto,
	)
	priceHandler.RegisterRoutes(mux, priceH)

	// Image handler
	imageH := imageHandler.NewHandler(
		s.container.ImageService,
		s.container.Crypto,
	)
	imageHandler.RegisterRoutes(mux, imageH)

	// Coupon handler
	couponH := couponHandler.NewHandler(
		s.container.CatalogAggregator,
		s.container.Crypto,
	)
	couponHandler.RegisterRoutes(mux, couponH)

	// Promotion code handler
	promotionCodeH := promotionCodeHandler.NewHandler(
		s.container.CatalogAggregator,
		s.container.Crypto,
	)
	promotionCodeHandler.RegisterRoutes(mux, promotionCodeH)
}

func (s *Server) applyMiddleware(handler http.Handler) http.Handler {
	// CORS middleware
	handler = httpx.EnableCORS(handler)

	// Logging middleware
	handler = middleware.Logger(s.logger)(handler)

	// Recovery middleware
	handler = middleware.Recovery(s.logger)(handler)

	return handler
}
