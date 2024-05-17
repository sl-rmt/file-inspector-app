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
	icons.Size()

	fileLabel := widget.NewLabel("File:")
	fileLabel.Hide()
	icons.Add(fileLabel)
	fileIcon := widget.NewFileIcon(nil)
	fileIcon.Hide()
	icons.Add(fileIcon)

	completeLabel := widget.NewLabel("Complete:")
	completeLabel.Hide()
	icons.Add(completeLabel)
	completeIcon := widget.NewIcon(theme.ConfirmIcon())
	completeIcon.Hide()
	icons.Add(completeIcon)

	errorLabel := widget.NewLabel("Error:")
	errorLabel.Hide()
	icons.Add(errorLabel)
	errorIcon := widget.NewIcon(theme.ErrorIcon())
	errorIcon.Hide()
	icons.Add(errorIcon)

	dangerLabel := widget.NewLabel("Danger:")
	dangerLabel.Hide()
	icons.Add(dangerLabel)
	dangerIcon := widget.NewIcon(theme.WarningIcon())
	dangerIcon.Hide()
	icons.Add(dangerIcon)

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
			log.Printf("chosen: %v", f.URI())
			fileIcon.SetURI(f.URI())
			fileLabel.Show()
			fileIcon.Show()
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
					completeLabel.Show()
					completeIcon.Show()
				}

				if result.Error != nil {
					launchErrorDialog(result.Error, window)
					analysisText.Set(result.Error.Error())
					
					errorLabel.Show()
					errorIcon.Show()
				}

				log.Println("File processing done")
				analysisText.Set(result.Analysis)

				if result.Dangerous {
					dangerLabel.Show()
					dangerIcon.Show()
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
		fileIcon.SetURI(nil)
		completeLabel.Hide()
		completeIcon.Hide()
		fileLabel.Hide()
		fileIcon.Hide()
		errorLabel.Hide()
		errorIcon.Hide()
		dangerLabel.Hide()
		dangerIcon.Hide()
	}))

	buttonsAndIcons := container.NewVBox()
	buttonsAndIcons.Add(widget.NewSeparator())
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
