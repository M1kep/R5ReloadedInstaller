package page

import (
	"R5ReloadedInstaller/internal/app"
	"R5ReloadedInstaller/internal/ui/interaction"
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// checkboxes is a list of all checkboxes
//var checkboxes = []widget.Bool{
//	installSDKToggle,
//	installLatestFLowToggle,
//	cleanScriptsToggle,
//	installLatestR5ScriptsToggle,
//	installPreReleaseSDKToggle,
//}

var (
	// Checkboxes
	installSDKToggle             widget.Bool
	installLatestFLowToggle      widget.Bool
	installLatestR5ScriptsToggle widget.Bool
	installPreReleaseSDKToggle   widget.Bool

	checkboxes = []map[*widget.Bool]string{
		{&installSDKToggle: "Install SDK"},
		{&installLatestFLowToggle: "Install Latest Flow"},
		{&installLatestR5ScriptsToggle: "Install Latest R5 Scripts"},
		{&installPreReleaseSDKToggle: "Install Pre-Release SDK"},
	}

	// checkBoxList is a list of all checkboxes
	checkBoxList []layout.FlexChild
)

type MainPage struct {
	interactionManager *interaction.Manager
	windowManager      app.WindowManager

	aboutPageButton widget.Clickable

	// Buttons
	startButton           widget.Clickable
	troubleShootingButton widget.Clickable
}

func (m *MainPage) ID() string {
	return "Main Page"
}

func (m *MainPage) Layout(gtx layout.Context) layout.Dimensions {
	startButtonLayout := layout.Rigid(
		func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(theme, &m.startButton, "Start")
			return btn.Layout(gtx)
		})

	troubleshootingButtonLayout := layout.Rigid(
		func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(theme, &m.troubleShootingButton, "Troubleshooting")
			return btn.Layout(gtx)
		},
	)

	aboutPageButtonLayout := layout.Rigid(
		func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(theme, &m.aboutPageButton, "About")
			return btn.Layout(gtx)
		},
	)

	spacerLayout := layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout)

	combined := append(checkBoxList, []layout.FlexChild{startButtonLayout, spacerLayout, troubleshootingButtonLayout, spacerLayout, aboutPageButtonLayout, spacerLayout}...)

	return layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceStart,
	}.Layout(gtx,
		combined...,
	)
}

func checkboxHandler(msg string) func() {
	return func() {
		fmt.Println(msg)
	}
}

func (m *MainPage) OnPageLoad(windowManager app.WindowManager) {
	// If the page has already been loaded, don't do anything
	if m.interactionManager != nil {
		return
	}

	m.interactionManager = interaction.NewManager()
	m.windowManager = windowManager
	// Buttons
	m.interactionManager.RegisterClickableHandler(&m.startButton, func() {
		fmt.Println("startButtonClicked")
	})

	m.interactionManager.RegisterClickableHandler(&m.aboutPageButton, func() {
		if m.windowManager != nil {
			if m.windowManager.KnowsPage("About Page") {
				m.windowManager.SetCurrentPageByID("About Page")
			} else {
				aboutPage := &AboutPage{}
				m.windowManager.SetCurrentPage(aboutPage)
			}
		}
	})

	m.interactionManager.RegisterClickableHandler(&m.troubleShootingButton, func() {
		if m.windowManager != nil {
			if m.windowManager.KnowsPage("Troubleshooting Page") {
				m.windowManager.SetCurrentPageByID("Troubleshooting Page")
			} else {
				troubleshootingPage := &TroubleshootingPage{}
				m.windowManager.SetCurrentPage(troubleshootingPage)
			}
		}
	})

	// Checkboxes
	m.interactionManager.RegisterBoolHandler(&installSDKToggle, checkboxHandler("installSDKToggle changed"))
	m.interactionManager.RegisterBoolHandler(&installLatestFLowToggle, checkboxHandler("installLatestFLowToggle changed"))
	m.interactionManager.RegisterBoolHandler(&installLatestR5ScriptsToggle, checkboxHandler("installLatestR5ScriptsToggle changed"))
	m.interactionManager.RegisterBoolHandler(&installPreReleaseSDKToggle, checkboxHandler("installPreReleaseSDKToggle changed"))

	checkBoxList = checkBoxListBuilder(checkboxes)
}

func (m *MainPage) OnPageUnload() {
	return
}

func (m *MainPage) HandleUserInteractions() {
	m.interactionManager.RunHandlers()
}
