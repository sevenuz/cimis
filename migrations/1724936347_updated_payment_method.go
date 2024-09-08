package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("z7j1snfgyn31sid")
		if err != nil {
			return err
		}

		collection.Name = "bar_payment_method"

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("z7j1snfgyn31sid")
		if err != nil {
			return err
		}

		collection.Name = "payment_method"

		return dao.SaveCollection(collection)
	})
}
