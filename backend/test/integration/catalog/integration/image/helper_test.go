package image_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"path/filepath"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func clearDatabaseAndS3(t *testing.T, ctx context.Context) {
	t.Helper()
	td.ClearBucket(t, ctx, s3Client)
	td.ClearImagesTable(t, ctx, testPool)
	td.ClearCategoriesTable(t, ctx, testPool)
	td.ClearProductsTable(t, ctx, testPool)
}

// createCategoryParent helper function to insert a parent category into the test database.
func createCategoryParent(t *testing.T, ctx context.Context, parentID uuid.UUID) {
	t.Helper()
	_, err := testPool.Exec(ctx, `INSERT INTO catalog.categories (id, name, description, status) VALUES ($1, $2, $3, $4)`,
		parentID, "Test Category", "A test category.", "published")
	require.NoError(t, err, "Failed to create parent category in database")
}

// createCategoryParent helper function to insert a parent category into the test database.
func createCategoryParentWithName(t *testing.T, ctx context.Context, parentID uuid.UUID, name string) {
	t.Helper()
	_, err := testPool.Exec(ctx, `INSERT INTO catalog.categories (id, name, description, status) VALUES ($1, $2, $3, $4)`,
		parentID, name, "A test category.", "published")
	require.NoError(t, err, "Failed to create parent category in database")
}

// createProductParent helper function to insert a parent product into the test database.
// This function also creates a category because of the foreign key constraint.
func createProductParent(t *testing.T, ctx context.Context, productID uuid.UUID) {
	t.Helper()
	categoryID := uuid.New()
	createCategoryParent(t, ctx, categoryID)
	_, err := testPool.Exec(ctx, `
		INSERT INTO catalog.products (
			id, name, description, category_id, duration, availability, buffer_time, cancellation_hours, stripe_product_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		productID, "Test Product", "A test product.", categoryID, 60, "online", 15, 24, "stripe_prod_12345")
	require.NoError(t, err, "Failed to create parent product in database")
}

// createMultipartForm creates a multipart form with the given fields and file.
// It now explicitly takes fileContentType.
func createMultipartForm(t *testing.T, fields map[string]string, filename string, fileContent string, fileContentType string) (*bytes.Buffer, string) {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	for key, value := range fields {
		err := writer.WriteField(key, value)
		require.NoError(t, err)
	}

	// Add file part
	if filename != "" && fileContent != "" {
		// Use CreatePart to manually set Content-Type for the file
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "image", filepath.Base(filename)))
		h.Set("Content-Type", fileContentType) // Explicitly set the content type

		// Corrected: Only use CreatePart and write to its returned io.Writer
		part, err := writer.CreatePart(h)
		require.NoError(t, err)
		_, err = io.WriteString(part, fileContent)
		require.NoError(t, err)
	}

	writer.Close()
	return body, writer.FormDataContentType()
}

// Helper to insert an image directly into the database for test setup
func insertImageIntoDB(t *testing.T, ctx context.Context, img *domain.Image) {
	t.Helper()
	_, err := testPool.Exec(ctx, `
		INSERT INTO catalog.images (id, parent_id, parent_type, title, s3_key, size, content_type, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		img.ID, img.ParentID, img.ParentType, img.Title, img.S3Key, img.Size, img.ContentType, img.IsActive, img.CreatedAt, img.UpdatedAt)
	require.NoError(t, err, "Failed to insert image into DB for test setup")
}

// Helper to check if an image exists in the database
func imageExistsInDB(t *testing.T, ctx context.Context, imageID uuid.UUID) bool {
	t.Helper()
	var count int
	err := testPool.QueryRow(ctx, "SELECT COUNT(*) FROM catalog.images WHERE id = $1", imageID).Scan(&count)
	require.NoError(t, err)
	return count > 0
}

// Helper to check if a file exists in S3
func fileExistsInS3(t *testing.T, ctx context.Context, key string) bool {
	t.Helper()
	_, err := s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(td.BUCKETNAME),
		Key:    aws.String(key),
	})
	if err != nil {
		var noSuchKey *types.NotFound
		if assert.ErrorAs(t, err, &noSuchKey) {
			return false // Object not found
		}
		require.Fail(t, fmt.Sprintf("Failed to check object existence for key %s: %v", key, err))
	}
	return true // Object found
}
