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

		// remove
		collection.Schema.RemoveField("ivrlyie1")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("xek65r5tfn9s7vz")
		if err != nil {
			return err
		}

		// add
		del_total := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "ivrlyie1",
			"name": "total",
			"type": "number",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null
			}
		}`), del_total)
		collection.Schema.AddField(del_total)

		return dao.SaveCollection(collection)
	})
}
