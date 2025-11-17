package image_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/application/image"
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"
	th "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestUploadImage TEST_PATH=test/integration/catalog/image/upload_image_test.go

func TestUploadImage(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// Create a minimal valid JPEG file (JPEG magic bytes + minimal structure)
	fileContent := []byte{
		0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01,
		0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0xFF, 0xD9,
	}

	t.Run("should successfully upload a new image for a category parent with valid admin token", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		parentID := uuid.New()
		createCategoryParent(t, ctx, parentID) // Create a parent for the image to attach to.

		additionalFields := map[string]string{
			"parent_id":   parentID.String(),
			"parent_type": string(domain.CategoryType),
			"title":       "My Test Image",
		}

		req := th.NewUploadImageRequest(t, ctx, testServerURL, fileContent, "image", "test-image.jpg", additionalFields, accessToken)

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
		assert.NoError(t, err)
		_, err = s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(td.BUCKETNAME),
			Key:    aws.String(key),
		})
		assert.NoError(t, err)
	})

	t.Run("should successfully upload a new image for a product parent with valid admin token", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		parentID := uuid.New()
		createProductParent(t, ctx, parentID) // Create a product parent for the image.

		product_type := domain.ProductType
		additionalFields := map[string]string{
			"parent_id":   parentID.String(),
			"parent_type": string(product_type),
			"title":       "My Product Image",
			"is_active":   "true",
		}

		req := th.NewUploadImageRequest(t, ctx, testServerURL, fileContent, "image", "product-image.jpg", additionalFields, accessToken)

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
		assert.NoError(t, err)

		// Verify the file was uploaded to S3 with the correct key format for a product
		_, err = s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(td.BUCKETNAME),
			Key:    aws.String(key),
		})
		assert.NoError(t, err)
	})

	t.Run("should fail if required form fields are missing", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		additionalFields := map[string]string{
			"parent_id": uuid.New().String(),
			// "parent_type" is missing
			"title":     "My Test Image",
			"is_active": "false",
		}

		req := th.NewUploadImageRequest(t, ctx, testServerURL, fileContent, "image", "test-image.jpg", additionalFields, accessToken)

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
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		additionalFields := map[string]string{
			"parent_id":   uuid.New().String(),
			"parent_type": string(domain.CategoryType),
			"title":       "My Test Image",
		}

		req := th.NewUploadImageRequest(t, ctx, testServerURL, []byte{}, "image", "", additionalFields, accessToken) // No file content

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail if parent ID does not exist", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		nonExistentParentID := uuid.New() // A valid UUID that does not exist in the DB

		additionalFields := map[string]string{
			"parent_id":   nonExistentParentID.String(),
			"parent_type": string(domain.CategoryType),
			"title":       "My Test Image",
			"is_active":   "true",
		}

		req := th.NewUploadImageRequest(t, ctx, testServerURL, fileContent, "image", "test-image.jpg", additionalFields, accessToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Contains(t, respBody.Error, errs.ErrRepositoryNotFound.Error())
	})

	t.Run("should succeed if a second active image is created for the same parent", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		parentID := uuid.New()
		createCategoryParent(t, ctx, parentID)

		// Upload the first active image successfully
		additionalFields1 := map[string]string{
			"parent_id":   parentID.String(),
			"parent_type": string(domain.CategoryType),
			"title":       "First Active Image",
			"is_active":   "true",
		}

		req1 := th.NewUploadImageRequest(t, ctx, testServerURL, fileContent, "image", "image1.jpg", additionalFields1, accessToken)
		resp1, _ := client.Do(req1)
		resp1.Body.Close()
		assert.Equal(t, http.StatusCreated, resp1.StatusCode)

		title2 := "Second Active Image"
		additionalFields2 := map[string]string{
			"parent_id":   parentID.String(),
			"parent_type": string(domain.CategoryType),
			"title":       title2,
			"is_active":   "true",
		}

		req2 := th.NewUploadImageRequest(t, ctx, testServerURL, fileContent, "image", "image2.jpg", additionalFields2, accessToken)
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

	t.Run("should return 401 when access token is missing", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)

		additionalFields := map[string]string{
			"parent_id":   uuid.New().String(),
			"parent_type": string(domain.CategoryType),
			"title":       "My Test Image",
		}

		req := th.NewUploadImageRequest(t, ctx, testServerURL, fileContent, "image", "test-image.jpg", additionalFields, "")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when session is expired", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create expired admin session
		accessToken := tu.SetupExpiredUserWithRole(t, ctx, identity.Administrator, authCtx)

		additionalFields := map[string]string{
			"parent_id":   uuid.New().String(),
			"parent_type": string(domain.CategoryType),
			"title":       "My Test Image",
		}

		req := th.NewUploadImageRequest(t, ctx, testServerURL, fileContent, "image", "test-image.jpg", additionalFields, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 403 when user has insufficient role", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create standard user (not admin)
		accessToken := tu.SetupStandardUser(t, ctx, authCtx)

		additionalFields := map[string]string{
			"parent_id":   uuid.New().String(),
			"parent_type": string(domain.CategoryType),
			"title":       "My Test Image",
		}

		req := th.NewUploadImageRequest(t, ctx, testServerURL, fileContent, "image", "test-image.jpg", additionalFields, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		clearDatabaseAndS3(t, ctx)

		additionalFields := map[string]string{
			"parent_id":   uuid.New().String(),
			"parent_type": string(domain.CategoryType),
			"title":       "My Test Image",
		}

		req := th.NewUploadImageRequest(t, ctx, testServerURL, fileContent, "image", "test-image.jpg", additionalFields, "invalid-token-12345")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
