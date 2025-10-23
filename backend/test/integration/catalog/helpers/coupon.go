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

// ClearCouponsTable cleans the coupons table for isolation
func ClearCouponsTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, `TRUNCATE catalog.coupons CASCADE;`)
	require.NoError(t, err)
}

// NewValidPercentOffCoupon creates a valid percent-off coupon for tests
func NewValidPercentOffCoupon(name string) *domain.Coupon {
	percentOff := 25.0
	maxRedemptions := 100

	return &domain.Coupon{
		ID:             uuid.New(),
		StripeCouponID: fmt.Sprintf("coupon_%s", uuid.New().String()[:12]),
		Name:           name,
		PercentOff:     &percentOff,
		Duration:       domain.CouponDurationOnce,
		MaxRedemptions: &maxRedemptions,
		TimesRedeemed:  0,
		IsValid:        true,
		CreatedAt:      time.Now().UTC().Truncate(time.Microsecond),
		UpdatedAt:      time.Now().UTC().Truncate(time.Microsecond),
		Metadata:       map[string]string{"test": "true", "type": "percent"},
	}
}

// NewValidAmountOffCoupon creates a valid amount-off coupon for tests
func NewValidAmountOffCoupon(name string, currency string) *domain.Coupon {
	amountOff := 500 // $5.00
	maxRedemptions := 50

	return &domain.Coupon{
		ID:             uuid.New(),
		StripeCouponID: fmt.Sprintf("coupon_%s", uuid.New().String()[:12]),
		Name:           name,
		AmountOff:      &amountOff,
		Currency:       &currency,
		Duration:       domain.CouponDurationOnce,
		MaxRedemptions: &maxRedemptions,
		TimesRedeemed:  0,
		IsValid:        true,
		CreatedAt:      time.Now().UTC().Truncate(time.Microsecond),
		UpdatedAt:      time.Now().UTC().Truncate(time.Microsecond),
		Metadata:       map[string]string{"test": "true", "type": "amount"},
	}
}

// NewValidRepeatingCoupon creates a valid repeating coupon for tests
func NewValidRepeatingCoupon(name string, durationMonths int) *domain.Coupon {
	percentOff := 15.0

	return &domain.Coupon{
		ID:               uuid.New(),
		StripeCouponID:   fmt.Sprintf("coupon_%s", uuid.New().String()[:12]),
		Name:             name,
		PercentOff:       &percentOff,
		Duration:         domain.CouponDurationRepeating,
		DurationInMonths: &durationMonths,
		TimesRedeemed:    0,
		IsValid:          true,
		CreatedAt:        time.Now().UTC().Truncate(time.Microsecond),
		UpdatedAt:        time.Now().UTC().Truncate(time.Microsecond),
		Metadata:         map[string]string{"test": "true", "type": "repeating"},
	}
}

// NewValidForeverCoupon creates a valid forever coupon for tests
func NewValidForeverCoupon(name string) *domain.Coupon {
	percentOff := 10.0
	maxRedemptions := 1000

	return &domain.Coupon{
		ID:             uuid.New(),
		StripeCouponID: fmt.Sprintf("coupon_%s", uuid.New().String()[:12]),
		Name:           name,
		PercentOff:     &percentOff,
		Duration:       domain.CouponDurationForever,
		MaxRedemptions: &maxRedemptions,
		TimesRedeemed:  0,
		IsValid:        true,
		CreatedAt:      time.Now().UTC().Truncate(time.Microsecond),
		UpdatedAt:      time.Now().UTC().Truncate(time.Microsecond),
		Metadata:       map[string]string{"test": "true", "type": "forever"},
	}
}

// NewExpiredCoupon creates an expired coupon for tests
func NewExpiredCoupon(name string) *domain.Coupon {
	percentOff := 20.0
	redeemBy := time.Now().Add(-24 * time.Hour) // 1 day ago

	return &domain.Coupon{
		ID:             uuid.New(),
		StripeCouponID: fmt.Sprintf("coupon_%s", uuid.New().String()[:12]),
		Name:           name,
		PercentOff:     &percentOff,
		Duration:       domain.CouponDurationOnce,
		TimesRedeemed:  0,
		IsValid:        true,
		RedeemBy:       &redeemBy,
		CreatedAt:      time.Now().UTC().Truncate(time.Microsecond),
		UpdatedAt:      time.Now().UTC().Truncate(time.Microsecond),
		Metadata:       map[string]string{"test": "true", "type": "expired"},
	}
}

// Helper function to insert a coupon for test setup
func InsertCoupon(t *testing.T, ctx context.Context, pool *pgxpool.Pool, coupon *domain.Coupon) {
	t.Helper()
	query := `
	INSERT INTO catalog.coupons (
		id, stripe_coupon_id, name, percent_off, amount_off, currency, 
		duration, duration_in_months, max_redemptions, times_redeemed, 
		valid, redeem_by, metadata, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
	)`

	// Handle metadata marshaling
	var metadataJSON []byte
	if coupon.Metadata != nil {
		var err error
		metadataJSON, err = json.Marshal(coupon.Metadata)
		require.NoError(t, err, "Failed to marshal metadata for coupon insertion")
	}

	_, err := pool.Exec(ctx, query,
		coupon.ID,
		coupon.StripeCouponID,
		coupon.Name,
		coupon.PercentOff,
		coupon.AmountOff,
		coupon.Currency,
		coupon.Duration,
		coupon.DurationInMonths,
		coupon.MaxRedemptions,
		coupon.TimesRedeemed,
		coupon.IsValid,
		coupon.RedeemBy,
		metadataJSON,
		coupon.CreatedAt,
		coupon.UpdatedAt,
	)
	require.NoError(t, err, fmt.Sprintf("Failed to pre-insert coupon '%s'", coupon.Name))
}

// GetCouponByID fetches a coupon from the database by its ID.
// This is a test helper function used to verify the state of the database.
func GetCouponByID(t *testing.T, ctx context.Context, couponID uuid.UUID, pool *pgxpool.Pool) (*domain.Coupon, error) {
	t.Helper()

	query := `
		SELECT id, stripe_coupon_id, name, percent_off, amount_off, currency, 
		       duration, duration_in_months, max_redemptions, times_redeemed, 
		       valid, redeem_by, metadata, created_at, updated_at
		FROM catalog.coupons 
		WHERE id = $1
	`

	var coupon domain.Coupon
	var metadataJSON []byte

	err := pool.QueryRow(ctx, query, couponID).Scan(
		&coupon.ID,
		&coupon.StripeCouponID,
		&coupon.Name,
		&coupon.PercentOff,
		&coupon.AmountOff,
		&coupon.Currency,
		&coupon.Duration,
		&coupon.DurationInMonths,
		&coupon.MaxRedemptions,
		&coupon.TimesRedeemed,
		&coupon.IsValid,
		&coupon.RedeemBy,
		&metadataJSON,
		&coupon.CreatedAt,
		&coupon.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errs.NewRepositoryNotFoundErr(err, "coupon")
		}
		return nil, errs.ClassifyPgError("get coupon by ID", err)
	}

	// Unmarshal metadata if present
	if metadataJSON != nil {
		err = json.Unmarshal(metadataJSON, &coupon.Metadata)
		require.NoError(t, err, "Failed to unmarshal metadata from database")
	}

	return &coupon, nil
}

// GetCouponByStripeID fetches a coupon from the database by its Stripe ID.
func GetCouponByStripeID(t *testing.T, ctx context.Context, stripeCouponID string, pool *pgxpool.Pool) (*domain.Coupon, error) {
	t.Helper()

	query := `
		SELECT id, stripe_coupon_id, name, percent_off, amount_off, currency, 
		       duration, duration_in_months, max_redemptions, times_redeemed, 
		       valid, redeem_by, metadata, created_at, updated_at
		FROM catalog.coupons 
		WHERE stripe_coupon_id = $1
	`

	var coupon domain.Coupon
	var metadataJSON []byte

	err := pool.QueryRow(ctx, query, stripeCouponID).Scan(
		&coupon.ID,
		&coupon.StripeCouponID,
		&coupon.Name,
		&coupon.PercentOff,
		&coupon.AmountOff,
		&coupon.Currency,
		&coupon.Duration,
		&coupon.DurationInMonths,
		&coupon.MaxRedemptions,
		&coupon.TimesRedeemed,
		&coupon.IsValid,
		&coupon.RedeemBy,
		&metadataJSON,
		&coupon.CreatedAt,
		&coupon.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errs.NewRepositoryNotFoundErr(err, "coupon")
		}
		return nil, errs.ClassifyPgError("get coupon by Stripe ID", err)
	}

	// Unmarshal metadata if present
	if metadataJSON != nil {
		err = json.Unmarshal(metadataJSON, &coupon.Metadata)
		require.NoError(t, err, "Failed to unmarshal metadata from database")
	}

	return &coupon, nil
}

// GetCouponValidStatus is a helper function to query the database directly and retrieve
// the valid status of a given coupon ID.
func GetCouponValidStatus(t *testing.T, ctx context.Context, couponID uuid.UUID, pool *pgxpool.Pool) bool {
	t.Helper()
	var isValid bool
	err := pool.QueryRow(ctx, "SELECT valid FROM catalog.coupons WHERE id = $1", couponID).Scan(&isValid)
	require.NoError(t, err, "Failed to get coupon valid status for ID %s", couponID)
	return isValid
}

// GetCouponTimesRedeemed is a helper function to get the current redemption count
func GetCouponTimesRedeemed(t *testing.T, ctx context.Context, couponID uuid.UUID, pool *pgxpool.Pool) int {
	t.Helper()
	var timesRedeemed int
	err := pool.QueryRow(ctx, "SELECT times_redeemed FROM catalog.coupons WHERE id = $1", couponID).Scan(&timesRedeemed)
	require.NoError(t, err, "Failed to get coupon times redeemed for ID %s", couponID)
	return timesRedeemed
}

// NewInvalidCoupon creates an invalid coupon for tests
func NewInvalidCoupon(name string) *domain.Coupon {
	percentOff := 25.0

	return &domain.Coupon{
		ID:             uuid.New(),
		StripeCouponID: fmt.Sprintf("coupon_%s", uuid.New().String()[:12]),
		Name:           name,
		PercentOff:     &percentOff,
		Duration:       domain.CouponDurationOnce,
		TimesRedeemed:  0,
		IsValid:        false, // This makes it invalid
		CreatedAt:      time.Now().UTC().Truncate(time.Microsecond),
		UpdatedAt:      time.Now().UTC().Truncate(time.Microsecond),
		Metadata:       map[string]string{"test": "true", "type": "invalid"},
	}
}

// GetCouponByIDOrNil is a helper function that returns nil if the coupon is not found instead of an error
func GetCouponByIDOrNil(t *testing.T, ctx context.Context, couponID uuid.UUID, pool *pgxpool.Pool) *domain.Coupon {
	t.Helper()
	coupon, err := GetCouponByID(t, ctx, couponID, pool)
	if err != nil {
		return nil
	}
	return coupon
}

// NewValidCouponWithRedeemBy creates a valid coupon with specific redeem by date for tests
func NewValidCouponWithRedeemBy(name string, redeemBy time.Time) *domain.Coupon {
	percentOff := 20.0
	maxRedemptions := 100

	return &domain.Coupon{
		ID:             uuid.New(),
		StripeCouponID: fmt.Sprintf("coupon_%s", uuid.New().String()[:12]),
		Name:           name,
		PercentOff:     &percentOff,
		Duration:       domain.CouponDurationOnce,
		MaxRedemptions: &maxRedemptions,
		TimesRedeemed:  0,
		IsValid:        true,
		RedeemBy:       &redeemBy,
		CreatedAt:      time.Now().UTC().Truncate(time.Microsecond),
		UpdatedAt:      time.Now().UTC().Truncate(time.Microsecond),
		Metadata:       map[string]string{"test": "true", "type": "with_redeem_by"},
	}
}

// NewValidPercentOffCouponWithRedemptionLimits creates a valid percent-off coupon with specific redemption limits for tests
func NewValidPercentOffCouponWithRedemptionLimits(name string, maxRedemptions int) *domain.Coupon {
	percentOff := 25.0

	return &domain.Coupon{
		ID:             uuid.New(),
		StripeCouponID: fmt.Sprintf("coupon_%s", uuid.New().String()[:12]),
		Name:           name,
		PercentOff:     &percentOff,
		Duration:       domain.CouponDurationOnce,
		MaxRedemptions: &maxRedemptions,
		TimesRedeemed:  0,
		IsValid:        true,
		CreatedAt:      time.Now().UTC().Truncate(time.Microsecond),
		UpdatedAt:      time.Now().UTC().Truncate(time.Microsecond),
		Metadata:       map[string]string{"test": "true", "type": "limited_redemptions"},
	}
}
