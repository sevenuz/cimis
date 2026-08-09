package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

// Grain: (event, product, 15-minute bucket of created). Powers the
// per-product sales-over-time line chart (which product sells at which time
// of day).
func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `{
			"id": "btniomjju8rk4hq",
			"created": "2026-08-09 00:00:00.000Z",
			"updated": "2026-08-09 00:00:00.000Z",
			"name": "bar_stats_product_time",
			"type": "view",
			"system": false,
			"schema": [],
			"listRule": "@request.auth.admin = true",
			"viewRule": "@request.auth.admin = true",
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {
				"query": "SELECT (bo.event || '_' || bs.product || '_' || datetime((strftime('%s', bo.created)/900)*900, 'unixepoch')) AS id, bo.event AS event, bs.product AS product, CAST(datetime((strftime('%s', bo.created)/900)*900, 'unixepoch') AS TEXT) AS bucket, CAST(SUM(bs.amount) AS REAL) AS units, CAST(SUM(bs.amount * bs.price) AS REAL) AS revenue FROM bar_serving bs JOIN bar_order bo ON bo.id = bs.\"order\" GROUP BY bo.event, bs.product, bucket"
			}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("btniomjju8rk4hq")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
