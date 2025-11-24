package availability

import (
	"fmt"
	"math"
	"sort"

	catalogDomain "github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

// ValidBlock represents a valid availability duration with its rationale
type ValidBlock struct {
	DurationMinutes int
	Rationale       string
}

// calculateValidBlocks computes all valid availability durations based on partner's products
// For shared rooms, availabilities must be exact multiples of (product duration + buffer time)
// to prevent fragmentation.
func calculateValidBlocks(products []*catalogDomain.ProductRes) []ValidBlock {
	if len(products) == 0 {
		return []ValidBlock{}
	}

	blocksMap := make(map[int]string)

	for _, product := range products {
		// Calculate total session time (service duration + buffer/cleanup time)
		sessionTime := product.Duration + product.BufferTime

		if sessionTime <= 0 {
			continue // Skip invalid products
		}

		// Single session
		blocksMap[sessionTime] = fmt.Sprintf("1×%s", product.Name)

		// Multiple sessions (2x, 3x, 4x)
		// We limit to 4x to keep suggestions reasonable
		for multiplier := 2; multiplier <= 4; multiplier++ {
			duration := sessionTime * multiplier
			// Only add if not already present with a simpler rationale
			if _, exists := blocksMap[duration]; !exists {
				blocksMap[duration] = fmt.Sprintf("%d×%s", multiplier, product.Name)
			}
		}
	}

	// Convert map to sorted slice
	blocks := make([]ValidBlock, 0, len(blocksMap))
	for duration, rationale := range blocksMap {
		blocks = append(blocks, ValidBlock{
			DurationMinutes: duration,
			Rationale:       rationale,
		})
	}

	// Sort by duration ascending
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].DurationMinutes < blocks[j].DurationMinutes
	})

	return blocks
}

// isValidDuration checks if the requested duration matches any valid block
// We allow a small tolerance (1 minute) to account for rounding differences
func isValidDuration(requestedMinutes float64, validBlocks []ValidBlock) bool {
	const toleranceMinutes = 1.0

	for _, block := range validBlocks {
		diff := math.Abs(requestedMinutes - float64(block.DurationMinutes))
		if diff <= toleranceMinutes {
			return true
		}
	}

	return false
}

// findClosestValidBlocks returns the 3 closest valid durations to the requested duration
// This is useful for suggesting alternatives when validation fails
func findClosestValidBlocks(requestedMinutes int, validBlocks []ValidBlock) []ValidBlock {
	if len(validBlocks) == 0 {
		return []ValidBlock{}
	}

	// Sort blocks by distance from requested duration
	sorted := make([]ValidBlock, len(validBlocks))
	copy(sorted, validBlocks)

	sort.Slice(sorted, func(i, j int) bool {
		distI := math.Abs(float64(sorted[i].DurationMinutes - requestedMinutes))
		distJ := math.Abs(float64(sorted[j].DurationMinutes - requestedMinutes))
		return distI < distJ
	})

	// Return top 3 closest matches
	maxResults := 3
	if len(sorted) < maxResults {
		maxResults = len(sorted)
	}

	return sorted[:maxResults]
}
