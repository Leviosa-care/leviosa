package category_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetAllPublishedCategories TEST_PATH=test/integration/catalog/category/get_all_published_categories_test.go

func TestGetAllPublishedCategories(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get only published categories", func(t *testing.T) {
		// Clean the database to ensure isolation for this test
		td.ClearCategoriesTable(t, ctx, testPool)

		// Setup: Insert a mix of published and draft categories
		publishedCat1 := td.NewValidCategory("Published Apple Products")
		publishedCat1.Status = domain.Published
		td.InsertCategory(t, ctx, publishedCat1, testPool)

		publishedCat2 := td.NewValidCategory("Published Samsung Products")
		publishedCat2.Status = domain.Published
		td.InsertCategory(t, ctx, publishedCat2, testPool)

		// These should NOT be returned
		draftCat := td.NewValidCategory("Draft Category")
		draftCat.Status = domain.Draft
		td.InsertCategory(t, ctx, draftCat, testPool)

		// Setup: Insert active images for publishedCat1
		img1 := td.NewValidImage(publishedCat1.ID)
		img1.Title = "Apple Image"
		img1.IsActive = true
		td.InsertImage(t, ctx, img1, testPool)

		// Setup: Insert an inactive image for publishedCat2 to ensure it's not returned
		img2 := td.NewValidImage(publishedCat2.ID)
		img2.Title = "Inactive Samsung Image"
		img2.IsActive = false
		td.InsertImage(t, ctx, img2, testPool)

		// Setup: Insert an image for the draft category to ensure it's not returned
		imgDraft := td.NewValidImage(draftCat.ID)
		imgDraft.Title = "Draft Image"
		imgDraft.IsActive = true
		td.InsertImage(t, ctx, imgDraft, testPool)
		req := newGetAllPublishedCategoriesRequest(t, ctx)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []*domain.CategoryWithImage
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		// Post-conditions: Verify the response content
		assert.Len(t, response, 2, "Expected 2 published categories in the response")

		// The categories should be ordered by name ASC
		assert.Equal(t, publishedCat1.Name, response[0].Category.Name)
		assert.NotNil(t, response[0].Image, "Expected an image for Published Apple Products")
		assert.Equal(t, img1.ID, response[0].Image.ID)

		assert.Equal(t, publishedCat2.Name, response[1].Category.Name)
		assert.Nil(t, response[1].Image, "Expected no active image for Published Samsung Products")

	})

	t.Run("should return an empty array when no published categories exist", func(t *testing.T) {
		// Clean the database to ensure a fresh start
		td.ClearCategoriesTable(t, ctx, testPool)

		// Setup: Insert a few categories, but none of them are published
		draftCat1 := td.NewValidCategory("Draft A")
		draftCat1.Status = domain.Draft
		td.InsertCategory(t, ctx, draftCat1, testPool)

		draftCat2 := td.NewValidCategory("Draft B")
		draftCat2.Status = domain.Draft
		td.InsertCategory(t, ctx, draftCat2, testPool)

		req := newGetAllPublishedCategoriesRequest(t, ctx)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var categories []*domain.CategoryWithImage
		err = json.NewDecoder(resp.Body).Decode(&categories)
		assert.NoError(t, err)

		// Post-conditions: Verify the response is an empty array
		assert.NotNil(t, categories, "Expected an empty but non-nil slice")
		assert.Len(t, categories, 0, "Expected 0 categories in the response")
	})

	t.Run("should return an empty array when the database is empty", func(t *testing.T) {
		// Clean the database to ensure it is completely empty
		td.ClearCategoriesTable(t, ctx, testPool)

		req := newGetAllPublishedCategoriesRequest(t, ctx)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var categories []*domain.CategoryWithImage
		err = json.NewDecoder(resp.Body).Decode(&categories)
		assert.NoError(t, err)

		// Post-conditions: Verify the response is an empty array
		assert.NotNil(t, categories, "Expected an empty but non-nil slice")
		assert.Len(t, categories, 0, "Expected 0 categories in the response")
	})
}

// newGetAllPublishedCategoriesRequest creates a new HTTP request for the public handler.
func newGetAllPublishedCategoriesRequest(t *testing.T, ctx context.Context) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+"/categories", nil)
	require.NoError(t, err)
	return req
}
