package auth

import (
	"context"

	"github.com/hashicorp/vault/api"
)

// SecretData represents the data returned from a secret store.
type SecretData struct {
	Data map[string]interface{}
}

// SecretReader abstracts reading secrets from a secrets store (e.g., Vault).
// This interface enables unit testing without a real Vault instance.
type SecretReader interface {
	// Read reads the secret at the given path and returns its top-level data map.
	// Returns nil SecretData when the path does not exist.
	Read(ctx context.Context, path string) (*SecretData, error)
}

// vaultSecretReader adapts the HashiCorp Vault client to the SecretReader interface.
type vaultSecretReader struct {
	client *api.Client
}

// NewVaultSecretReader creates a SecretReader backed by a real Vault client.
func NewVaultSecretReader(client *api.Client) SecretReader {
	return &vaultSecretReader{client: client}
}

func (r *vaultSecretReader) Read(ctx context.Context, path string) (*SecretData, error) {
	secret, err := r.client.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, nil
	}
	return &SecretData{Data: secret.Data}, nil
}
