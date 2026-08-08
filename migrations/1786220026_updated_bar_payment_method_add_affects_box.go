package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

// expected_in_box previously summed earned+booked_out across every payment
// method, but only cash physically sits in the box - card revenue never
// touches it. affects_box lets bookkeeping filter to the payment methods
// that actually count toward the box total; existing methods default to
// false and need to be flipped on for cash via the admin UI.
func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("wsn264572am4ub4")
		if err != nil {
			return err
		}

		field := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "goq5yrd0",
			"name": "affects_box",
			"type": "bool",
			"required": false,
			"unique": false,
			"options": {}
		}`), field)
		collection.Schema.AddField(field)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("wsn264572am4ub4")
		if err != nil {
			return err
		}

		collection.Schema.RemoveField("goq5yrd0")

		return dao.SaveCollection(collection)
	})
}
