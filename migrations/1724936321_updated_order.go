package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("7rp45edh7j2olq6")
		if err != nil {
			return err
		}

		collection.Name = "bar_order"

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("7rp45edh7j2olq6")
		if err != nil {
			return err
		}

		collection.Name = "order"

		return dao.SaveCollection(collection)
	})
}
