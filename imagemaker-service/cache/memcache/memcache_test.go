package memcache

import (
	"strings"
	"testing"

	memcache_impl "github.com/bradfitz/gomemcache/memcache"
	"github.com/iZIVer/imagemaker/imagemaker-service"
	imageserver_cache "github.com/iZIVer/imagemaker/imagemaker-service/cache"
	cachetest "github.com/iZIVer/imagemaker/imagemaker-service/cache/_test"
	"github.com/iZIVer/imagemaker/src"
)

var _ imageserver_cache.Cache = &Cache{}

func TestGetSet(t *testing.T) {
	cache := newTestCache(t)
	cachetest.TestGetSet(t, cache)
}

func TestGetMiss(t *testing.T) {
	cache := newTestCache(t)
	cachetest.TestGetMiss(t, cache)
}

func TestGetErrorServer(t *testing.T) {
	cache := newTestCacheInvalidServer()
	_, err := cache.Get(cachetest.KeyValid, imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestSetErrorServer(t *testing.T) {
	cache := newTestCacheInvalidServer()
	err := cache.Set(cachetest.KeyValid, testdata.Medium, imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetErrorUnmarshal(t *testing.T) {
	cache := newTestCache(t)
	data, err := testdata.Medium.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	data = data[:len(data)-1]
	err = cache.setData(cachetest.KeyValid, data)
	if err != nil {
		t.Fatal(err)
	}
	_, err = cache.Get(cachetest.KeyValid, imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ImageError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestSetErrorMarshal(t *testing.T) {
	cache := newTestCache(t)
	im := &imageserver.Image{
		Format: strings.Repeat("a", imageserver.ImageFormatMaxLen+1),
	}
	err := cache.Set(cachetest.KeyValid, im, imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ImageError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func newTestCache(tb testing.TB) *Cache {
	cache := newTestCacheWithClient(newTestClient("localhost:11211"))
	checkTestCacheAvailable(tb, cache)
	return cache
}

func newTestCacheInvalidServer() *Cache {
	return newTestCacheWithClient(newTestClient("localhost:11311"))
}

func newTestCacheWithClient(client *memcache_impl.Client) *Cache {
	return &Cache{
		Client: client,
	}
}

func newTestClient(server string) *memcache_impl.Client {
	return memcache_impl.New(server)
}

func checkTestCacheAvailable(tb testing.TB, cache *Cache) {
	err := cache.Client.Set(&memcache_impl.Item{
		Key:   "ping",
		Value: []byte("ping"),
	})
	if err != nil {
		tb.Skip(err)
	}
}
