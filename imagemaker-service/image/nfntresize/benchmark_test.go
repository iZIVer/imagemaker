package nfntresize

import (
	"testing"

	"github.com/iZIVer/imagemaker/imagemaker-service"
	imageserver_image "github.com/iZIVer/imagemaker/imagemaker-service/image"
	"github.com/iZIVer/imagemaker/src"
)

func BenchmarkSize(b *testing.B) {
	for _, tc := range []struct {
		name string
		im   *imageserver.Image
	}{
		{"Small", testdata.Small},
		{"Medium", testdata.Medium},
		{"Large", testdata.Large},
		{"Huge", testdata.Huge},
	} {
		benchmark(b, tc.name, tc.im, imageserver.Params{})
	}
}

func BenchmarkInterpolation(b *testing.B) {
	for _, it := range []string{
		"nearest_neighbor",
		"bilinear",
		"bicubic",
		"mitchell_netravali",
		"lanczos2",
		"lanczos3",
	} {
		benchmark(b, it, testdata.Medium, imageserver.Params{
			"interpolation": it,
		})
	}
}

func benchmark(b *testing.B, name string, im *imageserver.Image, params imageserver.Params) {
	nim, err := imageserver_image.Decode(im)
	if err != nil {
		b.Fatal(err)
	}
	params.Set("width", 100)
	params = imageserver.Params{
		param: params,
	}
	proc := &Processor{}
	b.Run(name, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := proc.Process(nim, params)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
