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

		// remove
		collection.Schema.RemoveField("nxqvnhsk")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("dkhhviph3mim48w")
		if err != nil {
			return err
		}

		// add
		del_amount := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "nxqvnhsk",
			"name": "amount",
			"type": "number",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null
			}
		}`), del_amount)
		collection.Schema.AddField(del_amount)

		return dao.SaveCollection(collection)
	})
}
