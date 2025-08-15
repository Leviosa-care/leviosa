package healthHandler

import (
	"github.com/hengadev/leviosa/internal/server/app"
)

type handler struct {
	*app.App
}

func New(appCtx *app.App) *handler {
	return &handler{appCtx}
}
