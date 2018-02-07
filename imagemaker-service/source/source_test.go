package source

import (
	"fmt"
	"testing"

	"github.com/iZIVer/imagemaker/imagemaker-service"
)

var _ imageserver.Server = &Server{}

func TestServer(t *testing.T) {
	srv := &Server{
		Server: imageserver.ServerFunc(func(params imageserver.Params) (*imageserver.Image, error) {
			if !params.Has(Param) {
				t.Fatal("no source param")
			}
			if params.Has("foo") {
				t.Fatal("unexpected param")
			}
			return &imageserver.Image{}, nil
		}),
	}
	_, err := srv.Get(imageserver.Params{
		Param: "source",
		"foo": "bar",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestServerErrorServer(t *testing.T) {
	srv := &Server{
		Server: imageserver.ServerFunc(func(params imageserver.Params) (*imageserver.Image, error) {
			return nil, fmt.Errorf("error")
		}),
	}
	_, err := srv.Get(imageserver.Params{Param: "source"})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestServerErrorNoSource(t *testing.T) {
	srv := &Server{}
	_, err := srv.Get(imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ParamError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}
