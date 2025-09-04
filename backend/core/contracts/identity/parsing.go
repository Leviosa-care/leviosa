package identity

import (
	"fmt"
	"strings"

	"github.com/Leviosa-care/core/errs"
)

// ParseRole converts a string to a Role with proper error handling (case-insensitive)
func ParseRole(roleStr string) (Role, error) {
	normalized := strings.ToLower(strings.TrimSpace(roleStr))
	if role, ok := stringRoles[normalized]; ok {
		return role, nil
	}
	return Visitor, errs.NewInvalidValueErr(fmt.Sprintf("invalid role: %q", roleStr))
}

// ConvertToRole converts a string to a Role (deprecated: use ParseRole instead)
// Deprecated: Use ParseRole for better error handling
func ConvertToRole(role string) (Role, bool) {
	r, err := ParseRole(role)
	return r, err == nil
}
