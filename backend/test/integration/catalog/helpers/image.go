package helpers

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func ClearImagesTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, "TRUNCATE TABLE catalog.images RESTART IDENTITY CASCADE;")
	require.NoError(t, err)
}

func NewValidImageRequest() *domain.CreateImageRequest {
	image := &domain.CreateImageRequest{
		ParentID:   uuid.New().String(),
		ParentType: domain.CategoryType,
		Title:      "title",
		IsActive:   BoolPtr(false),
	}
	return image
}

func NewValidImage(parentID uuid.UUID) *domain.Image {
	image := &domain.Image{
		ID:          uuid.New(),
		ParentID:    parentID,
		ParentType:  domain.ParentType("category"),
		Title:       "title",
		S3Key:       uuid.New().String(), // will be handled by the user later
		Size:        12000,
		ContentType: "",
		IsActive:    false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return image
}

func InsertImage(t *testing.T, ctx context.Context, img *domain.Image, pool *pgxpool.Pool) {
	t.Helper()
	query := `
		INSERT INTO catalog.images (
			id, parent_id, parent_type, title, s3_key, size, content_type, is_active, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := pool.Exec(
		ctx,
		query,
		img.ID,
		img.ParentID,
		img.ParentType,
		img.Title,
		img.S3Key,
		img.Size,
		img.ContentType,
		img.IsActive,
		img.CreatedAt,
	)
	require.NoError(t, err, "Failed to insert mock image for test setup")
}

// getImageStatus retrieves the is_active status of an image.
func GetImageStatus(t *testing.T, ctx context.Context, id uuid.UUID, db *pgxpool.Pool) bool {
	t.Helper()
	var isActive bool
	err := db.QueryRow(ctx, "SELECT is_active FROM catalog.images WHERE id = $1", id).Scan(&isActive)
	require.NoError(t, err, "Failed to get image status for ID %s", id)
	return isActive
}
