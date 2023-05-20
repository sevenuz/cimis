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
			"id": "6ogbpmksv47h2tl",
			"created": "2023-05-20 14:31:02.457Z",
			"updated": "2023-05-20 14:31:02.457Z",
			"name": "feature",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "nuegcptt",
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
					"id": "6lx7mhgy",
					"name": "files",
					"type": "file",
					"required": false,
					"unique": false,
					"options": {
						"maxSelect": 5,
						"maxSize": 5242880,
						"mimeTypes": [],
						"thumbs": []
					}
				},
				{
					"system": false,
					"id": "s53h1ran",
					"name": "data",
					"type": "json",
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

		collection, err := dao.FindCollectionByNameOrId("6ogbpmksv47h2tl")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
