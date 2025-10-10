package auth

import (
	"testing"

	"github.com/hengadev/encx"
)

// NewTestCrypto creates a test crypto instance using encx.NewTestCrypto
func NewTestCrypto(t *testing.T) (encx.CryptoService, error) {
	return encx.NewTestCrypto(t)
}