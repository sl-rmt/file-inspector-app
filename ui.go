package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	errorIconType   = "error"
	confirmIconType = "confirm"
	warningIconType = "warning"

	filePropertiesText = "File Properties"
	fileNameText       = "File Name:\t\t"
	hashLabelText      = "SHA256 Hash:\t"
	fileTypeText       = "File Type:\t\t"
	fileSizeText       = "File Size:\t\t"
	fileAnalysisText   = "File Analysis"

	metadataTableNumColumns       = 2
	metadataTableFieldColumnID    = 0
	metadataTableFieldColumnWidth = 200
)

func buildUI() *fyne.Container {
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

	// create the metadata table
	metadataTableData = [][]string{
		{"Field", "Value"},
	}

	metadataTable = widget.NewTable(
		func() (int, int) {
			return len(metadataTableData), metadataTableNumColumns
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Metadata")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if o.(*widget.Label) != nil {
				o.(*widget.Label).SetText(metadataTableData[i.Row][i.Col])
			}
		},
	)

	metadataTable.SetColumnWidth(metadataTableFieldColumnID, metadataTableFieldColumnWidth)
	metadataBox := container.NewScroll(metadataTable)

	// add text for the middle tabs

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

	return content
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

// Container box
func getPropertiesContainer(headingStyle fyne.TextStyle) *fyne.Container {
	props := container.NewVBox()
	propsHeading := widget.NewLabelWithStyle(filePropertiesText, fyne.TextAlignCenter, headingStyle)
	props.Add(widget.NewSeparator())
	props.Add(propsHeading)
	props.Add(widget.NewSeparator())
	return props
}

func getBoundStringAndLabelContainer(labelText string, boundString binding.String) *fyne.Container {
	titleStyle := fyne.TextStyle{
		Bold: true,
	}

	typeAndLabel := container.NewHBox()
	typeAndLabel.Add(widget.NewLabelWithStyle(labelText, fyne.TextAlignLeading, titleStyle))
	typeAndLabel.Add(widget.NewLabelWithData(boundString))
	return typeAndLabel
}

func getScrollContainer(labelText string, boundString binding.String) *container.Scroll {
	boundString.Set(labelText)
	textBox := widget.NewLabelWithData(boundString)
	textBox.Wrapping = fyne.TextWrapBreak
	container := container.NewScroll(textBox)
	return container
}

func getIconAndLabel(labelText string, hidden bool, styleType string) (*widget.Icon, *widget.Label, *widget.Separator) {
	// make the label
	label := widget.NewLabel(labelText)
	label.TextStyle.Bold = true

	// make the icon
	var icon *widget.Icon

	// change the default icon type depending on what's requested
	switch styleType {
	case errorIconType:
		icon = widget.NewIcon(theme.NewErrorThemedResource(theme.ErrorIcon()))
	case confirmIconType:
		icon = widget.NewIcon(theme.NewSuccessThemedResource(theme.ConfirmIcon()))
	case warningIconType:
		icon = widget.NewIcon(theme.NewWarningThemedResource(theme.WarningIcon()))
	default:
		// TODO what's a sensible default icon?
		icon = widget.NewIcon(theme.NewPrimaryThemedResource(theme.FileIcon()))
	}

	// make the separator
	separator := widget.NewSeparator()

	// hide them all
	if hidden {
		hideIconAndLabel(icon, label, separator)
	}

	return icon, label, separator
}
