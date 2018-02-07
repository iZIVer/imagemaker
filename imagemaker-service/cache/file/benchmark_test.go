package file

import (
	"testing"

	"github.com/iZIVer/imagemaker/imagemaker-service"
	cachetest "github.com/iZIVer/imagemaker/imagemaker-service/cache/_test"
	"github.com/iZIVer/imagemaker/src"
)

func BenchmarkGet(b *testing.B) {
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
			cch := newTestCache()
			cachetest.BenchmarkGet(b, cch, 1, tc.im)
		})
	}
}
