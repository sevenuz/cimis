package routes

import (
	"cimis/ui"
	"net/http"
	"testing"

	"github.com/pocketbase/pocketbase/tests"
)

func TestUiEndpoint(t *testing.T) {
	setupTestApp := func() (*tests.TestApp, error) {
		testApp, err := tests.NewTestApp("../test_pb_data")
		if err != nil {
			return nil, err
		}
		// no need to cleanup since scenario.Test() will do that for us
		// defer testApp.Cleanup()

		Register(testApp)

		return testApp, nil
	}

	scenarios := []tests.ApiScenario{
		{
			Name:            "POST index",
			Method:          http.MethodPost,
			Url:             "/",
			ExpectedStatus:  405,
			ExpectedContent: []string{""},
			TestAppFactory:  setupTestApp,
		},
		{
			Name:            "GET index",
			Method:          http.MethodGet,
			Url:             "/",
			ExpectedStatus:  200,
			ExpectedContent: []string{ui.Index_file},
			TestAppFactory:  setupTestApp,
		},
		{
			Name:            "GET index with random path",
			Method:          http.MethodGet,
			Url:             "/xyz/zyx/yzx",
			ExpectedStatus:  200,
			ExpectedContent: []string{ui.Index_file},
			TestAppFactory:  setupTestApp,
		},
		{
			Name:            "GET index with random path 2",
			Method:          http.MethodGet,
			Url:             "/page/partypeople/__4322?reset_buffer=1",
			ExpectedStatus:  200,
			ExpectedContent: []string{ui.Index_file},
			TestAppFactory:  setupTestApp,
		},
		{
			Name:            "GET icon",
			Method:          http.MethodGet,
			Url:             "/icon.png",
			ExpectedStatus:  200,
			ExpectedContent: []string{ui.Icon_file},
			TestAppFactory:  setupTestApp,
		},
		{
			Name:            "GET static normal",
			Method:          http.MethodGet,
			Url:             "/_app/version.json",
			ExpectedStatus:  200,
			ExpectedContent: []string{""},
			TestAppFactory:  setupTestApp,
		},
		{
			Name:            "GET static unknown",
			Method:          http.MethodGet,
			Url:             "/_app/somefile",
			ExpectedStatus:  404,
			ExpectedContent: []string{`{"code":404,"message":"Not Found.","data":{}}`},
			TestAppFactory:  setupTestApp,
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
