package domain

import "strings"

// AvailabilityType defines allowed values for availability
type AvailabilityType string

const (
	Online   AvailabilityType = "online"
	InPerson AvailabilityType = "in-person"
	Hybrid   AvailabilityType = "hybrid"
)

// IsValid checks if the AvailabilityType is one of the defined constants.
func (it AvailabilityType) IsValid() bool {
	switch strings.ToLower(string(it)) { // Case-insensitive check
	case strings.ToLower(string(Online)), strings.ToLower(string(InPerson)), strings.ToLower(string(Hybrid)):
		return true
	}
	return false
}
