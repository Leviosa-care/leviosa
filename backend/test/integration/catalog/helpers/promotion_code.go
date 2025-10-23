package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// ClearPromotionCodesTable cleans the promotion codes table for isolation
func ClearPromotionCodesTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, `TRUNCATE catalog.promotion_codes CASCADE;`)
	require.NoError(t, err)
}

// NewValidPromotionCode creates a valid promotion code for tests
func NewValidPromotionCode(code string, couponID uuid.UUID) *domain.PromotionCode {
	maxRedemptions := 50

	return &domain.PromotionCode{
		ID:                   uuid.New(),
		StripePromotionID:    fmt.Sprintf("promo_%s", uuid.New().String()[:12]),
		CouponID:             couponID,
		Code:                 code,
		Active:               true,
		MaxRedemptions:       &maxRedemptions,
		TimesRedeemed:        0,
		FirstTimeTransaction: false,
		CreatedAt:            time.Now().UTC().Truncate(time.Microsecond),
		UpdatedAt:            time.Now().UTC().Truncate(time.Microsecond),
		Metadata:             map[string]string{"test": "true", "type": "standard"},
	}
}

// NewValidPromotionCodeWithExpiry creates a valid promotion code with expiry for tests
func NewValidPromotionCodeWithExpiry(code string, couponID uuid.UUID, expiresAt time.Time, maxRedemptions ...interface{}) *domain.PromotionCode {
	var maxRed *int
	var restrictions *domain.PromotionCodeRestrictions

	// Handle optional parameters
	if len(maxRedemptions) > 0 {
		if redCount, ok := maxRedemptions[0].(int); ok {
			maxRed = &redCount
		}
	}
	if len(maxRedemptions) > 1 {
		if currencies, ok := maxRedemptions[1].([]string); ok && len(currencies) > 0 {
			restrictions = &domain.PromotionCodeRestrictions{
				CurrencyOptions: currencies,
			}
		}
	}

	// Default max redemptions if not specified
	if maxRed == nil {
		defaultMax := 25
		maxRed = &defaultMax
	}

	return &domain.PromotionCode{
		ID:                   uuid.New(),
		StripePromotionID:    fmt.Sprintf("promo_%s", uuid.New().String()[:12]),
		CouponID:             couponID,
		Code:                 code,
		Active:               true,
		MaxRedemptions:       maxRed,
		TimesRedeemed:        0,
		ExpiresAt:            &expiresAt,
		FirstTimeTransaction: true,
		Restrictions:         restrictions,
		CreatedAt:            time.Now().UTC().Truncate(time.Microsecond),
		UpdatedAt:            time.Now().UTC().Truncate(time.Microsecond),
		Metadata:             map[string]string{"test": "true", "type": "expiring"},
	}
}

// NewValidPromotionCodeWithMinAmount creates a valid promotion code with minimum amount for tests
func NewValidPromotionCodeWithMinAmount(code string, couponID uuid.UUID, minAmount int, currency string) *domain.PromotionCode {
	return &domain.PromotionCode{
		ID:                    uuid.New(),
		StripePromotionID:     fmt.Sprintf("promo_%s", uuid.New().String()[:12]),
		CouponID:              couponID,
		Code:                  code,
		Active:                true,
		TimesRedeemed:         0,
		FirstTimeTransaction:  false,
		MinimumAmount:         &minAmount,
		MinimumAmountCurrency: &currency,
		CreatedAt:             time.Now().UTC().Truncate(time.Microsecond),
		UpdatedAt:             time.Now().UTC().Truncate(time.Microsecond),
		Metadata:              map[string]string{"test": "true", "type": "min_amount"},
	}
}

// NewValidPromotionCodeWithRestrictions creates a valid promotion code with currency restrictions for tests
func NewValidPromotionCodeWithRestrictions(code string, couponID uuid.UUID, currencies []string) *domain.PromotionCode {
	return &domain.PromotionCode{
		ID:                   uuid.New(),
		StripePromotionID:    fmt.Sprintf("promo_%s", uuid.New().String()[:12]),
		CouponID:             couponID,
		Code:                 code,
		Active:               true,
		TimesRedeemed:        0,
		FirstTimeTransaction: false,
		Restrictions: &domain.PromotionCodeRestrictions{
			CurrencyOptions: currencies,
		},
		CreatedAt: time.Now().UTC().Truncate(time.Microsecond),
		UpdatedAt: time.Now().UTC().Truncate(time.Microsecond),
		Metadata:  map[string]string{"test": "true", "type": "restricted"},
	}
}

// NewValidPromotionCodeWithRedemptionLimits creates a valid promotion code with specific redemption limit for tests
func NewValidPromotionCodeWithRedemptionLimits(code string, couponID uuid.UUID, maxRedemptions int) *domain.PromotionCode {
	return &domain.PromotionCode{
		ID:                   uuid.New(),
		StripePromotionID:    fmt.Sprintf("promo_%s", uuid.New().String()[:12]),
		CouponID:             couponID,
		Code:                 code,
		Active:               true,
		MaxRedemptions:       &maxRedemptions,
		TimesRedeemed:        0,
		FirstTimeTransaction: false,
		CreatedAt:            time.Now().UTC().Truncate(time.Microsecond),
		UpdatedAt:            time.Now().UTC().Truncate(time.Microsecond),
		Metadata:             map[string]string{"test": "true", "type": "limited"},
	}
}

// NewExpiredPromotionCode creates an expired promotion code for tests
func NewExpiredPromotionCode(code string, couponID uuid.UUID) *domain.PromotionCode {
	maxRedemptions := 10
	expiresAt := time.Now().Add(-24 * time.Hour) // 1 day ago

	return &domain.PromotionCode{
		ID:                   uuid.New(),
		StripePromotionID:    fmt.Sprintf("promo_%s", uuid.New().String()[:12]),
		CouponID:             couponID,
		Code:                 code,
		Active:               true,
		MaxRedemptions:       &maxRedemptions,
		TimesRedeemed:        0,
		ExpiresAt:            &expiresAt,
		FirstTimeTransaction: false,
		CreatedAt:            time.Now().UTC().Truncate(time.Microsecond),
		UpdatedAt:            time.Now().UTC().Truncate(time.Microsecond),
		Metadata:             map[string]string{"test": "true", "type": "expired"},
	}
}

// NewInactivePromotionCode creates an inactive promotion code for tests
func NewInactivePromotionCode(code string, couponID uuid.UUID) *domain.PromotionCode {
	return &domain.PromotionCode{
		ID:                   uuid.New(),
		StripePromotionID:    fmt.Sprintf("promo_%s", uuid.New().String()[:12]),
		CouponID:             couponID,
		Code:                 code,
		Active:               false,
		TimesRedeemed:        0,
		FirstTimeTransaction: false,
		CreatedAt:            time.Now().UTC().Truncate(time.Microsecond),
		UpdatedAt:            time.Now().UTC().Truncate(time.Microsecond),
		Metadata:             map[string]string{"test": "true", "type": "inactive"},
	}
}

// Helper function to insert a promotion code for test setup
func InsertPromotionCode(t *testing.T, ctx context.Context, pool *pgxpool.Pool, promotionCode *domain.PromotionCode) {
	t.Helper()
	query := `
	INSERT INTO catalog.promotion_codes (
		id, stripe_promotion_id, coupon_id, code, active, 
		max_redemptions, times_redeemed, expires_at, first_time_transaction,
		minimum_amount, minimum_amount_currency, restrictions, 
		metadata, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
	)`

	// Handle metadata marshaling
	var metadataJSON []byte
	if promotionCode.Metadata != nil {
		var err error
		metadataJSON, err = json.Marshal(promotionCode.Metadata)
		require.NoError(t, err, "Failed to marshal metadata for promotion code insertion")
	}

	// Marshal restrictions to JSON
	var restrictionsJSON []byte
	if promotionCode.Restrictions != nil {
		var err error
		restrictionsJSON, err = json.Marshal(promotionCode.Restrictions)
		require.NoError(t, err, "Failed to marshal restrictions for promotion code insertion")
	}

	_, err := pool.Exec(ctx, query,
		promotionCode.ID,
		promotionCode.StripePromotionID,
		promotionCode.CouponID,
		promotionCode.Code,
		promotionCode.Active,
		promotionCode.MaxRedemptions,
		promotionCode.TimesRedeemed,
		promotionCode.ExpiresAt,
		promotionCode.FirstTimeTransaction,
		promotionCode.MinimumAmount,
		promotionCode.MinimumAmountCurrency,
		restrictionsJSON,
		metadataJSON,
		promotionCode.CreatedAt,
		promotionCode.UpdatedAt,
	)
	require.NoError(t, err, fmt.Sprintf("Failed to pre-insert promotion code '%s'", promotionCode.Code))
}

// GetPromotionCodeByID fetches a promotion code from the database by its ID.
// This is a test helper function used to verify the state of the database.
func GetPromotionCodeByID(t *testing.T, ctx context.Context, promotionCodeID uuid.UUID, pool *pgxpool.Pool) (*domain.PromotionCode, error) {
	t.Helper()

	query := `
		SELECT id, stripe_promotion_id, coupon_id, code, active, 
		       max_redemptions, times_redeemed, expires_at, first_time_transaction,
		       minimum_amount, minimum_amount_currency, restrictions, 
		       metadata, created_at, updated_at
		FROM catalog.promotion_codes 
		WHERE id = $1
	`

	var promotionCode domain.PromotionCode
	var metadataJSON []byte
	var restrictionsJSON []byte

	err := pool.QueryRow(ctx, query, promotionCodeID).Scan(
		&promotionCode.ID,
		&promotionCode.StripePromotionID,
		&promotionCode.CouponID,
		&promotionCode.Code,
		&promotionCode.Active,
		&promotionCode.MaxRedemptions,
		&promotionCode.TimesRedeemed,
		&promotionCode.ExpiresAt,
		&promotionCode.FirstTimeTransaction,
		&promotionCode.MinimumAmount,
		&promotionCode.MinimumAmountCurrency,
		&restrictionsJSON,
		&metadataJSON,
		&promotionCode.CreatedAt,
		&promotionCode.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errs.NewRepositoryNotFoundErr(err, "promotion code")
		}
		return nil, errs.ClassifyPgError("get promotion code by ID", err)
	}

	// Parse restrictions JSON if present
	if len(restrictionsJSON) > 0 {
		err = json.Unmarshal(restrictionsJSON, &promotionCode.Restrictions)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal restrictions from database: %w", err)
		}
		require.NoError(t, err, "Failed to unmarshal restrictions from database")
	}

	// Unmarshal metadata if present
	if metadataJSON != nil {
		err = json.Unmarshal(metadataJSON, &promotionCode.Metadata)
		require.NoError(t, err, "Failed to unmarshal metadata from database")
	}

	return &promotionCode, nil
}

// GetPromotionCodeByIDOrNil is a helper function that returns nil if the promotion code is not found instead of an error
func GetPromotionCodeByIDOrNil(t *testing.T, ctx context.Context, promotionCodeID uuid.UUID, pool *pgxpool.Pool) *domain.PromotionCode {
	t.Helper()
	promotionCode, err := GetPromotionCodeByID(t, ctx, promotionCodeID, pool)
	if err != nil {
		return nil
	}
	return promotionCode
}

// GetPromotionCodeByCode fetches a promotion code from the database by its code.
func GetPromotionCodeByCode(t *testing.T, ctx context.Context, code string, pool *pgxpool.Pool) (*domain.PromotionCode, error) {
	t.Helper()

	query := `
		SELECT id, stripe_promotion_id, coupon_id, code, active, 
		       max_redemptions, times_redeemed, expires_at, first_time_transaction,
		       minimum_amount, minimum_amount_currency, restrictions, 
		       metadata, created_at, updated_at
		FROM catalog.promotion_codes 
		WHERE code = $1
	`

	var promotionCode domain.PromotionCode
	var metadataJSON []byte
	var restrictionsJSON []byte

	err := pool.QueryRow(ctx, query, code).Scan(
		&promotionCode.ID,
		&promotionCode.StripePromotionID,
		&promotionCode.CouponID,
		&promotionCode.Code,
		&promotionCode.Active,
		&promotionCode.MaxRedemptions,
		&promotionCode.TimesRedeemed,
		&promotionCode.ExpiresAt,
		&promotionCode.FirstTimeTransaction,
		&promotionCode.MinimumAmount,
		&promotionCode.MinimumAmountCurrency,
		&restrictionsJSON,
		&metadataJSON,
		&promotionCode.CreatedAt,
		&promotionCode.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errs.NewRepositoryNotFoundErr(err, "promotion code")
		}
		return nil, errs.ClassifyPgError("get promotion code by code", err)
	}

	// Parse restrictions JSON if present
	if len(restrictionsJSON) > 0 {
		err = json.Unmarshal(restrictionsJSON, &promotionCode.Restrictions)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal restrictions from database: %w", err)
		}
	}

	// Unmarshal metadata if present
	if metadataJSON != nil {
		err = json.Unmarshal(metadataJSON, &promotionCode.Metadata)
		require.NoError(t, err, "Failed to unmarshal metadata from database")
	}

	return &promotionCode, nil
}

// GetPromotionCodeActiveStatus is a helper function to query the database directly and retrieve
// the active status of a given promotion code ID.
func GetPromotionCodeActiveStatus(t *testing.T, ctx context.Context, promotionCodeID uuid.UUID, pool *pgxpool.Pool) bool {
	t.Helper()
	var active bool
	err := pool.QueryRow(ctx, "SELECT active FROM catalog.promotion_codes WHERE id = $1", promotionCodeID).Scan(&active)
	require.NoError(t, err, "Failed to get promotion code active status for ID %s", promotionCodeID)
	return active
}

// GetPromotionCodeTimesRedeemed is a helper function to get the current redemption count
func GetPromotionCodeTimesRedeemed(t *testing.T, ctx context.Context, promotionCodeID uuid.UUID, pool *pgxpool.Pool) int {
	t.Helper()
	var timesRedeemed int
	err := pool.QueryRow(ctx, "SELECT times_redeemed FROM catalog.promotion_codes WHERE id = $1", promotionCodeID).Scan(&timesRedeemed)
	require.NoError(t, err, "Failed to get promotion code times redeemed for ID %s", promotionCodeID)
	return timesRedeemed
}

// UpdatePromotionCodeTimesRedeemed is a helper function to update the times_redeemed field for a promotion code
func UpdatePromotionCodeTimesRedeemed(t *testing.T, ctx context.Context, promotionCodeID uuid.UUID, timesRedeemed int, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, "UPDATE catalog.promotion_codes SET times_redeemed = $1, updated_at = NOW() WHERE id = $2", timesRedeemed, promotionCodeID)
	require.NoError(t, err, "Failed to update promotion code times redeemed for ID %s", promotionCodeID)
}

// GetPromotionCodesByCouponID is a helper function to get all promotion codes for a specific coupon
func GetPromotionCodesByCouponID(t *testing.T, ctx context.Context, couponID uuid.UUID, pool *pgxpool.Pool) []*domain.PromotionCode {
	t.Helper()

	query := `
		SELECT id, stripe_promotion_id, coupon_id, code, active, 
		       max_redemptions, times_redeemed, expires_at, first_time_transaction,
		       minimum_amount, minimum_amount_currency, restrictions, 
		       metadata, created_at, updated_at
		FROM catalog.promotion_codes 
		WHERE coupon_id = $1
		ORDER BY created_at DESC
	`

	rows, err := pool.Query(ctx, query, couponID)
	require.NoError(t, err, "Failed to query promotion codes by coupon ID")
	defer rows.Close()

	var promotionCodes []*domain.PromotionCode
	for rows.Next() {
		var promotionCode domain.PromotionCode
		var metadataJSON []byte
		var restrictionsJSON []byte

		err := rows.Scan(
			&promotionCode.ID,
			&promotionCode.StripePromotionID,
			&promotionCode.CouponID,
			&promotionCode.Code,
			&promotionCode.Active,
			&promotionCode.MaxRedemptions,
			&promotionCode.TimesRedeemed,
			&promotionCode.ExpiresAt,
			&promotionCode.FirstTimeTransaction,
			&promotionCode.MinimumAmount,
			&promotionCode.MinimumAmountCurrency,
			&restrictionsJSON,
			&metadataJSON,
			&promotionCode.CreatedAt,
			&promotionCode.UpdatedAt,
		)
		require.NoError(t, err, "Failed to scan promotion code")

		// Parse restrictions JSON if present
		if len(restrictionsJSON) > 0 {
			err = json.Unmarshal(restrictionsJSON, &promotionCode.Restrictions)
			require.NoError(t, err, "Failed to unmarshal restrictions from database")
		}

		// Unmarshal metadata if present
		if metadataJSON != nil {
			err = json.Unmarshal(metadataJSON, &promotionCode.Metadata)
			require.NoError(t, err, "Failed to unmarshal metadata from database")
		}

		promotionCodes = append(promotionCodes, &promotionCode)
	}

	require.NoError(t, rows.Err(), "Error iterating over promotion code rows")
	return promotionCodes
}
