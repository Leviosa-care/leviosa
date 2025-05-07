package userService_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/user"
	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/pkg/config"

	"github.com/google/uuid"
	"github.com/hengadev/test-assert"
)

func TestGetUserSessionData(t *testing.T) {
	email := "john.doe@hotmail.com"
	userID := uuid.NewString()
	tests := []struct {
		name          string
		email         string
		mockRepo      func() *MockRepo
		expectedError error
		expectedID    string
		expectedRole  models.Role
	}{
		{
			name:  "invalid email",
			email: "",
			mockRepo: func() *MockRepo {
				return &MockRepo{}
			},
			expectedError: domain.ErrInvalidValue,
			expectedID:    "",
			expectedRole:  models.VISITOR,
		},
		{
			name:  "database error",
			email: email,
			mockRepo: func() *MockRepo {
				return &MockRepo{
					GetUserSessionDataFunc: func(ctx context.Context, email string) (string, models.Role, error) {
						return "", models.VISITOR, rp.ErrDatabase
					},
				}
			},
			expectedError: domain.ErrQueryFailed,
			expectedID:    "",
			expectedRole:  models.VISITOR,
		},
		{
			name:  "context error",
			email: email,
			mockRepo: func() *MockRepo {
				return &MockRepo{
					GetUserSessionDataFunc: func(ctx context.Context, email string) (string, models.Role, error) {
						return "", models.VISITOR, rp.ErrContext
					},
				}
			},
			expectedError: rp.ErrContext,
			expectedID:    "",
			expectedRole:  models.VISITOR,
		},
		{
			name:  "unexpected error",
			email: email,
			mockRepo: func() *MockRepo {
				return &MockRepo{
					GetUserSessionDataFunc: func(ctx context.Context, email string) (string, models.Role, error) {
						return "", models.VISITOR, errors.New("unexpected error")
					},
				}
			},
			expectedError: domain.ErrUnexpectedType,
			expectedID:    "",
			expectedRole:  models.VISITOR,
		},
		{
			name:  "successful case",
			email: email,
			mockRepo: func() *MockRepo {
				return &MockRepo{
					GetUserSessionDataFunc: func(ctx context.Context, email string) (string, models.Role, error) {
						return userID, models.STANDARD, nil
					},
				}
			},
			expectedError: nil,
			expectedID:    userID,
			expectedRole:  models.STANDARD,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := tt.mockRepo()
			// TODO: I need to find the config for, or some sort of mock config just for the test
			config := &config.SecurityConfig{}
			service := userService.New(repo, config)
			id, role, err := service.GetUserSessionData(context.Background(), tt.email)
			assert.EqualError(t, err, tt.expectedError)
			assert.Equal(t, id, tt.expectedID)
			assert.Equal(t, role, tt.expectedRole)
		})
	}
}
