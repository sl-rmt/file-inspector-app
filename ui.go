package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

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
