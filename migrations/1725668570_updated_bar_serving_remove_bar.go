package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

// bar_serving.bar duplicated bar_order.bar for no reason - a serving's bar is
// always its order's bar (both are stamped from the same current_bar_id at
// creation time, they can never diverge in this design). Never actually
// queried by bar anywhere in the app; PocketBase's relation dot-notation
// filters (e.g. `order.bar="X"`) cover that if it's ever needed.
func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("addem50uzdhgc96")
		if err != nil {
			return err
		}

		collection.Schema.RemoveField("2yjtjvij")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("addem50uzdhgc96")
		if err != nil {
			return err
		}

		field := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "2yjtjvij",
			"name": "bar",
			"type": "relation",
			"required": true,
			"unique": false,
			"options": {
				"collectionId": "u8x3rq8d3znqqqb",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": []
			}
		}`), field)
		collection.Schema.AddField(field)

		return dao.SaveCollection(collection)
	})
}
