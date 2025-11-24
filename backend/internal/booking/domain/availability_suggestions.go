package domain

import (
	"github.com/google/uuid"
)

// AvailabilitySuggestions contains recommended availability block durations for a partner
type AvailabilitySuggestions struct {
	PartnerID         uuid.UUID
	RoomID            uuid.UUID
	AllocationType    AllocationType
	RecommendedBlocks []BlockSuggestion
}

// BlockSuggestion represents a recommended availability duration with product combinations
type BlockSuggestion struct {
	DurationMinutes     int
	Rationale           string
	ProductCombinations []ProductCombo
	Priority            int // 1 = highest priority (standard durations), 2 = good, 3 = alternative
}

// ProductCombo represents a specific combination of products that fit in a time block
type ProductCombo struct {
	Products      []ProductInfo
	TotalDuration int
	SessionCount  int
}

// ProductInfo contains simplified product information for suggestions
type ProductInfo struct {
	ID         uuid.UUID
	Name       string
	Duration   int
	BufferTime int
}
