package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"file-inspector/files"
)

const (
	appName = "File-Inspector"

	filePropertiesText = "File Properties"
	fileNameText       = "File Name:\t\t"
	hashLabelText      = "SHA256 Hash:\t"
	fileTypeText       = "File Type:\t\t"
	fileSizeText       = "File Size:\t\t"
	fileSAnalysisText = "File Analysis"
	defaultSelectText = "\n\n\t\tSelect a file to analyse..."
)

func main() {
	// create an app and window instance
	myApp := app.New()

	//deprecated: instead "export FYNE_THEME=light"
	//myApp.Settings().SetTheme(theme.LightTheme())

	window := myApp.NewWindow(appName)

	// setup styles
	headingStyle := fyne.TextStyle{
		Bold: true,
	}

	// Build file properties section
	// add properties in vertical box
	props := getPropertiesContainer(headingStyle)

	// file name
	fileNameBS := binding.NewString()
	nameAndLabel := getBoundStringAndLabelContainer(fileNameText, fileNameBS)
	props.Add(nameAndLabel)

	// hashBS
	hashBS := binding.NewString()
	hashAndLabel := getBoundStringAndLabelContainer(hashLabelText, hashBS)
	props.Add(hashAndLabel)

	// file mime type
	fileTypeBS := binding.NewString()
	typeAndLabel := getBoundStringAndLabelContainer(fileTypeText, fileTypeBS)
	props.Add(typeAndLabel)

	// file size
	fileSizeBS := binding.NewString()
	sizeAndLabel := getBoundStringAndLabelContainer(fileSizeText, fileSizeBS)
	props.Add(sizeAndLabel)


	// add file analysis section
	props.Add(widget.NewSeparator())
	props.Add(widget.NewLabelWithStyle(fileSAnalysisText, fyne.TextAlignCenter, headingStyle))
	props.Add(widget.NewSeparator())

	// add text for the middle tabs
	metadataTextBS := binding.NewString()
	metadataBox := getScrollContainer(defaultSelectText, metadataTextBS)

	analysisTextBS := binding.NewString()
	analysisBox := getScrollContainer(defaultSelectText, analysisTextBS)

	// add icons in a horizontal box
	icons := container.NewHBox()

	// Processed icon - if successfully processed
	processedLabel := widget.NewLabel("Processed")
	processedLabel.TextStyle.Bold = true
	processedIcon := widget.NewIcon(theme.NewSuccessThemedResource(theme.ConfirmIcon()))
	processedSeparator := widget.NewSeparator()
	hideIconAndLabel(processedIcon, processedLabel, processedSeparator)
	icons.Add(processedIcon)
	icons.Add(processedLabel)
	icons.Add(processedSeparator)

	// Error icon - if fille processing fails
	errorLabel := widget.NewLabel("Error")
	errorLabel.TextStyle.Bold = true
	// make a red error icon
	errorIcon := widget.NewIcon(theme.NewErrorThemedResource(theme.ErrorIcon()))
	errorSeparator := widget.NewSeparator()
	hideIconAndLabel(errorIcon, errorLabel, errorSeparator)
	icons.Add(errorIcon)
	icons.Add(errorLabel)
	icons.Add(errorSeparator)

	// Danger icon - only show where we find something suspicious
	dangerLabel := widget.NewLabel("Danger")
	dangerLabel.TextStyle.Bold = true
	// make an orange danger icon
	dangerIcon := widget.NewIcon(theme.NewWarningThemedResource(theme.WarningIcon()))
	dangerSeparator := widget.NewSeparator()
	hideIconAndLabel(dangerIcon, dangerLabel, dangerSeparator)
	icons.Add(dangerIcon)
	icons.Add(dangerLabel)
	icons.Add(dangerSeparator)
	iconSeparator := widget.NewSeparator()

	// add buttons in horizontal box
	buttons := container.NewHBox()
	var openButton *widget.Button

	openButton = widget.NewButtonWithIcon("Select File", theme.FileIcon(), func() {
		log.Println("Select file was clicked!")

		// lock open button until we're done processing
		openButton.Disable()

		// TODO make sure the view is reset

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
				hashBS.Set(properties.Hash)
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
				}
			}

			progress.Hide()
			openButton.Enable()
		}

		dialog.ShowFileOpen(onChosen, window)
	})

	buttons.Add(openButton)

	buttons.Add(widget.NewButtonWithIcon("Reset", theme.MediaReplayIcon(), func() {
		log.Println("Reset was clicked!")

		// blank all the fields
		analysisTextBS.Set("Select a file...")
		metadataTextBS.Set("Select a file...")
		fileNameBS.Set("")
		fileTypeBS.Set("")
		hashBS.Set("")
		fileSizeBS.Set("")

		// clear and hide icons
		iconSeparator.Hide()
		hideIconAndLabel(processedIcon, processedLabel, processedSeparator)
		hideIconAndLabel(errorIcon, errorLabel, errorSeparator)
		hideIconAndLabel(dangerIcon, dangerLabel, dangerSeparator)
	}))

	buttonsAndIcons := container.NewVBox()
	iconSeparator.Hide()
	buttonsAndIcons.Add(iconSeparator)
	buttonsAndIcons.Add(icons)
	buttonsAndIcons.Add(widget.NewSeparator())
	buttonsAndIcons.Add(buttons)
	buttonsAndIcons.Add(widget.NewSeparator())

	centreBox := container.NewAppTabs(
		container.NewTabItem("Content", analysisBox),
		container.NewTabItem("Metadata", metadataBox),
	)

	// set layout to borders
	content := container.NewBorder(props, buttonsAndIcons, nil, nil, centreBox)

	// set default size
	window.Resize(fyne.NewSize(700, 900))

	// run
	window.SetContent(content)
	window.ShowAndRun()
}

func hideIconAndLabel(icon *widget.Icon, label *widget.Label, sep *widget.Separator) {
	icon.Hide()
	label.Hide()
	sep.Hide()
}

func showIconAndLabel(icon *widget.Icon, label *widget.Label, sep *widget.Separator) {
	icon.Show()
	label.Show()
	sep.Show()
}
