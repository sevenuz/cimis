package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("n28qbs9mimev1m5")
		if err != nil {
			return err
		}

		// add
		new_admin := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "3udenc0e",
			"name": "admin",
			"type": "bool",
			"required": false,
			"unique": false,
			"options": {}
		}`), new_admin)
		collection.Schema.AddField(new_admin)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("n28qbs9mimev1m5")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("3udenc0e")

		return dao.SaveCollection(collection)
	})
}
