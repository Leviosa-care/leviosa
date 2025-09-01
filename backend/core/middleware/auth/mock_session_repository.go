package auth

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

func (m *MockSessionRepository) FindSessionByAccessToken(ctx context.Context, accessTokenHash string) ([]byte, error) {
	args := m.Called(ctx, accessTokenHash)
	return args.Get(0).([]byte), args.Error(1)
}
func (m *MockSessionRepository) FindSessionByRefreshToken(ctx context.Context, refreshTokenHash string) ([]byte, error) {
	args := m.Called(ctx, refreshTokenHash)
	return args.Get(0).([]byte), args.Error(1)
}
