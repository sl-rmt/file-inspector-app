package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func launchErrorDialog(err error, window fyne.Window) {
	d := dialog.NewError(err, window)
	d.Show()
}

func launchInfoDialog(title, message string, window *fyne.Window) {
	d := dialog.NewInformation(title, message, *window)
	d.Show()
}

func launchProcessingDialog(window *fyne.Window) *dialog.CustomDialog {
	// just label above a progress bar
	content := container.NewVBox()
	content.Add(widget.NewLabel("Please wait..."))
	progressBar := widget.NewProgressBarInfinite()
	content.Add(progressBar)
	d := dialog.NewCustomWithoutButtons("Processing", content, *window)

	// start the progress bar and show it
	progressBar.Start()
	d.Show()

	return d
}
