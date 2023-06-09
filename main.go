package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	"cimis/hooks"
	_ "cimis/migrations"
	"cimis/routes"
)

func main() {
	app := pocketbase.New()

	migratecmd.MustRegister(app, app.RootCmd, &migratecmd.Options{
		Automigrate: true, // auto creates migration files when making collection changes
	})

	hooks.Register(app)
	routes.Register(app)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
