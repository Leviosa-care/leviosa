package image_test

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestUploadImage(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully upload a new image for a category parent", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)

		parentID := uuid.New()
		createCategoryParent(t, ctx, parentID) // Create a parent for the image to attach to.

		fileContent := "test image content"
		formBody, contentType := createMultipartForm(t, map[string]string{
			"parent_id":   parentID.String(),
			"parent_type": string(domain.CategoryType),
			"title":       "My Test Image",
		}, "test-image.jpg", fileContent, "image/jpeg")

		req := newUploadImageRequest(t, ctx, formBody)
		req.Header.Set("Content-Type", contentType)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var responseBody struct {
			ID string `json:"id"`
		}
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err)
		assert.NotEqual(t, "", responseBody.ID)

		imageID, err := uuid.Parse(responseBody.ID)
		assert.NoError(t, err)

		// Verify the image was created in the database
		var count int
		err = testPool.QueryRow(ctx, "SELECT count(*) FROM catalog.images WHERE id = $1", responseBody.ID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)

		key, err := image.CreateParentImagePrefix(parentID, imageID, domain.CategoryType, "image/jpeg")
		require.NoError(t, err)

		require.NoError(t, err)
		_, err = s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(td.BUCKETNAME),
			Key:    aws.String(key),
		})
		assert.NoError(t, err)
	})

	t.Run("should successfully upload a new image for a product parent", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)

		parentID := uuid.New()
		createProductParent(t, ctx, parentID) // Create a product parent for the image.

		fileContent := "test image content"
		product_type := domain.ProductType
		formBody, contentType := createMultipartForm(t, map[string]string{
			"parent_id":   parentID.String(),
			"parent_type": string(product_type),
			"title":       "My Product Image",
			"is_active":   "true",
		}, "product-image.jpg", fileContent, "image/jpeg")

		req := newUploadImageRequest(t, ctx, formBody)
		req.Header.Set("Content-Type", contentType)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var responseBody struct {
			ID string `json:"id"`
		}
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err)
		assert.NotEqual(t, "", responseBody.ID)

		imageID, err := uuid.Parse(responseBody.ID)
		assert.NoError(t, err)

		// Verify the image was created in the database
		var count int
		err = testPool.QueryRow(ctx, "SELECT count(*) FROM catalog.images WHERE id = $1", responseBody.ID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)

		key, err := image.CreateParentImagePrefix(parentID, imageID, product_type, "image/jpeg")
		require.NoError(t, err)

		// Verify the file was uploaded to S3 with the correct key format for a product
		_, err = s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(td.BUCKETNAME),
			Key:    aws.String(key),
		})
		assert.NoError(t, err)
	})

	t.Run("should fail if required form fields are missing", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)

		fileContent := "test image content"
		formBody, contentType := createMultipartForm(t, map[string]string{
			"parent_id": uuid.New().String(),
			// "parent_type" is missing
			"title":     "My Test Image",
			"is_active": "false",
		}, "test-image.jpg", fileContent, "image/jpeg")

		req := newUploadImageRequest(t, ctx, formBody)
		req.Header.Set("Content-Type", contentType)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Contains(t, respBody.Error, "required form field 'parent_type' is missing")
	})

	t.Run("should fail if image file is missing", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)

		formBody, contentType := createMultipartForm(t, map[string]string{
			"parent_id":   uuid.New().String(),
			"parent_type": string(domain.CategoryType),
			"title":       "My Test Image",
		}, "", "", "image/jpeg") // No file content

		req := newUploadImageRequest(t, ctx, formBody)
		req.Header.Set("Content-Type", contentType)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail if parent ID does not exist", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)

		nonExistentParentID := uuid.New() // A valid UUID that does not exist in the DB

		fileContent := "test image content"
		formBody, contentType := createMultipartForm(t, map[string]string{
			"parent_id":   nonExistentParentID.String(),
			"parent_type": string(domain.CategoryType),
			"title":       "My Test Image",
			"is_active":   "true",
		}, "test-image.jpg", fileContent, "image/jpeg")

		req := newUploadImageRequest(t, ctx, formBody)
		req.Header.Set("Content-Type", contentType)

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

	t.Run("should succeed if a second active image is created for the same parent", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)

		parentID := uuid.New()
		createCategoryParent(t, ctx, parentID)

		// Upload the first active image successfully
		fileContent1 := "first image"
		formBody1, contentType1 := createMultipartForm(t, map[string]string{
			"parent_id":   parentID.String(),
			"parent_type": string(domain.CategoryType),
			"title":       "First Active Image",
			"is_active":   "true",
		}, "image1.jpg", fileContent1, "image/jpeg")

		req1 := newUploadImageRequest(t, ctx, formBody1)
		req1.Header.Set("Content-Type", contentType1)
		resp1, _ := client.Do(req1)
		resp1.Body.Close()
		assert.Equal(t, http.StatusCreated, resp1.StatusCode)

		// Now attempt to upload a second active image for the same parent
		fileContent2 := "second image"
		title2 := "Second Active Image"
		formBody2, contentType2 := createMultipartForm(t, map[string]string{
			"parent_id":   parentID.String(),
			"parent_type": string(domain.CategoryType),
			"title":       title2,
			"is_active":   "true",
		}, "image2.jpg", fileContent2, "image/jpeg")

		req2 := newUploadImageRequest(t, ctx, formBody2)
		req2.Header.Set("Content-Type", contentType2)
		resp2, _ := client.Do(req2)
		defer resp2.Body.Close()

		assert.Equal(t, http.StatusCreated, resp2.StatusCode)
		var respBody struct {
			ID string `json:"id"`
		}
		err := json.NewDecoder(resp2.Body).Decode(&respBody)
		assert.NoError(t, err)

		err = uuid.Validate(respBody.ID)
		assert.NoError(t, err)

		// Verify that the second image record was rolled back from the database
		var count int
		err = testPool.QueryRow(ctx, "SELECT count(*) FROM catalog.images WHERE title = $1", title2).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count, "The second image record should have been rolled back")
	})
}

func newUploadImageRequest(t *testing.T, ctx context.Context, formBody *bytes.Buffer) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/admin/images", formBody)
	require.NoError(t, err)
	return req
}
