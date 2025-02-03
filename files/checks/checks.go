package checks

import (
	"os"
)

// FolderExists checks the folder exists and is a folder, return boolean result
func FolderExists(folderPath string) bool {
	stat, err := os.Stat(folderPath)

	if os.IsNotExist(err) {
		return false
	}

	if !stat.IsDir() {
		return false
	}

	return true
}

// FileExists checks the folder exists and is a folder, return boolean result
func FileExists(filePath string) bool {
	stat, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		return false
	}

	if stat.IsDir() {
		return false
	}

	return true
}
