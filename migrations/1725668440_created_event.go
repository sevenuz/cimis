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
			"id": "dax2xrm503yoqmh",
			"created": "2026-08-06 00:00:00.000Z",
			"updated": "2026-08-06 00:00:00.000Z",
			"name": "event",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "92g8u8xt",
					"name": "name",
					"type": "text",
					"required": true,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				},
				{
					"system": false,
					"id": "x5wfauhf",
					"name": "active",
					"type": "bool",
					"required": false,
					"unique": false,
					"options": {}
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

		collection, err := dao.FindCollectionByNameOrId("dax2xrm503yoqmh")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
