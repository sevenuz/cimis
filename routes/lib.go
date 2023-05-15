package routes

import (
	"github.com/pocketbase/pocketbase/core"
)

func Register(app core.App) {
	addUiServingRoutes(app)
}
