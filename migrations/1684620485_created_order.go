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
			"id": "7rp45edh7j2olq6",
			"created": "2023-05-20 22:08:05.680Z",
			"updated": "2023-05-20 22:08:05.680Z",
			"name": "order",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "ne6ewiye",
					"name": "user",
					"type": "relation",
					"required": true,
					"unique": false,
					"options": {
						"collectionId": "n28qbs9mimev1m5",
						"cascadeDelete": false,
						"minSelect": null,
						"maxSelect": 1,
						"displayFields": []
					}
				},
				{
					"system": false,
					"id": "w6dy4rul",
					"name": "payment_method",
					"type": "relation",
					"required": true,
					"unique": false,
					"options": {
						"collectionId": "z7j1snfgyn31sid",
						"cascadeDelete": false,
						"minSelect": null,
						"maxSelect": 1,
						"displayFields": []
					}
				},
				{
					"system": false,
					"id": "uvkfk37b",
					"name": "total",
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

		collection, err := dao.FindCollectionByNameOrId("7rp45edh7j2olq6")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
