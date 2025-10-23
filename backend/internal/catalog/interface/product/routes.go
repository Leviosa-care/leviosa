package productHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /products", middleware.EnableCORS(h.GetAllPublishedProducts))
	router.HandleFunc("GET /products/{id}", middleware.EnableCORS(h.GetProductByID))
	router.HandleFunc("GET /admin/products", middleware.EnableCORS(h.GetAdminAllProducts))
	router.HandleFunc("POST /admin/products", middleware.EnableCORS(h.CreateProductWithPrice))
	router.HandleFunc("PATCH /admin/products/{id}", middleware.EnableCORS(h.ModifyProduct))
	router.HandleFunc("DELETE /admin/products/{id}", middleware.EnableCORS(h.RemoveProduct))
}
