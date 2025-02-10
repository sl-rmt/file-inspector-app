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
)

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
