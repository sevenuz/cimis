package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

// Grain: (event, bar, 15-minute bucket of created), bookouts excluded. Baked
// at the finest granularity the stats UI offers - a 15/30/60 min toggle
// re-buckets client-side by summing consecutive rows, so no dynamic SQL or
// extra views are needed for that. Powers per-bar-per-hour revenue/order
// count (busiest hour/bar) and, summed across bars per bucket, the overall
// event revenue curve.
func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `{
			"id": "v9o0mmadsfxxgn6",
			"created": "2026-08-09 00:00:00.000Z",
			"updated": "2026-08-09 00:00:00.000Z",
			"name": "bar_stats_time_bucket",
			"type": "view",
			"system": false,
			"schema": [],
			"listRule": "@request.auth.admin = true",
			"viewRule": "@request.auth.admin = true",
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {
				"query": "SELECT (bo.event || '_' || bo.bar || '_' || datetime((strftime('%s', bo.created)/900)*900, 'unixepoch')) AS id, bo.event AS event, bo.bar AS bar, CAST(datetime((strftime('%s', bo.created)/900)*900, 'unixepoch') AS TEXT) AS bucket, CAST(COUNT(*) AS REAL) AS order_count, CAST(SUM(bo.total) AS REAL) AS revenue FROM bar_order bo WHERE bo.is_bookout = 0 GROUP BY bo.event, bo.bar, bucket"
			}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("v9o0mmadsfxxgn6")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
