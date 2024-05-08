package main

import (
	"fmt"
	"log"
	"path"

	"fyne.io/fyne/v2"
)

func processFile(filePath string, window *fyne.Window) {
	log.Printf("Processing file %q\n", filePath)

	fileExt := path.Ext(filePath)

	switch fileExt {
	case ".msg":
		fallthrough
	case ".eml":
		log.Println("Parsing email file")
		processEmailFile(filePath, window)
	case ".docx":
		log.Println("Parsing document file")
	default:
		launchInfoDialog("Unknown file", fmt.Sprintf("Can only inspect known file types, not %s", fileExt), window)
	}
}

func processEmailFile(filePath string, window *fyne.Window) {
	panic("unimplemented")
}
