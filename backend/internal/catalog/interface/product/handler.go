package productHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware/auth"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
	CreateProductWithPrice(w http.ResponseWriter, r *http.Request)
	GetAdminAllProducts(w http.ResponseWriter, r *http.Request)
	GetAllPublishedProducts(w http.ResponseWriter, r *http.Request)
	GetProductByID(w http.ResponseWriter, r *http.Request)
	ModifyProduct(w http.ResponseWriter, r *http.Request)
	RemoveProduct(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	productService ports.ProductService
	aggr           ports.ProductAggregatorService
	authmw         auth.AuthMiddleware
}

func New(productService ports.ProductService, aggr ports.ProductAggregatorService, authmw auth.AuthMiddleware) Handler {
	return &handler{
		productService: productService,
		aggr:           aggr,
		authmw:         authmw,
	}
}

// NOTE: the old way of doing
// type handler struct {
// 	productService ports.ProductService
// 	imageService   ports.ImageParentService
// 	priceService   ports.PriceService
// }
//
// func New(productService ports.ProductService, imageService ports.ImageParentService, priceService ports.PriceService) Handler {
// 	return &handler{
// 		productService: productService,
// 		imageService:   imageService,
// 		priceService:   priceService,
// 	}
// }
