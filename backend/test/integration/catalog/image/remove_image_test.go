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
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestRemoveImage TEST_PATH=test/integration/catalog/image/remove_image_test.go

func TestRemoveImage(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}
	t.Run("should successfully remove an existing image", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)

		parentID := uuid.New()
		createCategoryParent(t, ctx, parentID) // Ensure parent exists

		imageID := uuid.New()

		contentType := "image/jpeg"

		s3Key, err := image.CreateParentImagePrefix(parentID, imageID, domain.CategoryType, contentType)
		require.NoError(t, err)
		mockImage := &domain.Image{
			ID: imageID, ParentID: parentID, ParentType: domain.CategoryType,
			Title: "Image to delete", S3Key: s3Key, Size: 100, ContentType: contentType,
			IsActive: false, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}
		insertImageIntoDB(t, ctx, mockImage)
		// Upload a dummy file to S3
		_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(td.BUCKETNAME), Key: aws.String(s3Key), Body: bytes.NewReader([]byte("dummy content")),
			ContentLength: aws.Int64(int64(len("dummy content"))), ContentType: aws.String("image/jpeg"),
		})
		require.NoError(t, err, "Failed to upload dummy S3 file for setup")

		// Pre-conditions
		require.True(t, imageExistsInDB(t, ctx, imageID), "Image should exist in DB before deletion")
		require.True(t, fileExistsInS3(t, ctx, s3Key), "File should exist in S3 before deletion")

		requestBody := domain.ImageModifierRequest{
			ImageID:    imageID.String(),
			ParentID:   parentID.String(),
			ParentType: string(domain.CategoryType),
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newRemoveImageRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Post-conditions
		assert.False(t, imageExistsInDB(t, ctx, imageID), "Image should be deleted from DB")
		assert.False(t, fileExistsInS3(t, ctx, s3Key), "File should be deleted from S3")
	})

	t.Run("should return 400 Bad Request if image_id is missing", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)

		requestBody := domain.ImageModifierRequest{
			ImageID:    "", // Missing
			ParentID:   uuid.New().String(),
			ParentType: string(domain.ProductType),
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newRemoveImageRequest(t, ctx, jsonBody)
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
		assert.Contains(t, respBody.Error, "image ID: invalid value, must be a valid UUID.")
	})

	t.Run("should return 400 Bad Request if parent_id is missing", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)

		requestBody := domain.ImageModifierRequest{
			ImageID:    uuid.New().String(),
			ParentID:   "", // Missing
			ParentType: string(domain.ProductType),
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newRemoveImageRequest(t, ctx, jsonBody)
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
		assert.Contains(t, respBody.Error, "parent ID: invalid value, must be a valid UUID.")
	})

	t.Run("should return 400 Bad Request if parent_type is invalid", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)

		requestBody := domain.ImageModifierRequest{
			ImageID:    uuid.New().String(),
			ParentID:   uuid.New().String(),
			ParentType: "invalid_type", // Invalid
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newRemoveImageRequest(t, ctx, jsonBody)
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
		assert.Contains(t, respBody.Error, "parent type: invalid parent type value.")
	})

	t.Run("should return 404 Not Found if image does not exist", func(t *testing.T) {
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

		req := newRemoveImageRequest(t, ctx, jsonBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		fmt.Println("here the status code is:", resp.StatusCode)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Contains(t, respBody.Error, errs.ErrDomainNotFound.Error(), nonExistentImageID)
	})

	t.Run("should return 400 Bad Request if image exists but parent_id mismatch", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)

		correctParentID := uuid.New()
		mismatchedParentID := uuid.New()
		createCategoryParent(t, ctx, correctParentID) // Ensure correct parent exists

		imageID := uuid.New()

		contentType := "image/jpeg"

		// s3Key := fmt.Sprintf("images/category/%s.jpg", imageID)
		s3Key, err := image.CreateParentImagePrefix(correctParentID, imageID, domain.CategoryType, contentType)
		require.NoError(t, err)

		mockImage := &domain.Image{
			ID: imageID, ParentID: correctParentID, ParentType: domain.CategoryType,
			Title: "Image to delete", S3Key: s3Key, Size: 100, ContentType: contentType,
			IsActive: false, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}
		insertImageIntoDB(t, ctx, mockImage)
		// Upload a dummy file to S3
		_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(td.BUCKETNAME),
			Key:    aws.String(s3Key), Body: bytes.NewReader([]byte("dummy content")),
			ContentLength: aws.Int64(int64(len("dummy content"))), ContentType: aws.String("image/jpeg"),
		})
		require.NoError(t, err, "Failed to upload dummy S3 file for setup")

		requestBody := domain.ImageModifierRequest{
			ImageID:    imageID.String(),
			ParentID:   mismatchedParentID.String(), // Mismatch
			ParentType: string(domain.CategoryType),
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := newRemoveImageRequest(t, ctx, jsonBody)
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
		assert.Contains(t, respBody.Error, fmt.Sprintf("image ID %s does not belong to category with ID %s", imageID, mismatchedParentID))

		// Verify image and file are NOT deleted
		assert.True(t, imageExistsInDB(t, ctx, imageID), "Image should NOT be deleted from DB on mismatch")
		assert.True(t, fileExistsInS3(t, ctx, s3Key), "File should NOT be deleted from S3 on mismatch")
	})
}

func newRemoveImageRequest(t *testing.T, ctx context.Context, jsonBody []byte) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, testServerURL+"/admin/images", bytes.NewReader(jsonBody))
	require.NoError(t, err)
	return req
}
