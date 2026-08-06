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
			"id": "u8x3rq8d3znqqqb",
			"created": "2026-08-06 00:00:00.000Z",
			"updated": "2026-08-06 00:00:00.000Z",
			"name": "bar",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "en6ea2h8",
					"name": "name",
					"type": "text",
					"required": true,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
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

		collection, err := dao.FindCollectionByNameOrId("u8x3rq8d3znqqqb")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
