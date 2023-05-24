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

		collection, err := dao.FindCollectionByNameOrId("7rp45edh7j2olq6")
		if err != nil {
			return err
		}

		// add
		new_tip := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "ze8j8xht",
			"name": "tip",
			"type": "number",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null
			}
		}`), new_tip)
		collection.Schema.AddField(new_tip)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("7rp45edh7j2olq6")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("ze8j8xht")

		return dao.SaveCollection(collection)
	})
}
