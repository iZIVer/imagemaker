// Package _test provides utilities for imageserver/cache.Cache testing.
package _test

import (
	"sync"
	"testing"

	"github.com/iZIVer/imagemaker/compare"
	"github.com/iZIVer/imagemaker/imagemaker-service"
	imageserver_cache "github.com/iZIVer/imagemaker/imagemaker-service/cache"
	"github.com/iZIVer/imagemaker/src"
)

const (
	// KeyValid is a valid Cache key (with content)
	KeyValid = "test"
	// KeyMiss is an invalid Cache key (without content)
	KeyMiss = "unknown"
)

// TestGetSet is a helper to test imageserver/cache.Cache.Get()/Set().
func TestGetSet(t *testing.T, cache imageserver_cache.Cache) {
	err := cache.Set(KeyValid, testdata.Medium, imageserver.Params{})
	if err != nil {
		t.Fatal(err)
	}
	im, err := cache.Get(KeyValid, imageserver.Params{})
	if err != nil {
		t.Fatal(err)
	}
	if im == nil {
		t.Fatal("image nil")
	}
	diff := compare.Compare(im, testdata.Medium)
	if len(diff) != 0 {
		t.Fatalf("images not equal, diff:\n%+v", diff)
	}
}

// TestGetMiss is a helper to test imageserver/cache.Cache.Get() with a "cache miss".
func TestGetMiss(t *testing.T, cache imageserver_cache.Cache) {
	im, err := cache.Get(KeyMiss, imageserver.Params{})
	if err != nil {
		t.Fatal(err)
	}
	if im != nil {
		t.Fatal("image not nil")
	}
}

// MapCache is a simple imageserver/cache.Cache implementation (it wraps a map) for tests.
type MapCache struct {
	mutex sync.RWMutex
	data  map[string]*imageserver.Image
}

// NewMapCache creates a new CacheMap.
func NewMapCache() *MapCache {
	return &MapCache{
		data: make(map[string]*imageserver.Image),
	}
}

// Get implements imageserver/cache.Cache.
func (cache *MapCache) Get(key string, params imageserver.Params) (*imageserver.Image, error) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	return cache.data[key], nil
}

// Set implements imageserver/cache.Cache.
func (cache *MapCache) Set(key string, im *imageserver.Image, params imageserver.Params) error {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	cache.data[key] = im
	return nil
}
