package settings

import (
	"github.com/Leviosa-care/leviosa/backend/internal/settings/ports"

	"github.com/hengadev/encx"
	// amqp "github.com/rabbitmq/amqp091-go" // Commented out: RabbitMQ integration disabled in favor of EventPublisher interface
)

type SettingsService struct {
	repo      ports.SettingsRepository
	media     ports.SettingsMedia
	crypto    encx.CryptoService
	publisher ports.EventPublisher
	// mq     *amqp.Connection // Replaced with publisher interface
}

func New(repo ports.SettingsRepository, media ports.SettingsMedia, crypto encx.CryptoService, publisher ports.EventPublisher) ports.SettingsService {
	return &SettingsService{
		repo:      repo,
		media:     media,
		crypto:    crypto,
		publisher: publisher,
		// mq:     rabbitMQConn, // Replaced with publisher interface
	}
}
