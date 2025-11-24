package availability

import (
	"fmt"
	"sort"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	catalogDomain "github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

// generateBlockSuggestions creates availability block suggestions based on products and room allocation type.
//
// Algorithm:
//   1. For each product, calculate session time (Duration + BufferTime)
//   2. Generate single session blocks (1x product)
//   3. Generate multi-session blocks (2x, 3x, 4x product)
//   4. For shared rooms, also suggest standard durations (60, 90, 120, 180, 240 minutes)
//      if they divide evenly by session time
//
// Multiple products may suggest the same duration (e.g., 60-min product and 120-min product
// both suggest 120-min blocks). These are consolidated into a single suggestion with multiple
// product combinations, allowing practitioners to see all options for a given block duration.
func generateBlockSuggestions(
	products []*catalogDomain.ProductRes,
	allocationType domain.AllocationType,
) []domain.BlockSuggestion {
	suggestions := make(map[int]*domain.BlockSuggestion)

	for _, product := range products {
		sessionTime := product.Duration + product.BufferTime

		// Single session blocks
		addSuggestion(suggestions, sessionTime, 1, product)

		// Multiple session blocks (2x, 3x, 4x)
		for multiplier := 2; multiplier <= 4; multiplier++ {
			duration := sessionTime * multiplier
			addSuggestion(suggestions, duration, multiplier, product)
		}

		// If shared room, also suggest common standard durations
		if allocationType == domain.AllocationTypeShared {
			commonDurations := []int{60, 90, 120, 180, 240}
			for _, commonDuration := range commonDurations {
				if sessionTime < commonDuration {
					sessions := commonDuration / sessionTime
					// Only suggest if it divides evenly
					if sessionTime*sessions == commonDuration {
						addSuggestion(suggestions, commonDuration, sessions, product)
					}
				}
			}
		}
	}

	// Convert map to slice
	result := make([]domain.BlockSuggestion, 0, len(suggestions))
	for _, suggestion := range suggestions {
		result = append(result, *suggestion)
	}

	return result
}

// addSuggestion adds or updates a block suggestion for a given duration.
//
// If a suggestion for this duration already exists, adds a new product combination to it.
// If not, creates a new suggestion with the given product as the first combination.
//
// This allows multiple products to suggest the same duration (e.g., 2x 60-min sessions
// and 1x 120-min session both suggest 120-minute blocks).
func addSuggestion(
	suggestions map[int]*domain.BlockSuggestion,
	duration int,
	sessionCount int,
	product *catalogDomain.ProductRes,
) {
	productInfo := domain.ProductInfo{
		ID:         product.ID,
		Name:       product.Name,
		Duration:   product.Duration,
		BufferTime: product.BufferTime,
	}

	combo := domain.ProductCombo{
		Products: []domain.ProductInfo{productInfo},
		TotalDuration: duration,
		SessionCount:  sessionCount,
	}

	if existing, ok := suggestions[duration]; ok {
		// Add to existing suggestion's combinations
		existing.ProductCombinations = append(existing.ProductCombinations, combo)
	} else {
		// Create new suggestion
		rationale := fmt.Sprintf("%d×%s session(s)", sessionCount, product.Name)
		suggestions[duration] = &domain.BlockSuggestion{
			DurationMinutes:     duration,
			Rationale:           rationale,
			ProductCombinations: []domain.ProductCombo{combo},
		}
	}
}

// rankSuggestions assigns priority levels and sorts suggestions by priority and duration.
//
// Priority Levels:
//   Priority 1 (Highest): Standard durations (60, 90, 120, 180 minutes)
//     - These are industry-standard time blocks that work well for scheduling
//     - Easy for clients to understand and book
//
//   Priority 2 (Good): 30-minute increments (30, 150, 210, etc.)
//     - Flexible for scheduling and calendar display
//     - Compatible with most booking systems
//
//   Priority 3 (Alternative): Non-standard durations (45, 75, etc.)
//     - Valid but may create scheduling complexity
//     - Suggested only when they match product offerings exactly
//
// Sorting: Priority (ascending), then duration (ascending) within each priority level.
func rankSuggestions(suggestions []domain.BlockSuggestion) {
	// Standard durations that work well for scheduling
	standardDurations := map[int]bool{
		60:  true, // 1 hour
		90:  true, // 1.5 hours
		120: true, // 2 hours
		180: true, // 3 hours
	}

	for i := range suggestions {
		if standardDurations[suggestions[i].DurationMinutes] {
			// Highest priority for standard durations
			suggestions[i].Priority = 1
		} else if suggestions[i].DurationMinutes%30 == 0 {
			// Good priority for 30-minute increments
			suggestions[i].Priority = 2
		} else {
			// Lower priority for non-standard durations
			suggestions[i].Priority = 3
		}
	}

	// Sort by priority (ascending), then by duration (ascending)
	sort.Slice(suggestions, func(i, j int) bool {
		if suggestions[i].Priority != suggestions[j].Priority {
			return suggestions[i].Priority < suggestions[j].Priority
		}
		return suggestions[i].DurationMinutes < suggestions[j].DurationMinutes
	})
}
