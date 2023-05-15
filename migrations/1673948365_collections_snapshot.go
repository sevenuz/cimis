package migrations

import (
	"encoding/json"

	"cimis/models"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	pb_models "github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `[
		` + models.JsonLanguageKey + `,
		` + models.JsonRest + `
		]`
		collections := []*pb_models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collections); err != nil {
			return err
		}

		if err := daos.New(db).ImportCollections(collections, true, nil); err != nil {
			return err
		}

		query := db.NewQuery("CREATE UNIQUE INDEX language_index ON language_value(language_key, iso);")
		if _, err := query.Execute(); err != nil {
			return err
		}
		query2 := db.NewQuery("CREATE UNIQUE INDEX page_index ON page(path, iso);")
		if _, err := query2.Execute(); err != nil {
			return err
		}
		return nil
	}, func(db dbx.Builder) error {
		return nil
	})
}
