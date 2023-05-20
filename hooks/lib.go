package hooks

import (
	"github.com/pocketbase/pocketbase/core"
)

func Register(app core.App) {
	addPageValidationHook(app)
}
