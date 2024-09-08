package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("dkhhviph3mim48w")
		if err != nil {
			return err
		}

		collection.Name = "bar_inventory"

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("dkhhviph3mim48w")
		if err != nil {
			return err
		}

		collection.Name = "inventory"

		return dao.SaveCollection(collection)
	})
}
