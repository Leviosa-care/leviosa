package sessionService_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/session"
	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"

	"github.com/google/uuid"
	"github.com/hengadev/test-assert"
)

func TestCreateSession(t *testing.T) {
	userID := uuid.NewString()
	tests := []struct {
		name          string
		userID        string
		role          models.Role
		mockRepo      func() *MockRepo
		expectedError error
	}{
		{
			name:          "user ID not UUID",
			userID:        "user123",
			role:          models.STANDARD,
			mockRepo:      func() *MockRepo { return &MockRepo{} },
			expectedError: domain.ErrInvalidValue,
		},
		{
			name:          "invalid role",
			userID:        userID,
			role:          models.VISITOR,
			mockRepo:      func() *MockRepo { return &MockRepo{} },
			expectedError: domain.ErrInvalidValue,
		},
		// TODO: how to get that case, since the function will never return an error
		// {
		// 	name:   "JSON marshal error",
		// 	userID: userID,
		// 	role:   models.BASIC,
		// 	mockRepo: func() *MockRepo {
		// 		return &MockRepo{
		// 			CreateSessionFunc: func(ctx context.Context, sessionID string, sessionEncoded []byte) error {
		// 				return rp.
		// 			},
		// 		}
		// 	},
		// 	expectedError: domain.ErrInvalidValue,
		// },
		{
			name:   "database error",
			userID: userID,
			role:   models.BASIC,
			mockRepo: func() *MockRepo {
				return &MockRepo{
					CreateSessionFunc: func(ctx context.Context, sessionID string, sessionEncoded []byte) error {
						return rp.ErrDatabase
					},
				}
			},
			expectedError: domain.ErrQueryFailed,
		},
		{
			name:   "context deadline exceeded error",
			userID: userID,
			role:   models.BASIC,
			mockRepo: func() *MockRepo {
				return &MockRepo{
					CreateSessionFunc: func(ctx context.Context, sessionID string, sessionEncoded []byte) error {
						return rp.ErrContext
					},
				}
			},
			expectedError: rp.ErrContext,
		},
		{
			name:   "context deadline exceeded error",
			userID: userID,
			role:   models.BASIC,
			mockRepo: func() *MockRepo {
				return &MockRepo{
					CreateSessionFunc: func(ctx context.Context, sessionID string, sessionEncoded []byte) error {
						return errors.New("unexpect type error")
					},
				}
			},
			expectedError: domain.ErrUnexpectedType,
		},
		{
			name:   "succcessful case",
			userID: userID,
			role:   models.BASIC,
			mockRepo: func() *MockRepo {
				return &MockRepo{
					CreateSessionFunc: func(ctx context.Context, sessionID string, sessionEncoded []byte) error {
						return nil
					},
				}
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			repo := tt.mockRepo()
			service := sessionService.New(repo)
			sessionID, err := service.CreateSession(ctx, tt.userID, tt.role)
			assert.EqualError(t, err, tt.expectedError)
			if tt.expectedError != nil {
				assert.Equal(t, sessionID, "")
				return
			}
			if err := uuid.Validate(sessionID); err != nil {
				t.Errorf("expect session ID to be UUID: %s", err)
			}
		})
	}
}
