package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	appName = "File-Inspector"

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
	

	// setup styles
	headingStyle := fyne.TextStyle{
		Bold: true,
	}

	// Build file properties section
	// add properties in vertical box
	props := getPropertiesContainer(headingStyle)

	// file name
	fileNameBS = binding.NewString()
	nameAndLabel := getBoundStringAndLabelContainer(fileNameText, fileNameBS)
	props.Add(nameAndLabel)

	// hashBS
	fileHashBS = binding.NewString()
	hashAndLabel := getBoundStringAndLabelContainer(hashLabelText, fileHashBS)
	props.Add(hashAndLabel)

	// file mime type
	fileTypeBS = binding.NewString()
	typeAndLabel := getBoundStringAndLabelContainer(fileTypeText, fileTypeBS)
	props.Add(typeAndLabel)

	// file size
	fileSizeBS = binding.NewString()
	sizeAndLabel := getBoundStringAndLabelContainer(fileSizeText, fileSizeBS)
	props.Add(sizeAndLabel)

	// add file analysis section
	props.Add(widget.NewSeparator())
	props.Add(widget.NewLabelWithStyle(fileAnalysisText, fyne.TextAlignCenter, headingStyle))
	props.Add(widget.NewSeparator())

	// add text for the middle tabs
	metadataTextBS = binding.NewString()
	metadataBox := getScrollContainer(defaultSelectText, metadataTextBS)

	analysisTextBS = binding.NewString()

	analysisBox := getScrollContainer(defaultSelectText, analysisTextBS)

	// add icons in a horizontal box
	icons := container.NewHBox()

	// Processed icon - if successfully processed
	processedIcon, processedLabel, processedSeparator = getIconAndLabel("Processed", true, confirmIconType)
	icons.Add(processedIcon)
	icons.Add(processedLabel)
	icons.Add(processedSeparator)

	// Error icon - if file processing fails
	errorIcon, errorLabel, errorSeparator = getIconAndLabel("Error", true, errorIconType)
	icons.Add(errorIcon)
	icons.Add(errorLabel)
	icons.Add(errorSeparator)

	// Danger icon - only show where we find something suspicious
	dangerIcon, dangerLabel, dangerSeparator = getIconAndLabel("Danger", true, warningIconType)
	icons.Add(dangerIcon)
	icons.Add(dangerLabel)
	icons.Add(dangerSeparator)

	iconSeparator = widget.NewSeparator()

	// add buttons in horizontal box
	buttons := container.NewHBox()
	openButton = widget.NewButtonWithIcon("Select File", theme.FileIcon(), onOpenButtonClicked)

	buttons.Add(openButton)

	buttons.Add(widget.NewButtonWithIcon("Reset", theme.MediaReplayIcon(), onResetButtonClicked))

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

	// TODO set up drag and drop
	//window.SetOnDropped(onFileDroppedin)

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

// TODO get this working
// // this is called if a file is dropped into the UI
// func onFileDroppedin(_ fyne.Position, uris []fyne.URI) {
// 	log.Printf("%d files dropped in\n", len(uris))

// 	onFileChosen(uris[0], nil)
// }
