package domain

import (
	"embed"
)

var (
	//go:embed ../../assets/instagram.png
	InstagramImage []byte

	//go:embed ../../templates/*.html
	EmailTemplates embed.FS
)

