package gif

import (
	"testing"

	"github.com/iZIVer/imagemaker/imagemaker-service"
	imageserver_image_test "github.com/iZIVer/imagemaker/imagemaker-service/image/_test"
	_ "github.com/iZIVer/imagemaker/imagemaker-service/image/jpeg"
	"github.com/iZIVer/imagemaker/src"
)

func BenchmarkEncoder(b *testing.B) {
	enc := &Encoder{}
	params := imageserver.Params{}
	for _, tc := range []struct {
		name string
		im   *imageserver.Image
	}{
		{"Small", testdata.Small},
		{"Medium", testdata.Medium},
		{"Large", testdata.Large},
		{"Huge", testdata.Huge},
	} {
		b.Run(tc.name, func(b *testing.B) {
			imageserver_image_test.BenchmarkEncoder(b, enc, tc.im, params)
		})
	}
}
