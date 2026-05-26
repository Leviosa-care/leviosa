package user_test

import (
	"context"
	"io"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/application/user"
	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/hengadev/encx"
	stripe "github.com/stripe/stripe-go/v82"
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

func (m *mockUserRepository) GetUserByEmailHash(ctx context.Context, emailHash string) (*domain.UserEncx, error) {
	args := m.Called(ctx, emailHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserEncx), args.Error(1)
}

func (m *mockUserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.UserEncx, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserEncx), args.Error(1)
}

func (m *mockUserRepository) GetPendingUsers(ctx context.Context) ([]*domain.UserEncx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.UserEncx), args.Error(1)
}

func (m *mockUserRepository) GetAllUsers(ctx context.Context) ([]*domain.UserEncx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.UserEncx), args.Error(1)
}

func (m *mockUserRepository) CreateUser(ctx context.Context, u *domain.UserEncx) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *mockUserRepository) UpdateUser(ctx context.Context, u *domain.UserEncx) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *mockUserRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockUserRepository) GetUserByGoogleID(ctx context.Context, googleID string) (*domain.UserEncx, error) {
	args := m.Called(ctx, googleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserEncx), args.Error(1)
}

func (m *mockUserRepository) GetUserByAppleID(ctx context.Context, appleID string) (*domain.UserEncx, error) {
	args := m.Called(ctx, appleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserEncx), args.Error(1)
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

func (m *mockCryptoService) GetPepper() []byte                  { return nil }
func (m *mockCryptoService) GetArgon2Params() *encx.Argon2Params { return nil }
func (m *mockCryptoService) GetAlias() string                   { return "" }
func (m *mockCryptoService) GenerateDEK() ([]byte, error)       { return nil, nil }
func (m *mockCryptoService) EncryptData(_ context.Context, _ []byte, _ []byte) ([]byte, error) {
	return nil, nil
}
func (m *mockCryptoService) DecryptData(_ context.Context, _ []byte, _ []byte) ([]byte, error) {
	return nil, nil
}
func (m *mockCryptoService) EncryptDEK(_ context.Context, _ []byte) ([]byte, error) { return nil, nil }
func (m *mockCryptoService) DecryptDEKWithVersion(_ context.Context, _ []byte, _ int) ([]byte, error) {
	return nil, nil
}
func (m *mockCryptoService) RotateKEK(_ context.Context) error { return nil }
func (m *mockCryptoService) HashBasic(ctx context.Context, data []byte) string {
	args := m.Called(ctx, data)
	return args.String(0)
}
func (m *mockCryptoService) HashSecure(_ context.Context, _ []byte) (string, error) { return "", nil }
func (m *mockCryptoService) CompareSecureHashAndValue(_ context.Context, _ any, _ string) (bool, error) {
	return false, nil
}
func (m *mockCryptoService) CompareBasicHashAndValue(_ context.Context, _ any, _ string) (bool, error) {
	return false, nil
}
func (m *mockCryptoService) EncryptStream(_ context.Context, _ io.Reader, _ io.Writer, _ []byte) error {
	return nil
}
func (m *mockCryptoService) DecryptStream(_ context.Context, _ io.Reader, _ io.Writer, _ []byte) error {
	return nil
}
func (m *mockCryptoService) GetCurrentKEKVersion(_ context.Context, _ string) (int, error) {
	return 0, nil
}
func (m *mockCryptoService) GetKMSKeyIDForVersion(_ context.Context, _ string, _ int) (string, error) {
	return "", nil
}

type mockStripeService struct {
	mock.Mock
}

func (m *mockStripeService) CreateCustomer(ctx context.Context, userID uuid.UUID, email, firstName, lastName string) (*stripe.Customer, error) {
	args := m.Called(ctx, userID, email, firstName, lastName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.Customer), args.Error(1)
}

func (m *mockStripeService) RetrieveCustomer(_ context.Context, _ string) (*stripe.Customer, error) {
	return nil, nil
}

func (m *mockStripeService) UpdateCustomer(_ context.Context, _ string, _ *stripe.CustomerUpdateParams) (*stripe.Customer, error) {
	return nil, nil
}

func (m *mockStripeService) DeleteCustomer(_ context.Context, _ string) (*stripe.Customer, error) {
	return nil, nil
}

func (m *mockStripeService) FindCustomerByUserID(_ context.Context, _ uuid.UUID) (*stripe.Customer, error) {
	return nil, nil
}

func (m *mockStripeService) CreateConnectedAccount(_ context.Context, _ uuid.UUID) (string, error) {
	return "", nil
}

func TestUserService_GetOrCreateOAuthUser(t *testing.T) {
	ctx := context.Background()

	t.Run("should return existing OAuth user when Google ID exists", func(t *testing.T) {
		mockRepo := &mockUserRepository{}
		mockCrypto := &mockCryptoService{}
		mockStripe := &mockStripeService{}

		userService := user.New(mockRepo, mockCrypto, mockStripe)

		existingUserEncx := &domain.UserEncx{
			ID: uuid.New(),
		}

		// Mock that the Google ID exists
		mockRepo.On("GetUserByGoogleID", ctx, "google123").Return(existingUserEncx, nil)

		_, isNewUser, err := userService.GetOrCreateOAuthUser(ctx, "google", "google123", "test@example.com", "John", "Doe")

		require.NoError(t, err)
		assert.False(t, isNewUser)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should create new OAuth user when neither OAuth ID nor email exists", func(t *testing.T) {
		mockRepo := &mockUserRepository{}
		mockCrypto := &mockCryptoService{}
		mockStripe := &mockStripeService{}

		userService := user.New(mockRepo, mockCrypto, mockStripe)

		// Mock that neither OAuth ID nor email exists
		mockRepo.On("GetUserByGoogleID", ctx, "google456").Return(nil, errs.ErrRepositoryNotFound)
		mockRepo.On("GetUserByEmailHash", ctx, mock.AnythingOfType("string")).Return(nil, errs.ErrRepositoryNotFound)

		// Mock user creation
		mockRepo.On("CreateUser", ctx, mock.AnythingOfType("*domain.UserEncx")).Return(nil)

		// Mock Stripe customer creation
		mockStripe.On("CreateCustomer", ctx, mock.AnythingOfType("uuid.UUID"), "new@example.com", "Jane", "Smith").Return(nil, nil)

		result, isNewUser, err := userService.GetOrCreateOAuthUser(ctx, "google", "google456", "new@example.com", "Jane", "Smith")

		require.NoError(t, err)
		assert.True(t, isNewUser)
		assert.Equal(t, "new@example.com", result.Email)
		assert.Equal(t, "Jane", result.FirstName)
		assert.Equal(t, "Smith", result.LastName)

		mockRepo.AssertExpectations(t)
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
