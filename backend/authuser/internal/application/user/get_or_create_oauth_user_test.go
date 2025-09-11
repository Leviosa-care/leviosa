package user_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/authuser/internal/application/user"
	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock implementations for testing
type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) ExistsByEmailHash(ctx context.Context, emailHash string) (bool, error) {
	args := m.Called(ctx, emailHash)
	return args.Bool(0), args.Error(1)
}

func (m *mockUserRepository) GetUserByEmailHash(ctx context.Context, emailHash string) (*domain.User, error) {
	args := m.Called(ctx, emailHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepository) GetPendingUsers(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *mockUserRepository) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *mockUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockUserRepository) GetUserByGoogleID(ctx context.Context, googleID string) (*domain.User, error) {
	args := m.Called(ctx, googleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepository) GetUserByAppleID(ctx context.Context, appleID string) (*domain.User, error) {
	args := m.Called(ctx, appleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepository) ExistsByGoogleID(ctx context.Context, googleID string) (bool, error) {
	args := m.Called(ctx, googleID)
	return args.Bool(0), args.Error(1)
}

func (m *mockUserRepository) ExistsByAppleID(ctx context.Context, appleID string) (bool, error) {
	args := m.Called(ctx, appleID)
	return args.Bool(0), args.Error(1)
}

type mockCryptoService struct {
	mock.Mock
}

func (m *mockCryptoService) HashBasic(ctx context.Context, data []byte) string {
	args := m.Called(ctx, data)
	return args.String(0)
}

func (m *mockCryptoService) ProcessStruct(ctx context.Context, data interface{}) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *mockCryptoService) DecryptStruct(ctx context.Context, data interface{}) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

type mockStripeService struct {
	mock.Mock
}

func (m *mockStripeService) CreateCustomer(ctx context.Context, userID uuid.UUID, email, firstName, lastName string) (*mockStripeCustomer, error) {
	args := m.Called(ctx, userID, email, firstName, lastName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mockStripeCustomer), args.Error(1)
}

type mockStripeCustomer struct {
	ID string
}

func TestUserService_GetOrCreateOAuthUser(t *testing.T) {
	ctx := context.Background()

	t.Run("should return existing OAuth user when Google ID exists", func(t *testing.T) {
		mockRepo := &mockUserRepository{}
		mockCrypto := &mockCryptoService{}
		mockStripe := &mockStripeService{}
		
		userService := user.New(mockRepo, mockCrypto, mockStripe)

		existingUser := &domain.User{
			ID:        uuid.New(),
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
			GoogleID:  "google123",
			State:     domain.Active,
			Role:      identity.Standard.String(),
		}

		// Mock that the Google ID exists
		mockRepo.On("GetUserByGoogleID", ctx, "google123").Return(existingUser, nil)
		mockCrypto.On("DecryptStruct", ctx, existingUser).Return(nil)

		result, isNewUser, err := userService.GetOrCreateOAuthUser(ctx, "google", "google123", "test@example.com", "John", "Doe")

		require.NoError(t, err)
		assert.False(t, isNewUser)
		assert.Equal(t, existingUser.ID, result.ID)
		assert.Equal(t, existingUser.Email, result.Email)
		
		mockRepo.AssertExpectations(t)
		mockCrypto.AssertExpectations(t)
	})

	t.Run("should create new OAuth user when neither OAuth ID nor email exists", func(t *testing.T) {
		mockRepo := &mockUserRepository{}
		mockCrypto := &mockCryptoService{}
		mockStripe := &mockStripeService{}
		
		userService := user.New(mockRepo, mockCrypto, mockStripe)

		// Mock that neither OAuth ID nor email exists
		mockRepo.On("GetUserByGoogleID", ctx, "google456").Return(nil, errs.ErrRepositoryNotFound)
		mockRepo.On("GetUserByEmailHash", ctx, "hashed_email").Return(nil, errs.ErrRepositoryNotFound)
		
		// Mock crypto operations
		mockCrypto.On("HashBasic", ctx, []byte("new@example.com")).Return("hashed_email")
		mockCrypto.On("ProcessStruct", ctx, mock.AnythingOfType("*domain.User")).Return(nil)
		mockCrypto.On("DecryptStruct", ctx, mock.AnythingOfType("*domain.User")).Return(nil)
		
		// Mock user creation
		mockRepo.On("CreateUser", ctx, mock.AnythingOfType("*domain.User")).Return(nil)
		
		// Mock Stripe customer creation
		mockStripeCustomer := &mockStripeCustomer{ID: "stripe_customer_123"}
		mockStripe.On("CreateCustomer", ctx, mock.AnythingOfType("uuid.UUID"), "new@example.com", "Jane", "Smith").Return(mockStripeCustomer, nil)
		mockRepo.On("UpdateUser", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

		result, isNewUser, err := userService.GetOrCreateOAuthUser(ctx, "google", "google456", "new@example.com", "Jane", "Smith")

		require.NoError(t, err)
		assert.True(t, isNewUser)
		assert.Equal(t, "new@example.com", result.Email)
		assert.Equal(t, "Jane", result.FirstName)
		assert.Equal(t, "Smith", result.LastName)
		
		mockRepo.AssertExpectations(t)
		mockCrypto.AssertExpectations(t)
		mockStripe.AssertExpectations(t)
	})

	t.Run("should return error with invalid provider", func(t *testing.T) {
		mockRepo := &mockUserRepository{}
		mockCrypto := &mockCryptoService{}
		mockStripe := &mockStripeService{}
		
		userService := user.New(mockRepo, mockCrypto, mockStripe)

		_, _, err := userService.GetOrCreateOAuthUser(ctx, "invalid_provider", "oauth123", "test@example.com", "John", "Doe")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported provider")
	})

	t.Run("should return error with missing required parameters", func(t *testing.T) {
		mockRepo := &mockUserRepository{}
		mockCrypto := &mockCryptoService{}
		mockStripe := &mockStripeService{}
		
		userService := user.New(mockRepo, mockCrypto, mockStripe)

		// Test missing provider
		_, _, err := userService.GetOrCreateOAuthUser(ctx, "", "oauth123", "test@example.com", "John", "Doe")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "provider, OAuth user ID, and email are required")

		// Test missing OAuth ID
		_, _, err = userService.GetOrCreateOAuthUser(ctx, "google", "", "test@example.com", "John", "Doe")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "provider, OAuth user ID, and email are required")

		// Test missing email
		_, _, err = userService.GetOrCreateOAuthUser(ctx, "google", "oauth123", "", "John", "Doe")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "provider, OAuth user ID, and email are required")
	})
}