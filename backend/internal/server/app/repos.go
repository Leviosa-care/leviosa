package app

import (
	"database/sql"

	"github.com/hengadev/leviosa/internal/domain/event"
	"github.com/hengadev/leviosa/internal/domain/message"
	"github.com/hengadev/leviosa/internal/domain/otp"
	"github.com/hengadev/leviosa/internal/domain/product"
	"github.com/hengadev/leviosa/internal/domain/register"
	"github.com/hengadev/leviosa/internal/domain/session"
	"github.com/hengadev/leviosa/internal/domain/settings"
	"github.com/hengadev/leviosa/internal/domain/throttler"
	"github.com/hengadev/leviosa/internal/domain/user"
	"github.com/hengadev/leviosa/internal/domain/vote"

	"github.com/redis/go-redis/v9"
)

type Repos struct {
	User        userService.Reader
	Session     sessionService.Reader
	Event       eventService.Reader
	Vote        vote.Reader
	Register    registerService.Reader
	Throttler   throttlerService.Reader
	Product     productService.Reader
	OTP         otpService.Reader
	Message     messageService.Reader
	Settings    settings.Reader
	SQLiteDB    *sql.DB
	RedisClient *redis.Client
}
