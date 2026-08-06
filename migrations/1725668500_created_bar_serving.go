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
			"id": "addem50uzdhgc96",
			"created": "2026-08-06 00:00:00.000Z",
			"updated": "2026-08-06 00:00:00.000Z",
			"name": "bar_serving",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "ntzx57hs",
					"name": "order",
					"type": "relation",
					"required": true,
					"unique": false,
					"options": {
						"collectionId": "g7x7a85mhqk42t6",
						"cascadeDelete": true,
						"minSelect": null,
						"maxSelect": 1,
						"displayFields": []
					}
				},
				{
					"system": false,
					"id": "0f36ono2",
					"name": "product",
					"type": "relation",
					"required": true,
					"unique": false,
					"options": {
						"collectionId": "nnad5hvrhvlw5bn",
						"cascadeDelete": false,
						"minSelect": null,
						"maxSelect": 1,
						"displayFields": []
					}
				},
				{
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
				},
				{
					"system": false,
					"id": "baf39mbd",
					"name": "amount",
					"type": "number",
					"required": true,
					"unique": false,
					"options": {
						"min": null,
						"max": null
					}
				},
				{
					"system": false,
					"id": "n6bjmo3b",
					"name": "price",
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
					"id": "mnha1vts",
					"name": "free",
					"type": "bool",
					"required": false,
					"unique": false,
					"options": {}
				}
			],
			"listRule": "@request.auth.id != \"\"",
			"viewRule": "@request.auth.id != \"\"",
			"createRule": "@request.auth.id != \"\"",
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
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("addem50uzdhgc96")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
