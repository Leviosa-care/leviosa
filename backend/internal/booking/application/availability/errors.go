package availability

import (
	"fmt"
	"strings"
)

// InvalidDurationError represents an availability duration that doesn't align with product offerings
type InvalidDurationError struct {
	RequestedDuration int
	ValidBlocks       []ValidBlock
}

// Error returns a user-friendly error message with suggestions
func (e *InvalidDurationError) Error() string {
	if len(e.ValidBlocks) == 0 {
		return fmt.Sprintf(
			"Availability duration of %d minutes does not align with your product offerings. "+
				"No valid durations available - please configure your products first.",
			e.RequestedDuration,
		)
	}

	suggestions := make([]string, 0, len(e.ValidBlocks))
	for _, block := range e.ValidBlocks {
		suggestions = append(suggestions,
			fmt.Sprintf("%d min (%s)", block.DurationMinutes, block.Rationale))
	}

	return fmt.Sprintf(
		"Availability duration of %d minutes does not align with your product offerings. "+
			"Suggested durations: %s",
		e.RequestedDuration,
		strings.Join(suggestions, ", "),
	)
}

// ToJSON converts the error to a structured JSON-friendly format
func (e *InvalidDurationError) ToJSON() map[string]interface{} {
	suggestions := make([]map[string]interface{}, 0, len(e.ValidBlocks))
	for _, block := range e.ValidBlocks {
		suggestions = append(suggestions, map[string]interface{}{
			"minutes":   block.DurationMinutes,
			"rationale": block.Rationale,
		})
	}

	return map[string]interface{}{
		"error":               "Availability duration does not align with product offerings",
		"requested_duration":  e.RequestedDuration,
		"suggested_durations": suggestions,
	}
}
