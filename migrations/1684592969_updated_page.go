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

		collection, err := dao.FindCollectionByNameOrId("yf7bx6hp7amhg5h")
		if err != nil {
			return err
		}

		// add
		new_addon := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "9sb2nhnt",
			"name": "addon",
			"type": "select",
			"required": false,
			"unique": false,
			"options": {
				"maxSelect": 4,
				"values": [
					"particle",
					"bar",
					"pretix",
					"external_link"
				]
			}
		}`), new_addon)
		collection.Schema.AddField(new_addon)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("yf7bx6hp7amhg5h")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("9sb2nhnt")

		return dao.SaveCollection(collection)
	})
}
