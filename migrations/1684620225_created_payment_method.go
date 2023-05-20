package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `{
			"id": "z7j1snfgyn31sid",
			"created": "2023-05-20 22:03:45.887Z",
			"updated": "2023-05-20 22:03:45.887Z",
			"name": "payment_method",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "m8p8cdw8",
					"name": "name",
					"type": "relation",
					"required": true,
					"unique": false,
					"options": {
						"collectionId": "5rxmpfs22bhonfh",
						"cascadeDelete": false,
						"minSelect": null,
						"maxSelect": 1,
						"displayFields": []
					}
				},
				{
					"system": false,
					"id": "6ozdsmbh",
					"name": "fee_percentage",
					"type": "number",
					"required": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null
					}
				},
				{
					"system": false,
					"id": "pxltgwmv",
					"name": "fee",
					"type": "number",
					"required": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null
					}
				}
			],
			"listRule": null,
			"viewRule": null,
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("z7j1snfgyn31sid")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
