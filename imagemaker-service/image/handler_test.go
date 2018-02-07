package image

import (
	"fmt"
	"image"
	"testing"

	"github.com/iZIVer/imagemaker/imagemaker-service"
	"github.com/iZIVer/imagemaker/src"
)

var _ imageserver.Handler = &Handler{}

func TestHandler(t *testing.T) {
	hdr := &Handler{}
	_, err := hdr.Handle(testdata.Medium, imageserver.Params{"quality": 85})
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandlerNoChange(t *testing.T) {
	hdr := &Handler{}
	im, err := hdr.Handle(testdata.Medium, imageserver.Params{})
	if err != nil {
		t.Fatal(err)
	}
	if im != testdata.Medium {
		t.Fatal("not equal")
	}
}

func TestHandlerFormat(t *testing.T) {
	hdr := &Handler{}
	_, err := hdr.Handle(testdata.Medium, imageserver.Params{"format": "jpeg"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandlerProcessor(t *testing.T) {
	hdr := &Handler{
		Processor: ProcessorFunc(func(nim image.Image, params imageserver.Params) (image.Image, error) {
			return nim, nil
		}),
	}
	_, err := hdr.Handle(testdata.Medium, imageserver.Params{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandlerErrorFormatParam(t *testing.T) {
	hdr := &Handler{}
	_, err := hdr.Handle(testdata.Medium, imageserver.Params{"format": "unknown"})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ParamError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestHandlerErrorFormatImage(t *testing.T) {
	hdr := &Handler{}
	_, err := hdr.Handle(&imageserver.Image{Format: "unknown"}, imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ImageError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestHandlerErrorDecode(t *testing.T) {
	hdr := &Handler{}
	_, err := hdr.Handle(testdata.Invalid, imageserver.Params{"format": "jpeg"})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ImageError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestHandlerErrorProcessor(t *testing.T) {
	hdr := &Handler{
		Processor: ProcessorFunc(func(nim image.Image, params imageserver.Params) (image.Image, error) {
			return nil, fmt.Errorf("error")
		}),
	}
	_, err := hdr.Handle(testdata.Medium, imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestHandlerErrorEncode(t *testing.T) {
	hdr := &Handler{}
	_, err := hdr.Handle(testdata.Medium, imageserver.Params{"quality": 9001})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ParamError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}
