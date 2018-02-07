package crop

import (
	"net/http"
	"testing"

	"github.com/iZIVer/imagemaker/compare"
	"github.com/iZIVer/imagemaker/imagemaker-service"
	imageserver_http "github.com/iZIVer/imagemaker/imagemaker-service/http"
)

var _ imageserver_http.Parser = &Parser{}

func TestParse(t *testing.T) {
	ps := &Parser{}
	for _, tc := range []struct {
		name               string
		url                string
		expectedParams     imageserver.Params
		expectedParamError string
	}{
		{
			name:           "Empty",
			url:            "http://localhost",
			expectedParams: imageserver.Params{},
		},
		{
			name: "Valid",
			url:  "http://localhost?crop=1,2|3,4",
			expectedParams: imageserver.Params{param: imageserver.Params{
				"min_x": 1,
				"min_y": 2,
				"max_x": 3,
				"max_y": 4,
			}},
		},
		{
			name:               "Invalid",
			url:                "http://localhost?crop=invalid",
			expectedParamError: "crop",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tc.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			params := imageserver.Params{}
			err = ps.Parse(req, params)
			if err != nil {
				if err, ok := err.(*imageserver.ParamError); ok && tc.expectedParamError == err.Param {
					return
				}
				t.Fatal(err)
			}
			if tc.expectedParamError != "" {
				t.Fatalf("no error, expected: %s", tc.expectedParamError)
			}
			diff := compare.Compare(params, tc.expectedParams)
			if len(diff) != 0 {
				t.Fatalf("unexpected params:\ngot: %v\nwant: %v\ndiff:\n%+v", params, tc.expectedParams, diff)
			}
		})
	}
}

func TestResolve(t *testing.T) {
	ps := &Parser{}
	for _, tc := range []struct {
		name     string
		param    string
		expected string
	}{
		{
			name:     "Param",
			param:    param,
			expected: param,
		},
		{
			name:     "MinX",
			param:    param + ".min_x",
			expected: param,
		},
		{
			name:     "Other",
			param:    "foobar",
			expected: "",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			httpParam := ps.Resolve(tc.param)
			if httpParam != tc.expected {
				t.Logf("param %s", tc.param)
				t.Fatalf("unexpected result: got '%s', want %s''", httpParam, tc.expected)
			}
		})
	}
}
