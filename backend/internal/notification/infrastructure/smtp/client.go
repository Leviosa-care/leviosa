package smtp

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/notification/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/notification/ports"

	"gopkg.in/gomail.v2"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type SMTPClient struct {
	config SMTPConfig
}

func NewSMTPClient(config SMTPConfig) ports.EmailService {
	return &SMTPClient{
		config: config,
	}
}

func (c *SMTPClient) SendOTPEmail(ctx context.Context, req domain.OTPEmailRequest) error {
	data := struct {
		OTP     string
		Year    int
		LogoURL string
	}{
		OTP:     req.OTP,
		Year:    time.Now().Year(),
		LogoURL: req.LogoURL,
	}

	emailReq := &domain.EmailRequest{
		To:       req.ToEmail,
		Subject:  "Your OTP Code",
		Template: "otp",
		Data:     data,
	}

	return c.sendEmail(ctx, emailReq, req.FromEmail)
}

func (c *SMTPClient) SendWelcomeEmail(ctx context.Context, req domain.WelcomeEmailRequest) error {
	data := struct {
		FirstName   string
		LastName    string
		CompanyName string
		Year        int
		LogoURL     string
	}{
		FirstName:   req.ToFirstName,
		LastName:    req.ToLastName,
		CompanyName: req.CompanyName,
		Year:        time.Now().Year(),
		LogoURL:     req.LogoURL,
	}

	emailReq := &domain.EmailRequest{
		To:       req.ToEmail,
		Subject:  fmt.Sprintf("Welcome to %s", req.CompanyName),
		Template: "welcome",
		Data:     data,
	}

	return c.sendEmail(ctx, emailReq, req.FromEmail)
}

func (c *SMTPClient) SendVerifyEmailEmail(ctx context.Context, req domain.VerifyEmailRequest) error {
	data := struct {
		FirstName   string
		LastName    string
		CompanyName string
		Year        int
		LogoURL     string
	}{
		FirstName:   req.ToFirstName,
		LastName:    req.ToLastName,
		CompanyName: req.CompanyName,
		Year:        time.Now().Year(),
		LogoURL:     req.LogoURL,
	}

	emailReq := &domain.EmailRequest{
		To:       req.ToEmail,
		Subject:  "Verify Your Email",
		Template: "verify_email",
		Data:     data,
	}

	return c.sendEmail(ctx, emailReq, req.FromEmail)
}

func (c *SMTPClient) SendEventNotificationEmail(ctx context.Context, req domain.EventNotificationRequest) error {
	data := struct {
		FirstName   string
		LastName    string
		Event       string
		Details     string
		CompanyName string
		Year        int
		LogoURL     string
	}{
		FirstName:   req.ToFirstName,
		LastName:    req.ToLastName,
		Event:       req.Event,
		Details:     req.Details,
		CompanyName: req.CompanyName,
		Year:        time.Now().Year(),
		LogoURL:     req.LogoURL,
	}

	emailReq := &domain.EmailRequest{
		To:       req.ToEmail,
		Subject:  fmt.Sprintf("Event Notification: %s", req.Event),
		Template: "event_notification",
		Data:     data,
	}

	return c.sendEmail(ctx, emailReq, req.FromEmail)
}

func (c *SMTPClient) SendPaymentNotificationEmail(ctx context.Context, req domain.PaymentNotificationRequest) error {
	data := struct {
		FirstName   string
		LastName    string
		Amount      string
		Product     string
		PaymentDate string
		CompanyName string
		Year        int
		LogoURL     string
	}{
		FirstName:   req.ToFirstName,
		LastName:    req.ToLastName,
		Amount:      req.Amount,
		Product:     req.Product,
		PaymentDate: req.PaymentDate,
		CompanyName: req.CompanyName,
		Year:        time.Now().Year(),
		LogoURL:     req.LogoURL,
	}

	emailReq := &domain.EmailRequest{
		To:       req.ToEmail,
		Subject:  "Payment Confirmation",
		Template: "payment",
		Data:     data,
	}

	return c.sendEmail(ctx, emailReq, req.FromEmail)
}

func (c *SMTPClient) SendPaymentFailedEmail(ctx context.Context, req domain.PaymentNotificationRequest) error {
	data := struct {
		FirstName   string
		LastName    string
		Amount      string
		Product     string
		PaymentDate string
		CompanyName string
		Year        int
		LogoURL     string
	}{
		FirstName:   req.ToFirstName,
		LastName:    req.ToLastName,
		Amount:      req.Amount,
		Product:     req.Product,
		PaymentDate: req.PaymentDate,
		CompanyName: req.CompanyName,
		Year:        time.Now().Year(),
		LogoURL:     req.LogoURL,
	}

	emailReq := &domain.EmailRequest{
		To:       req.ToEmail,
		Subject:  "Payment Failed",
		Template: "payment_failed",
		Data:     data,
	}

	return c.sendEmail(ctx, emailReq, req.FromEmail)
}

func (c *SMTPClient) SendBookingConfirmationEmail(ctx context.Context, req domain.BookingConfirmationRequest) error {
	data := struct {
		FirstName       string
		LastName        string
		BookingID       string
		ProductName     string
		RoomName        string
		Building        string
		Address         string
		Date            string
		Time            string
		PartnerName     string
		Amount          string
		Year            int
		BookingTokenURL string
	}{
		FirstName:       req.ToFirstName,
		LastName:        req.ToLastName,
		BookingID:       req.BookingID,
		ProductName:     req.ProductName,
		RoomName:        req.RoomName,
		Building:        req.Building,
		Address:         req.Address,
		Date:            req.Date,
		Time:            req.Time,
		PartnerName:     req.PartnerName,
		Amount:          req.Amount,
		Year:            req.Year,
		BookingTokenURL: req.BookingTokenURL,
	}

	emailReq := &domain.EmailRequest{
		To:       req.ToEmail,
		Subject:  "Confirmation de réservation",
		Template: "booking_confirmation",
		Data:     data,
	}

	return c.sendEmail(ctx, emailReq, "") // from email handled by SMTP config
}

func (c *SMTPClient) SendBookingCancellationEmail(ctx context.Context, req domain.BookingCancellationRequest) error {
	data := struct {
		FirstName   string
		LastName    string
		BookingID   string
		ProductName string
		RoomName    string
		Date        string
		Time        string
		Reason      string
		Year        int
	}{
		FirstName:   req.ToFirstName,
		LastName:    req.ToLastName,
		BookingID:   req.BookingID,
		ProductName: req.ProductName,
		RoomName:    req.RoomName,
		Date:        req.Date,
		Time:        req.Time,
		Reason:      req.Reason,
		Year:        req.Year,
	}

	emailReq := &domain.EmailRequest{
		To:       req.ToEmail,
		Subject:  "Annulation de réservation",
		Template: "booking_cancellation",
		Data:     data,
	}

	return c.sendEmail(ctx, emailReq, "")
}

func (c *SMTPClient) SendBookingReminderEmail(ctx context.Context, req domain.BookingReminderRequest) error {
	data := struct {
		FirstName   string
		LastName    string
		BookingID   string
		ProductName string
		RoomName    string
		Building    string
		Address     string
		Date        string
		Time        string
		PartnerName string
		Year        int
	}{
		FirstName:   req.ToFirstName,
		LastName:    req.ToLastName,
		BookingID:   req.BookingID,
		ProductName: req.ProductName,
		RoomName:    req.RoomName,
		Building:    req.Building,
		Address:     req.Address,
		Date:        req.Date,
		Time:        req.Time,
		PartnerName: req.PartnerName,
		Year:        req.Year,
	}

	emailReq := &domain.EmailRequest{
		To:       req.ToEmail,
		Subject:  "Rappel de réservation",
		Template: "booking_reminder",
		Data:     data,
	}

	return c.sendEmail(ctx, emailReq, "")
}

func (c *SMTPClient) sendEmail(ctx context.Context, request *domain.EmailRequest, fromEmail string) error {
	if request == nil {
		return errs.ErrInvalidValue
	}

	if request.To == "" {
		return fmt.Errorf("recipient email cannot be empty: %w", errs.ErrInvalidValue)
	}

	if request.Subject == "" {
		return fmt.Errorf("email subject cannot be empty: %w", errs.ErrInvalidValue)
	}

	if fromEmail == "" {
		fromEmail = c.config.Username
	}

	m := gomail.NewMessage()
	m.SetHeader("From", fromEmail)
	m.SetHeader("To", request.To)
	m.SetHeader("Subject", request.Subject)

	if len(request.CarbonCopy) > 0 {
		addresses := make([]string, 0, len(request.CarbonCopy))
		for email, name := range request.CarbonCopy {
			addresses = append(addresses, m.FormatAddress(email, name))
		}
		m.SetHeader("Cc", addresses...)
	}

	// Embed Instagram image
	instagramPath, err := writeTempFile(domain.InstagramImage, "instagram.png")
	if err != nil {
		return fmt.Errorf("write instagram temp file: %w", err)
	}
	defer os.Remove(instagramPath)
	m.Embed(instagramPath, gomail.Rename("instagram.png"))

	// Embed additional images if specified
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

	d := gomail.NewDialer(c.config.Host, c.config.Port, c.config.Username, c.config.Password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("send email via SMTP: %w", err)
	}

	return nil
}

func (c *SMTPClient) renderTemplate(templateName string, data interface{}) (string, error) {
	templatePath := fmt.Sprintf("templates/%s.html", templateName)

	tmpl, err := template.ParseFS(domain.EmailTemplates, templatePath)
	if err != nil {
		return "", fmt.Errorf("parse email template %s: %w", templateName, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute email template %s: %w", templateName, err)
	}

	return buf.String(), nil
}

func writeTempFile(data []byte, filename string) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("file data cannot be empty: %w", errs.ErrInvalidValue)
	}

	ext := filepath.Ext(filename)
	stem := filename[:len(filename)-len(ext)]

	f, err := os.CreateTemp("", stem+"*"+ext)
	if err != nil {
		return "", fmt.Errorf("create temp file %s: %w", filename, err)
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		os.Remove(f.Name())
		return "", fmt.Errorf("write temp file %s: %w", filename, err)
	}

	return f.Name(), nil
}
