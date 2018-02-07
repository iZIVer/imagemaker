package gamma

import (
	"image"
	"testing"

	"github.com/iZIVer/imagemaker/imagemaker-service"
	imageserver_image "github.com/iZIVer/imagemaker/imagemaker-service/image"
	"github.com/iZIVer/imagemaker/src"
)

func BenchmarkProcessor(b *testing.B) {
	nim, err := imageserver_image.Decode(testdata.Medium)
	if err != nil {
		b.Fatal(err)
	}
	prc := NewProcessor(2.2, false)
	params := imageserver.Params{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := prc.Process(nim, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProcessorHighQuality(b *testing.B) {
	nim, err := imageserver_image.Decode(testdata.Medium)
	if err != nil {
		b.Fatal(err)
	}
	prc := NewProcessor(2.2, true)
	params := imageserver.Params{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := prc.Process(nim, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCorrectionProcessor(b *testing.B) {
	nim, err := imageserver_image.Decode(testdata.Medium)
	if err != nil {
		b.Fatal(err)
	}
	prc := NewCorrectionProcessor(
		imageserver_image.ProcessorFunc(func(nim image.Image, params imageserver.Params) (image.Image, error) {
			return nim, nil
		}),
		true,
	)
	params := imageserver.Params{}
	for i := 0; i < b.N; i++ {
		_, err := prc.Process(nim, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}
