package file_scan_domain

import (
	"os"
	"path/filepath"
	"strings"
)

func ScanFiles(rootPath string) []string {
	var files []string
	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || strings.HasPrefix(info.Name(), ".") {
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files
}
