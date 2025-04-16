package mailService

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"gopkg.in/gomail.v2"
)

const TEMPLATE_PATH = "/internal/mail/templates/"

func sendMail(from, to, subject, templateFilename, password string, data any, images map[string]string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", subject)

	//Embed the images
	for path, rename := range images {
		m.Embed(path, gomail.Rename(rename))
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %s", err)
	}
	t, err := template.ParseFiles(wd + templateFilename)
	if err != nil {
		return fmt.Errorf("parsing template: %s", err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return fmt.Errorf("execute template: %s", err)
	}
	m.SetBody("text/html", tpl.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, from, password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("dial and sent mail: %s", err)
	}
	return nil
}

func getCompanyCredentials() (string, string) {
	return os.Getenv("GMAIL_EMAIL"), os.Getenv("GMAIL_PASSWORD")
}
