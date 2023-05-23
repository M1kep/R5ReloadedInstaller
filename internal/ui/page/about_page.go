package page

import (
	"R5ReloadedInstaller/internal/app"
	"R5ReloadedInstaller/internal/ui/interaction"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type AboutPage struct {
	interactionManager *interaction.Manager
	windowManager      app.WindowManager
}

func (a *AboutPage) OnPageUnload() {
	return
}

var (
	backButton widget.Clickable
)

// Name returns the name of the page
func (a *AboutPage) ID() string {
	return "About Page"
}

// OnPageLoad is called when the page is loaded
func (a *AboutPage) OnPageLoad(windowManager app.WindowManager) {
	// Skip if already initialized
	if a.interactionManager != nil {
		return
	}

	a.interactionManager = interaction.NewManager()
	a.windowManager = windowManager

	a.interactionManager.RegisterClickableHandler(&backButton, func() {
		if a.windowManager != nil {
			a.windowManager.SetCurrentPageByID("Main Page")
		}
	})
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
func (a *AboutPage) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, material.H1(theme, "About").Layout)
			},
		),
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, material.Body1(theme, "This is the about page").Layout)
			},
		),
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(theme, &backButton, "Back")
				return btn.Layout(gtx)
			},
		),
	)
}

func (a *AboutPage) HandleUserInteractions() {
	a.interactionManager.RunHandlers()
}
