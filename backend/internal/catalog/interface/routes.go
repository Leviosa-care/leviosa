package handler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/category"
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/image"
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/price"
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/product"
)

func RegisterRoutes(
	router *http.ServeMux,
	categoryHandler categoryHandler.Handler,
	productHandler productHandler.Handler,
	priceHandler priceHandler.Handler,
	imageHandler imageHandler.Handler,
) {
	productHandler.RegisterRoutes(router)
	categoryHandler.RegisterRoutes(router)
	priceHandler.RegisterRoutes(router)
	imageHandler.RegisterRoutes(router)

}
