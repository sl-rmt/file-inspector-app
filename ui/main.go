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
	appName = "File Inspector"
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

	// each property is | label | value | in horizontal box
	// file name
	nameAndLabel := container.NewHBox()
	fileName := binding.NewString()
	nameAndLabel.Add(widget.NewLabel("File Name:\t\t"))
	nameAndLabel.Add(widget.NewLabelWithData(fileName))
	props.Add(nameAndLabel)

	// hash
	hashAndLabel := container.NewHBox()
	hash := binding.NewString()
	hashAndLabel.Add(widget.NewLabel("SHA256 Hash:\t"))
	hashAndLabel.Add(widget.NewLabelWithData(hash))
	props.Add(hashAndLabel)

	// file mime type
	typeAndLabel := container.NewHBox()
	fileType := binding.NewString()
	typeAndLabel.Add(widget.NewLabel("File Type:\t\t"))
	typeAndLabel.Add(widget.NewLabelWithData(fileType))
	props.Add(typeAndLabel)

	// file size
	sizeAndLabel := container.NewHBox()
	size := binding.NewString()
	sizeAndLabel.Add(widget.NewLabel("File Size:\t\t"))
	sizeAndLabel.Add(widget.NewLabelWithData(size))
	props.Add(sizeAndLabel)

	analysisHeading := widget.NewLabelWithStyle("File Analysis", fyne.TextAlignCenter, headingStyle)
	props.Add(widget.NewSeparator())
	props.Add(analysisHeading)
	props.Add(widget.NewSeparator())

	// add text for the middle
	analysisText := binding.NewString()
	analysisText.Set("Select a file to analyse...")
	analysisTextBox := widget.NewLabelWithData(analysisText)
	analysisTextBox.Wrapping = fyne.TextWrapBreak
	analysisBox := container.NewScroll(analysisTextBox)

	// add icons in a horizontal box
	icons := container.NewHBox()

	processedLabel := widget.NewLabel("Processed")
	processedIcon := widget.NewIcon(theme.ConfirmIcon())
	processedSeparator := widget.NewSeparator()
	hideIconAndLabel(processedIcon, processedLabel, processedSeparator)
	icons.Add(processedIcon)
	icons.Add(processedLabel)
	icons.Add(processedSeparator)

	errorLabel := widget.NewLabel("Error")
	errorIcon := widget.NewIcon(theme.ErrorIcon())
	errorSeparator := widget.NewSeparator()
	hideIconAndLabel(errorIcon, errorLabel, errorSeparator)
	icons.Add(errorIcon)
	icons.Add(errorLabel)
	icons.Add(errorSeparator)

	dangerLabel := widget.NewLabel("Danger")
	dangerIcon := widget.NewIcon(theme.WarningIcon())
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

		// lock openbutton
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
			err = files.SetFileProperties(f.URI().Path(), fileName, hash, fileType, size, &window)

			if err != nil {
				analysisText.Set(fmt.Sprintf("Error processing file: %q\n", err.Error()))
				errorLabel.Show()
				errorIcon.Show()
			} else {
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

	// set layout to borders
	content := container.NewBorder(props, buttonsAndIcons, nil, nil, analysisBox)

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
