package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/hengadev/leviosa/internal/server/app"
	eventHandler "github.com/hengadev/leviosa/internal/server/handler/event"
	healthHandler "github.com/hengadev/leviosa/internal/server/handler/health"
	productHandler "github.com/hengadev/leviosa/internal/server/handler/product"
	settingsHandler "github.com/hengadev/leviosa/internal/server/handler/settings"
	userHandler "github.com/hengadev/leviosa/internal/server/handler/user"
	"github.com/hengadev/leviosa/internal/server/handler/vote"
	mw "github.com/hengadev/leviosa/internal/server/middleware"
	"github.com/hengadev/leviosa/pkg/ctxutil"
)

func (s *Server) addRoutes(h *app.App) {
	router := http.NewServeMux()

	// basic
	router.HandleFunc("GET /healthz", healthz)
	router.HandleFunc("GET /hello", hello)

	// handlers declaration
	hh := healthHandler.New(h)
	uh := userHandler.New(h)
	vh := vote.NewHandler(h)
	eh := eventHandler.New(h)
	ph := productHandler.New(h)
	sh := settingsHandler.New(h)

	uh.RegisterRoutes(router)
	sh.RegisterRoutes(router)
	hh.RegisterRoutes(router)

	// middlewares declaration
	rateLimit := mw.PerIPRateLimit(1, 1)

	// user
	router.HandleFunc("GET /users/me", uh.GetUser)
	router.HandleFunc("PUT /users/me", uh.UpdateUser)
	router.HandleFunc("DELETE /users/me", uh.DeleteUser)

	router.HandleFunc("POST /users/exists", uh.CheckUserExists)

	// auth
	router.HandleFunc("POST /auth/signin", rateLimit(uh.Signin))
	router.HandleFunc("POST /auth/email", rateLimit(uh.VerifyEmail))
	router.HandleFunc("POST /auth/otp", rateLimit(uh.ValidateUserOTP))
	// TODO:
	// auth/general
	// auth/address
	// auth/password
	// auth/pending
	router.HandleFunc("POST /auth/register", rateLimit(uh.RegisterUserOTP))
	router.HandleFunc("GET /auth/approve-user", rateLimit(uh.GetUsersToApprove))
	router.HandleFunc("POST /auth/approve-user", rateLimit(uh.ApproveUserRegistration))
	router.HandleFunc("POST /auth/signout", uh.Signout)

	router.HandleFunc("POST /oauth/{provider}", uh.HandleOAuth)

	// vote
	router.HandleFunc("GET /votes/{month}/{year}", vh.GetVotesByUserID)

	// register
	// NOTE: the old way to do the reservation thing
	// mux.Handle("POST /register", registerHandler.MakeRegistration())
	// TODO: the better way to do the reservation thing
	// mux.Handle("POST /register/event", registerHandler.MakeRegistration())
	// mux.Handle("POST /register/consultation", registerHandler.MakeRegistration())

	// products
	router.HandleFunc("POST /admin/products", ph.CreateProduct)
	router.HandleFunc("GET /products/{id}", ph.GetProduct)
	router.HandleFunc("DELETE /admin/products/{id}", ph.DeleteProduct)
	router.HandleFunc("PUT /admin/products/{id}", ph.UpdateProduct)

	// categories
	router.HandleFunc("POST /product-types", ph.CreateOffer)
	router.HandleFunc("DELETE /product-types/{id}", ph.DeleteOffer)

	// events
	router.HandleFunc("GET /events/{id}", eh.FindEventByID)
	router.HandleFunc("POST /events", eh.CreateEvent)
	router.HandleFunc("PUT /events/{id}", eh.ModifyEvent)
	router.HandleFunc("DELETE /events/{id}", eh.FindEventByID)
	router.HandleFunc("GET /events/users", eh.FindEventsForUser)
	router.HandleFunc("POST /upload-image", handleImage)

	s.srv.Handler = router
}

// TODO: how can I make groups for that thing and make sure that I can add as much middleware to a group as I want ?

func healthz(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "logger not found in context", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.ErrorContext(ctx, "Just hit the simple 'hello' endpoint")
	w.Write([]byte("OK"))
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hit server")
	message := struct {
		Message string `json:"message"`
	}{
		Message: "Hello world",
	}
	json.NewEncoder(w).Encode(message)
}

func handleImage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("here we are in the handle image handler")
	err := r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to retrieve image file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Println("the filename that I uploaded is:", handler.Filename)

	w.WriteHeader(http.StatusOK)
}
