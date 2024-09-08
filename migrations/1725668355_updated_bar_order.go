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
		new_is_bookout := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "e1ce8ufl",
			"name": "is_bookout",
			"type": "bool",
			"required": false,
			"unique": false,
			"options": {}
		}`), new_is_bookout)
		collection.Schema.AddField(new_is_bookout)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("7rp45edh7j2olq6")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("e1ce8ufl")

		return dao.SaveCollection(collection)
	})
}
