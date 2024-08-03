package cache

import (
	"os"
	"testing"
)

func TestCacheCRUD(t *testing.T) {
	cachePath := "/tmp/cache.json"
	defer func() {
		os.Remove(cachePath)
	}()
	var cache *Cache
	t.Run("TestCacheLoad", func(t *testing.T) {
		cache = NewCache(cachePath)
		if len(cache.processedFiles) != 0 {
			t.Errorf("Cache should be empty, but got %d", len(cache.processedFiles))
		}
	})
	t.Run("TestCacheSave", func(t *testing.T) {
		cache.AddProcessedFile("test")
		err := cache.save()
		if err != nil {
			t.Errorf("Failed to save cache.")
		}
		if cache.IsProcessed("test") == false {
			t.Errorf("Cache should have test.")
		}
	})
	t.Run("TestCacheLoad", func(t *testing.T) {
		cache2 := NewCache(cachePath)
		if cache2.IsProcessed("test") == false {
			t.Errorf("Cache should have test.")
		}
	})
}
