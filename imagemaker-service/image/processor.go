package image

import (
	"image"

	"github.com/iZIVer/imagemaker/imagemaker-service"
)

// Processor processes a Go Image.
//
// It can return the given Image, but should not modify it.
type Processor interface {
	Process(image.Image, imageserver.Params) (image.Image, error)
	Changer
}

// ProcessorFunc is a Processor func.
type ProcessorFunc func(image.Image, imageserver.Params) (image.Image, error)

// Process implements Processor.
func (f ProcessorFunc) Process(nim image.Image, params imageserver.Params) (image.Image, error) {
	return f(nim, params)
}

// Change implements Processor.
func (f ProcessorFunc) Change(params imageserver.Params) bool {
	return true
}

// ListProcessor is a Processor implementation that wrap a list of Processor.
type ListProcessor []Processor

// Process implements Processor.
func (prc ListProcessor) Process(nim image.Image, params imageserver.Params) (image.Image, error) {
	for _, p := range prc {
		var err error
		nim, err = p.Process(nim, params)
		if err != nil {
			return nil, err
		}
	}
	return nim, nil
}

// Change implements Processor.
func (prc ListProcessor) Change(params imageserver.Params) bool {
	for _, p := range prc {
		if p.Change(params) {
			return true
		}
	}
	return false
}

// ChangeProcessor is a Processor implementation that alway return true for the Change method.
type ChangeProcessor struct {
	Processor
}

// Change implements Processor.
func (prc *ChangeProcessor) Change(params imageserver.Params) bool {
	return true
}
