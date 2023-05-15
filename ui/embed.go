// Package ui handles the PocketBase Admin frontend embedding.
package ui

import (
	"embed"

	"github.com/labstack/echo/v5"
)

//go:embed build/index.html
var Index_file string

//go:embed build/icon.png
var Icon_file []byte

//go:embed build/particles.min.js
var Particles_file []byte

//go:embed all:build
var distDir embed.FS

// DistDirFS contains the embedded dist directory files (without the "dist" prefix)
var distDirFS = echo.MustSubFS(distDir, "build")
var AssetsDirFS = echo.MustSubFS(distDirFS, "_app")
