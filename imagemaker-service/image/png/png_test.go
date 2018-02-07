package png

import (
	"testing"

	"github.com/iZIVer/imagemaker/imagemaker-service"
	imageserver_image "github.com/iZIVer/imagemaker/imagemaker-service/image"
	imageserver_image_test "github.com/iZIVer/imagemaker/imagemaker-service/image/_test"
)

var _ imageserver_image.Encoder = &Encoder{}

func TestEncoder(t *testing.T) {
	imageserver_image_test.TestEncoder(t, &Encoder{}, "png")
}

func TestEncoderChange(t *testing.T) {
	c := (&Encoder{}).Change(imageserver.Params{})
	if c {
		t.Fatal("not false")
	}
}
