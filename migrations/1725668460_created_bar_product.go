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
			"id": "nnad5hvrhvlw5bn",
			"created": "2026-08-06 00:00:00.000Z",
			"updated": "2026-08-06 00:00:00.000Z",
			"name": "bar_product",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "6kewtihb",
					"name": "slug",
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
					"id": "ura9ssyz",
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
					"id": "o32zgc0b",
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
					"id": "y8berow4",
					"name": "type",
					"type": "select",
					"required": true,
					"unique": false,
					"options": {
						"maxSelect": 1,
						"values": [
							"product",
							"discount_percentage",
							"discount",
							"deposit"
						]
					}
				},
				{
					"system": false,
					"id": "ptjvgfak",
					"name": "color",
					"type": "text",
					"required": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				},
				{
					"system": false,
					"id": "kb767cls",
					"name": "order",
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
					"id": "fcyz8bh4",
					"name": "deactivated",
					"type": "bool",
					"required": false,
					"unique": false,
					"options": {}
				},
				{
					"system": false,
					"id": "zbtd6gjm",
					"name": "admin_only",
					"type": "bool",
					"required": false,
					"unique": false,
					"options": {}
				},
				{
					"system": false,
					"id": "uozwz4u8",
					"name": "instructions",
					"type": "text",
					"required": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				},
				{
					"system": false,
					"id": "k48p0a5f",
					"name": "requires_deposit",
					"type": "bool",
					"required": false,
					"unique": false,
					"options": {}
				},
				{
					"system": false,
					"id": "zf6jaqbr",
					"name": "bars",
					"type": "relation",
					"required": false,
					"unique": false,
					"options": {
						"collectionId": "u8x3rq8d3znqqqb",
						"cascadeDelete": false,
						"minSelect": null,
						"maxSelect": null,
						"displayFields": []
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

		collection, err := dao.FindCollectionByNameOrId("nnad5hvrhvlw5bn")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
