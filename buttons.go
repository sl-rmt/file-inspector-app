package main

import (
	"file-inspector/files"
	"fmt"
	"log"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

const (
	defaultSelectText = "\n\n\n\n\n\t\t\t\t\t\tSelect a file to analyse..."
)

func onOpenButtonClicked() {
	log.Println("Select file was clicked!")

	// lock open button until we're done processing
	openButton.Disable()
	onChosen := onFileChosen
	dialog.ShowFileOpen(onChosen, window)
}

func onFileChosen(f fyne.URIReadCloser, err error) {
	// make sure the view is reset
	onResetButtonClicked()

	if err != nil {
		log.Printf("Error from file picker: %s\n", err.Error())
		return
	}
	if f == nil {
		log.Println("Nil result from file picker")
		return
	}

	// file chosen - update UI
	iconSeparator.Show() // show the separator above the row of icons

	log.Printf("chosen: %v", f.URI())

	filePathString := f.URI().Path()
	fileExtension := filepath.Ext(filePathString)

	// check chosen file
	if !fileOkayToProcess(fileExtension) {
		launchInfoDialog("Unsupported File Type", "File is not currently supported", &window)
		openButton.Enable()
		return
	}

	// launch progress bad dialog
	progress := launchProcessingDialog(&window)

	// get and set the file properties
	properties, err := files.GetFileProperties(filePathString)

	if err != nil {
		analysisTextBS.Set(fmt.Sprintf("Error processing file: %q\n", err.Error()))
		errorLabel.Show()
		errorIcon.Show()
	} else {
		// set the values
		fileNameBS.Set(properties.FileName)
		fileTypeBS.Set(properties.FileType)
		fileHashBS.Set(properties.Hash)
		fileSizeBS.Set(properties.Size)

		matches, explanation := files.CheckMime(fileExtension, properties.FileType)

		if !matches {
			showIconAndLabel(dangerIcon, dangerLabel, dangerSeparator)

			log.Printf("Mismatched extension and MIME type. %s", explanation)
			analysisTextBS.Set(fmt.Sprintf("Mismatched extension and MIME type.\n\n%s", explanation))

			// tidy up and return, as we don't want to parse the file as the wrong type
			progress.Hide()
			openButton.Enable()
			return
		}

		// process the file and show the analysis
		result := files.ProcessFile(filePathString)

		if result.Completed {
			showIconAndLabel(processedIcon, processedLabel, processedSeparator)
		}

		if result.Error != nil {
			log.Printf("Processing complete with error: %q\n", result.Error.Error())

			// notify the user
			launchErrorDialog(result.Error, window)
			analysisTextBS.Set(result.Error.Error())
			showIconAndLabel(errorIcon, errorLabel, errorSeparator)
		}

		log.Println("File processing done")
		metadataTextBS.Set(result.Metadata)
		analysisTextBS.Set(result.Analysis)

		if result.Dangerous {
			showIconAndLabel(dangerIcon, dangerLabel, dangerSeparator)
			launchInfoDialog("Potentially Dangerous File", "Warning: dangerous file found", &window)
		}
	}

	progress.Hide()
	openButton.Enable()
}

func onResetButtonClicked() {
	log.Println("Reset was clicked!")

	// blank all the fields
	analysisTextBS.Set(defaultSelectText)
	metadataTextBS.Set(defaultSelectText)
	fileNameBS.Set("")
	fileTypeBS.Set("")
	fileHashBS.Set("")
	fileSizeBS.Set("")

	// clear and hide icons
	iconSeparator.Hide()
	hideIconAndLabel(processedIcon, processedLabel, processedSeparator)
	hideIconAndLabel(errorIcon, errorLabel, errorSeparator)
	hideIconAndLabel(dangerIcon, dangerLabel, dangerSeparator)
}
