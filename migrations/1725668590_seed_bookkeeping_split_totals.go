package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

var bookkeepingSplitTranslations = map[string][2]string{
	"ui_earned":          {"Erwirtschaftet (Alle Bezahlmethoden)", "Earned (all payment methods)"},
	"ui_booked_out":      {"Ausgebucht", "Booked out"},
	"ui_expected_in_box": {"Erwartet in der Kasse", "Expected in the box"},
}

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

		for name, t := range bookkeepingSplitTranslations {
			key_record, err := dao.FindFirstRecordByData("language_key", "name", name)
			if err != nil {
				key_record = models.NewRecord(key_collection)
				key_record.Set("name", name)
				if err := dao.SaveRecord(key_record); err != nil {
					return err
				}
			}

			if err := ensureLanguageValue(dao, value_collection, key_record.Id, "de", t[0]); err != nil {
				return err
			}
			if err := ensureLanguageValue(dao, value_collection, key_record.Id, "en", t[1]); err != nil {
				return err
			}
		}

		return nil
	}, func(db dbx.Builder) error {
		return nil
	})
}
