package nfntresize

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/iZIVer/imagemaker/imagemaker-service"
	imageserver_http "github.com/iZIVer/imagemaker/imagemaker-service/http"
)

var _ imageserver_http.Parser = &Parser{}

func TestParse(t *testing.T) {
	p := &Parser{}
	for _, tc := range []struct {
		name               string
		query              url.Values
		expectedParams     imageserver.Params
		expectedParamError string
	}{
		{
			name: "Empty",
		},
		{
			name:  "Width",
			query: url.Values{"width": {"100"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"width": 100,
			}},
		},
		{
			name:  "Height",
			query: url.Values{"height": {"100"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"height": 100,
			}},
		},
		{
			name:  "Interpolation",
			query: url.Values{"interpolation": {"lanczos3"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"interpolation": "lanczos3",
			}},
		},
		{
			name:  "Mode",
			query: url.Values{"mode": {"resize"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"mode": "resize",
			}},
		},
		{
			name:               "WidthInvalid",
			query:              url.Values{"width": {"invalid"}},
			expectedParamError: globalParam + ".width",
		},
		{
			name:               "HeightInvalid",
			query:              url.Values{"height": {"invalid"}},
			expectedParamError: globalParam + ".height",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			u := &url.URL{
				Scheme:   "http",
				Host:     "localhost",
				RawQuery: tc.query.Encode(),
			}
			req, err := http.NewRequest("GET", u.String(), nil)
			if err != nil {
				t.Fatal(err)
			}
			params := imageserver.Params{}
			err = p.Parse(req, params)
			if err != nil {
				if err, ok := err.(*imageserver.ParamError); ok && tc.expectedParamError == err.Param {
					return
				}
				t.Fatal(err)
			}
			if params.String() != tc.expectedParams.String() {
				t.Fatalf("unexpected params: got %s, want %s", params, tc.expectedParams)
			}
		})
	}
}

func TestResolve(t *testing.T) {
	p := &Parser{}
	httpParam := p.Resolve(globalParam + ".width")
	if httpParam != "width" {
		t.Fatal("not equal")
	}
}

func TestResolveNoMatch(t *testing.T) {
	p := &Parser{}
	httpParam := p.Resolve("foo")
	if httpParam != "" {
		t.Fatal("not equal")
	}
}
