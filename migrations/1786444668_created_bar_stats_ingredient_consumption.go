package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

// Grain: (event, ingredient). Replaces the O(servings x recipe_ingredients)
// nested loop in inventory/+page.svelte's compute() - that loop scanned every
// recipe_ingredient row for every serving row client-side; this pushes the
// same join+aggregation into SQL. Same auth rule as bar_inventory/bar_order
// (any authenticated user, not admin-only) since the inventory page itself is
// viewable by any logged-in user, only the add-forms are admin-gated.
func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `{
			"id": "3dgkqt5t4iwxydu",
			"created": "2026-08-11 00:00:00.000Z",
			"updated": "2026-08-11 00:00:00.000Z",
			"name": "bar_stats_ingredient_consumption",
			"type": "view",
			"system": false,
			"schema": [],
			"listRule": "@request.auth.id != \"\"",
			"viewRule": "@request.auth.id != \"\"",
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {
				"query": "SELECT (bo.event || '_' || ri.ingredient) AS id, bo.event AS event, ri.ingredient AS ingredient, CAST(SUM(bs.amount * ri.quantity) AS REAL) AS consumed FROM bar_serving bs JOIN bar_order bo ON bo.id = bs.\"order\" JOIN bar_recipe_ingredient ri ON ri.product = bs.product GROUP BY bo.event, ri.ingredient"
			}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("3dgkqt5t4iwxydu")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
