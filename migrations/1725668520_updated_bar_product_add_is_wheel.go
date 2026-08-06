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
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("nnad5hvrhvlw5bn")
		if err != nil {
			return err
		}

		new_is_wheel := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "1zdklx1y",
			"name": "is_wheel",
			"type": "bool",
			"required": false,
			"unique": false,
			"options": {}
		}`), new_is_wheel)
		collection.Schema.AddField(new_is_wheel)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("nnad5hvrhvlw5bn")
		if err != nil {
			return err
		}

		collection.Schema.RemoveField("1zdklx1y")

		return dao.SaveCollection(collection)
	})
}
