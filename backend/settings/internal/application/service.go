package settings

import (
	"github.com/Leviosa-care/settings/internal/ports"

	"github.com/hengadev/encx"
	amqp "github.com/rabbitmq/amqp091-go"
)

type SettingsService struct {
	repo   ports.SettingsRepository
	media  ports.SettingsMedia
	crypto encx.CryptoService
	mq     *amqp.Connection
}

func New(repo ports.SettingsRepository, media ports.SettingsMedia, crypto encx.CryptoService, rabbitMQConn *amqp.Connection) ports.SettingsService {
	return &SettingsService{
		repo:   repo,
		media:  media,
		crypto: crypto,
		mq:     rabbitMQConn,
	}
}
