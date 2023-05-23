package app

import (
	"gioui.org/layout"
)

type Page interface {
	// ID returns the page ID
	ID() string

	// HandleUserInteractions handles user interactions is called just before the page is rendered
	HandleUserInteractions()

	// Layout draws the page UI components into the provided layout context
	// to be eventually drawn on screen.
	Layout(layout.Context) layout.Dimensions

	// OnPageLoad is called when the page is loaded
	OnPageLoad(manager WindowManager)

	// OnPageUnload is called when the page is unloaded
	OnPageUnload()

	// SetWindowManager sets the window manager
	//SetWindowManager(manager WindowManager)
}

type WindowManager interface {
	// SetCurrentPage sets the current page
	SetCurrentPage(page Page)

	// SetCurrentPageByID sets the current page by ID
	SetCurrentPageByID(id string)

	// KnowsPage returns true if the window manager knows the page
	KnowsPage(id string) bool

	// GetCurrentPage returns the current page
	GetCurrentPage() Page

	// PreviousPage goes to the previous page
	PreviousPage()
}
