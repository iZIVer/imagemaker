package image

import (
	"github.com/iZIVer/imagemaker/imagemaker-service"
)

// Server is a imageserver.Server implementation that gets the Image from a Provider.
//
// It uses the "format" param or the DefaultFormat variable to determine which Encoder is used.
type Server struct {
	Provider      Provider
	DefaultFormat string
}

// Get implements imageserver.Server.
func (srv *Server) Get(params imageserver.Params) (*imageserver.Image, error) {
	enc, format, err := getEncoderFormat(srv.DefaultFormat, params)
	if err != nil {
		return nil, err
	}
	nim, err := srv.Provider.Get(params)
	if err != nil {
		return nil, err
	}
	im, err := encode(nim, format, enc, params)
	if err != nil {
		return nil, err
	}
	return im, nil
}
