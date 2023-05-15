package routes

import (
	"cimis/ui"
	"github.com/labstack/echo/v5/middleware"

	// For route extension
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"net/http"
)

func addUiServingRoutes(app core.App) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {

		e.Router.GET(
			"/*",
			func(c echo.Context) error {
				return c.HTML(http.StatusOK, ui.Index_file)
			},
			middleware.Gzip(),
		)
		e.Router.GET(
			"/icon.png",
			func(c echo.Context) error {
				return c.Blob(http.StatusOK, "image/png", ui.Icon_file)
			},
			middleware.Gzip(),
		)
		e.Router.GET(
			"/particles.min.js",
			func(c echo.Context) error {
				return c.Blob(http.StatusOK, "text/javascript", ui.Particles_file)
			},
			middleware.Gzip(),
		)
		e.Router.GET(
			"/_app/*",
			apis.StaticDirectoryHandler(ui.AssetsDirFS, false),
			middleware.Gzip(),
		)

		return nil
	})
}
