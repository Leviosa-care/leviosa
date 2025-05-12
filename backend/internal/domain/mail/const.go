package mailService

import (
	"embed"
)

var (
	//go:embed assets/instagram.png
	instagramImage []byte
	//go:embed templates/*.html
	emailTemplates embed.FS
)
