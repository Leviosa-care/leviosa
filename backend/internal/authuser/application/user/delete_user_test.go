package user_test

import (
	"context"
	"errors"
	"testing"

	userService "github.com/Leviosa-care/leviosa/backend/internal/authuser/application/user"
	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stripe/stripe-go/v82"
)

// Mock implementations
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.UserEncx, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserEncx), args.Error(1)
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// Implement other required interface methods as no-ops for this test
func (m *MockUserRepository) ExistsByEmailHash(ctx context.Context, emailHash string) (bool, error) {
	args := m.Called(ctx, emailHash)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmailHash(ctx context.Context, emailHash string) (*domain.UserEncx, error) {
	args := m.Called(ctx, emailHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserEncx), args.Error(1)
}

func (m *MockUserRepository) GetPendingUsers(ctx context.Context) ([]*domain.UserEncx, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.UserEncx), args.Error(1)
}

func (m *MockUserRepository) GetAllUsers(ctx context.Context) ([]*domain.UserEncx, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.UserEncx), args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *domain.UserEncx) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *domain.UserEncx) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) ExistsByAppleID(ctx context.Context, appleID string) (bool, error) {
	args := m.Called(ctx, appleID)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) ExistsByGoogleID(ctx context.Context, googleID string) (bool, error) {
	args := m.Called(ctx, googleID)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) GetUserByAppleID(ctx context.Context, appleID string) (*domain.UserEncx, error) {
	args := m.Called(ctx)
	return args.Get(0).(*domain.UserEncx), args.Error(1)
}

func (m *MockUserRepository) GetUserByGoogleID(ctx context.Context, googleID string) (*domain.UserEncx, error) {
	args := m.Called(ctx)
	return args.Get(0).(*domain.UserEncx), args.Error(1)
}

type MockCryptoService struct {
	mock.Mock
}

func (m *MockCryptoService) DecryptStruct(ctx context.Context, data interface{}) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockCryptoService) ProcessStruct(ctx context.Context, data interface{}) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockCryptoService) EncryptStruct(ctx context.Context, data interface{}) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

type MockStripeService struct {
	mock.Mock
}

func (m *MockStripeService) DeleteCustomer(ctx context.Context, customerID string) (*stripe.Customer, error) {
	args := m.Called(ctx, customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.Customer), args.Error(1)
}

// Implement other required interface methods as no-ops for this test
func (m *MockStripeService) CreateCustomer(ctx context.Context, userID uuid.UUID, email, firstName, lastName string) (*stripe.Customer, error) {
	args := m.Called(ctx, userID, email, firstName, lastName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.Customer), args.Error(1)
}

func (m *MockStripeService) RetrieveCustomer(ctx context.Context, customerID string) (*stripe.Customer, error) {
	args := m.Called(ctx, customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.Customer), args.Error(1)
}

func (m *MockStripeService) UpdateCustomer(ctx context.Context, customerID string, params *stripe.CustomerUpdateParams) (*stripe.Customer, error) {
	args := m.Called(ctx, customerID, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.Customer), args.Error(1)
}

func (m *MockStripeService) FindCustomerByUserID(ctx context.Context, userID uuid.UUID) (*stripe.Customer, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.Customer), args.Error(1)
}

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	stripeCustomerID := "cus_test123"

	t.Run("should successfully delete user with Stripe customer", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		mockCrypto := new(MockCryptoService)
		mockStripe := new(MockStripeService)

		testUser := &domain.User{
			ID:               userID,
			Email:            "test@example.com",
			StripeCustomerID: stripeCustomerID,
		}

		// Mock expectations
		mockRepo.On("GetUserByID", ctx, userID).Return(testUser, nil)
		mockCrypto.On("DecryptStruct", ctx, testUser).Return(nil)
		mockStripe.On("DeleteCustomer", ctx, stripeCustomerID).Return(&stripe.Customer{ID: stripeCustomerID}, nil)
		mockRepo.On("DeleteUser", ctx, userID).Return(nil)

		service := userService.New(mockRepo, mockCrypto, mockStripe)

		// Act
		err := service.DeleteUser(ctx, userID)

		// Assert
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockCrypto.AssertExpectations(t)
		mockStripe.AssertExpectations(t)
	})

	t.Run("should successfully delete user without Stripe customer", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		mockCrypto := new(MockCryptoService)
		mockStripe := new(MockStripeService)

		testUser := &domain.User{
			ID:               userID,
			Email:            "test@example.com",
			StripeCustomerID: "", // No Stripe customer
		}

		// Mock expectations
		mockRepo.On("GetUserByID", ctx, userID).Return(testUser, nil)
		mockCrypto.On("DecryptStruct", ctx, testUser).Return(nil)
		// Stripe should not be called when no customer ID
		mockRepo.On("DeleteUser", ctx, userID).Return(nil)

		service := userService.New(mockRepo, mockCrypto, mockStripe)

		// Act
		err := service.DeleteUser(ctx, userID)

		// Assert
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockCrypto.AssertExpectations(t)
		mockStripe.AssertNotCalled(t, "DeleteCustomer")
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		mockCrypto := new(MockCryptoService)
		mockStripe := new(MockStripeService)

		mockRepo.On("GetUserByID", ctx, userID).Return(nil, errs.ErrRepositoryNotFound)

		service := userService.New(mockRepo, mockCrypto, mockStripe)

		// Act
		err := service.DeleteUser(ctx, userID)

		// Assert
		require.Error(t, err)
		assert.True(t, errors.Is(err, errs.ErrDomainNotFound))
		mockRepo.AssertExpectations(t)
	})

	t.Run("should handle Stripe customer not found gracefully", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		mockCrypto := new(MockCryptoService)
		mockStripe := new(MockStripeService)

		testUser := &domain.User{
			ID:               userID,
			Email:            "test@example.com",
			StripeCustomerID: stripeCustomerID,
		}

		// Mock expectations
		mockRepo.On("GetUserByID", ctx, userID).Return(testUser, nil)
		mockCrypto.On("DecryptStruct", ctx, testUser).Return(nil)
		mockStripe.On("DeleteCustomer", ctx, stripeCustomerID).Return(nil, errs.ErrInvalidValue) // Stripe customer not found
		mockRepo.On("DeleteUser", ctx, userID).Return(nil)                                       // Should still delete user

		service := userService.New(mockRepo, mockCrypto, mockStripe)

		// Act
		err := service.DeleteUser(ctx, userID)

		// Assert
		require.NoError(t, err) // Should succeed despite Stripe customer not found
		mockRepo.AssertExpectations(t)
		mockCrypto.AssertExpectations(t)
		mockStripe.AssertExpectations(t)
	})

	t.Run("should return error when Stripe service fails", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		mockCrypto := new(MockCryptoService)
		mockStripe := new(MockStripeService)

		testUser := &domain.User{
			ID:               userID,
			Email:            "test@example.com",
			StripeCustomerID: stripeCustomerID,
		}

		// Mock expectations
		mockRepo.On("GetUserByID", ctx, userID).Return(testUser, nil)
		mockCrypto.On("DecryptStruct", ctx, testUser).Return(nil)
		mockStripe.On("DeleteCustomer", ctx, stripeCustomerID).Return(nil, errs.NewExternalServiceErr(errors.New("stripe error"), "stripe unavailable"))

		service := userService.New(mockRepo, mockCrypto, mockStripe)

		// Act
		err := service.DeleteUser(ctx, userID)

		// Assert
		require.Error(t, err)
		assert.True(t, errors.Is(err, errs.ErrExternalService))
		mockRepo.AssertExpectations(t)
		mockCrypto.AssertExpectations(t)
		mockStripe.AssertExpectations(t)
		// User should not be deleted if Stripe fails
		mockRepo.AssertNotCalled(t, "DeleteUser")
	})

	t.Run("should return error when database deletion fails", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		mockCrypto := new(MockCryptoService)
		mockStripe := new(MockStripeService)

		testUser := &domain.User{
			ID:               userID,
			Email:            "test@example.com",
			StripeCustomerID: stripeCustomerID,
		}

		// Mock expectations
		mockRepo.On("GetUserByID", ctx, userID).Return(testUser, nil)
		mockCrypto.On("DecryptStruct", ctx, testUser).Return(nil)
		mockStripe.On("DeleteCustomer", ctx, stripeCustomerID).Return(&stripe.Customer{ID: stripeCustomerID}, nil)
		mockRepo.On("DeleteUser", ctx, userID).Return(errs.ErrRepositoryNotFound)

		service := userService.New(mockRepo, mockCrypto, mockStripe)

		// Act
		err := service.DeleteUser(ctx, userID)

		// Assert
		require.Error(t, err)
		assert.True(t, errors.Is(err, errs.ErrDomainNotFound))
		mockRepo.AssertExpectations(t)
		mockCrypto.AssertExpectations(t)
		mockStripe.AssertExpectations(t)
	})

	t.Run("should return error when decryption fails", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		mockCrypto := new(MockCryptoService)
		mockStripe := new(MockStripeService)

		testUser := &domain.User{
			ID:    userID,
			Email: "test@example.com",
		}

		// Mock expectations
		mockRepo.On("GetUserByID", ctx, userID).Return(testUser, nil)
		mockCrypto.On("DecryptStruct", ctx, testUser).Return(errors.New("decryption failed"))

		service := userService.New(mockRepo, mockCrypto, mockStripe)

		// Act
		err := service.DeleteUser(ctx, userID)

		// Assert
		require.Error(t, err)
		assert.True(t, errors.Is(err, errs.ErrNotDecrypted))
		mockRepo.AssertExpectations(t)
		mockCrypto.AssertExpectations(t)
		// Neither Stripe nor database operations should be called
		mockStripe.AssertNotCalled(t, "DeleteCustomer")
		mockRepo.AssertNotCalled(t, "DeleteUser")
	})
}

