package aggregatorHandler

import "strings"

// maskEmail masks email for GDPR-compliant logging
func maskEmail(email string) string {
	if email == "" {
		return "[empty]"
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "[invalid]"
	}

	local, domain := parts[0], parts[1]
	if len(local) == 0 {
		return "[invalid]"
	}

	// Show first character + asterisks + last character if long enough
	if len(local) <= 2 {
		return string(local[0]) + "*@" + domain
	}

	return string(local[0]) + strings.Repeat("*", len(local)-2) + string(local[len(local)-1]) + "@" + domain
}
