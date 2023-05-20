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
			"id": "dkhhviph3mim48w",
			"created": "2023-05-20 21:56:39.319Z",
			"updated": "2023-05-20 21:56:39.319Z",
			"name": "inventory",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "wsitdvtp",
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
					"id": "teelib8h",
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
					"id": "nxqvnhsk",
					"name": "amount",
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
					"id": "gbhlkngy",
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
					"id": "to4rinss",
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
					"id": "sskersyh",
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
					"id": "zuu1k3i4",
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
					"id": "dxro3sov",
					"name": "deactivated",
					"type": "bool",
					"required": false,
					"unique": false,
					"options": {}
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

		collection, err := dao.FindCollectionByNameOrId("dkhhviph3mim48w")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
