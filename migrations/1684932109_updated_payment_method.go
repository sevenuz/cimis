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

		collection, err := dao.FindCollectionByNameOrId("z7j1snfgyn31sid")
		if err != nil {
			return err
		}

		// add
		new_color := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "clf1brin",
			"name": "color",
			"type": "text",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_color)
		collection.Schema.AddField(new_color)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("z7j1snfgyn31sid")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("clf1brin")

		return dao.SaveCollection(collection)
	})
}
