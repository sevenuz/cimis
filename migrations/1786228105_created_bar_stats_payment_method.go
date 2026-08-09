package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

// Grain: (event, bar, payment_method), bookouts excluded (they're manual box
// adjustments, not sales). fee_cost is the charged total minus the raw sum of
// the underlying bar_serving line items - this measures the surcharge actually
// collected without inverting bar_payment_method's fee_percentage/fee formula
// (total *= fee_percentage; total += fee, see get_serving_total in bar.ts),
// which would risk a rounding-mismatch bug.
func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `{
			"id": "qjtqaw1x0nwuvlp",
			"created": "2026-08-09 00:00:00.000Z",
			"updated": "2026-08-09 00:00:00.000Z",
			"name": "bar_stats_payment_method",
			"type": "view",
			"system": false,
			"schema": [],
			"listRule": "@request.auth.admin = true",
			"viewRule": "@request.auth.admin = true",
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {
				"query": "SELECT (bo.event || '_' || bo.bar || '_' || bo.payment_method) AS id, bo.event AS event, bo.bar AS bar, bo.payment_method AS payment_method, CAST(COUNT(*) AS REAL) AS order_count, CAST(SUM(bo.total) AS REAL) AS volume, CAST(AVG(bo.total) AS REAL) AS avg_order_value, CAST(COALESCE(SUM(ot.order_total), 0) AS REAL) AS raw_total, CAST(SUM(bo.total) - COALESCE(SUM(ot.order_total), 0) AS REAL) AS fee_cost FROM bar_order bo LEFT JOIN (SELECT \"order\", SUM(amount * price) AS order_total FROM bar_serving GROUP BY \"order\") ot ON ot.\"order\" = bo.id WHERE bo.is_bookout = 0 GROUP BY bo.event, bo.bar, bo.payment_method"
			}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("qjtqaw1x0nwuvlp")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
