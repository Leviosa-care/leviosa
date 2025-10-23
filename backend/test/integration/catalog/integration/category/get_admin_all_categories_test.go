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

// newGetAllCategoriesRequest creates a new HTTP request for the GetAdminAllCategories handler.
func newGetAllCategoriesRequest(t *testing.T, ctx context.Context) *http.Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServerURL+"/admin/categories", nil)
	require.NoError(t, err)
	return req
}

func TestGetAdminAllCategories(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get all categories when they exist", func(t *testing.T) {
		// Clean the database to ensure isolation for this test
		td.ClearCategoriesTable(t, ctx, testPool)

		// Setup: Insert multiple categories
		cat1 := td.NewValidCategory("Apple Products")
		td.InsertCategory(t, ctx, cat1, testPool)

		cat2 := td.NewValidCategory("Samsung Products")
		td.InsertCategory(t, ctx, cat2, testPool)

		cat3 := td.NewValidCategory("Category Z")
		td.InsertCategory(t, ctx, cat3, testPool)

		// Setup: Insert active images
		img1 := td.NewValidImage(cat1.ID)
		img1.Title = "Apple Image"
		img1.IsActive = true
		td.InsertImage(t, ctx, img1, testPool)

		img2 := td.NewValidImage(cat2.ID)
		img2.Title = "Samsung Image"
		img2.IsActive = true
		td.InsertImage(t, ctx, img2, testPool)

		// Setup: Insert an inactive image for category 3 to ensure it's not returned
		img3 := td.NewValidImage(cat3.ID)
		img3.Title = "Inactive Z Image"
		img3.IsActive = false
		td.InsertImage(t, ctx, img3, testPool)

		req := newGetAllCategoriesRequest(t, ctx)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var response []*domain.CategoryWithImage
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		// Post-conditions: Verify the response content
		require.Len(t, response, 3, "Expected 3 categories in the response")

		// The categories should be ordered by name ASC
		assert.Equal(t, cat1.Name, response[0].Category.Name)
		assert.NotNil(t, response[0].Image, "Expected an image for Apple Products")
		assert.Equal(t, img1.ID, response[0].Image.ID)

		assert.Equal(t, cat3.Name, response[1].Category.Name)
		assert.Nil(t, response[1].Image, "Expected no active image for Category Z")

		assert.Equal(t, cat2.Name, response[2].Category.Name)
		assert.NotNil(t, response[2].Image, "Expected an image for Samsung Products")
		assert.Equal(t, img2.ID, response[2].Image.ID)

		// Verify the full data of one of the returned categories
		assert.Equal(t, cat1.ID, response[0].Category.ID)
		assert.Equal(t, cat1.Description, response[0].Category.Description)
	})

	t.Run("should return an empty array when no categories exist", func(t *testing.T) {
		// Clean the database to ensure it is empty
		td.ClearCategoriesTable(t, ctx, testPool)

		req := newGetAllCategoriesRequest(t, ctx)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []*domain.CategoryWithImage
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		// Post-conditions: Verify the response is an empty array
		assert.NotNil(t, response, "Expected an empty but non-nil slice")
		assert.Len(t, response, 0, "Expected 0 categories in the response")
	})
}

// NOTE: this is the old way
// func TestGetAdminAllCategories(t *testing.T) {
// 	ctx := context.Background()
// 	client := &http.Client{Timeout: 10 * time.Second}
//
// 	t.Run("should successfully get all categories when they exist", func(t *testing.T) {
// 		// Clean the database to ensure isolation for this test
// 		td.ClearCategoriesTable(t, ctx, testPool)
//
// 		// Setup: Insert multiple categories
// 		cat1 := td.NewValidCategory("Apple Products")
// 		td.InsertCategory(t, ctx, cat1, testPool)
//
// 		cat2 := td.NewValidCategory("Samsung Products")
// 		td.InsertCategory(t, ctx, cat2, testPool)
//
// 		cat3 := td.NewValidCategory("Category Z")
// 		td.InsertCategory(t, ctx, cat3, testPool)
//
// 		img1 := td.NewValidImage(cat1.ID)
// 		img1.Title = "Apple Image"
// 		img1.IsActive = true
// 		td.InsertImage(t, ctx, img1, testPool)
//
// 		img2 := td.NewValidImage(cat2.ID)
// 		img2.Title = "Samsung Image"
// 		img2.IsActive = true
// 		td.InsertImage(t, ctx, img2, testPool)
//
// 		// Setup: Insert an inactive image for category 3 to ensure it's not returned
// 		img3 := td.NewValidImage(cat3.ID)
// 		img3.Title = "Inactive Z Image"
// 		img3.IsActive = false
// 		td.InsertImage(t, ctx, img3, testPool)
//
// 		req := newGetAllCategoriesRequest(t, ctx)
// 		resp, err := client.Do(req)
// 		assert.NoError(t, err)
// 		defer resp.Body.Close()
//
// 		assert.Equal(t, http.StatusOK, resp.StatusCode)
//
// 		var response []*domain.CategoryWithImage
// 		err = json.NewDecoder(resp.Body).Decode(&response)
// 		assert.NoError(t, err)
//
// 		// Post-conditions: Verify the response content
// 		require.Len(t, response, 3, "Expected 3 categories in the response")
//
// 		// The categories should be ordered by name ASC
// 		assert.Equal(t, cat1.Name, response[0].Category.Name)
// 		assert.NotNil(t, response[0].Image, "Expected an image for Apple Products")
// 		assert.Equal(t, img1.ID, response[0].Image.ID)
//
// 		assert.Equal(t, cat3.Name, response[1].Category.Name)
// 		assert.Nil(t, response[1].Image, "Expected no active image for Category Z")
//
// 		assert.Equal(t, cat2.Name, response[2].Category.Name)
// 		assert.NotNil(t, response[2].Image, "Expected an image for Samsung Products")
// 		assert.Equal(t, img2.ID, response[2].Image.ID)
//
// 		// Verify the full data of one of the returned categories
// 		assert.Equal(t, cat1.ID, response[0].Category.ID)
// 		assert.Equal(t, cat1.Description, response[0].Category.Description)
// 	})
//
// 	t.Run("should return an empty array when no categories exist", func(t *testing.T) {
// 		// Clean the database to ensure it is empty
// 		td.ClearCategoriesTable(t, ctx, testPool)
//
// 		req := newGetAllCategoriesRequest(t, ctx)
// 		resp, err := client.Do(req)
// 		assert.NoError(t, err)
// 		defer resp.Body.Close()
//
// 		assert.Equal(t, http.StatusOK, resp.StatusCode)
//
// 		var response []*domain.CategoryWithImage
// 		err = json.NewDecoder(resp.Body).Decode(&response)
// 		assert.NoError(t, err)
//
// 		// Post-conditions: Verify the response is an empty array
// 		assert.NotNil(t, response, "Expected an empty but non-nil slice")
// 		assert.Len(t, response, 0, "Expected 0 categories in the response")
// 	})
// }
