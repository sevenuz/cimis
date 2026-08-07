package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

// requires_deposit turned out unnecessary - the confirm-order warning is now a
// static always-shown reminder (see ui_bar_confirmation_alert) rather than
// conditional on a per-product flag.
func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("nnad5hvrhvlw5bn")
		if err != nil {
			return err
		}

		collection.Schema.RemoveField("k48p0a5f")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("nnad5hvrhvlw5bn")
		if err != nil {
			return err
		}

		field := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "k48p0a5f",
			"name": "requires_deposit",
			"type": "bool",
			"required": false,
			"unique": false,
			"options": {}
		}`), field)
		collection.Schema.AddField(field)

		return dao.SaveCollection(collection)
	})
}
