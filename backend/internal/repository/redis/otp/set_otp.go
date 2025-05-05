package otpRepository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/hengadev/leviosa/internal/domain/otp"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/redis/go-redis/v9"
)

func (o *Repository) setOTP(ctx context.Context, key string, data *otpService.OTP) error {
	encoded, err := json.Marshal(data)
	if err != nil {
		return rp.NewDatabaseErr(fmt.Errorf("failed to encode OTP data: %w", err))
	}

	// Use pipelines for atomic operations
	pipe := o.client.Pipeline()
	pipe.Set(ctx, key, encoded, otpService.OTPDURATION)

	// Execute pipeline
	if _, err := pipe.Exec(ctx); err != nil {
		switch {
		case errors.Is(err, redis.ErrClosed), errors.Is(err, &net.OpError{}):
			return rp.NewDatabaseErr(fmt.Errorf("redis connection error: %w", err))
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return rp.NewContextErr(err)
		default:
			return rp.NewDatabaseErr(fmt.Errorf("failed to store OTP: %w", err))
		}
	}

	return nil
}
