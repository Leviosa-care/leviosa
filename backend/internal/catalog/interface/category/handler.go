package categoryHandler

import (
	"net/http"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/middleware/auth"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
	CreateCategory(w http.ResponseWriter, r *http.Request)
	GetAdminAllCategories(w http.ResponseWriter, r *http.Request)
	GetAllPublishedCategories(w http.ResponseWriter, r *http.Request)
	GetCategoryByID(w http.ResponseWriter, r *http.Request)
	// TODO: add something to get all the images for a given category
	// something like GetCategoryImages
	ModifyCategory(w http.ResponseWriter, r *http.Request)
	RemoveCategory(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	svc    ports.CategoryService
	aggr   ports.CategoryImagesService
	authmw auth.AuthMiddleware
}

func New(categoryService ports.CategoryService, imageService ports.ImageParentService, categoryImagesService ports.CategoryImagesService, authmw auth.AuthMiddleware) Handler {
	return &handler{
		svc:    categoryService,
		aggr:   categoryImagesService,
		authmw: authmw,
	}
}

// NOTE: the old thing that I used to do
// type handler struct {
// 	svc   ports.CategoryService
// 	image ports.ImageParentService
// }
// func New(categoryService ports.CategoryService, imageService ports.ImageParentService) Handler {
// 	return &handler{
// 		svc:   categoryService,
// 		image: imageService,
// 	}
// }
