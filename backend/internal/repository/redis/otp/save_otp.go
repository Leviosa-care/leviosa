package otpRepository

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	rp "github.com/hengadev/leviosa/internal/repository"

	"github.com/redis/go-redis/v9"
)

func (o *Repository) SaveOTP(ctx context.Context, emailHash string, otpEncoded []byte, duration time.Duration) error {
	key := formatOTPKey(emailHash)

	result := o.client.Set(ctx, key, otpEncoded, duration)
	if err := result.Err(); err != nil {
		switch {
		case errors.Is(err, redis.ErrClosed), errors.Is(err, &net.OpError{}):
			return rp.NewDatabaseErr(fmt.Errorf("redis connection error: %w", err))
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return rp.NewContextErr(err)
		default:
			return rp.NewDatabaseErr(fmt.Errorf("failed to store OTP: %w", err))
		}
	}
	if result.Val() == "" {
		return rp.NewNotCreatedErr(nil, "session")
	}

	return nil
}
