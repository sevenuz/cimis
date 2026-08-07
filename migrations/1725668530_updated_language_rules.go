package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

// The new /bar/recipes page lets an app-admin (user.admin = true, not a PocketBase
// superuser) create translated product names on the fly instead of requiring the
// raw PocketBase dashboard. That needs create/update access on language_key and
// language_value, which previously had those rules unset (superuser-only) - a plain
// user.admin flag never satisfied that, regardless of how correctly it was set.
func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)
		rule := "@request.auth.admin = true"

		key_collection, err := dao.FindCollectionByNameOrId("5rxmpfs22bhonfh")
		if err != nil {
			return err
		}
		key_collection.CreateRule = &rule
		key_collection.UpdateRule = &rule
		if err := dao.SaveCollection(key_collection); err != nil {
			return err
		}

		value_collection, err := dao.FindCollectionByNameOrId("ilf6q6610e4ellf")
		if err != nil {
			return err
		}
		value_collection.CreateRule = &rule
		value_collection.UpdateRule = &rule
		return dao.SaveCollection(value_collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		key_collection, err := dao.FindCollectionByNameOrId("5rxmpfs22bhonfh")
		if err != nil {
			return err
		}
		key_collection.CreateRule = nil
		key_collection.UpdateRule = nil
		if err := dao.SaveCollection(key_collection); err != nil {
			return err
		}

		value_collection, err := dao.FindCollectionByNameOrId("ilf6q6610e4ellf")
		if err != nil {
			return err
		}
		value_collection.CreateRule = nil
		value_collection.UpdateRule = nil
		return dao.SaveCollection(value_collection)
	})
}
