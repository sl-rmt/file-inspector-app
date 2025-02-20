package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	appName = "File-Inspector"

	windowWidth  = 700
	windowHeight = 900

	filePropertiesText = "File Properties"
	fileNameText       = "File Name:\t\t"
	hashLabelText      = "SHA256 Hash:\t"
	fileTypeText       = "File Type:\t\t"
	fileSizeText       = "File Size:\t\t"
	fileAnalysisText   = "File Analysis"
	defaultSelectText  = "\n\n\n\n\n\t\t\t\t\t\tSelect a file to analyse..."
)

// These all need to be global to allow us to break up some of the functions
var (
	window        fyne.Window
	openButton    *widget.Button
	iconSeparator *widget.Separator

	analysisTextBS binding.String
	fileNameBS     binding.String
	fileTypeBS     binding.String
	fileSizeBS     binding.String
	fileHashBS     binding.String
	metadataTextBS binding.String

	errorLabel     *widget.Label
	errorIcon      *widget.Icon
	errorSeparator *widget.Separator

	processedIcon      *widget.Icon
	processedLabel     *widget.Label
	processedSeparator *widget.Separator

	dangerIcon      *widget.Icon
	dangerLabel     *widget.Label
	dangerSeparator *widget.Separator
)

func main() {
	// create an app and window instance
	myApp := app.New()
	myApp.Settings().SetTheme(&WindowTheme{Theme: theme.DefaultTheme()})
	window = myApp.NewWindow(appName)

	// get UI content
	content := buildUI()

	// TODO set up drag and drop
	//window.SetOnDropped(onFileDroppedin)

	// set default size
	window.Resize(fyne.NewSize(windowWidth, windowHeight))

	// run
	window.SetContent(content)
	window.ShowAndRun()
}
