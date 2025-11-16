package product

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"

	"github.com/google/uuid"
)

// CreateProduct handles the creation of a new product, including validation and persistence.
func (s *ProductService) CreateProduct(ctx context.Context, request *domain.CreateProductRequest) (string, error) {
	categoryID, err := uuid.Parse(request.CategoryID)
	if err != nil {
		return "", errs.NewInvalidValueErr("category ID is required and must be a valid UUID")
	}

	_, err = s.sharedRepo.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return "", fmt.Errorf("get category with given ID: %w", err)
	}

	product := &domain.Product{
		Name:              strings.ToLower(request.Name),
		Description:       request.Description,
		CategoryID:        categoryID,
		Duration:          request.Duration,
		Status:            domain.Draft,
		Availability:      request.Availability,
		BufferTime:        request.BufferTime,
		CancellationHours: request.CancellationHours,
		Metadata:          request.Metadata,
	}

	if err := product.Valid(ctx); err != nil {
		return "", errs.NewInvalidValueErr(err.Error())
	}

	product.ID = uuid.New()
	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	req := domain.CreateStripeProductRequest{
		Name:        product.Name,
		Description: product.Description,
	}
	stripeProduct, err := s.stripe.CreateProduct(ctx, req)
	if err != nil {
		log.Printf("Service: Failed to create product/price in Stripe: %v", err)
		return "", errs.NewExternalServiceErr(err, "stripe product creation failed")
	}
	product.StripeProductID = stripeProduct.ID

	productID, err := s.repo.AddProduct(ctx, product)
	if err != nil {
		log.Printf("Service: Failed to create product in DB, attempting Stripe rollback: %v", err)
		if rollbackErr := s.stripe.DeactivateProduct(ctx, stripeProduct.ID); rollbackErr != nil {
			log.Printf("Service: Failed to rollback Stripe product %s. Data inconsistency detected! Rollback error: %v", stripeProduct.ID, rollbackErr)
		}
		return "", fmt.Errorf("failed to create product: %w", err)
	}
	return productID, nil
}
