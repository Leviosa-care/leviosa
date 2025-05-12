package mailService

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/gomail.v2"
)

func (s *service) sendMail(ctx context.Context, to, subject, templateFilename string, data any, carbonCopy, images map[string]string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	addresses := make([]string, 0, len(carbonCopy))
	for email, name := range carbonCopy {
		addresses = append(addresses, m.FormatAddress(email, name))
	}
	m.SetHeader("Cc", addresses...)

	logo, err := s.getLogo(ctx)
	if err != nil {
		return fmt.Errorf("get logo before sending email: %w", err)
	}
	logoPath, err := writeTempFile(logo, "logo.jpg")
	if err != nil {
		return err
	}
	defer os.Remove(logoPath)

	instaPath, err := writeTempFile(instagramImage, "instagram.png")
	if err != nil {
		return err
	}
	defer os.Remove(instaPath)

	m.Embed(logoPath, gomail.Rename("logo.jpg"))
	m.Embed(instaPath, gomail.Rename("instagram.png"))

	// Embed other images in mail
	for path, rename := range images {
		m.Embed(path, gomail.Rename(rename))
	}

	t, err := template.ParseFS(emailTemplates, fmt.Sprintf("templates/%s.html", templateFilename))
	if err != nil {
		return fmt.Errorf("parsing template: %s", err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return fmt.Errorf("execute template: %s", err)
	}

	m.SetBody("text/html", tpl.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, s.email, s.password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("dial and sent mail: %s", err)
	}
	return nil
}

func writeTempFile(data []byte, filename string) (string, error) {
	tmpDir := os.TempDir()
	tmpPath := filepath.Join(tmpDir, filename)

	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return "", err
	}
	return tmpPath, nil
}

func (s *service) getLogo(ctx context.Context) ([]byte, error) {
	url, err := s.repo.GetLogo(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching image: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
