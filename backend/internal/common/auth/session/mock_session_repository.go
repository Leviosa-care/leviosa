package session

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockSessionRepository implements the minimal SessionRepository interface for testing
type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) FindSessionByTokenHash(ctx context.Context, tokenHash string) ([]byte, error) {
	args := m.Called(ctx, tokenHash)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockSessionRepository) FindSessionByAccessTokenHash(ctx context.Context, accessTokenHash string) (string, []byte, error) {
	args := m.Called(ctx, accessTokenHash)
	return args.String(0), args.Get(1).([]byte), args.Error(2)
}
func (m *MockSessionRepository) FindSessionByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (string, []byte, error) {
	args := m.Called(ctx, refreshTokenHash)
	return args.String(0), args.Get(1).([]byte), args.Error(2)
}
