package cache

import (
	"encoding/json"
	"os"
	"sync"
)

// ProcessedFile is a struct that contains the path of a processed file, will be stored as json in the cache file.
type ProcessedFile struct {
	FilePath string `json:"filePath"`
}

type CacheFile struct {
	ProcessedFiles []ProcessedFile `json:"processedFiles"`
}

type Cache struct {
	processedFiles []ProcessedFile
	cacheFilePath  string
	lock           sync.Mutex
}

func NewCache(cacheFilePath string) *Cache {
	cache := &Cache{
		cacheFilePath:  cacheFilePath,
		processedFiles: make([]ProcessedFile, 0),
	}
	// try to load the cache file
	cache.Load()
	return cache
}

// Load reads the cache file and loads the processed files into the cache.
func (c *Cache) Load() error {
	cacheFile := CacheFile{}
	cacheFileContents, err := os.ReadFile(c.cacheFilePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(cacheFileContents, &cacheFile)
	if err != nil {
		return err
	}
	c.processedFiles = cacheFile.ProcessedFiles
	return nil
}

// Save writes the processed files in the cache to the cache file.
func (c *Cache) save() error {
	cacheFile := CacheFile{
		ProcessedFiles: c.processedFiles,
	}
	cacheFileContents, err := json.Marshal(cacheFile)
	if err != nil {
		return err
	}
	err = os.WriteFile(c.cacheFilePath, cacheFileContents, 0644)
	if err != nil {
		return err
	}
	return nil
}

// AddProcessedFile adds a processed file to the cache.
func (c *Cache) AddProcessedFile(filePath string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.processedFiles = append(c.processedFiles, ProcessedFile{filePath})
	c.save()
}

// IsProcessed checks if a file has been processed before.
func (c *Cache) IsProcessed(filePath string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	for _, processedFile := range c.processedFiles {
		if processedFile.FilePath == filePath {
			return true
		}
	}
	return false
}
