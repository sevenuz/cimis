package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

// Seeds default de/en translations for the generic app-chrome ui_ keys (login,
// buttons, bar/inventory/recipes/bookkeeping labels, etc). Deliberately excludes
// site-content keys (ui_title, ui_ticket_shop_url, social links, cookie banner,
// contact email, ...) since those hold real business content this migration has
// no business guessing at.
//
// Idempotent by design: every key/value is only created if missing, so running
// this against a DB that already has some of these translated does not clobber
// existing content. Because of that "only if missing" behavior, the down
// migration intentionally does nothing - deleting on rollback could not tell a
// pre-existing key from one this migration created, and leaving a few extra
// translation strings around is harmless, unlike guessing wrong and deleting
// content that predates this migration.
var defaultUiTranslations = map[string][2]string{
	// [0] = de, [1] = en
	"ui_add_inventory":                 {"Bestand hinzufügen", "Add inventory"},
	"ui_admin_only":                    {"Nur für Admins", "Admin only"},
	"ui_amount":                        {"Menge", "Amount"},
	"ui_back":                          {"Zurück", "Back"},
	"ui_bar":                           {"Bar", "Bar"},
	"ui_bar_confirmation_alert":        {"ACHTUNG: Pfand und Wechselgeld nicht vergessen!", "WARNING: Don't forget deposit and change!"},
	"ui_bar_selector":                  {"Bar", "Bar"},
	"ui_bookkeeping":                   {"Buchführungs-Möglichkeiten", "Bookkeeping"},
	"ui_bookout":                       {"Kassenschnitt", "Cash out"},
	"ui_cancel":                        {"Abbrechen", "Cancel"},
	"ui_cash_per_bar":                  {"Bargeld pro Bar", "Cash per bar"},
	"ui_color":                         {"Farbe", "Color"},
	"ui_date":                          {"Datum", "Date"},
	"ui_deactivated":                   {"Deaktiviert", "Deactivated"},
	"ui_discount_percentage_only_once": {"Rabatt kann nur einmal angewendet werden", "Discount can only be applied once"},
	"ui_edit":                          {"Bearbeiten", "Edit"},
	"ui_empty_selection":               {"Keine Auswahl getroffen", "No selection made"},
	"ui_event":                         {"Event", "Event"},
	"ui_free":                          {"Gratis", "Free"},
	"ui_ingredients":                   {"Zutaten", "Ingredients"},
	"ui_instructions":                  {"Anleitung", "Instructions"},
	"ui_invalid_date":                  {"Ungültiges Datum", "Invalid date"},
	"ui_inventory":                     {"Inventar", "Inventory"},
	"ui_is_wheel":                      {"Ist das Glücksrad", "Is the lucky wheel"},
	"ui_items":                         {"Artikel", "items"},
	"ui_logged_in_as":                  {"Angemeldet als", "Logged in as"},
	"ui_login":                         {"Anmelden", "Login"},
	"ui_logout":                        {"Abmelden", "Logout"},
	"ui_name":                          {"Name", "Name"},
	"ui_new":                           {"Neu", "New"},
	"ui_no_active_event":               {"Kein aktives Event gesetzt", "No active event set"},
	"ui_note":                          {"Notiz", "Note"},
	"ui_order":                         {"Reihenfolge", "Order"},
	"ui_order_saved":                   {"Bestellung gespeichert", "Order saved"},
	"ui_payment_method":                {"Zahlungsmethode", "Payment method"},
	"ui_pending":                       {"ausstehend", "pending"},
	"ui_price":                         {"Preis", "Price"},
	"ui_products":                      {"Produkte", "Products"},
	"ui_recipe_ingredients":            {"Rezeptzutaten", "Recipe ingredients"},
	"ui_recipes":                       {"Rezepte", "Recipes"},
	"ui_save":                          {"Speichern", "Save"},
	"ui_show":                          {"Anzeigen", "Show"},
	"ui_slug":                          {"Slug", "Slug"},
	"ui_total":                         {"Gesamt", "Total"},
	"ui_type":                          {"Typ", "Type"},
	"ui_unit":                          {"Einheit", "Unit"},
	"ui_use_custom_range":              {"Eigenen Zeitraum verwenden", "Use custom date range"},
	"ui_wheel_prompt":                  {"Welches Getränk wurde gewonnen?", "Which drink was won?"},
}

func ensureLanguageValue(dao *daos.Dao, collection *models.Collection, keyId string, iso string, value string) error {
	existing, err := dao.FindRecordsByExpr(
		"language_value",
		dbx.HashExp{"language_key": keyId, "iso": iso},
	)
	if err == nil && len(existing) > 0 {
		return nil // already translated, don't overwrite
	}

	record := models.NewRecord(collection)
	record.Set("language_key", keyId)
	record.Set("iso", iso)
	record.Set("value", value)
	return dao.SaveRecord(record)
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

		for name, t := range defaultUiTranslations {
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
		// intentionally a no-op, see the comment on defaultUiTranslations above
		return nil
	})
}
