package domain

import (
	"strings"
)

type PublishedStatus string

const (
	Published PublishedStatus = "published"
	Draft     PublishedStatus = "draft"
	Archived  PublishedStatus = "archived"
)

// IsValid checks if the PublishedStatus is one of the defined constants.
func (ps PublishedStatus) IsValid() bool {
	switch strings.ToLower(string(ps)) { // Case-insensitive check
	case strings.ToLower(string(Published)), strings.ToLower(string(Draft)), strings.ToLower(string(Archived)):
		return true
	}
	return false
}
