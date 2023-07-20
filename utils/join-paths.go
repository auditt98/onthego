package utils

import (
	"path/filepath"
)

func JoinPaths(elements ...string) string {
	return filepath.Join(elements...)
}
