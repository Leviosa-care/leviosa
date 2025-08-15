package authHandler

import (
	"github.com/hengadev/leviosa/internal/server/app"
)

// TODO: complete that thing brother

type handler struct {
	*app.App
}

func New(appCtx *app.App) *handler {
	return &handler{appCtx}
}
