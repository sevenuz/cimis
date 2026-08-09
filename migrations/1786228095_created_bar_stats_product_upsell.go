package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

// Grain: (event, product). avg_order_total is the average of the *per-order*
// total (not per-line) across every order that contains this product - the
// "average order size when this product is bought" upsell signal.
func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `{
			"id": "mf8qmbrmnq1ighz",
			"created": "2026-08-09 00:00:00.000Z",
			"updated": "2026-08-09 00:00:00.000Z",
			"name": "bar_stats_product_upsell",
			"type": "view",
			"system": false,
			"schema": [],
			"listRule": "@request.auth.admin = true",
			"viewRule": "@request.auth.admin = true",
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {
				"query": "SELECT (bo.event || '_' || bs.product) AS id, bo.event AS event, bs.product AS product, CAST(COUNT(DISTINCT bo.id) AS REAL) AS order_count, CAST(AVG(ot.order_total) AS REAL) AS avg_order_total FROM bar_serving bs JOIN bar_order bo ON bo.id = bs.\"order\" JOIN (SELECT \"order\", SUM(amount * price) AS order_total FROM bar_serving GROUP BY \"order\") ot ON ot.\"order\" = bs.\"order\" GROUP BY bo.event, bs.product"
			}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("mf8qmbrmnq1ighz")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
