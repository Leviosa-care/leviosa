package partner_test

import (
	"context"
	"errors"
	"io"
	"testing"

	partnerApp "github.com/Leviosa-care/leviosa/backend/internal/authuser/application/partner"
	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/authuser/ports"
	catalogDomain "github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	catalogPorts "github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
	"github.com/stripe/stripe-go/v82"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// --- Mock implementations ---

type mockPartnerRepository struct {
	mock.Mock
}

func (m *mockPartnerRepository) CreatePartner(ctx context.Context, partner *domain.PartnerEncx) error {
	args := m.Called(ctx, partner)
	return args.Error(0)
}

func (m *mockPartnerRepository) GetPartnerByID(ctx context.Context, partnerID uuid.UUID) (*domain.PartnerEncx, error) {
	args := m.Called(ctx, partnerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PartnerEncx), args.Error(1)
}

func (m *mockPartnerRepository) GetPartnerByUserID(ctx context.Context, userID uuid.UUID) (*domain.PartnerEncx, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PartnerEncx), args.Error(1)
}

func (m *mockPartnerRepository) GetAllPartners(ctx context.Context) ([]*domain.PartnerEncx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.PartnerEncx), args.Error(1)
}

func (m *mockPartnerRepository) GetAllPartnersByCategory(ctx context.Context, categoryID uuid.UUID) ([]*domain.PartnerEncx, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.PartnerEncx), args.Error(1)
}

func (m *mockPartnerRepository) GetAllPartnersByCategories(ctx context.Context, categoryIDs []uuid.UUID) ([]*domain.PartnerEncx, error) {
	args := m.Called(ctx, categoryIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.PartnerEncx), args.Error(1)
}

func (m *mockPartnerRepository) GetAllPartnersByProduct(ctx context.Context, productID uuid.UUID) ([]*domain.PartnerEncx, error) {
	args := m.Called(ctx, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.PartnerEncx), args.Error(1)
}

func (m *mockPartnerRepository) GetAllPartnersByProducts(ctx context.Context, productIDs []uuid.UUID) ([]*domain.PartnerEncx, error) {
	args := m.Called(ctx, productIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.PartnerEncx), args.Error(1)
}

func (m *mockPartnerRepository) UpdatePartner(ctx context.Context, partner *domain.PartnerEncx) error {
	args := m.Called(ctx, partner)
	return args.Error(0)
}

func (m *mockPartnerRepository) DeletePartner(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockPartnerRepository) VerifyPartner(ctx context.Context, userID uuid.UUID, verifiedByUserID uuid.UUID) error {
	args := m.Called(ctx, userID, verifiedByUserID)
	return args.Error(0)
}

func (m *mockPartnerRepository) GetAllPartnersWithStripeAccount(ctx context.Context) ([]*domain.PartnerEncx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.PartnerEncx), args.Error(1)
}

func (m *mockPartnerRepository) UpdatePartnerStripeStatus(ctx context.Context, partnerID uuid.UUID, status domain.StripeAccountStatus) error {
	args := m.Called(ctx, partnerID, status)
	return args.Error(0)
}

func (m *mockPartnerRepository) GetPublicPartners(ctx context.Context) ([]*domain.PublicPartnerRow, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.PublicPartnerRow), args.Error(1)
}

func (m *mockPartnerRepository) GetPublicPartnerByID(ctx context.Context, partnerID uuid.UUID) (*domain.PublicPartnerRow, error) {
	args := m.Called(ctx, partnerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PublicPartnerRow), args.Error(1)
}

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.UserEncx, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserEncx), args.Error(1)
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

func (m *mockUserRepository) GetPendingUsers(ctx context.Context) ([]*domain.UserEncx, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.UserEncx), args.Error(1)
}

func (m *mockUserRepository) GetAllUsers(ctx context.Context) ([]*domain.UserEncx, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.UserEncx), args.Error(1)
}

func (m *mockUserRepository) CreateUser(ctx context.Context, user *domain.UserEncx) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepository) UpdateUser(ctx context.Context, user *domain.UserEncx) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockUserRepository) ExistsByAppleID(ctx context.Context, appleID string) (bool, error) {
	args := m.Called(ctx, appleID)
	return args.Bool(0), args.Error(1)
}

func (m *mockUserRepository) ExistsByGoogleID(ctx context.Context, googleID string) (bool, error) {
	args := m.Called(ctx, googleID)
	return args.Bool(0), args.Error(1)
}

func (m *mockUserRepository) GetUserByAppleID(ctx context.Context, appleID string) (*domain.UserEncx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserEncx), args.Error(1)
}

func (m *mockUserRepository) GetUserByGoogleID(ctx context.Context, googleID string) (*domain.UserEncx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserEncx), args.Error(1)
}

type mockProductService struct {
	mock.Mock
}

func (m *mockProductService) GetProductByID(ctx context.Context, ID string) (*catalogDomain.ProductRes, error) {
	args := m.Called(ctx, ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*catalogDomain.ProductRes), args.Error(1)
}

func (m *mockProductService) GetAllPublishedProducts(ctx context.Context) ([]*catalogDomain.ProductRes, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*catalogDomain.ProductRes), args.Error(1)
}

func (m *mockProductService) GetAllProducts(ctx context.Context) ([]*catalogDomain.ProductRes, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*catalogDomain.ProductRes), args.Error(1)
}

type mockCategoryService struct {
	mock.Mock
}

func (m *mockCategoryService) GetCategoryByID(ctx context.Context, ID string) (*catalogDomain.Category, error) {
	args := m.Called(ctx, ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*catalogDomain.Category), args.Error(1)
}

func (m *mockCategoryService) GetAllPublishedCategories(ctx context.Context) ([]*catalogDomain.Category, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*catalogDomain.Category), args.Error(1)
}

func (m *mockCategoryService) GetAllCategories(ctx context.Context) ([]*catalogDomain.Category, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*catalogDomain.Category), args.Error(1)
}

type mockCryptoService struct {
	mock.Mock
}

func (m *mockCryptoService) GetPepper() []byte {
	args := m.Called()
	return args.Get(0).([]byte)
}

func (m *mockCryptoService) GetArgon2Params() *encx.Argon2Params {
	args := m.Called()
	if v := args.Get(0); v != nil {
		return v.(*encx.Argon2Params)
	}
	return nil
}

func (m *mockCryptoService) GetAlias() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockCryptoService) GenerateDEK() ([]byte, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockCryptoService) EncryptData(ctx context.Context, plaintext []byte, dek []byte) ([]byte, error) {
	args := m.Called(ctx, plaintext, dek)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockCryptoService) DecryptData(ctx context.Context, ciphertext []byte, dek []byte) ([]byte, error) {
	args := m.Called(ctx, ciphertext, dek)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockCryptoService) EncryptDEK(ctx context.Context, plaintextDEK []byte) ([]byte, error) {
	args := m.Called(ctx, plaintextDEK)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockCryptoService) DecryptDEKWithVersion(ctx context.Context, ciphertextDEK []byte, kekVersion int) ([]byte, error) {
	args := m.Called(ctx, ciphertextDEK, kekVersion)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockCryptoService) RotateKEK(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockCryptoService) HashBasic(ctx context.Context, data []byte) string {
	args := m.Called(ctx, data)
	return args.String(0)
}

func (m *mockCryptoService) HashSecure(ctx context.Context, value []byte) (string, error) {
	args := m.Called(ctx, value)
	return args.String(0), args.Error(1)
}

func (m *mockCryptoService) CompareSecureHashAndValue(ctx context.Context, value any, hashValue string) (bool, error) {
	args := m.Called(ctx, value, hashValue)
	return args.Bool(0), args.Error(1)
}

func (m *mockCryptoService) CompareBasicHashAndValue(ctx context.Context, value any, hashValue string) (bool, error) {
	args := m.Called(ctx, value, hashValue)
	return args.Bool(0), args.Error(1)
}

func (m *mockCryptoService) EncryptStream(ctx context.Context, reader io.Reader, writer io.Writer, dek []byte) error {
	args := m.Called(ctx, reader, writer, dek)
	return args.Error(0)
}

func (m *mockCryptoService) DecryptStream(ctx context.Context, reader io.Reader, writer io.Writer, dek []byte) error {
	args := m.Called(ctx, reader, writer, dek)
	return args.Error(0)
}

func (m *mockCryptoService) GetCurrentKEKVersion(ctx context.Context, alias string) (int, error) {
	args := m.Called(ctx, alias)
	return args.Int(0), args.Error(1)
}

func (m *mockCryptoService) GetKMSKeyIDForVersion(ctx context.Context, alias string, version int) (string, error) {
	args := m.Called(ctx, alias, version)
	return args.String(0), args.Error(1)
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

func (m *mockStripeService) RetrieveCustomer(ctx context.Context, customerID string) (*stripe.Customer, error) {
	args := m.Called(ctx, customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.Customer), args.Error(1)
}

func (m *mockStripeService) UpdateCustomer(ctx context.Context, customerID string, params *stripe.CustomerUpdateParams) (*stripe.Customer, error) {
	args := m.Called(ctx, customerID, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.Customer), args.Error(1)
}

func (m *mockStripeService) DeleteCustomer(ctx context.Context, customerID string) (*stripe.Customer, error) {
	args := m.Called(ctx, customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.Customer), args.Error(1)
}

func (m *mockStripeService) FindCustomerByUserID(ctx context.Context, userID uuid.UUID) (*stripe.Customer, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.Customer), args.Error(1)
}

func (m *mockStripeService) CreateConnectedAccount(ctx context.Context, userID uuid.UUID) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *mockStripeService) CreateAccountLink(ctx context.Context, accountID, returnType, returnURL, refreshURL string) (string, error) {
	args := m.Called(ctx, accountID, returnType, returnURL, refreshURL)
	return args.String(0), args.Error(1)
}

func (m *mockStripeService) VerifyConnectWebhookSignature(payload []byte, signature string) (string, bool, bool, error) {
	args := m.Called(payload, signature)
	return args.String(0), args.Bool(1), args.Bool(2), args.Error(3)
}

// --- Helper to create a PartnerService with mocks ---

func setupPartnerService(
	partnerRepo ports.PartnerRepository,
	userRepo ports.UserRepository,
	productSvc catalogPorts.PublicProductService,
	categorySvc catalogPorts.PublicCategoryService,
	crypto encx.CryptoService,
	stripeSvc ports.StripeService,
) ports.PartnerService {
	svc, _ := partnerApp.New(context.Background(), partnerRepo, userRepo, productSvc, categorySvc, crypto, stripeSvc)
	return svc
}

// --- Tests ---

func TestCreatePartner_StripeConnect(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	t.Run("should create Stripe connected account and store account ID on success", func(t *testing.T) {
		mockPartnerRepo := new(mockPartnerRepository)
		mockUserRepo := new(mockUserRepository)
		mockProductSvc := new(mockProductService)
		mockCategorySvc := new(mockCategoryService)
		mockCrypto := new(mockCryptoService)
		mockStripe := new(mockStripeService)

		stripeAccountID := "acct_test_connect_123"

		// Catalog returns empty lists (errgroup wraps context, so use mock.Anything)
		mockCategorySvc.On("GetAllPublishedCategories", mock.Anything).Return([]*catalogDomain.Category{}, nil)
		mockProductSvc.On("GetAllPublishedProducts", mock.Anything).Return([]*catalogDomain.ProductRes{}, nil)

		// Stripe Connect account creation succeeds
		mockStripe.On("CreateConnectedAccount", ctx, userID).Return(stripeAccountID, nil)

		// Crypto stubs for ProcessPartnerEncx
		mockCrypto.On("GenerateDEK").Return([]byte("test-dek-32-bytes-long-enough!!!"), nil)
		mockCrypto.On("EncryptData", ctx, mock.Anything, mock.Anything).Return([]byte("encrypted"), nil)
		mockCrypto.On("EncryptDEK", ctx, mock.Anything).Return([]byte("encrypted-dek"), nil)
		mockCrypto.On("GetAlias").Return("test-alias")
		mockCrypto.On("GetCurrentKEKVersion", ctx, "test-alias").Return(1, nil)

		// Partner repo CreatePartner captures the partner
		mockPartnerRepo.On("CreatePartner", ctx, mock.AnythingOfType("*domain.PartnerEncx")).Return(nil)

		svc := setupPartnerService(mockPartnerRepo, mockUserRepo, mockProductSvc, mockCategorySvc, mockCrypto, mockStripe)

		partner, err := svc.CreatePartner(ctx, userID, "bio", "experience", nil, nil)

		require.NoError(t, err)
		assert.NotNil(t, partner)
		assert.Equal(t, stripeAccountID, partner.StripeConnectedAccountID, "Stripe account ID should be stored on partner")
		assert.Equal(t, domain.StripeAccountStatusPending, partner.StripeAccountStatus, "Partner should start with pending Stripe status")

		mockStripe.AssertExpectations(t)
		mockPartnerRepo.AssertExpectations(t)
	})

	t.Run("should save partner with stripe_status=pending when Stripe fails", func(t *testing.T) {
		mockPartnerRepo := new(mockPartnerRepository)
		mockUserRepo := new(mockUserRepository)
		mockProductSvc := new(mockProductService)
		mockCategorySvc := new(mockCategoryService)
		mockCrypto := new(mockCryptoService)
		mockStripe := new(mockStripeService)

		// Catalog returns empty lists
		mockCategorySvc.On("GetAllPublishedCategories", mock.Anything).Return([]*catalogDomain.Category{}, nil)
		mockProductSvc.On("GetAllPublishedProducts", mock.Anything).Return([]*catalogDomain.ProductRes{}, nil)

		// Stripe Connect account creation fails
		mockStripe.On("CreateConnectedAccount", ctx, userID).Return("", errors.New("stripe api error"))

		// Crypto stubs for ProcessPartnerEncx
		mockCrypto.On("GenerateDEK").Return([]byte("test-dek-32-bytes-long-enough!!!"), nil)
		mockCrypto.On("EncryptData", ctx, mock.Anything, mock.Anything).Return([]byte("encrypted"), nil)
		mockCrypto.On("EncryptDEK", ctx, mock.Anything).Return([]byte("encrypted-dek"), nil)
		mockCrypto.On("GetAlias").Return("test-alias")
		mockCrypto.On("GetCurrentKEKVersion", ctx, "test-alias").Return(1, nil)

		// Partner repo should still be called (partner saved with pending status)
		mockPartnerRepo.On("CreatePartner", ctx, mock.AnythingOfType("*domain.PartnerEncx")).Return(nil)

		svc := setupPartnerService(mockPartnerRepo, mockUserRepo, mockProductSvc, mockCategorySvc, mockCrypto, mockStripe)

		partner, err := svc.CreatePartner(ctx, userID, "bio", "experience", nil, nil)

		require.NoError(t, err, "CreatePartner should not return error when Stripe fails")
		assert.NotNil(t, partner)
		assert.Empty(t, partner.StripeConnectedAccountID, "Stripe account ID should be empty when Stripe fails")
		assert.Equal(t, domain.StripeAccountStatusPending, partner.StripeAccountStatus, "Partner should be saved with pending Stripe status")

		mockStripe.AssertExpectations(t)
		mockPartnerRepo.AssertExpectations(t)
	})

	t.Run("should call CreateConnectedAccount exactly once with correct userID", func(t *testing.T) {
		mockPartnerRepo := new(mockPartnerRepository)
		mockUserRepo := new(mockUserRepository)
		mockProductSvc := new(mockProductService)
		mockCategorySvc := new(mockCategoryService)
		mockCrypto := new(mockCryptoService)
		mockStripe := new(mockStripeService)

		mockCategorySvc.On("GetAllPublishedCategories", mock.Anything).Return([]*catalogDomain.Category{}, nil)
		mockProductSvc.On("GetAllPublishedProducts", mock.Anything).Return([]*catalogDomain.ProductRes{}, nil)
		mockStripe.On("CreateConnectedAccount", ctx, userID).Return("acct_123", nil)
		mockCrypto.On("GenerateDEK").Return([]byte("test-dek-32-bytes-long-enough!!!"), nil)
		mockCrypto.On("EncryptData", ctx, mock.Anything, mock.Anything).Return([]byte("encrypted"), nil)
		mockCrypto.On("EncryptDEK", ctx, mock.Anything).Return([]byte("encrypted-dek"), nil)
		mockCrypto.On("GetAlias").Return("test-alias")
		mockCrypto.On("GetCurrentKEKVersion", ctx, "test-alias").Return(1, nil)
		mockPartnerRepo.On("CreatePartner", ctx, mock.AnythingOfType("*domain.PartnerEncx")).Return(nil)

		svc := setupPartnerService(mockPartnerRepo, mockUserRepo, mockProductSvc, mockCategorySvc, mockCrypto, mockStripe)

		_, err := svc.CreatePartner(ctx, userID, "bio", "experience", nil, nil)
		require.NoError(t, err)

		// Verify CreateConnectedAccount was called exactly once with the correct userID
		mockStripe.AssertCalled(t, "CreateConnectedAccount", ctx, userID)
		mockStripe.AssertNumberOfCalls(t, "CreateConnectedAccount", 1)
	})
}
