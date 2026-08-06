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
			"id": "5s33pb1flrqaauv",
			"created": "2026-08-06 00:00:00.000Z",
			"updated": "2026-08-06 00:00:00.000Z",
			"name": "bar_recipe_ingredient",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "5njs1qpl",
					"name": "product",
					"type": "relation",
					"required": true,
					"unique": false,
					"options": {
						"collectionId": "nnad5hvrhvlw5bn",
						"cascadeDelete": true,
						"minSelect": null,
						"maxSelect": 1,
						"displayFields": []
					}
				},
				{
					"system": false,
					"id": "c29y4dur",
					"name": "ingredient",
					"type": "relation",
					"required": true,
					"unique": false,
					"options": {
						"collectionId": "zohq2dguzpxkdjt",
						"cascadeDelete": false,
						"minSelect": null,
						"maxSelect": 1,
						"displayFields": []
					}
				},
				{
					"system": false,
					"id": "vpxcgy2h",
					"name": "quantity",
					"type": "number",
					"required": true,
					"unique": false,
					"options": {
						"min": null,
						"max": null
					}
				}
			],
			"listRule": "",
			"viewRule": "",
			"createRule": "@request.auth.admin = true",
			"updateRule": "@request.auth.admin = true",
			"deleteRule": "@request.auth.admin = true",
			"options": {}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("5s33pb1flrqaauv")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
