package main

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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

	// add text for the middle in a vertical scroller
	mainText := canvas.NewText("Select a file...", color.White)
	middle := container.NewVScroll(mainText)

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
			processFile(f.URI().String(), &window)
		}

		dialog.ShowFileOpen(onChosen, window)

	}))

	// TODO use NewButtonWithIcon
	buttons.Add(widget.NewButton("Reset", func() {
		log.Println("Reset was clicked!")
		mainText.Text = "Select a file..."
	}))

	// set layout to borders
	content := container.NewBorder(nil, buttons, nil, nil, middle)

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
