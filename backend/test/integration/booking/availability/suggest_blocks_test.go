package availability_test

// import (
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"testing"
// 	"time"
//
// 	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
// 	allocationHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/allocation"
// 	availabilityHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/availability"
// 	buildingHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/building"
// 	roomHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/booking/room"
// 	catalogHelpers "github.com/Leviosa-care/leviosa/backend/test/helpers/catalog"
//
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )
//
// func TestSuggestBlocks(t *testing.T) {
// 	ctx := context.Background()
// 	client := &http.Client{Timeout: 10 * time.Second}
//
// 	t.Run("should generate suggestions for dedicated room", func(t *testing.T) {
// 		// Clean state
// 		allocationHelpers.ClearAllocationTables(t, ctx, testPool)
// 		availabilityHelpers.ClearAvailabilityTable(t, ctx, testPool)
// 		roomHelpers.ClearRoomTables(t, ctx, testPool)
// 		buildingHelpers.ClearBuildingTable(t, ctx, testPool)
// 		catalogHelpers.ClearProductsTable(t, ctx, testPool)
//
// 		// Create test data
// 		partnerUser := authCtx.CreateTestUser(t, ctx, "partner", "partner@suggest.com")
// 		sessionToken := authCtx.CreateTestSession(t, ctx, partnerUser)
//
// 		building := buildingHelpers.CreateTestBuilding(t, ctx, testPool, crypto)
// 		room := roomHelpers.CreateTestRoom(t, ctx, testPool, crypto, building.ID)
//
// 		// Create dedicated room allocation
// 		startDate := time.Now()
// 		endDate := startDate.Add(365 * 24 * time.Hour) // 1 year
// 		allocation := allocationHelpers.CreateTestAllocation(t, ctx, testPool, crypto, partnerUser.ID, room.ID, domain.AllocationTypeDedicated, startDate, endDate)
// 		allocationHelpers.InsertAllocationEncx(t, ctx, allocation, testPool)
//
// 		// Create products
// 		categoryID := catalogHelpers.CreateTestCategory(t, ctx, testPool)
// 		products := catalogHelpers.CreateDefaultTestProducts(t, ctx, testPool, categoryID)
//
// 		// Make request
// 		req := availabilityHelpers.NewSuggestBlocksRequest(t, ctx, testServerURL, partnerUser.ID.String(), room.ID.String(), sessionToken)
// 		resp, err := client.Do(req)
// 		require.NoError(t, err)
// 		defer resp.Body.Close()
//
// 		// Assert response
// 		assert.Equal(t, http.StatusOK, resp.StatusCode)
//
// 		var response domain.GetAvailabilitySuggestionsResponse
// 		err = json.NewDecoder(resp.Body).Decode(&response)
// 		require.NoError(t, err)
//
// 		// Verify response structure
// 		assert.Equal(t, partnerUser.ID, response.PartnerID)
// 		assert.Equal(t, room.ID, response.RoomID)
// 		assert.Equal(t, string(domain.AllocationTypeDedicated), response.AllocationType)
// 		assert.NotEmpty(t, response.RecommendedBlocks, "Should have block suggestions")
//
// 		// Verify suggestions based on products
// 		// Products: 60min (75min with buffer), 90min (105min with buffer), 30min (40min with buffer)
// 		// Expected blocks: 40, 75, 80 (2x40), 105, 120 (3x40), 150 (2x75), etc.
//
// 		// Verify priority ranking
// 		priorities := make(map[int]int) // priority -> count
// 		for _, block := range response.RecommendedBlocks {
// 			priorities[block.Priority]++
// 			assert.NotEmpty(t, block.Rationale, "Should have rationale")
// 			assert.NotEmpty(t, block.ProductCombinations, "Should have product combinations")
// 		}
//
// 		// Should have Priority 1 suggestions (standard durations: 60, 90, 120, 180)
// 		assert.Greater(t, priorities[1], 0, "Should have priority 1 suggestions")
// 	})
//
// 	t.Run("should generate suggestions for shared room", func(t *testing.T) {
// 		// Clean state
// 		allocationHelpers.ClearAllocationTables(t, ctx, testPool)
// 		roomHelpers.ClearRoomTables(t, ctx, testPool)
// 		buildingHelpers.ClearBuildingTable(t, ctx, testPool)
//
// 		// Create test data
// 		partnerUser := authCtx.CreateTestUser(t, ctx, "partner", "partner2@suggest.com")
// 		sessionToken := authCtx.CreateTestSession(t, ctx, partnerUser)
//
// 		building := buildingHelpers.CreateTestBuilding(t, ctx, testPool, crypto)
// 		room := roomHelpers.CreateTestRoom(t, ctx, testPool, crypto, building.ID)
//
// 		// Create shared room allocation
// 		startDate := time.Now()
// 		endDate := startDate.Add(365 * 24 * time.Hour)
// 		allocation := allocationHelpers.CreateTestAllocation(t, ctx, testPool, crypto, partnerUser.ID, room.ID, domain.AllocationTypeShared, startDate, endDate)
// 		allocationHelpers.InsertAllocationEncx(t, ctx, allocation, testPool)
//
// 		// Make request
// 		req := availabilityHelpers.NewSuggestBlocksRequest(t, ctx, testServerURL, partnerUser.ID.String(), room.ID.String(), sessionToken)
// 		resp, err := client.Do(req)
// 		require.NoError(t, err)
// 		defer resp.Body.Close()
//
// 		// Assert response
// 		assert.Equal(t, http.StatusOK, resp.StatusCode)
//
// 		var response domain.GetAvailabilitySuggestionsResponse
// 		err = json.NewDecoder(resp.Body).Decode(&response)
// 		require.NoError(t, err)
//
// 		// Verify shared allocation
// 		assert.Equal(t, string(domain.AllocationTypeShared), response.AllocationType)
// 		assert.NotEmpty(t, response.RecommendedBlocks)
//
// 		// Shared rooms should include standard durations (60, 90, 120, 180, 240)
// 		standardDurations := map[int]bool{60: false, 90: false, 120: false, 180: false, 240: false}
// 		for _, block := range response.RecommendedBlocks {
// 			if _, exists := standardDurations[block.DurationMinutes]; exists {
// 				standardDurations[block.DurationMinutes] = true
// 			}
// 		}
//
// 		// At least some standard durations should be present
// 		foundStandard := false
// 		for _, found := range standardDurations {
// 			if found {
// 				foundStandard = true
// 				break
// 			}
// 		}
// 		assert.True(t, foundStandard, "Shared room should suggest some standard durations")
// 	})
//
// 	t.Run("should verify priority ranking order", func(t *testing.T) {
// 		// Clean state
// 		allocationHelpers.ClearAllocationTables(t, ctx, testPool)
// 		roomHelpers.ClearRoomTables(t, ctx, testPool)
// 		buildingHelpers.ClearBuildingTable(t, ctx, testPool)
//
// 		// Create test data
// 		partnerUser := authCtx.CreateTestUser(t, ctx, "partner", "partner3@suggest.com")
// 		sessionToken := authCtx.CreateTestSession(t, ctx, partnerUser)
//
// 		building := buildingHelpers.CreateTestBuilding(t, ctx, testPool, crypto)
// 		room := roomHelpers.CreateTestRoom(t, ctx, testPool, crypto, building.ID)
//
// 		startDate := time.Now()
// 		endDate := startDate.Add(365 * 24 * time.Hour)
// 		allocation := allocationHelpers.CreateTestAllocation(t, ctx, testPool, crypto, partnerUser.ID, room.ID, domain.AllocationTypeDedicated, startDate, endDate)
// 		allocationHelpers.InsertAllocationEncx(t, ctx, allocation, testPool)
//
// 		// Make request
// 		req := availabilityHelpers.NewSuggestBlocksRequest(t, ctx, testServerURL, partnerUser.ID.String(), room.ID.String(), sessionToken)
// 		resp, err := client.Do(req)
// 		require.NoError(t, err)
// 		defer resp.Body.Close()
//
// 		assert.Equal(t, http.StatusOK, resp.StatusCode)
//
// 		var response domain.GetAvailabilitySuggestionsResponse
// 		err = json.NewDecoder(resp.Body).Decode(&response)
// 		require.NoError(t, err)
//
// 		// Verify suggestions are sorted by priority (ascending), then duration (ascending)
// 		prevPriority := 0
// 		prevDuration := 0
// 		for _, block := range response.RecommendedBlocks {
// 			if block.Priority > prevPriority {
// 				prevPriority = block.Priority
// 				prevDuration = 0 // Reset duration when priority changes
// 			} else if block.Priority == prevPriority {
// 				assert.GreaterOrEqual(t, block.DurationMinutes, prevDuration,
// 					"Within same priority, duration should be ascending")
// 				prevDuration = block.DurationMinutes
// 			}
// 		}
// 	})
//
// 	t.Run("should return 400 for invalid partner ID", func(t *testing.T) {
// 		partnerUser := authCtx.CreateTestUser(t, ctx, "partner", "partner4@suggest.com")
// 		sessionToken := authCtx.CreateTestSession(t, ctx, partnerUser)
//
// 		roomID := uuid.New()
// 		req := availabilityHelpers.NewSuggestBlocksRequest(t, ctx, testServerURL, "invalid-uuid", roomID.String(), sessionToken)
// 		resp, err := client.Do(req)
// 		require.NoError(t, err)
// 		defer resp.Body.Close()
//
// 		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
// 	})
//
// 	t.Run("should return 400 for invalid room ID", func(t *testing.T) {
// 		partnerUser := authCtx.CreateTestUser(t, ctx, "partner", "partner5@suggest.com")
// 		sessionToken := authCtx.CreateTestSession(t, ctx, partnerUser)
//
// 		req := availabilityHelpers.NewSuggestBlocksRequest(t, ctx, testServerURL, partnerUser.ID.String(), "invalid-uuid", sessionToken)
// 		resp, err := client.Do(req)
// 		require.NoError(t, err)
// 		defer resp.Body.Close()
//
// 		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
// 	})
// }
