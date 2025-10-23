package imageHandler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/hengadev/errsx"
)

func (h *handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	const maxFormMemory = 32 << 20             // 32 MB
	err := r.ParseMultipartForm(maxFormMemory) // 32 MB max memory buffer
	if err != nil {
		log.Printf("Handler: Error parsing multipart form: %v", err)
		httpx.RespondWithError(w, errors.New("failed to parse form data, request too large or invalid"), http.StatusBadRequest)
		return
	}
	defer r.MultipartForm.RemoveAll()

	file, header, err := r.FormFile("image")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			httpx.RespondWithError(w, errors.New("required form file 'image' is missing"), http.StatusBadRequest)
		} else {
			log.Printf("Handler: Error retrieving form file: %v", err)
			httpx.RespondWithError(w, errors.New("failed to retrieve form file"), http.StatusBadRequest)
		}
		return
	}
	defer file.Close()

	// header.Size is the file size in bytes (int64).
	fileSize := header.Size
	contentType := header.Header.Get("Content-Type")

	// fields check
	var fieldErrs errsx.Map
	parentID := r.FormValue("parent_id")
	if parentID == "" {
		fieldErrs.Set("parent ID", "required form field 'parent_id' is missing")
	}

	parentType := r.FormValue("parent_type")
	if parentType == "" {
		fieldErrs.Set("parent type", "required form field 'parent_type' is missing")
	}

	title := r.FormValue("title")
	if title == "" {
		fieldErrs.Set("title", "required form field 'title' is missing")
	}

	isActiveStr := r.FormValue("is_active")
	isActive, err := strconv.ParseBool(isActiveStr)
	if err != nil {
		isActive = false
	}

	if fieldErrs != nil {
		httpx.RespondWithError(w, fieldErrs.AsError(), http.StatusBadRequest)
		return
	}

	req := &domain.CreateImageRequest{
		ParentID:   parentID,
		ParentType: domain.ParentType(parentType),
		Title:      title,
		IsActive:   &isActive, // & to make it a pointer
	}

	imageID, err := h.svc.AddImage(ctx, req, file, fileSize, contentType)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidValue):
			httpx.RespondWithError(w, err, http.StatusBadRequest)
		case errors.Is(err, errs.ErrDomainNotFound):
			httpx.RespondWithError(w, err, http.StatusNotFound)
		case errors.Is(err, errs.ErrExternalService):
			httpx.RespondWithError(w, fmt.Errorf("failed to upload image due to an external service error: %w", err), http.StatusServiceUnavailable)
		case errors.Is(err, errs.ErrUnexpectedError), errors.Is(err, errs.ErrUnexpectedError):
			log.Printf("Handler: Internal server error during image upload: %v", err)
			httpx.RespondWithError(w, errors.New("an internal server error occurred"), http.StatusInternalServerError)
		default:
			// A catch-all for any other unhandled errors.
			log.Printf("Handler: Unhandled error from service during image upload: %v", err)
			httpx.RespondWithError(w, errors.New("an unexpected error occurred"), http.StatusInternalServerError)
		}
		return
	}

	httpx.RespondWithJSON(w, struct {
		ID string `json:"id"`
	}{ID: imageID}, http.StatusCreated)

	log.Printf("Handler: Image created successfully. ID: %s", imageID)
}
