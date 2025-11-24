package availability

import (
	"fmt"
	"sort"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	catalogDomain "github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

// generateBlockSuggestions creates availability block suggestions based on products
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

// addSuggestion adds or updates a block suggestion
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

// rankSuggestions assigns priority levels and sorts suggestions
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
