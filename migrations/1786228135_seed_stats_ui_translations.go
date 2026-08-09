package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

var statsTranslations = map[string][2]string{
	"ui_stats":                   {"Statistik", "Statistics"},
	"ui_revenue":                 {"Umsatz", "Revenue"},
	"ui_all_events":              {"Alle Events", "All events"},
	"ui_per_product":             {"Pro Produkt", "Per product"},
	"ui_per_payment_method":      {"Pro Zahlungsmethode", "Per payment method"},
	"ui_per_bar":                 {"Pro Bar", "Per bar"},
	"ui_over_time":               {"Im Zeitverlauf", "Over time"},
	"ui_top_sellers":             {"Meistverkauft", "Top sellers"},
	"ui_bottom_sellers":          {"Wenigstverkauft", "Bottom sellers"},
	"ui_units_sold":              {"Verkaufte Einheiten", "Units sold"},
	"ui_free_units":              {"Gratis-Einheiten", "Free units"},
	"ui_paid_units":              {"Bezahlte Einheiten", "Paid units"},
	"ui_revenue_per_bar":         {"Umsatz pro Bar", "Revenue per bar"},
	"ui_avg_order_size":          {"Ø Bestellgröße", "Avg. order size"},
	"ui_order_count":             {"Anzahl Bestellungen", "Order count"},
	"ui_volume":                  {"Umsatzvolumen", "Volume"},
	"ui_avg_order_value":         {"Ø Bestellwert", "Avg. order value"},
	"ui_fee_cost":                {"Kartengebühren", "Card fee cost"},
	"ui_busiest_hour":            {"Umsatz pro Stunde und Bar", "Revenue per hour and bar"},
	"ui_revenue_over_time":       {"Umsatzverlauf", "Revenue over time"},
	"ui_product_sales_over_time": {"Produktverkäufe im Zeitverlauf", "Product sales over time"},
	"ui_bucket_size":             {"Intervall", "Bucket size"},
	"ui_product":                 {"Produkt", "Product"},
	"ui_wheel_results":           {"Glücksrad-Ergebnisse", "Lucky wheel results"},
	"ui_wheel_spins":             {"Radspins", "Wheel spins"},
	"ui_avg_reward_value":        {"Ø Belohnungswert pro Spin", "Avg. reward value per spin"},
	"ui_accumulated_revenue":     {"Kumulierter Umsatz", "Accumulated revenue"},
	"ui_avg_order_size_explanation": {
		"Durchschnittlicher Gesamtwert der kompletten Bestellung - für alle Bestellungen, die dieses Produkt enthielten, nicht der Preis des Produkts selbst. Ein hoher Wert bedeutet: Bestellungen mit diesem Produkt sind meist größere Warenkörbe (oft zusammen mit mehreren anderen Sachen bestellt) - ein nützliches Signal für Bundling/Zusatzverkäufe. Liegt der Wert nah am eigenen Preis des Produkts, wird es meist alleine gekauft.",
		"Average total value of the whole order - across every order that included this product, not the product's own price. A high number means orders containing it tend to be bigger baskets (often bought alongside several other things) - a useful bundle/upsell signal. A number close to the product's own price means it's mostly bought alone.",
	},
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

		for name, t := range statsTranslations {
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
