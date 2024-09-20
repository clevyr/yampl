package util

import (
	"path/filepath"
	"strings"
)

func IsYaml(path string) bool {
	ext := filepath.Ext(path)
	return strings.EqualFold(ext, ".yaml") || strings.EqualFold(ext, ".yml")
}
