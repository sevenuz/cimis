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

		collection, err := dao.FindCollectionByNameOrId("dkhhviph3mim48w")
		if err != nil {
			return err
		}

		// add
		new_admin_only := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "9hwhgred",
			"name": "admin_only",
			"type": "bool",
			"required": false,
			"unique": false,
			"options": {}
		}`), new_admin_only)
		collection.Schema.AddField(new_admin_only)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("dkhhviph3mim48w")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("9hwhgred")

		return dao.SaveCollection(collection)
	})
}
