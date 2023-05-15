package models

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

const N_LanguageKey = "language_key" // the name of the collection

var _ models.Model = (*LanguageKey)(nil)

type LanguageKey struct {
	models.BaseModel

	Name string `db:"name" json:"name"`
}

func (m *LanguageKey) TableName() string {
	return N_LanguageKey
}

func GetLanguageKey(app core.App, id string) (*LanguageKey, error) {
	s := &LanguageKey{}
	err := app.Dao().ModelQuery(&LanguageKey{}).
		AndWhere(dbx.HashExp{"id": id}).
		Limit(1).
		One(s)
	if err != nil {
		return nil, apis.NewNotFoundError("The language_key does not exist: " + id, err)
	}
	return s, nil
}

const JsonLanguageKey = `{
        "id": "5rxmpfs22bhonfh",
        "name": "` + N_LanguageKey + `",
        "type": "base",
        "system": false,
        "schema": [
            {
                "id": "4un9spdd",
                "name": "name",
                "type": "text",
                "system": false,
                "required": true,
                "unique": true,
                "options": {
                    "min": null,
                    "max": null,
                    "pattern": ""
                }
            }
        ],
        "listRule": null,
        "viewRule": "",
        "createRule": null,
        "updateRule": null,
        "deleteRule": null,
        "options": {}
    }`
