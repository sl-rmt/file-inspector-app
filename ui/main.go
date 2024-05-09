package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const (
	appName = "File Inspector"
)

func main() {
	// create an app and window instance
	myApp := app.New()
	window := myApp.NewWindow(appName)

	// add properties in vertical box
	props := container.NewVBox()
	propsHeading := widget.NewLabel("File Properties")
	propsHeading.Alignment = fyne.TextAlignCenter
	props.Add(propsHeading)

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

	analysisHeading := widget.NewLabel("File Analysis")
	analysisHeading.Alignment = fyne.TextAlignCenter
	props.Add(analysisHeading)

	// add text for the middle
	// TODO make this scrollable
	
	analysisText := binding.NewString()
	analysisText.Set("Select a file...")
	analysisTextBox := widget.NewLabelWithData(analysisText)
	analysisTextBox.Wrapping = fyne.TextWrapBreak
	analysisBox := container.NewScroll(analysisTextBox)

	// add buttons in horizontal box
	buttons := container.NewHBox()

	// TODO use NewButtonWithIcon
	buttons.Add(widget.NewButton("Select File", func() {
		log.Println("Select file was clicked!")

		onChosen := func(f fyne.URIReadCloser, err error) {
			if err != nil {
				log.Printf("Error from file picker: %s\n", err.Error())
				return
			}
			if f == nil {
				log.Println("Nil result from file picker")
				return
			}

			log.Printf("chosen: %v", f.URI())

			// get and set the file properties
			err = getProps(f.URI().Path(), fileName, hash, fileType, size, &window)

			if err == nil {
				// process the file and show the analysis
				processFile(f.URI().Path(), &window, analysisText)
			}
		}

		dialog.ShowFileOpen(onChosen, window)

	}))

	// TODO use NewButtonWithIcon
	buttons.Add(widget.NewButton("Reset", func() {
		log.Println("Reset was clicked!")
		analysisText.Set("Select a file...")
		fileName.Set("")
		fileType.Set("")
		hash.Set("")
		size.Set("")
	}))

	// set layout to borders
	content := container.NewBorder(props, buttons, nil, nil, analysisBox)

	// set default size
	window.Resize(fyne.NewSize(700, 900))

	// run
	window.SetContent(content)
	window.ShowAndRun()
}

func launchErrorDialog(err error, window *fyne.Window) {
	d := dialog.NewError(err, *window)
	d.Show()
}

func launchInfoDialog(title, message string, window *fyne.Window) {
	d := dialog.NewInformation(title, message, *window)
	d.Show()
}
