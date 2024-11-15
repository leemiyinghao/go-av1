package file_scan_domain

import (
	"testing"
	"path/filepath"

	"github.com/stretchr/testify/assert"
)

func TestScanFiles(t *testing.T) {
	rootPath := "./test_assets/"
	files := ScanFiles(rootPath)

	assert.Equal(t, files, []string{
		filepath.Join("test_assets", "file1.txt"),
		filepath.Join("test_assets", "subdir", "file2.txt"),
	})
}

func TestScanFiles_Directory(t *testing.T) {
	rootPath := "./test_assets/subdir"
	files := ScanFiles(rootPath)

	assert.Equal(t, files, []string{
		filepath.Join("test_assets", "subdir", "file2.txt"),
	})
}

