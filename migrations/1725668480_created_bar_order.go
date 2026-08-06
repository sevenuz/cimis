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
			"id": "g7x7a85mhqk42t6",
			"created": "2026-08-06 00:00:00.000Z",
			"updated": "2026-08-06 00:00:00.000Z",
			"name": "bar_order",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "cegkq25k",
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
					"id": "7n36jn3s",
					"name": "payment_method",
					"type": "relation",
					"required": true,
					"unique": false,
					"options": {
						"collectionId": "wsn264572am4ub4",
						"cascadeDelete": false,
						"minSelect": null,
						"maxSelect": 1,
						"displayFields": []
					}
				},
				{
					"system": false,
					"id": "fi8t1moa",
					"name": "total",
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
					"id": "0aykq2cx",
					"name": "is_bookout",
					"type": "bool",
					"required": false,
					"unique": false,
					"options": {}
				},
				{
					"system": false,
					"id": "9wih3ppp",
					"name": "client_id",
					"type": "text",
					"required": true,
					"unique": true,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				},
				{
					"system": false,
					"id": "rshkr98a",
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
					"id": "e3k73xqk",
					"name": "event",
					"type": "relation",
					"required": false,
					"unique": false,
					"options": {
						"collectionId": "dax2xrm503yoqmh",
						"cascadeDelete": false,
						"minSelect": null,
						"maxSelect": 1,
						"displayFields": []
					}
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

		collection, err := dao.FindCollectionByNameOrId("g7x7a85mhqk42t6")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
