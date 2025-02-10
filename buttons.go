package main

import (
	"file-inspector/files"
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func onOpenButtonClicked() {
	log.Println("Select file was clicked!")

	// lock open button until we're done processing
	openButton.Disable()

	onChosen := func(f fyne.URIReadCloser, err error) {
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

		// check chosen file
		if !fileOkayToProcess(f.URI().Path()) {
			launchInfoDialog("Unsupported File Type", "File is not currently supported", &window)
			return
		}

		// make sure the view is reset
		onResetButtonClicked()

		progress := launchProcessingDialog(&window)

		// get and set the file properties
		properties, err := files.GetFileProperties(f.URI().Path())

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

			// process the file and show the analysis
			result := files.ProcessFile(f.URI().Path())

			if result.Completed {
				showIconAndLabel(processedIcon, processedLabel, processedSeparator)
			}

			if result.Error != nil {
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

	dialog.ShowFileOpen(onChosen, window)
}

func onResetButtonClicked() {
	log.Println("Reset was clicked!")

	// blank all the fields
	analysisTextBS.Set("Select a file...")
	metadataTextBS.Set("Select a file...")
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
