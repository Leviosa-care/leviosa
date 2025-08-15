package userHandler

import (
	"github.com/hengadev/leviosa/internal/server/app"
)

type AppInstance struct {
	*app.App
}

func New(appCtx *app.App) *AppInstance {
	return &AppInstance{appCtx}
}
