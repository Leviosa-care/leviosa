package otpRepository

import (
	"context"
	"errors"
	"fmt"

	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/redis/go-redis/v9"
)

func (o *Repository) GetOTP(ctx context.Context, emailHash string) ([]byte, error) {
	key := formatOTPKey(emailHash)
	otpEncoded, err := o.client.Get(ctx, key).Bytes()
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			return nil, rp.NewNotFoundErr(err, "OTP with specified hash email")
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return nil, rp.NewContextErr(err)
		default:
			return nil, rp.NewDatabaseErr(fmt.Errorf("failed to check existing OTP: %w", err))
		}
	}

	return otpEncoded, nil
}
