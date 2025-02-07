package main

import "path/filepath"

// check the file extension is one we support
// return false if not
func fileOkayToProcess(filePathString string) bool {
	ext := filepath.Ext(filePathString)

	switch ext {
	case ".eml":
		fallthrough
	case ".msg":
		fallthrough
	case ".pdf":
		return true
	default:
		return false
	}
}
