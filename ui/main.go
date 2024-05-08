package main

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	appName = "File Inspector"
)

func main() {
	// create an instance
	myApp := app.New()
	window := myApp.NewWindow(appName)

	// add button
	// TODO use NewButtonWithIcon
	button := widget.NewButton("Choose file", func() {
		log.Println("Button pressed")
	})

	middle := canvas.NewText("Select a file...", color.White)

	// set layout to borders
	//top := canvas.NewText("Select a file...", color.White)
	content := container.NewBorder(nil, button, nil, nil, middle)

	// set default size
	window.Resize(fyne.NewSize(600, 900))

	// run
	window.SetContent(content)
	window.ShowAndRun()
}
