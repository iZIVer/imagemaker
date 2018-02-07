// Package nfntresize provides a imageserver/http.Parser implementation for imageserver/image/nfntresize.Processor.
package nfntresize

import (
	"net/http"
	"strings"

	"github.com/iZIVer/imagemaker/imagemaker-service"
	imageserver_http "github.com/iZIVer/imagemaker/imagemaker-service/http"
)

const (
	globalParam = "nfntresize"
)

// Parser is a imageserver/http.Parser implementation for imageserver/image/nfntresize.Processor.
//
// It takes the params from the HTTP URL query and stores them in a Params.
// This Params is added to the given Params at the key "nfntresize".
//
// See imageserver/image/nfntresize.Processor for params list.
type Parser struct{}

// Parse implements imageserver/http.Parser.
func (parser *Parser) Parse(req *http.Request, params imageserver.Params) error {
	p := imageserver.Params{}
	err := parse(req, p)
	if err != nil {
		if err, ok := err.(*imageserver.ParamError); ok {
			err.Param = globalParam + "." + err.Param
		}
		return err
	}
	if !p.Empty() {
		params.Set(globalParam, p)
	}
	return nil
}

func parse(req *http.Request, params imageserver.Params) error {
	if err := imageserver_http.ParseQueryInt("width", req, params); err != nil {
		return err
	}
	if err := imageserver_http.ParseQueryInt("height", req, params); err != nil {
		return err
	}
	imageserver_http.ParseQueryString("interpolation", req, params)
	imageserver_http.ParseQueryString("mode", req, params)
	return nil
}

// Resolve implements imageserver/http.Parser.
func (parser *Parser) Resolve(param string) string {
	if !strings.HasPrefix(param, globalParam+".") {
		return ""
	}
	return strings.TrimPrefix(param, globalParam+".")
}
