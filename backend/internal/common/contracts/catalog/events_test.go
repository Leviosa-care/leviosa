package catalog

import ( "encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoryEventSerialization(t *testing.T) {
	categoryID := uuid.New().String()
	now := time.Now().UTC()

	t.Run("CategoryCreatedEvent", func(t *testing.T) {
		// Create event
		event := CategoryCreatedEvent{
			ID:          categoryID,
			Name:        "Massage",
			Description: "Therapeutic massage services",
			Status:      "published",
			Metadata:    map[string]any{"color": "#FF5733"},
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		// Test validation
		require.NoError(t, event.Validate())

		// Serialize to JSON
		data, err := json.Marshal(event)
		require.NoError(t, err)

		// Deserialize from JSON
		var deserialized CategoryCreatedEvent
		err = json.Unmarshal(data, &deserialized)
		require.NoError(t, err)

		// Verify equality
		assert.Equal(t, event.ID, deserialized.ID)
		assert.Equal(t, event.Name, deserialized.Name)
		assert.Equal(t, event.Description, deserialized.Description)
		assert.Equal(t, event.Status, deserialized.Status)
		assert.Equal(t, event.Metadata, deserialized.Metadata)
		assert.True(t, event.CreatedAt.Equal(deserialized.CreatedAt))
		assert.True(t, event.UpdatedAt.Equal(deserialized.UpdatedAt))
	})

	t.Run("CategoryUpdatedEvent", func(t *testing.T) {
		event := CategoryUpdatedEvent{
			ID:          categoryID,
			Name:        "Updated Massage",
			Description: "Updated description",
			Status:      "draft",
			Metadata:    map[string]any{"color": "#123456"},
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		// Test validation
		require.NoError(t, event.Validate())

		// Test JSON roundtrip
		data, err := json.Marshal(event)
		require.NoError(t, err)

		var deserialized CategoryUpdatedEvent
		err = json.Unmarshal(data, &deserialized)
		require.NoError(t, err)

		assert.Equal(t, event.ID, deserialized.ID)
		assert.Equal(t, event.Name, deserialized.Name)
	})

	t.Run("CategoryDeletedEvent", func(t *testing.T) {
		event := CategoryDeletedEvent{
			ID: categoryID,
		}

		// Test validation
		require.NoError(t, event.Validate())

		// Test JSON roundtrip
		data, err := json.Marshal(event)
		require.NoError(t, err)

		var deserialized CategoryDeletedEvent
		err = json.Unmarshal(data, &deserialized)
		require.NoError(t, err)

		assert.Equal(t, event.ID, deserialized.ID)
	})
}

func TestProductEventSerialization(t *testing.T) {
	categoryID := uuid.New().String()
	productID := uuid.New().String()
	now := time.Now().UTC()

	t.Run("ProductCreatedEvent", func(t *testing.T) {
		event := ProductCreatedEvent{
			ID:                productID,
			Name:              "Swedish Massage - 60min",
			Description:       "Relaxing full body massage",
			CategoryID:        categoryID,
			Duration:          60,
			Status:            "published",
			Availability:      "in-person",
			BufferTime:        15,
			CancellationHours: 24,
			StripeProductID:   "prod_ABC123",
			Metadata:          map[string]any{"key": "value"},
			CreatedAt:         now,
			UpdatedAt:         now,
		}

		// Test validation
		require.NoError(t, event.Validate())

		// Serialize to JSON
		data, err := json.Marshal(event)
		require.NoError(t, err)

		// Deserialize from JSON
		var deserialized ProductCreatedEvent
		err = json.Unmarshal(data, &deserialized)
		require.NoError(t, err)

		// Verify equality
		assert.Equal(t, event.ID, deserialized.ID)
		assert.Equal(t, event.Name, deserialized.Name)
		assert.Equal(t, event.Description, deserialized.Description)
		assert.Equal(t, event.CategoryID, deserialized.CategoryID)
		assert.Equal(t, event.Duration, deserialized.Duration)
		assert.Equal(t, event.Status, deserialized.Status)
		assert.Equal(t, event.Availability, deserialized.Availability)
		assert.Equal(t, event.BufferTime, deserialized.BufferTime)
		assert.Equal(t, event.CancellationHours, deserialized.CancellationHours)
		assert.Equal(t, event.StripeProductID, deserialized.StripeProductID)
		assert.Equal(t, event.Metadata, deserialized.Metadata)
		assert.True(t, event.CreatedAt.Equal(deserialized.CreatedAt))
		assert.True(t, event.UpdatedAt.Equal(deserialized.UpdatedAt))
	})

	t.Run("ProductUpdatedEvent", func(t *testing.T) {
		event := ProductUpdatedEvent{
			ID:                productID,
			Name:              "Updated Product",
			Description:       "Updated description",
			CategoryID:        categoryID,
			Duration:          90,
			Status:            "draft",
			Availability:      "online",
			BufferTime:        30,
			CancellationHours: 48,
			StripeProductID:   "prod_UPDATED",
			Metadata:          map[string]any{"updated": true},
			CreatedAt:         now,
			UpdatedAt:         now,
		}

		// Test validation
		require.NoError(t, event.Validate())

		// Test JSON roundtrip
		data, err := json.Marshal(event)
		require.NoError(t, err)

		var deserialized ProductUpdatedEvent
		err = json.Unmarshal(data, &deserialized)
		require.NoError(t, err)

		assert.Equal(t, event.ID, deserialized.ID)
		assert.Equal(t, event.Name, deserialized.Name)
		assert.Equal(t, event.Duration, deserialized.Duration)
	})

	t.Run("ProductDeletedEvent", func(t *testing.T) {
		event := ProductDeletedEvent{
			ID: productID,
		}

		// Test validation
		require.NoError(t, event.Validate())

		// Test JSON roundtrip
		data, err := json.Marshal(event)
		require.NoError(t, err)

		var deserialized ProductDeletedEvent
		err = json.Unmarshal(data, &deserialized)
		require.NoError(t, err)

		assert.Equal(t, event.ID, deserialized.ID)
	})
}

func TestEventValidation(t *testing.T) {
	t.Run("invalid UUID", func(t *testing.T) {
		event := CategoryCreatedEvent{
			ID:        "invalid-uuid",
			Name:      "Test",
			Status:    "published",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := event.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid category ID")
	})

	t.Run("empty name", func(t *testing.T) {
		event := CategoryCreatedEvent{
			ID:        uuid.New().String(),
			Name:      "",
			Status:    "published",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := event.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "category name cannot be empty")
	})

	t.Run("negative duration", func(t *testing.T) {
		event := ProductCreatedEvent{
			ID:                uuid.New().String(),
			Name:              "Test Product",
			CategoryID:        uuid.New().String(),
			Duration:          -10,
			Status:            "published",
			Availability:      "online",
			BufferTime:        15,
			CancellationHours: 24,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		err := event.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product duration must be positive")
	})
}
