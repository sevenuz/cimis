package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

// Split out from 1725668540_seed_default_ui_translations.go: that migration may
// already have run (and be recorded as applied) before ui_bar_confirmation_alert
// replaced ui_deposit_warning/ui_requires_deposit in its map, and an already-applied
// migration file is never re-run - so this key needs its own migration to actually
// reach databases that already migrated past 1725668540.
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

		key_record, err := dao.FindFirstRecordByData("language_key", "name", "ui_bar_confirmation_alert")
		if err != nil {
			key_record = models.NewRecord(key_collection)
			key_record.Set("name", "ui_bar_confirmation_alert")
			if err := dao.SaveRecord(key_record); err != nil {
				return err
			}
		}

		if err := ensureLanguageValue(dao, value_collection, key_record.Id, "de", "ACHTUNG: Pfand und Wechselgeld nicht vergessen!"); err != nil {
			return err
		}
		return ensureLanguageValue(dao, value_collection, key_record.Id, "en", "WARNING: Don't forget deposit and change!")
	}, func(db dbx.Builder) error {
		// intentionally a no-op, same reasoning as 1725668540_seed_default_ui_translations.go
		return nil
	})
}
