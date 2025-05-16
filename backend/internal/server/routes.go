package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hengadev/leviosa/internal/server/app"
	"github.com/hengadev/leviosa/internal/server/handler/event"
	"github.com/hengadev/leviosa/internal/server/handler/health"
	"github.com/hengadev/leviosa/internal/server/handler/product"
	"github.com/hengadev/leviosa/internal/server/handler/settings"
	"github.com/hengadev/leviosa/internal/server/handler/user"
	"github.com/hengadev/leviosa/internal/server/handler/vote"
	mw "github.com/hengadev/leviosa/internal/server/middleware"
)

func (s *Server) addRoutes(h *app.App) {
	router := http.NewServeMux()

	// basic check health
	router.HandleFunc("GET /healthz", healthz)
	router.HandleFunc("GET /hello", hello)

	// handlers declaration
	healthHandler := healthHandler.New(h)
	userHandler := userHandler.New(h)
	voteHandler := vote.NewHandler(h)
	eventHandler := eventHandler.New(h)
	productHandler := productHandler.New(h)
	settingsHandler := settingsHandler.New(h)

	// middlewares declaration
	rateLimit := mw.PerIPRateLimit(1, 1)

	router.HandleFunc("GET /health", healthHandler.CheckHealth)

	// user
	router.HandleFunc("GET /users/me", userHandler.GetUser)
	router.HandleFunc("PUT /users/me", userHandler.UpdateUser)
	router.HandleFunc("DELETE /users/me", userHandler.DeleteUser)

	router.HandleFunc("POST /users/exists", userHandler.CheckUserExists)

	router.HandleFunc("POST /auth/signin", rateLimit(userHandler.Signin))

	router.HandleFunc("POST /auth/register", rateLimit(userHandler.RegisterUserOTP))
	router.HandleFunc("POST /auth/validate-otp", rateLimit(userHandler.ValidateUserOTP))
	router.HandleFunc("GET /auth/approve-user", rateLimit(userHandler.GetUsersToApprove))
	router.HandleFunc("POST /auth/approve-user", rateLimit(userHandler.ApproveUserRegistration))
	router.HandleFunc("POST /auth/signout", userHandler.Signout)

	router.HandleFunc("POST /oauth/{provider}", userHandler.HandleOAuth)

	// vote
	router.HandleFunc("GET /votes/{month}/{year}", voteHandler.GetVotesByUserID)

	// settings
	router.HandleFunc("POST /settings/logo", settingsHandler.AddLogo)

	// register
	// NOTE: the old way to do the reservation thing
	// mux.Handle("POST /register", registerHandler.MakeRegistration())
	// TODO: the better way to do the reservation thing
	// mux.Handle("POST /register/event", registerHandler.MakeRegistration())
	// mux.Handle("POST /register/consultation", registerHandler.MakeRegistration())

	// products
	router.HandleFunc("POST /products", productHandler.CreateProduct)
	router.HandleFunc("GET /products/{id}", productHandler.GetProduct)
	router.HandleFunc("DELETE /products/{id}", productHandler.DeleteProduct)
	router.HandleFunc("PUT /products/{id}", productHandler.UpdateProduct)

	// product types
	router.HandleFunc("POST /product-types", productHandler.CreateOffer)
	router.HandleFunc("DELETE /product-types/{id}", productHandler.DeleteOffer)

	// event
	router.HandleFunc("GET /events/{id}", eventHandler.FindEventByID)
	router.HandleFunc("POST /events", eventHandler.CreateEvent)
	router.HandleFunc("PUT /events/{id}", eventHandler.ModifyEvent)
	router.HandleFunc("DELETE /events/{id}", eventHandler.FindEventByID)
	router.HandleFunc("GET /events/users", eventHandler.FindEventsForUser)
	router.HandleFunc("POST /upload-image", handleImage)

	s.srv.Handler = router
}

// TODO: how can I make groups for that thing and make sure that I can add as much middleware to a group as I want ?

func healthz(w http.ResponseWriter, r *http.Request) {
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
