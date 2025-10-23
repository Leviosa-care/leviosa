package smtp

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/notification/internal/domain"
	"github.com/Leviosa-care/notification/internal/ports"

	"gopkg.in/gomail.v2"
)

type SMTPClient struct {
	username string
	password string
	host     string
	port     int
	cache    *domain.CompanyCache
}

func New(username, password, host string, port int, cache *domain.CompanyCache) ports.EmailService {
	return &SMTPClient{
		username: username,
		password: password,
		host:     host,
		port:     port,
		cache:    cache,
	}
}

func (c *SMTPClient) SendEmail(ctx context.Context, request *domain.EmailRequest) error {
	if request == nil {
		return errs.NewInvalidValueErr(fmt.Errorf("email request cannot be nil"))
	}

	if request.To == "" {
		return errs.NewInvalidValueErr(fmt.Errorf("recipient email cannot be empty"))
	}

	if request.Subject == "" {
		return errs.NewInvalidValueErr(fmt.Errorf("email subject cannot be empty"))
	}

	m := gomail.NewMessage()
	m.SetHeader("From", c.cache.GetCompanyEmail())
	m.SetHeader("To", request.To)
	m.SetHeader("Subject", request.Subject)

	if len(request.CarbonCopy) > 0 {
		addresses := make([]string, 0, len(request.CarbonCopy))
		for email, name := range request.CarbonCopy {
			addresses = append(addresses, m.FormatAddress(email, name))
		}
		m.SetHeader("Cc", addresses...)
	}

	logo := c.cache.GetLogo()
	if len(logo) > 0 {
		logoPath, err := writeTempFile(logo, "logo.jpg")
		if err != nil {
			return fmt.Errorf("write logo temp file: %w", err)
		}
		defer os.Remove(logoPath)
		m.Embed(logoPath, gomail.Rename("logo.jpg"))
	}

	instagramPath, err := writeTempFile(domain.InstagramImage, "instagram.png")
	if err != nil {
		return fmt.Errorf("write instagram temp file: %w", err)
	}
	defer os.Remove(instagramPath)
	m.Embed(instagramPath, gomail.Rename("instagram.png"))

	for path, rename := range request.Images {
		m.Embed(path, gomail.Rename(rename))
	}

	if request.Template != "" {
		htmlBody, err := c.renderTemplate(request.Template, request.Data)
		if err != nil {
			return fmt.Errorf("render email template: %w", err)
		}
		m.SetBody("text/html", htmlBody)
	}

	d := gomail.NewDialer(c.host, c.port, c.username, c.password)

	if err := d.DialAndSend(m); err != nil {
		return errs.NewConnectionFailureErr(fmt.Errorf("send email via SMTP: %w", err))
	}

	return nil
}

func (c *SMTPClient) renderTemplate(templateName string, data interface{}) (string, error) {
	templatePath := fmt.Sprintf("templates/%s.html", templateName)

	tmpl, err := template.ParseFS(domain.EmailTemplates, templatePath)
	if err != nil {
		return "", errs.NewInvalidValueErr(fmt.Errorf("parse email template %s: %w", templateName, err))
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", errs.NewInvalidValueErr(fmt.Errorf("execute email template %s: %w", templateName, err))
	}

	return buf.String(), nil
}

func writeTempFile(data []byte, filename string) (string, error) {
	if len(data) == 0 {
		return "", errs.NewInvalidValueErr(fmt.Errorf("file data cannot be empty"))
	}

	tmpDir := os.TempDir()
	tmpPath := filepath.Join(tmpDir, filename)

	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return "", errs.NewInternalErr(fmt.Errorf("write temp file %s: %w", filename, err))
	}

	return tmpPath, nil
}
