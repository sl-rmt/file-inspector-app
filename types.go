package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

const (
	windowFontSize = 18
)

type WindowTheme struct {
	fyne.Theme
}

func (c WindowTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return windowFontSize
	}
	return c.Theme.Size(name)
}