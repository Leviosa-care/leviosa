package mail

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Leviosa-care/notification/internal/domain"
	"github.com/Leviosa-care/notification/internal/ports"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	"github.com/Leviosa-care/leviosa/backend/internal/common/messaging/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

var _ ports.MailService = (*MailService)(nil)

type MailService struct {
	email    string
	password string
	// cache    *cache
	cache *domain.MailCache
}

func New(ctx context.Context, rabbitConn *amqp.Connection) (ports.MailService, error) {
	email := os.Getenv("GMAIL_EMAIL")
	if email == "" {
		return nil, fmt.Errorf("environment variable 'GMAIL_EMAIL' not found")
	}
	password := os.Getenv("GMAIL_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("environment variable 'GMAIL_PASSWORD' not found")
	}

	// Example: GET /settings/bulk?keys=company_email,company_logo
	keysList := []string{settings.CompanyEmail, settings.CompanyInstagram, settings.CompanyLegalAddress, settings.CompanyLogo}
	// TODO: use the bulk handler for this and find the best way to format the best url
	baseURL := "http://backend:3500"
	url := fmt.Sprintf("%s/settings/bulk?keys=%s", baseURL,
		strings.Join(keysList, ","))

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("settings service returned %d", resp.StatusCode)
	}

	var list []settings.SettingDTO
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, err
	}

	settingsMap := make(map[string]string)
	for _, s := range list {
		settingsMap[s.Key] = s.Value
	}

	cache := domain.NewMailCache(
		settingsMap[settings.CompanyEmail],
		settingsMap[settings.CompanyInstagram],
		settingsMap[settings.CompanyLegalAddress],
		[]byte(settingsMap[settings.CompanyLogo]),
	)

	mailService := &MailService{
		email:    email,
		password: password,
		cache:    cache,
	}
	ch, err := rabbitmq.NewChannel(rabbitConn)
	if err != nil {
		return nil, fmt.Errorf("consumer channel for mail service not created")
	}

	mailService.StartMailSettingConsumer(ctx, ch)

	return mailService, nil
}
