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
	// create an instance
	myApp := app.New()
	window := myApp.NewWindow(appName)

	// add properties in vertical box
	props := container.NewVBox()

	fileName := binding.NewString()
	props.Add(widget.NewLabelWithData(fileName))

	hash := binding.NewString()
	props.Add(widget.NewLabelWithData(hash))

	fileType := binding.NewString()
	props.Add(widget.NewLabelWithData(fileType))

	size := binding.NewString()
	props.Add(widget.NewLabelWithData(size))

	// add text for the middle in a vertical scroller
	mainText := binding.NewString()
	middle := container.NewVScroll(widget.NewLabelWithData(mainText))
	mainText.Set("Select a file...")

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
			getProps(f.URI().Path(), fileName, hash, fileType, size, &window)

			// process the file and show the analysis
			processFile(f.URI().Path(), &window, mainText)
		}

		dialog.ShowFileOpen(onChosen, window)

	}))

	// TODO use NewButtonWithIcon
	buttons.Add(widget.NewButton("Reset", func() {
		log.Println("Reset was clicked!")
		mainText.Set("Select a file...")
		fileName.Set("")
		fileType.Set("")
		hash.Set("")
		size.Set("")
	}))

	// set layout to borders
	content := container.NewBorder(props, buttons, nil, nil, middle)

	// set default size
	window.Resize(fyne.NewSize(600, 900))

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
