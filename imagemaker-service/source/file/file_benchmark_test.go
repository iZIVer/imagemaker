package file

import (
	"testing"

	"github.com/iZIVer/imagemaker/imagemaker-service"
	imageserver_source "github.com/iZIVer/imagemaker/imagemaker-service/source"
	"github.com/iZIVer/imagemaker/src"
)

func BenchmarkServerGet(b *testing.B) {
	srv := &Server{
		Root: testdata.Dir,
	}
	for _, tc := range []struct {
		name     string
		filename string
	}{
		{"Small", testdata.SmallFileName},
		{"Medium", testdata.MediumFileName},
		{"Large", testdata.LargeFileName},
		{"Huge", testdata.HugeFileName},
	} {
		b.Run(tc.name, func(b *testing.B) {
			params := imageserver.Params{
				imageserver_source.Param: tc.filename,
			}
			var bs int
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				im, err := srv.Get(params)
				if err != nil {
					b.Fatal(err)
				}
				bs = len(im.Data)
			}
			b.SetBytes(int64(bs))
		})
	}
}
