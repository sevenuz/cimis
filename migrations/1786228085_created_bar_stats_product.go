package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

// Grain: (event, bar, product, free). Powers per-product units/revenue (client
// sums across bar/free for per-event totals, across event for all-time),
// free-vs-paid split, revenue per product per bar, and top/bottom seller
// ranking - all without pulling every bar_serving row to the client.
func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `{
			"id": "1lq59uj9alor9tz",
			"created": "2026-08-09 00:00:00.000Z",
			"updated": "2026-08-09 00:00:00.000Z",
			"name": "bar_stats_product",
			"type": "view",
			"system": false,
			"schema": [],
			"listRule": "@request.auth.admin = true",
			"viewRule": "@request.auth.admin = true",
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {
				"query": "SELECT (bo.event || '_' || bo.bar || '_' || bs.product || '_' || bs.free) AS id, bo.event AS event, bo.bar AS bar, bs.product AS product, bs.free AS free, CAST(SUM(bs.amount) AS REAL) AS units, CAST(SUM(bs.amount * bs.price) AS REAL) AS revenue FROM bar_serving bs JOIN bar_order bo ON bo.id = bs.\"order\" GROUP BY bo.event, bo.bar, bs.product, bs.free"
			}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("1lq59uj9alor9tz")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
