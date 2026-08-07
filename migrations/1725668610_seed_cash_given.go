package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		key_collection, err := dao.FindCollectionByNameOrId("language_key")
		if err != nil {
			return err
		}
		value_collection, err := dao.FindCollectionByNameOrId("language_value")
		if err != nil {
			return err
		}

		key_record, err := dao.FindFirstRecordByData("language_key", "name", "ui_cash_given")
		if err != nil {
			key_record = models.NewRecord(key_collection)
			key_record.Set("name", "ui_cash_given")
			if err := dao.SaveRecord(key_record); err != nil {
				return err
			}
		}

		if err := ensureLanguageValue(dao, value_collection, key_record.Id, "de", "Erhaltenes Bargeld"); err != nil {
			return err
		}
		return ensureLanguageValue(dao, value_collection, key_record.Id, "en", "Cash given")
	}, func(db dbx.Builder) error {
		return nil
	})
}
