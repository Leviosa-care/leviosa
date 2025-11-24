package mocks

import (
	"context"

	catalogDomain "github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/google/uuid"
)

// MockProductService provides a test implementation of PublicProductService
type MockProductService struct {
	products []*catalogDomain.ProductRes
}

// NewMockProductService creates a mock product service with default test products
func NewMockProductService() *MockProductService {
	return &MockProductService{
		products: []*catalogDomain.ProductRes{
			{
				ID:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Name:       "60-Minute Massage",
				Duration:   60,
				BufferTime: 15,
			},
			{
				ID:         uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Name:       "90-Minute Massage",
				Duration:   90,
				BufferTime: 15,
			},
			{
				ID:         uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				Name:       "30-Minute Consultation",
				Duration:   30,
				BufferTime: 10,
			},
		},
	}
}

// GetAllActiveProductsByPartnerID returns mock products for testing
func (m *MockProductService) GetAllActiveProductsByPartnerID(ctx context.Context, partnerID uuid.UUID) ([]*catalogDomain.ProductRes, error) {
	return m.products, nil
}

// SetProducts allows tests to configure custom product data
func (m *MockProductService) SetProducts(products []*catalogDomain.ProductRes) {
	m.products = products
}
