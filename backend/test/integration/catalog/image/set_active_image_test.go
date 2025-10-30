package image_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/application/image"
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetActiveImage(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully set an inactive image as active and deactivate old one", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)

		parentID := uuid.New()
		createCategoryParent(t, ctx, parentID) // Ensure parent exists

		contentType := "image/jpeg"

		// Setup: Create an existing active image
		activeImageID := uuid.New()

		activeImageS3Key, err := image.CreateParentImagePrefix(parentID, activeImageID, domain.CategoryType, contentType)
		require.NoError(t, err)

		activeImage := &domain.Image{
			ID: activeImageID, ParentID: parentID, ParentType: domain.CategoryType,
			Title: "Old Active Image", S3Key: activeImageS3Key,
			Size: 100, ContentType: contentType, IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}
		insertImageIntoDB(t, ctx, activeImage)

		// Setup: Create a new inactive image that will be set active
		newActiveImageID := uuid.New()

		newActiveImageS3Key, err := image.CreateParentImagePrefix(parentID, newActiveImageID, domain.CategoryType, contentType)
		require.NoError(t, err)

		newActiveImage := &domain.Image{
			ID: newActiveImageID, ParentID: parentID, ParentType: domain.CategoryType,
			Title: "New Active Image", S3Key: newActiveImageS3Key,
			Size: 100, ContentType: contentType, IsActive: false, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}
		insertImageIntoDB(t, ctx, newActiveImage)

		// Pre-conditions
		assert.True(t, td.GetImageStatus(t, ctx, activeImageID, testPool), "Old image should be active initially")
		assert.False(t, td.GetImageStatus(t, ctx, newActiveImageID, testPool), "New image should be inactive initially")

		requestBody := domain.ImageModifierRequest{
			ImageID:    newActiveImageID.String(),
			ParentID:   parentID.String(),
			ParentType: string(domain.CategoryType),
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newSetActiveImageRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Post-conditions
		assert.True(t, td.GetImageStatus(t, ctx, newActiveImageID, testPool), "New image should be active after request")
		assert.False(t, td.GetImageStatus(t, ctx, activeImageID, testPool), "Old image should be inactive after request")
	})

	t.Run("should successfully set an inactive image as active when no other is active", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)

		parentID := uuid.New()
		createCategoryParent(t, ctx, parentID) // Ensure parent exists

		contentType := "image/jpeg"

		// Setup: Create only one inactive image
		imageID := uuid.New()

		s3Key, err := image.CreateParentImagePrefix(parentID, imageID, domain.CategoryType, contentType)
		require.NoError(t, err)

		mockImage := &domain.Image{
			ID: imageID, ParentID: parentID, ParentType: domain.CategoryType,
			Title: "Single Image", S3Key: s3Key,
			Size: 100, ContentType: contentType, IsActive: false, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}
		insertImageIntoDB(t, ctx, mockImage)

		// Pre-condition
		assert.False(t, td.GetImageStatus(t, ctx, imageID, testPool), "Image should be inactive initially")

		requestBody := domain.ImageModifierRequest{
			ImageID:    imageID.String(),
			ParentID:   parentID.String(),
			ParentType: string(domain.CategoryType),
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newSetActiveImageRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Post-condition
		assert.True(t, td.GetImageStatus(t, ctx, imageID, testPool), "Image should be active after request")
	})

	t.Run("should return 400 Bad Request if image_id is missing", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)

		requestBody := domain.ImageModifierRequest{
			ImageID:    "", // Missing
			ParentID:   uuid.New().String(),
			ParentType: string(domain.ProductType),
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newSetActiveImageRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		// assert.Contains(t, respBody.Error, "image ID: invalid value, must be a valid UUID.")
		assert.Contains(t, respBody.Error, errs.ErrInvalidValue.Error())
	})

	t.Run("should return 400 Bad Request if parent_id is missing", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)

		requestBody := domain.ImageModifierRequest{
			ImageID:    uuid.New().String(),
			ParentID:   "", // Missing
			ParentType: string(domain.ProductType),
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newSetActiveImageRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Contains(t, respBody.Error, errs.ErrInvalidValue.Error())
	})

	t.Run("should return 400 Bad Request if parent_type is invalid", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)

		requestBody := domain.ImageModifierRequest{
			ImageID:    uuid.New().String(),
			ParentID:   uuid.New().String(),
			ParentType: "invalid_type", // Invalid
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newSetActiveImageRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Contains(t, respBody.Error, errs.ErrInvalidValue.Error())
	})

	t.Run("should return 404 Not Found if image to activate does not exist", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)

		parentID := uuid.New()
		createCategoryParent(t, ctx, parentID) // Ensure parent exists

		nonExistentImageID := uuid.New()

		requestBody := domain.ImageModifierRequest{
			ImageID:    nonExistentImageID.String(),
			ParentID:   parentID.String(),
			ParentType: string(domain.CategoryType),
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newSetActiveImageRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Contains(t, respBody.Error, errs.ErrDomainNotFound.Error())
	})

	t.Run("should return 404 Not Found if parent for image does not exist", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)

		nonExistentParentID := uuid.New()
		imageID := uuid.New()

		// Setup: Create the image, but its parent does NOT exist in the DB
		mockImage := &domain.Image{
			ID: imageID, ParentID: nonExistentParentID, ParentType: domain.CategoryType,
			Title: "Image with non-existent parent", S3Key: fmt.Sprintf("images/category/%s.jpg", imageID),
			Size: 100, ContentType: "image/jpeg", IsActive: false, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}
		insertImageIntoDB(t, ctx, mockImage) // Image exists, but its parent doesn't

		requestBody := domain.ImageModifierRequest{
			ImageID:    imageID.String(),
			ParentID:   nonExistentParentID.String(),
			ParentType: string(domain.CategoryType),
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newSetActiveImageRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Contains(t, respBody.Error, errs.ErrDomainNotFound.Error())
	})

	t.Run("should return 404 Not Found if image exists but parent_id/parent_type mismatch", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)

		correctParentID := uuid.New()
		mismatchedParentID := uuid.New()
		createCategoryParent(t, ctx, correctParentID)                                   // Ensure correct parent exists
		createCategoryParentWithName(t, ctx, mismatchedParentID, "Other test category") // Ensure mismatched parent also exists

		contentType := "image/jpeg"

		imageID := uuid.New()

		s3Key, err := image.CreateParentImagePrefix(correctParentID, imageID, domain.CategoryType, contentType)
		require.NoError(t, err)

		mockImage := &domain.Image{
			ID: imageID, ParentID: correctParentID, ParentType: domain.CategoryType, // Image belongs to correctParentID
			Title: "Image with correct parent", S3Key: s3Key,
			Size: 100, ContentType: contentType, IsActive: false, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}
		insertImageIntoDB(t, ctx, mockImage)

		requestBody := domain.ImageModifierRequest{
			ImageID:    imageID.String(),
			ParentID:   mismatchedParentID.String(), // Mismatch: trying to activate with wrong parent
			ParentType: string(domain.CategoryType),
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newSetActiveImageRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		// TODO: find the right message to check
		assert.Contains(t, respBody.Error, errs.ErrDomainNotFound.Error())

		// Verify image status hasn't changed
		assert.False(t, td.GetImageStatus(t, ctx, imageID, testPool), "Image status should not change on mismatch")
	})
}

func newSetActiveImageRequest(t *testing.T, ctx context.Context, jsonBody []byte) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/admin/images/set-active", bytes.NewReader(jsonBody))
	require.NoError(t, err)
	return req
}
