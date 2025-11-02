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

	// Catalog HTTP handlers
	categoryHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/category"
	couponHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/coupon"
	imageHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/image"
	priceHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/price"
	productHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/product"
	promotionCodeHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/promotion_code"

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

func (s *Server) setupAuthuserRoutes(router *http.ServeMux) {
	// Auth handler
	authH := authHandler.New(
		s.container.AuthAggregator,
		s.container.Crypto,
	)
	authH.RegisterRoutes(router)

	// User handler
	userH := userHandler.New(
		s.container.UserService,
		s.container.Crypto,
	)
	userH.RegisterRoutes(router)

	// Partner handler
	partnerH := partnerHandler.New(
		s.container.PartnerService,
		s.container.Crypto,
	)
	partnerH.RegisterRoutes(router)
}

func (s *Server) setupCatalogRoutes(router *http.ServeMux) {
	// Category handler
	categoryH := categoryHandler.New(
		s.container.CatalogAggregator,
		&s.container.ImageService,
		s.container.Crypto,
	)
	categoryH.RegisterRoutes(router)

	// Product handler
	productH := productHandler.New(
		s.container.CatalogAggregator,
		s.container.Crypto,
	)
	productH.RegisterRoutes(router)

	// Price handler
	priceH := priceHandler.New(
		s.container.CatalogAggregator,
		s.container.Crypto,
	)
	priceH.RegisterRoutes(router)

	// Image handler
	imageH := imageHandler.New(
		s.container.ImageService,
		s.container.Crypto,
	)
	imageH.RegisterRoutes(router)

	// Coupon handler
	couponH := couponHandler.New(
		s.container.CatalogAggregator,
		s.container.Crypto,
	)
	couponH.RegisterRoutes(router)

	// Promotion code handler
	promotionCodeH := promotionCodeHandler.New(
		s.container.CatalogAggregator,
		s.container.Crypto,
	)
	promotionCodeH.RegisterRoutes(router)
}

// func (s *Server) applyMiddleware(handler http.Handler) http.Handler {
func (s *Server) applyMiddleware(handler middleware.Handler) middleware.Handler {
	// CORS middleware
	handler = middleware.EnableCORS(handler)

	// Logging middleware
	handler = middleware.AttachLogger(envmode.Prod, s.logger)(handler)

	// Recovery middleware
	handler = middleware.Recovery(s.logger)(handler)

	return handler
}
