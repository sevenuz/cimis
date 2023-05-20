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
			"id": "xek65r5tfn9s7vz",
			"created": "2023-05-20 22:13:49.045Z",
			"updated": "2023-05-20 22:13:49.045Z",
			"name": "serving",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "khwvit8y",
					"name": "order",
					"type": "relation",
					"required": true,
					"unique": false,
					"options": {
						"collectionId": "7rp45edh7j2olq6",
						"cascadeDelete": false,
						"minSelect": null,
						"maxSelect": 1,
						"displayFields": []
					}
				},
				{
					"system": false,
					"id": "cji1eb51",
					"name": "product",
					"type": "relation",
					"required": true,
					"unique": false,
					"options": {
						"collectionId": "dkhhviph3mim48w",
						"cascadeDelete": false,
						"minSelect": null,
						"maxSelect": 1,
						"displayFields": []
					}
				},
				{
					"system": false,
					"id": "55jvsvah",
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
					"id": "ccr7h6mh",
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
					"id": "ivrlyie1",
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

		collection, err := dao.FindCollectionByNameOrId("xek65r5tfn9s7vz")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
