package models_test

import (
	"testing"

	"github.com/hengadev/leviosa/internal/domain/user/models"
	"github.com/hengadev/leviosa/tests/utils"

	"github.com/hengadev/test-assert"
)

func TestString(t *testing.T) {
	tests := []struct {
		role     models.Role
		expected string
		name     string
	}{
		{role: models.VISITOR, expected: "visitor", name: "Get string unknown"},
		{role: models.STANDARD, expected: "basic", name: "Get string basic"},
		{role: models.PREMIUM, expected: "premium", name: "Get string premium"},
		{role: models.GUEST, expected: "guest", name: "Get string guest"},
		{role: models.PARTNER, expected: "partner", name: "Get string partner"},
		{role: models.ADMINISTRATOR, expected: "admin", name: "Get string admin"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.role.String()
			assert.Equal(t, got, tt.expected)
		})
	}
}

func TestConvertToRole(t *testing.T) {
	tests := []struct {
		roleStr  string
		expected models.Role
		name     string
	}{
		{roleStr: "visitor", expected: models.VISITOR, name: "Convert to VISITOR"},
		{roleStr: "standard", expected: models.STANDARD, name: "Convert to STANDARD"},
		{roleStr: "premium", expected: models.PREMIUM, name: "Convert to PREMIUM"},
		{roleStr: "guest", expected: models.GUEST, name: "Convert to GUEST"},
		{roleStr: "partner", expected: models.PARTNER, name: "Convert to PARTNER"},
		{roleStr: "administrator", expected: models.ADMINISTRATOR, name: "Convert to ADMINISTRATOR"},
		{roleStr: test.GenerateRandomString(5), expected: models.VISITOR, name: "Convert to VISITOR"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := models.ConvertToRole(tt.roleStr)
			assert.Equal(t, got, tt.expected)
		})
	}
}
