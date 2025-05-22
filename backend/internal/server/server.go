package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/hengadev/leviosa/internal/server/app"
	mw "github.com/hengadev/leviosa/internal/server/middleware"
	"github.com/hengadev/leviosa/pkg/envmode"
)

type Server struct {
	srv *http.Server
}

func New(
	appCtx *app.App,
	env envmode.Mode,
	slogHandler slog.Handler,
	opts ...ServerOption,
) *Server {

	// build server with default options.
	server := &Server{
		srv: &http.Server{
			// might need to change the values of the timeouts
			Addr:              ":5000",
			IdleTimeout:       120 * time.Second,
			ReadTimeout:       1 * time.Second,
			WriteTimeout:      1 * time.Second,
			ReadHeaderTimeout: 1 * time.Second,
		},
	}
	// add server options
	for _, opt := range opts {
		opt(server)
	}

	// create router and add routes
	server.addRoutes(appCtx)

	// add middlewares common to all routes. [Order important]
	server.Use(
		mw.SetUserContext(appCtx.Svcs.Session.GetSession),
		mw.AttachLogger(env, slogHandler),
		mw.SetOrigin,
	)
	return server
}

func (s *Server) ListenAndServe() error {
	return s.srv.ListenAndServe()
}

func (s *Server) ListenAndServeTLS(certFile, keyFile string) error {
	return s.srv.ListenAndServeTLS(certFile, keyFile)
}

func (s *Server) Close() error {
	return s.srv.Close()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

// A function to add middleware to all the routes of the multiplexer.
func (s *Server) Use(middlewares ...func(http.Handler) http.Handler) {
	for _, mw := range middlewares {
		s.srv.Handler = mw(s.srv.Handler)
	}
}
