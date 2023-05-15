package models

import (
	"github.com/pocketbase/pocketbase/tests"
	"testing"
)

func TestSurveyEndpoint(t *testing.T) {
	setupTestApp := func() *tests.TestApp {
		testApp, err := tests.NewTestApp("../test_pb_data")
		if err != nil {
			t.Errorf("Should not fail while creating testapp, but %d", err)
		}

		return testApp
	}
	app := setupTestApp()

	lk := LanguageKey{}
	if lk.TableName() != N_LanguageKey {
		t.Fail()
	}

	_, err := GetLanguageKey(app, "unknown")
	if err == nil {
		t.Errorf("Should fail for unknown key")
	}

	_, err = GetLanguageKey(app, "kutizv8q6d9cbb8")
	if err != nil {
		t.Errorf("Should find key, but %d", err)
	}

}
