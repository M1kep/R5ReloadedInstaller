package app

import "gioui.org/layout"

type Page interface {
	// Name returns the name of the page
	Name() string

	// HandleUserInteractions handles user interactions is called just before the page is rendered
	HandleUserInteractions()

	// Layout draws the page UI components into the provided layout context
	// to be eventually drawn on screen.
	Layout(layout.Context) layout.Dimensions

	// OnPageLoad is called when the page is loaded
	OnPageLoad()

	// OnPageUnload is called when the page is unloaded
	OnPageUnload()
}
