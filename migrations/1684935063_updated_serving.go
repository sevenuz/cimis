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

		collection, err := dao.FindCollectionByNameOrId("xek65r5tfn9s7vz")
		if err != nil {
			return err
		}

		// add
		new_type := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "qxk1ub4w",
			"name": "type",
			"type": "select",
			"required": true,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"values": [
					"product",
					"discount_percentage",
					"discount",
					"deposit"
				]
			}
		}`), new_type)
		collection.Schema.AddField(new_type)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("xek65r5tfn9s7vz")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("qxk1ub4w")

		return dao.SaveCollection(collection)
	})
}
