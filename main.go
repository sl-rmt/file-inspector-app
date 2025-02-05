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

	// ColorRed is the red primary color name
	ColorRed = "red"
	ColorNameForegroundOnWarning fyne.ThemeColorName = "foregroundOnWarning"
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

	// add properties in vertical box
	props := container.NewVBox()
	propsHeading := widget.NewLabelWithStyle("File Properties", fyne.TextAlignCenter, headingStyle)
	props.Add(widget.NewSeparator())
	props.Add(propsHeading)
	props.Add(widget.NewSeparator())

	// file name
	nameAndLabel := container.NewHBox()
	fileName := binding.NewString()
	titleStyle := fyne.TextStyle{
		Bold: true,
	}
	nameAndLabel.Add(widget.NewLabelWithStyle("File Name:\t\t", fyne.TextAlignLeading, titleStyle))
	nameAndLabel.Add(widget.NewLabelWithData(fileName))
	props.Add(nameAndLabel)

	// hash
	hashAndLabel := container.NewHBox()
	hash := binding.NewString()
	hashAndLabel.Add(widget.NewLabelWithStyle("SHA256 Hash:\t", fyne.TextAlignLeading, titleStyle))
	hashAndLabel.Add(widget.NewLabelWithData(hash))
	props.Add(hashAndLabel)

	// file mime type
	typeAndLabel := container.NewHBox()
	fileType := binding.NewString()
	typeAndLabel.Add(widget.NewLabelWithStyle("File Type:\t\t", fyne.TextAlignLeading, titleStyle))
	typeAndLabel.Add(widget.NewLabelWithData(fileType))
	props.Add(typeAndLabel)

	// file size
	sizeAndLabel := container.NewHBox()
	size := binding.NewString()
	sizeAndLabel.Add(widget.NewLabelWithStyle("File Size:\t\t", fyne.TextAlignLeading, titleStyle))
	sizeAndLabel.Add(widget.NewLabelWithData(size))
	props.Add(sizeAndLabel)

	props.Add(widget.NewSeparator())
	props.Add(widget.NewLabelWithStyle("File Analysis", fyne.TextAlignCenter, headingStyle))
	props.Add(widget.NewSeparator())

	// add text for the middle tabs
	metadataText := binding.NewString()
	metadataText.Set("Select a file to analyse...")
	metadataTextBox := widget.NewLabelWithData(metadataText)
	metadataTextBox.Wrapping = fyne.TextWrapBreak
	metadataBox := container.NewScroll(metadataTextBox)

	analysisText := binding.NewString()
	analysisText.Set("Select a file to analyse...")
	analysisTextBox := widget.NewLabelWithData(analysisText)
	analysisTextBox.Wrapping = fyne.TextWrapBreak
	analysisBox := container.NewScroll(analysisTextBox)

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
				analysisText.Set(fmt.Sprintf("Error processing file: %q\n", err.Error()))
				errorLabel.Show()
				errorIcon.Show()
			} else {
				// set the values
				fileName.Set(properties.FileName)
				fileType.Set(properties.FileType)
				hash.Set(properties.Hash)
				size.Set(properties.Size)

				// process the file and show the analysis
				result := files.ProcessFile(f.URI().Path())

				if result.Completed {
					showIconAndLabel(processedIcon, processedLabel, processedSeparator)
				}

				if result.Error != nil {
					launchErrorDialog(result.Error, window)
					analysisText.Set(result.Error.Error())
					showIconAndLabel(errorIcon, errorLabel, errorSeparator)
				}

				log.Println("File processing done")
				metadataText.Set(result.Metadata)
				analysisText.Set(result.Analysis)

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
		analysisText.Set("Select a file...")
		metadataText.Set("Select a file...")
		fileName.Set("")
		fileType.Set("")
		hash.Set("")
		size.Set("")

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
