package hooks

import (
	//"cimis/models"
	// "fmt"
	// "log"
	"strings"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	pb_models "github.com/pocketbase/pocketbase/models"
)

func contains(value string, array []string) bool {
	for _, a := range array {
		if a == value {
			return true
		}
	}
	return false
}

func validatePage(collection *pb_models.Collection, record *pb_models.Record) error {
	if collection.Name == "page" {
		// if record.GetString("type") == "container" {
		path, ok := record.Get("path").(string)
		if !ok {
			return apis.NewBadRequestError("Something wrent wrong reading the path", nil)
		}
		addons, ok := record.Get("addon").([]string)
		if !ok {
			return apis.NewBadRequestError("Something wrent wrong reading the addons", nil)
		}

		if strings.HasPrefix(path, "bar") && !contains("bar", addons) {
			return apis.NewBadRequestError("Path \"bar\" is reserved for the addon \"bar\".", nil)
		}

		if contains("bar", addons) {
			if len(addons) > 1 {
				return apis.NewBadRequestError("On addon bar, no other addon is selectable.", nil)
			}
			if path != "bar" {
				return apis.NewBadRequestError("Addon bar only works with path \"bar\".", nil)
			}
		}
		if contains("external_link", addons) {
			if len(addons) > 1 {
				return apis.NewBadRequestError("On addon external_link, no other addon is selectable.", nil)
			}
			if !strings.HasPrefix(path, "https://") {
				return apis.NewBadRequestError("On addon external_link, an external link with https is required.", nil)
			}
		}
		if contains("particle", addons) {
			// TODO create/update feature entry
		}
		// shows value and type of field
		// i := record.Get("addon")
		// log.Println(fmt.Sprintf("(%v, %T)\n", i, i))
	}
	return nil
}

func addPageValidationHook(app core.App) {
	app.OnRecordBeforeCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		return validatePage(e.Collection, e.Record)
	})
	app.OnRecordBeforeUpdateRequest().Add(func(e *core.RecordUpdateEvent) error {
		return validatePage(e.Collection, e.Record)
	})
}
