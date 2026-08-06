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
			"id": "58x3s2vtsq5jcn9",
			"created": "2026-08-06 00:00:00.000Z",
			"updated": "2026-08-06 00:00:00.000Z",
			"name": "bar_inventory",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "8yg3uks3",
					"name": "event",
					"type": "relation",
					"required": true,
					"unique": false,
					"options": {
						"collectionId": "dax2xrm503yoqmh",
						"cascadeDelete": false,
						"minSelect": null,
						"maxSelect": 1,
						"displayFields": []
					}
				},
				{
					"system": false,
					"id": "7my5lthh",
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
					"id": "nnlf3m45",
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
					"id": "685c0czu",
					"name": "note",
					"type": "text",
					"required": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				}
			],
			"listRule": "@request.auth.id != \"\"",
			"viewRule": "@request.auth.id != \"\"",
			"createRule": "@request.auth.admin = true",
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

		collection, err := dao.FindCollectionByNameOrId("58x3s2vtsq5jcn9")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
